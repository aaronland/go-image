package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aaronland/go-image/app/transform"
	_ "github.com/aaronland/go-image/colour"
	_ "github.com/aaronland/go-image/resize"
	"github.com/sfomuseum/go-flags/flagset"
	_ "gocloud.dev/blob/fileblob"
)

func main() {

	var max int
	var profile string

	ctx := context.Background()
	logger := log.Default()

	fs := flagset.NewFlagSet("resize")

	fs.IntVar(&max, "max", 0, "The maximum dimension of the resized image")
	fs.StringVar(&profile, "profile", "", "...")

	flagset.Parse(fs)

	resize_uri := fmt.Sprintf("resize://?max=%d", max)
	suffix := fmt.Sprintf("-%d", max)

	transformation_uris := []string{
		resize_uri,
	}

	if profile != "" {
		profile_uri := fmt.Sprintf("%s://", profile)
		transformation_uris = append(transformation_uris, profile_uri)
	}

	paths := fs.Args()

	opts := &transform.RunOptions{
		TransformationURIs: transformation_uris,
		ApplySuffix:        suffix,
		BucketURI:          "file:///",
		Logger:             logger,
	}

	err := transform.RunWithOptions(ctx, opts, paths...)

	if err != nil {
		logger.Fatalf("Failed to transform images, %v", err)
	}
}
