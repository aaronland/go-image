package halftone

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var process string
var scale_factor int
var source_uri string
var target_uri string
var preserve_exif bool
var rotate bool

var extra_transformations multi.MultiCSVString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("halftone")

	fs.IntVar(&scale_factor, "scale-factor", 2, "The scale factor to use for the halftone process.")
	fs.StringVar(&process, "process", "atkinson", "The halftone process to use.")

	fs.StringVar(&source_uri, "source-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are read from.")
	fs.StringVar(&target_uri, "target-uri", "file:///", "A valid gocloud.dev/blob.Bucket URI where images are written to.")
	fs.Var(&extra_transformations, "transformation-uri", "Zero or more additional `transform.Transformation` URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).")

	fs.BoolVar(&preserve_exif, "preserve-exif", false, "Copy EXIF data from source image final target image.")
	fs.BoolVar(&rotate, "rotate", true, `Automatically rotate based on EXIF orientation. This does NOT update any of the original EXIF data with one exception: If the -rotate flag is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1".`)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Apply a halftone (dithering) process to one or more images.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s uri(N) uri(N)\n", os.Args[0])
		fs.PrintDefaults()
	}
	
	return fs
}
