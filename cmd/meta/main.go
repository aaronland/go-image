package main

import (
	"flag"
	"log"
	"os"

	"github.com/mandykoh/prism/meta/autometa"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			panic(err)
		}

		defer r.Close()

		md, _, err := autometa.Load(r)

		if err != nil {
			panic(err)
		}

		profile, err := md.ICCProfile()

		if err != nil {
			panic(err)
		}

		if profile != nil {

			// also: https://pkg.go.dev/github.com/mandykoh/prism@v0.35.2/meta/pngmeta

			description, err := profile.Description()

			if err != nil {
				panic(err)
			}

			log.Println(path, description)
		}

		r.Seek(0, 0)

		exif.RegisterParsers(mknote.All...)

		x, err := exif.Decode(r)

		if err != nil {
			log.Println("POO")
			// panic(err)
		}

		sp, _ := x.Get(exif.ColorSpace)

		val, _ := sp.StringVal()

		log.Println("SP", val)
	}

}
