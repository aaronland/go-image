package encode

import (
	"context"
	"image"
	"image/jpeg"
	"io"
)

type JPEGEncoder struct {
	Encoder
	options *jpeg.Options
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewJPEGEncoder, "jpg", "jpeg")
}

func NewJPEGEncoder(ctx context.Context, uri string) (Encoder, error) {

	opts := &jpeg.Options{Quality: 100}

	e := &JPEGEncoder{
		options: opts,
	}

	return e, nil
}

func (e *JPEGEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {
	return jpeg.Encode(wr, im, e.options)
}
