package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/aaronland/go-image/rotate"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-heic-exif-extractor/v2"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/dsoprea/go-png-image-structure/v2"
	"github.com/gabriel-vasile/mimetype"
	"github.com/strukturag/libheif-go"
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

	slog.Info("ROTATE", "orientation", orientation)
	return true, im, nil

	ctx := context.Background()
	r_im, err := rotate.RotateImageWithOrientation(ctx, im, orientation)

	if err != nil {
		return false, nil, err
	}

	return true, r_im, nil
}

func ImageFromHEIC(body []byte) (image.Image, error) {

	// https://github.com/spacestation93/heif_howto

	// First decode the HEIC image

	im_ctx, err := libheif.NewContext()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new libheif context, %w", err)
	}

	err = im_ctx.ReadFromMemory(body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read input data, %w", err)
	}

	im_handle, err := im_ctx.GetPrimaryImageHandle()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive primary image handler, %w", err)
	}

	h_im, err := im_handle.DecodeImage(libheif.ColorspaceUndefined, libheif.ChromaUndefined, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image, %w", err)
	}

	im, err := h_im.GetImage()

	if err != nil {
		return nil, fmt.Errorf("Failed to create image.Image, %w", err)
	}

	return im, nil
}

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
	slog.Info("m", "type", mtype)

	switch im_fmt {
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

			// Ideally remove "Orientation" tag here but the following
			// causes a panic below. Maybe there is an easier way?
			// The dsoprea/go-{IMAGE} packages are  still a bit of mystery
			// to me. See also:
			// https://github.com/strukturag/libheif/issues/227

			/*
				ib := exif.NewIfdBuilderFromExistingChain(heic_ifd)

				ifdPath := "IFD0"

				ifd_ib, err := exif.GetOrCreateIbFromRootIb(ib, ifdPath)

				if err != nil {
					return nil, nil, nil, err
				}

				oint, _ := strconv.Atoi("1")	// top left
				oint16 := uint16(oint)

				err = ifd_ib.SetStandardWithName("Orientation", []uint16{oint16})

				if err != nil {
					return nil, nil, nil, err
				}

				slog.Info("SET ORIENTATION TO 1")

				var im_buf bytes.Buffer
				im_wr := bufio.NewWriter(&im_buf)

				err = jpeg.Encode(im_wr, im, nil)

				if err != nil {
					return nil, nil, nil, err
				}

				im_wr.Flush()

				jpeg_jmp := jpegstructure.NewJpegMediaParser()
				jpeg_mc, err := jpeg_jmp.ParseBytes(im_buf.Bytes())

				if err != nil {
					return nil, nil, nil, err
				}

				jpeg_ib := exif.NewIfdBuilderFromExistingChain(ifd)
				jpeg_sl := jpeg_mc.(*jpegstructure.SegmentList)

				err = jpeg_sl.SetExif(jpeg_ib)

				if err != nil {
					return nil, nil, nil, err
				}

				var jpeg_buf bytes.Buffer
				jpeg_wr := bufio.NewWriter(&jpeg_buf)

				err = jpeg_sl.Write(jpeg_wr)

				if err != nil {
					return nil, nil, nil, err
				}

				jpeg_wr.Flush()

				buf_mp, err := jpeg_jmp.ParseBytes(jpeg_buf.Bytes())

				if err != nil {
					return nil, nil, nil, err
				}

				jpg_ifd, _, err := buf_mp.Exif()

				if err != nil {
					slog.Warn("Failed to derive EXIF from update JPEG ifd", "error", err)
				} else {
					ifd = jpg_ifd
				}

			*/

		default:
			return nil, nil, nil, fmt.Errorf("Unsupported media type")
		}
	}

	return im, mtype, ifd, nil
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

		im, mtype, ifd, err := DecodeImage(body)

		if err != nil {
			log.Fatal(err)
		}

		// Rotate image if necessary

		rotated, r_im, err := RotateFromOrientation(im, mtype, ifd)

		if err != nil {
			log.Fatal(err)
		}

		im = r_im

		// Write new image w/ EXIF

		new_path := fmt.Sprintf("%s.jpg", path)
		wr, err := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		if err != nil {
			log.Fatal(err)
		}

		err = WriteJpeg(wr, im, ifd, rotated)

		err = wr.Close()

		if err != nil {
			log.Fatal(err)
		}

		log.Println("WROTE", new_path)
	}
}

func WriteJpeg(wr io.Writer, im image.Image, ifd *exif.Ifd, rotated bool) error {

	jpeg_opts := &jpeg.Options{
		Quality: 100,
	}

	if ifd == nil {
		return jpeg.Encode(wr, im, jpeg_opts)
	}

	// Do EXIF dance

	var im_buf bytes.Buffer
	im_wr := bufio.NewWriter(&im_buf)

	err := jpeg.Encode(im_wr, im, jpeg_opts)

	if err != nil {
		return err
	}

	im_wr.Flush()

	// Write EXIF back to JPEG

	jmp := jpegstructure.NewJpegMediaParser()

	mp, err := jmp.ParseBytes(im_buf.Bytes())

	if err != nil {
		return err
	}

	ib := exif.NewIfdBuilderFromExistingChain(ifd)
	sl := mp.(*jpegstructure.SegmentList)

	if rotated {

		ifdPath := "IFD0"

		ifd_ib, err := exif.GetOrCreateIbFromRootIb(ib, ifdPath)

		if err != nil {
			return err
		}

		oint, _ := strconv.Atoi("1") // top left

		oint16 := uint16(oint)

		err = ifd_ib.SetStandardWithName("Orientation", []uint16{oint16})

		if err != nil {
			log.Fatal(err)
		}

	}

	err = sl.SetExif(ib)

	if err != nil {
		return err
	}

	err = sl.Write(wr)

	if err != nil {
		return err
	}

	return nil
}
