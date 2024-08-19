package decode

import (
	"bytes"
	"bufio"
	"context"
	"image"
	"image/jpeg"
	"io"

	"github.com/jdeng/goheif"
	"github.com/aaronland/go-image/colour"
)

// HEICDecoder is a struct that implements the `Decoder` interface for
// decoding HEIC image.
type HEICDecoder struct {
	Decoder
}

func init() {
	ctx := context.Background()
	RegisterDecoder(ctx, NewHEICDecoder, "heic")
}

// NewHEICDecoder returns a new `HEICDecoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.heic
func NewHEICDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &HEICDecoder{}
	return e, nil
}

// Decode will decode the body of 'r' in to an `image.Image` instance using the `image/heic` package.
func (e *HEICDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {

	exif, err := goheif.ExtractExif(r)
	
	if err != nil {
		return nil, "", err
	}

	im, err := goheif.Decode(r)
	
	if err != nil {
		return nil, "", err
	}

	im = colour.ToDisplayP3(im)
	
	var buf bytes.Buffer
	wr := bufio.NewWriter(buf)
	
	exif_wr, _ := newWriterExif(wr, exif)
	
	err = jpeg.Encode(exif_wr, img, nil)
	
	if err != nil {
		return nil, "", err
	}

	jpeg_r := bytes.NewReader(buf.Bytes())
	return image.Decode(r)	
}

// Skip Writer for exif writing
type writerSkipper struct {
	w           io.Writer
	bytesToSkip int
}

func (w *writerSkipper) Write(data []byte) (int, error) {
	if w.bytesToSkip <= 0 {
		return w.w.Write(data)
	}

	if dataLen := len(data); dataLen < w.bytesToSkip {
		w.bytesToSkip -= dataLen
		return dataLen, nil
	}

	if n, err := w.w.Write(data[w.bytesToSkip:]); err == nil {
		n += w.bytesToSkip
		w.bytesToSkip = 0
		return n, nil
	} else {
		return n, err
	}
}

func newWriterExif(w io.Writer, exif []byte) (io.Writer, error) {
	writer := &writerSkipper{w, 2}
	soi := []byte{0xff, 0xd8}
	if _, err := w.Write(soi); err != nil {
		return nil, err
	}

	if exif != nil {
		app1Marker := 0xe1
		markerlen := 2 + len(exif)
		marker := []byte{0xff, uint8(app1Marker), uint8(markerlen >> 8), uint8(markerlen & 0xff)}
		if _, err := w.Write(marker); err != nil {
			return nil, err
		}

		if _, err := w.Write(exif); err != nil {
			return nil, err
		}
	}

	return writer, nil
}
