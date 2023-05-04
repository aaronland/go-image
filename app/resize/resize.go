package resize

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/aaronland/go-image/app/transform"
	"github.com/sfomuseum/go-flags/flagset"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

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

	opts := &transform.RunOptions{
		TransformationURIs: transformation_uris,
		ApplySuffix:        suffix,
		SourceURI:          source_uri,
		Logger:             logger,
	}

	paths := fs.Args()

	return transform.RunWithOptions(ctx, opts, paths...)
}
