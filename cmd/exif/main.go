package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-image/encode"
	"github.com/aaronland/go-image/rotate"
	"github.com/dsoprea/go-exif/v3"
	"github.com/gabriel-vasile/mimetype"
)

func RotateFromOrientation(im image.Image, mtype *mimetype.MIME, ifd *exif.Ifd) (bool, image.Image, error) {

	if ifd == nil {
		return false, im, nil
	}

	// Ignore EXIF Orientation tags in libheif, kthxbye...
	// https://github.com/strukturag/libheif/issues/227

	if mtype.String() == "image/heic" {
		return true, im, nil
	}

	results, err := ifd.FindTagWithName("Orientation")

	if err != nil {
		return false, nil, err
	}

	ite := results[0]
	orientation, err := ite.FormatFirst()

	if err != nil {
		return false, nil, err
	}

	slog.Info("rotate", "orientation", orientation)

	// Rotate

	if orientation == "1" {
		return false, im, nil
	}

	ctx := context.Background()
	r_im, err := rotate.RotateImageWithOrientation(ctx, im, orientation)

	if err != nil {
		return false, nil, err
	}

	return true, r_im, nil
}

func main() {

	flag.Parse()

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

		im, mtype, ifd, err := decode.DecodeImage(body)

		if err != nil {
			log.Fatal(err)
		}

		// Rotate image if necessary

		rotated, r_im, err := RotateFromOrientation(im, mtype, ifd)

		if err != nil {
			log.Fatal(err)
		}

		im = r_im

		// Update EXIF if necessary

		ib, err := newIfdBuilder(ifd, rotated)

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

func newIfdBuilder(ifd *exif.Ifd, rotated bool) (*exif.IfdBuilder, error) {

	var ib *exif.IfdBuilder

	if ifd != nil {

		ib = exif.NewIfdBuilderFromExistingChain(ifd)

		if rotated {

			ifdPath := "IFD0"

			ifd_ib, err := exif.GetOrCreateIbFromRootIb(ib, ifdPath)

			if err != nil {
				return nil, err
			}

			oint, _ := strconv.Atoi("1") // top left
			oint16 := uint16(oint)

			err = ifd_ib.SetStandardWithName("Orientation", []uint16{oint16})

			if err != nil {
				return nil, err
			}
		}
	}

	return ib, nil
}
