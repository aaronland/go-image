package contour

import (
	"flag"
	"fmt"
	"os"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var source_uri string
var target_uri string

var n int
var scale float64

var preserve_exif bool
var rotate bool

var extra_transformations multi.MultiCSVString

// DefaultFlagSet returns a `flag.FlagSet` instance configured with the default flags
// for running an image contour-ing application.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("contour")

	fs.IntVar(&n, "n", 12, "The number of iterations used to generate contours.")
	fs.Float64Var(&scale, "scale", 1.0, "The scale of the final output relative to the input image.")

	fs.StringVar(&source_uri, "source-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are read from.")
	fs.StringVar(&target_uri, "target-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are written to.")

	fs.BoolVar(&preserve_exif, "preserve-exif", false, "Copy EXIF data from source image final target image.")
	fs.BoolVar(&rotate, "rotate", true, `Automatically rotate based on EXIF orientation. This does NOT update any of the original EXIF data with one exception: If the -rotate flag is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1".`)

	fs.Var(&extra_transformations, "transformation-uri", "Zero or more additional `transform.Transformation` URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Apply a \"contouring\" process to one or more images.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}

	return fs
}
