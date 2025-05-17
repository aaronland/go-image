package transform

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
)

// RunOptions is a struct containing configuration details for running an image transformation application.
type RunOptions struct {
	// TranformationURIs is one or more `transform.Tranformation` URIs used to apply transformations to an image.
	TransformationURIs []string
	// SourceURI is a `gocloud.dev/blob.Bucket` URI specifying the location where images are read from.
	SourceURI string
	// SourceURI is a `gocloud.dev/blob.Bucket` URI specifying the location where images are written to.
	TargetURI string
	// ApplySuffix is an optional suffix to apply to the final image filename.
	ApplySuffix string
	// ImageFormat is an optional image format used to encode the final image.
	ImageFormat string
	// Copy EXIF data from source image to the final image. This does NOT update any of the original EXIF data with one exception:
	// If the `Rotate` property is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1".
	PreserveExif bool
	// Automatically rotate based on EXIF orientation.
	Rotate bool
	// One or more image URIs (paths) to transform.
	Paths []string
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	paths := fs.Args()

	opts := &RunOptions{
		TransformationURIs: transformation_uris,
		SourceURI:          source_uri,
		TargetURI:          source_uri,
		ApplySuffix:        apply_suffix,
		ImageFormat:        image_format,
		Rotate:             rotate,
		PreserveExif:       preserve_exif,
		Paths:              paths,
	}

	return opts, nil
}
