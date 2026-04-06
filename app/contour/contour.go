package contour

import (
	"context"
	"flag"
	"fmt"

	_ "github.com/aaronland/go-image/v2/contour"

	"github.com/aaronland/go-image/v2/app/transform"
	"github.com/sfomuseum/go-flags/flagset"
)

// Run invokes the image contour-ing application using the default flags.
func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

// Run invokes the image contour-ing application using 'fs' to derive flag values.
func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	contour_uri := fmt.Sprintf("contour://?n=%d&scale=%f", n, scale)
	suffix := fmt.Sprintf("-contour-%d", n)

	transformation_uris := []string{
		contour_uri,
	}

	for _, e := range extra_transformations {
		transformation_uris = append(transformation_uris, e)
	}

	paths := fs.Args()

	opts := &transform.RunOptions{
		TransformationURIs: transformation_uris,
		ApplySuffix:        suffix,
		SourceURI:          source_uri,
		TargetURI:          target_uri,
		Paths:              paths,
	}

	return transform.RunWithOptions(ctx, opts)
}
