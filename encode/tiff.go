package encode

import (
	"github.com/dsoprea/go-exif/v3"
	"image"
	"io"
	"log/slog"

	"golang.org/x/image/tiff"
)

func EncodeTIFF(wr io.Writer, im image.Image, ib *exif.IfdBuilder, tiff_opts *tiff.Options) error {

	if ib == nil {
		return tiff.Encode(wr, im, tiff_opts)
	}

	slog.Warn("WriteTIFF method does not support writing EXIF data.")
	return tiff.Encode(wr, im, tiff_opts)
}
