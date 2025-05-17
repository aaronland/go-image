package resize

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var max int
var profile string
var source_uri string
var target_uri string
var preserve_exif bool

var extra_transformations multi.MultiCSVString

// DefaultFlagSet returns a `flag.FlagSet` instance configured with the default flags
// for running an image resizing application.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("resize")

	fs.IntVar(&max, "max", 0, "The maximum dimension of the resized image")
	fs.StringVar(&profile, "profile", "", "An optional colour profile to apply to the resized image. Valid options are: adobergb, displayp3.")
	fs.BoolVar(&preserve_exif, "preserve-exif", true, "Copy EXIF data from source image final target image.")

	fs.StringVar(&source_uri, "source-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are read from.")
	fs.StringVar(&target_uri, "target-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are written to.")
	fs.Var(&extra_transformations, "transformation-uri", "Zero or more additional `transform.Transformation` URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).")

	return fs
}
