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
	"github.com/whosonfirst/go-ioutil"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	tr, err := transform.NewMultiTransformationWithURIs(ctx, transformation_uris...)

	if err != nil {
		return fmt.Errorf("Failed to create transformation, %w", err)
	}

	b, err := bucket.OpenBucket(ctx, bucket_uri)

	if err != nil {
		return fmt.Errorf("Failed to open bucket, %w", err)
	}

	defer b.Close()

	paths := fs.Args()

	for _, key := range paths {

		r, err := b.NewReader(ctx, key, nil)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %v", key, err)
		}

		rs, err := ioutil.NewReadSeekCloser(r)

		if err != nil {
			return fmt.Errorf("Failed to create read seek closer for %s, %w", key, err)
		}

		defer rs.Close()

		dec, err := decode.NewDecoder(ctx, key)

		if err != nil {
			return fmt.Errorf("Failed to create decoder for %s, %w", key, err)
		}

		im, im_format, err := dec.Decode(ctx, rs)

		if err != nil {
			return fmt.Errorf("Failed to decode %s, %v", key, err)
		}

		new_im, err := tr.Transform(ctx, im)

		if err != nil {
			return fmt.Errorf("Failed to transform %s, %v", key, err)
		}

		new_key := key
		new_ext := filepath.Ext(key)

		if image_format != "" && image_format != im_format {

			old_ext := new_ext
			new_ext = fmt.Sprintf(".%s", image_format)

			new_key = strings.Replace(new_key, old_ext, new_ext, 1)
		}

		if apply_suffix != "" {

			key_root := filepath.Dir(new_key)
			key_name := filepath.Base(new_key)
			key_ext := filepath.Ext(new_key)

			new_keyname := strings.Replace(key_name, key_ext, "", 1)
			new_keyname = fmt.Sprintf("%s%s%s", new_keyname, apply_suffix, key_ext)

			new_key = filepath.Join(key_root, new_keyname)
		}

		wr, err := b.NewWriter(ctx, new_key, nil)

		if err != nil {
			return fmt.Errorf("Failed to create new writer for %s, %v", new_key, err)
		}

		if image_format == "" {
			image_format = im_format
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
