package transform

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-image/encode"
	"github.com/aaronland/go-image/transform"
	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/sfomuseum/go-flags/flagset"
	"gocloud.dev/blob"
)

type RunOptions struct {
	TransformationURIs []string
	SourceURI          string
	TargetURI          string
	ApplySuffix        string
	ImageFormat        string
	Logger             *log.Logger
}

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	opts := &RunOptions{
		TransformationURIs: transformation_uris,
		SourceURI:          source_uri,
		TargetURI:          source_uri,
		ApplySuffix:        apply_suffix,
		ImageFormat:        image_format,
		Logger:             logger,
	}

	paths := fs.Args()

	return RunWithOptions(ctx, opts, paths...)
}

func RunWithOptions(ctx context.Context, opts *RunOptions, paths ...string) error {

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

	for _, key := range paths {

		go func(key string) {

			err := applyTransformation(ctx, opts, tr, source_b, target_b, key)

			if err != nil {
				err_ch <- err
			}

			done_ch <- true
		}(key)
	}

	remaining := len(paths)

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

	r, err := bucket.NewReadSeekCloser(ctx, source_b, key, nil)

	if err != nil {
		return fmt.Errorf("Failed to open %s for reading, %v", key, err)
	}

	defer r.Close()

	dec, err := decode.NewDecoder(ctx, key)

	if err != nil {
		return fmt.Errorf("Failed to create decoder for %s, %w", key, err)
	}

	im, im_format, err := dec.Decode(ctx, r)

	if err != nil {
		return fmt.Errorf("Failed to decode %s, %v", key, err)
	}

	new_im, err := tr.Transform(ctx, im)

	if err != nil {
		return fmt.Errorf("Failed to transform %s, %v", key, err)
	}

	new_key := key
	new_ext := filepath.Ext(key)

	if opts.ImageFormat != "" && opts.ImageFormat != im_format {

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

	enc, err := encode.NewEncoder(ctx, new_key)

	if err != nil {
		return fmt.Errorf("Failed to create new encoder, %w", err)
	}

	err = enc.Encode(ctx, wr, new_im)

	if err != nil {
		return fmt.Errorf("Failed to encode %s, %w", new_key, err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close writer for %s, %v", new_key, err)
	}

	return nil
}
