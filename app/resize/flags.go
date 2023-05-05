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

var extra_transformations multi.MultiCSVString

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("resize")

	fs.IntVar(&max, "max", 0, "The maximum dimension of the resized image")
	fs.StringVar(&profile, "profile", "", "...")

	fs.StringVar(&source_uri, "source-uri", "file:///", "")
	fs.StringVar(&target_uri, "target-uri", "file:///", "")
	fs.Var(&extra_transformations, "transformation-uri", "")

	return fs
}
