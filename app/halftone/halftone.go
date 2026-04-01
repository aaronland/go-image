package halftone

import (
	"context"
	"flag"
	"fmt"

	_ "github.com/aaronland/go-image/v2/halftone"
	
	"github.com/aaronland/go-image/v2/app/transform"
	"github.com/sfomuseum/go-flags/flagset"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	flagset.Parse(fs)

	resize_uri := fmt.Sprintf("halftone://?process=%s&scale-factor=%d", process, scale_factor)
	suffix := fmt.Sprintf("-%s-%d", process, scale_factor)

	transformation_uris := []string{
		resize_uri,
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
		Rotate: rotate,
		PreserveExif: preserve_exif,
		Paths: paths,
	}

	return transform.RunWithOptions(ctx, opts)
}
