// Package resize provides methods for running a base image resizing application
// that can be imported alongside custom `transform.Transformation` and `gocloud.dev/blob`
// packages.
package resize

import (
	"context"
	"flag"
	"fmt"

	"github.com/aaronland/go-image/v2/app/transform"
	"github.com/sfomuseum/go-flags/flagset"
)

// Run invokes the image resizing application using the default flags.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// Run invokes the image resizing application using a custom `flag.FlagSet` instance.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	resize_uri := fmt.Sprintf("resize://?max=%d", max)
	suffix := fmt.Sprintf("-%d", max)

	transformation_uris := []string{
		resize_uri,
	}

	for _, e := range extra_transformations {
		transformation_uris = append(transformation_uris, e)
	}

	// Always do colour profile stuff last

	if profile != "" {
		profile_uri := fmt.Sprintf("%s://", profile)
		transformation_uris = append(transformation_uris, profile_uri)
	}

	paths := fs.Args()

	opts := &transform.RunOptions{
		TransformationURIs: transformation_uris,
		ApplySuffix:        suffix,
		SourceURI:          source_uri,
		TargetURI:          target_uri,
		Rotate:             rotate,
		PreserveExif:       preserve_exif,
		Paths:              paths,
	}

	return transform.RunWithOptions(ctx, opts)
}
