package svg

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image/v2/contour"
	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/gocloud/blob/bucket"
	"github.com/sfomuseum/go-flags/flagset"
	"gocloud.dev/blob"
)

// RunOptions is a struct containing configuration details for running an image transformation application.
type RunOptions struct {
	// SourceURI is a `gocloud.dev/blob.Bucket` URI specifying the location where images are read from.
	SourceURI string
	// SourceURI is a `gocloud.dev/blob.Bucket` URI specifying the location where images are written to.
	TargetURI string
	// One or more image URIs (paths) to process.
	Paths []string
}

// Run invokes the application to generate SVG derived from an image's contours using the default flags.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// Run invokes the application to generate SVG derived from an image's contours using a custom `flag.FlagSet` instance.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	paths := fs.Args()

	opts := &RunOptions{
		SourceURI: source_uri,
		TargetURI: source_uri,
		Paths:     paths,
	}

	return RunWithOptions(ctx, opts)
}

// Run invokes the application to generate SVG derived from an image's contours configured using 'opts'.
func RunWithOptions(ctx context.Context, opts *RunOptions) error {

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

			err := deriveSVG(ctx, opts, source_b, target_b, key)

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

func deriveSVG(ctx context.Context, opts *RunOptions, source_b *blob.Bucket, target_b *blob.Bucket, key string) error {

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

	im, _, _, err := decode.DecodeImage(ctx, r)

	if err != nil {
		return fmt.Errorf("Failed to decode %s, %v", key, err)
	}

	old_ext := filepath.Ext(key)
	new_ext := ".svg"

	new_key := key
	new_key = strings.Replace(new_key, old_ext, "", 1)

	new_key = fmt.Sprintf("%s-%03d%s", new_key, n, new_ext)

	wr, err := target_b.NewWriter(ctx, new_key, nil)

	if err != nil {
		return fmt.Errorf("Failed to create new writer for %s, %w", new_key, err)
	}

	err = contour.ContourImageSVG(ctx, wr, im, n, scale)

	if err != nil {
		return fmt.Errorf("Failed to derive SVG from %s, %w", key, err)
	}

	err = wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close writer for %s, %w", new_key, err)
	}

	return nil
}
