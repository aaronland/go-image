package colour

import (
	"context"
	"image"
	"runtime"

	"github.com/aaronland/go-image/transform"
	"github.com/mandykoh/prism"
	"github.com/mandykoh/prism/adobergb"
	"github.com/mandykoh/prism/srgb"
)

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "adobergb", NewAdobeRGBTransformation)
}

type AdobeRGBTransformation struct {
	transform.Transformation
}

func NewAdobeRGBTransformation(ctx context.Context, uri string) (transform.Transformation, error) {
	tr := &AdobeRGBTransformation{}
	return tr, nil
}

func (tr *AdobeRGBTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	new_im := ToAdobeRGB(im)
	return new_im, nil
}

// ToAdobeRGB converts all the coloura in 'im' to match the Adobe RGB colour profile.
func ToAdobeRGB(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := adobergb.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}
