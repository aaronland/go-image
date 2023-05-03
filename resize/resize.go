package resize

import (
	"context"
	"fmt"
	"image"
	"math"
	"net/url"
	"strconv"

	"github.com/aaronland/go-image/transform"
	nfnt_resize "github.com/nfnt/resize"
)

type ResizeTransformation struct {
	transform.Transformation
	max int
}

func init() {

	ctx := context.Background()
	transform.RegisterTransformation(ctx, "resize", NewResizeTransformation)
}

func NewResizeTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	query := parsed.Query()
	str_max := query.Get("max")

	if str_max == "" {
		return nil, fmt.Errorf("Missing parameter: max")
	}

	max, err := strconv.Atoi(str_max)

	if err != nil {
		return nil, fmt.Errorf("Failed to convert ?max= parameter, %w", err)
	}

	tr := &ResizeTransformation{
		max: max,
	}

	return tr, nil
}

func (tr *ResizeTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return ResizeImage(ctx, im, tr.max)
}

func ResizeImage(ctx context.Context, im image.Image, max int) (image.Image, error) {

	// calculating w,h is probably unnecessary since we're
	// calling resize.Thumbnail but it will do for now...
	// (20180708/thisisaaronland)

	bounds := im.Bounds()
	dims := bounds.Max

	width := dims.X
	height := dims.Y

	ratio_w := float64(max) / float64(width)
	ratio_h := float64(max) / float64(height)

	ratio := math.Min(ratio_w, ratio_h)

	w := uint(float64(width) * ratio)
	h := uint(float64(height) * ratio)

	sm := nfnt_resize.Thumbnail(w, h, im, nfnt_resize.Lanczos3)

	return sm, nil
}
