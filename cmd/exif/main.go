package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	_ "log/slog"
	"os"

	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-image/encode"
	"github.com/aaronland/go-image/exif"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	for _, path := range flag.Args() {

		r, err := os.Open(path)

		if err != nil {
			log.Fatalf("Failed to open %s for reading, %v", path, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			log.Fatalf("Failed to read %s, %v", path, err)
		}

		// Open with EXIF

		im, _, ifd, err := decode.DecodeImage(ctx, body)

		if err != nil {
			log.Fatal(err)
		}

		// Update EXIF if necessary

		ib, err := exif.NewIfdBuilderWithOrientation(ifd, "1")

		// Write new image w/ EXIF

		new_path := fmt.Sprintf("%s.jpg", path)
		wr, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		err = encode.EncodeJPEG(wr, im, ib, nil)

		err = wr.Close()

		if err != nil {
			log.Fatal(err)
		}

		log.Println("WROTE", new_path)
	}
}
