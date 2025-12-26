package multiply

import (
	"image"
	"image/color"
)

// Multiply returns a new image that is the result of blending src1 and src2
// with a Photoshop-style “multiply” blend mode.
func Multiply(im1, im2 image.Image) *image.RGBA {

	if im1 == nil || im2 == nil {
		return nil
	}

	bounds := im1.Bounds().Intersect(im2.Bounds())
	if bounds.Empty() {
		return nil
	}

	dst := image.NewRGBA(bounds)

	// Pre‑allocate variables to reduce allocations inside the loop.
	var c1, c2 color.RGBA
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 = color.RGBAModel.Convert(im1.At(x, y)).(color.RGBA)
			c2 = color.RGBAModel.Convert(im2.At(x, y)).(color.RGBA)

			r := uint8((uint16(c1.R) * uint16(c2.R)) / 255)
			g := uint8((uint16(c1.G) * uint16(c2.G)) / 255)
			b := uint8((uint16(c1.B) * uint16(c2.B)) / 255)
			a := uint8((uint16(c1.A) * uint16(c2.A)) / 255)

			dst.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}

	return dst
}
