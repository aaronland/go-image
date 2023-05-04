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
)

type RunOptions struct {
	TransformationURIs []string
	BucketURI          string
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
		BucketURI:          bucket_uri,
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

	b, err := bucket.OpenBucket(ctx, opts.BucketURI)

	if err != nil {
		return fmt.Errorf("Failed to open bucket, %w", err)
	}

	defer b.Close()

	for _, key := range paths {

		if opts.BucketURI == "file:///" {

			abs_key, err := filepath.Abs(key)

			if err != nil {
				return fmt.Errorf("Failed to derive absolute path for %s, %w", key, err)
			}

			key = abs_key
		}

		r, err := bucket.NewReadSeekCloser(ctx, b, key, nil)

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

		wr, err := b.NewWriter(ctx, new_key, nil)

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
	}

	return nil
}
