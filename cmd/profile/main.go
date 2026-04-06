package main

import (
	"flag"
	"fmt"
	_ "image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {

	flag.Parse()

	exif.RegisterParsers(mknote.All...)

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()

		x, err := exif.Decode(r)

		if err != nil {
			log.Fatal(err)
		}

		v, err := x.Get(exif.ColorSpace)

		if err != nil {
			log.Fatal(err)
		}

		str_v := v.String()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s '%s'\n", path, str_v)
	}
}
