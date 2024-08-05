package main

import (
	"flag"
	"os"
	"log"

	"github.com/mandykoh/prism/meta/autometa"
)

func main() {

	flag.Parse()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer r.Close()

		md, _, err := autometa.Load(r)

		if err != nil {
			log.Fatal(err)
		}
		

		pr, err := md.ICCProfile()

		if err != nil {
			log.Fatal(err)
		}

		d, err := pr.Description()
		log.Println(path, d, err)
	}
}
