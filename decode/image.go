package decode

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-heic-exif-extractor/v2"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/dsoprea/go-png-image-structure/v2"
	"github.com/dsoprea/go-tiff-image-structure/v2"
	"github.com/gabriel-vasile/mimetype"
	_ "golang.org/x/image/tiff"
)

func DecodeImage(body []byte) (image.Image, *mimetype.MIME, *exif.Ifd, error) {

	var ifd *exif.Ifd
	var im image.Image

	br := bytes.NewReader(body)

	im, im_fmt, err := image.Decode(br)

	if err != nil {
		// Check error here...
		slog.Warn("Failed to decode image natively", "error", err)
	}

	mtype := mimetype.Detect(body)

	switch im_fmt {
	case "gif":
		// pass
	case "jpeg":

		jmp := jpegstructure.NewJpegMediaParser()
		mc, err := jmp.ParseBytes(body)

		if err != nil {
			return nil, nil, nil, err
		}

		jpg_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Warn("Failed to derive EXIF", "error", err)
		} else {
			ifd = jpg_ifd
		}

	case "png":

		mp := pngstructure.NewPngMediaParser()

		mc, err := mp.ParseBytes(body)

		if err != nil {
			return nil, nil, nil, err
		}

		png_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Warn("Failed to derive EXIF", "error", err)
		} else {
			ifd = png_ifd
		}

	case "tiff":

		mp := tiffstructure.NewTiffMediaParser()

		mc, err := mp.ParseBytes(body)

		if err != nil {
			return nil, nil, nil, err
		}

		tiff_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Warn("Failed to derive EXIF", "error", err)
		} else {
			ifd = tiff_ifd
		}

	default:

		switch mtype.String() {
		case "image/heic":

			heic_im, err := ImageFromHEIC(body)

			if err != nil {
				return nil, nil, nil, err
			}

			im = heic_im

			mp := heicexif.NewHeicExifMediaParser()
			mc, err := mp.ParseBytes(body)

			if err != nil {
				return nil, nil, nil, err
			}

			heic_ifd, _, err := mc.Exif()

			if err != nil {
				slog.Warn("Failed to derive EXIF", "error", err)
			} else {
				ifd = heic_ifd
			}

			// Note: We are NOT removing or updating the Orientation tag
			// (which is assigned but incorrect) in libheif because I can
			// not figure out hwo to do that using the dsoprea packages
			// without causing everything to panic later in the code. Instead
			// we are accounting for this in RotateFromOrientation
			// https://github.com/strukturag/libheif/issues/227

		default:
			return nil, nil, nil, fmt.Errorf("Unsupported media type")
		}
	}

	return im, mtype, ifd, nil
}
