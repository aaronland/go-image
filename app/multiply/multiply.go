package multiply

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

	images := fs.Args()
	count := len(images)
	
	if count < 2{
		return fmt.Errorf("Insufficient images")
	}

	current := images[0]

	for i := 1; i < count; i++ {

		paths := []string{
			current,
			images[i],
		}

		suffix = i
		multiply_uri := fmt.Sprintf("multiply://", other)

		transformation_uris := []string{
			multiply_uri,
		}
		
		opts := &transform.RunOptions{
			TransformationURIs: transformation_uris,
			ApplySuffix:        suffix,
			SourceURI:          source_uri,
			TargetURI:          target_uri,
			Rotate:             rotate,
			PreserveExif:       preserve_exif,
			Paths:              paths,
		}
		
		if format != "" {
			opts.ImageFormat = format
		}
		
		return transform.RunWithOptions(ctx, opts)
	}

	// colour profile stuff
}
