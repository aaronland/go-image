package colour

import (
	"context"
	"image"
	"runtime"

	"github.com/aaronland/go-image/transform"
	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/displayp3"
	"github.com/mandykoh/prism/srgb"
)

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "displayp3", NewDisplayP3Transformation)
}

type DisplayP3Transformation struct {
	transform.Transformation
}

func NewDisplayP3Transformation(ctx context.Context, uri string) (transform.Transformation, error) {
	tr := &DisplayP3Transformation{}
	return tr, nil
}

func (tr *DisplayP3Transformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	new_im := ToDisplayP3(im)
	return new_im, nil
}

// ToDisplayP3 converts all the coloura in 'im' to match the Apple Display P3 colour profile.
func ToDisplayP3(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := displayp3.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}
