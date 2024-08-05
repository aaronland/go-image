package main

import (
	"flag"
	"log"
	"os"

	"github.com/aaronland/go-image/colour"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()

		sp, err := colour.Colorspace(r)

		log.Println(path, sp, err)
	}
}
