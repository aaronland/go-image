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

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("transform")
	fs.Var(&transformation_uris, "transformation-uri", "")

	fs.StringVar(&source_uri, "source-uri", "file:///", "")
	fs.StringVar(&target_uri, "target-uri", "file:///", "")
	fs.StringVar(&apply_suffix, "apply-suffix", "", "")
	fs.StringVar(&image_format, "format", "", "")

	return fs
}
