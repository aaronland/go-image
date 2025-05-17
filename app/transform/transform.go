// Package transform provides methods for running a base image transformation application
// that can be imported alongside custom `transform.Transformation` and `gocloud.dev/blob`
// packages.
package transform

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-image/v2/encode"
	aa_exif "github.com/aaronland/go-image/v2/exif"
	"github.com/aaronland/go-image/v2/transform"
	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/dsoprea/go-exif/v3"
	"github.com/gabriel-vasile/mimetype"
	"gocloud.dev/blob"
)

// Run invokes the image transformation application using the default flags.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// Run invokes the image transformation application using a custom `flag.FlagSet` instance.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

// Run invokes the image transformation application configured using 'opts'.
func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	tr, err := transform.NewMultiTransformationWithURIs(ctx, opts.TransformationURIs...)

	if err != nil {
		return fmt.Errorf("Failed to create transformation, %w", err)
	}

	source_b, err := bucket.OpenBucket(ctx, opts.SourceURI)

	if err != nil {
		return fmt.Errorf("Failed to open source, %w", err)
	}

	defer source_b.Close()

	target_b, err := bucket.OpenBucket(ctx, opts.TargetURI)

	if err != nil {
		return fmt.Errorf("Failed to open target, %w", err)
	}

	defer target_b.Close()

	done_ch := make(chan bool)
	err_ch := make(chan error)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, key := range opts.Paths {

		go func(key string) {

			err := applyTransformation(ctx, opts, tr, source_b, target_b, key)

			if err != nil {
				err_ch <- err
			}

			done_ch <- true
		}(key)
	}

	remaining := len(opts.Paths)

	for remaining > 0 {
		select {
		case <-ctx.Done():
			return nil
		case <-done_ch:
			remaining -= 1
		case err := <-err_ch:
			return err
		}
	}

	return nil
}

func applyTransformation(ctx context.Context, opts *RunOptions, tr transform.Transformation, source_b *blob.Bucket, target_b *blob.Bucket, key string) error {

	if opts.SourceURI == "file:///" {

		abs_key, err := filepath.Abs(key)

		if err != nil {
			return fmt.Errorf("Failed to derive absolute path for %s, %w", key, err)
		}

		key = abs_key
	}

	im_r, err := bucket.NewReadSeekCloser(ctx, source_b, key, nil)

	if err != nil {
		return fmt.Errorf("Failed to open %s for reading, %v", key, err)
	}

	defer im_r.Close()

	decode_opts := &decode.DecodeImageOptions{
		Rotate: opts.Rotate,
	}

	im, im_fmt, ifd, err := decode.DecodeImageWithOptions(ctx, im_r, decode_opts)

	if err != nil {
		return fmt.Errorf("Failed to decode %s, %v", key, err)
	}

	new_im, err := tr.Transform(ctx, im)

	if err != nil {
		return fmt.Errorf("Failed to transform %s, %v", key, err)
	}

	var ib *exif.IfdBuilder

	if opts.PreserveExif {

		new_ib, err := aa_exif.NewIfdBuilderWithOrientation(ifd, "1")

		if err != nil {
			return fmt.Errorf("Failed to create new IFD builder, %w", err)
		}

		ib = new_ib
	}

	new_key := key
	new_ext := filepath.Ext(key)

	if opts.ImageFormat != "" && opts.ImageFormat != im_fmt {

		old_ext := new_ext
		new_ext = fmt.Sprintf(".%s", opts.ImageFormat)

		new_key = strings.Replace(new_key, old_ext, new_ext, 1)
	}

	if opts.ApplySuffix != "" {

		key_root := filepath.Dir(new_key)
		key_name := filepath.Base(new_key)
		key_ext := filepath.Ext(new_key)

		new_keyname := strings.Replace(key_name, key_ext, "", 1)
		new_keyname = fmt.Sprintf("%s%s%s", new_keyname, opts.ApplySuffix, key_ext)

		new_key = filepath.Join(key_root, new_keyname)
	}

	wr, err := target_b.NewWriter(ctx, new_key, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new writer for %s, %v", new_key, err)
	}

	if opts.ImageFormat == "" {

		_, err := im_r.Seek(0, 0)

		if err != nil {
			return fmt.Errorf("Failed to rewind image reader to determine filetype, %w", err)
		}

		mtype, err := mimetype.DetectReader(im_r)

		if err != nil {
			return fmt.Errorf("Failed to determine image from image reader, %w", err)
		}

		opts.ImageFormat = mtype.String()
	}

	switch opts.ImageFormat {
	case "jpg", "jpeg", "image/jpeg":
		err = encode.EncodeJPEG(ctx, wr, new_im, ib, nil)
	case "png", "image/png":
		err = encode.EncodePNG(ctx, wr, new_im, ib)
	case "tiff", "image/tiff":
		err = encode.EncodeTIFF(ctx, wr, new_im, ib, nil)
	default:
		return fmt.Errorf("Unsupported filetype (%s)", opts.ImageFormat)
	}

	if err != nil {
		return fmt.Errorf("Failed to encode %s, %w", new_key, err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close writer for %s, %v", new_key, err)
	}

	return nil
}
