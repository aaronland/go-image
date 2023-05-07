package transform

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var transformation_uris multi.MultiCSVString
var source_uri string
var target_uri string
var apply_suffix string
var image_format string

// DefaultFlagSet returns a `flag.FlagSet` instance configured with the default flags
// for running an image transformation application.
func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("transform")
	fs.Var(&transformation_uris, "transformation-uri", "One or more additional `transform.Transformation` URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).")

	fs.StringVar(&source_uri, "source-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are read from.")
	fs.StringVar(&target_uri, "target-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are written to.")
	fs.StringVar(&apply_suffix, "apply-suffix", "", "An optional suffix to apply to the final image filename.")
	fs.StringVar(&image_format, "format", "", "An optional image format used to encode the final image.")

	return fs
}
