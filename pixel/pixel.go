// Package pixels provides methods for altering images at a per-pixel level.
package pixel

import (
	"image"
	"image/color"
	"sync"
)

// ReplacePixelKey is a struct that defines candidate pixel colours to replace with another colour
type ReplacePixelKey struct {
	// Zero or more `color.Color` instances whose pixel values will be replaced.
	Candidates []color.Color
	// Replacement is a `color.Color` instance that will be used to replace specific pixels.
	Replacement color.Color
}

// PixelFunc defines a function that will return a new colour for a given pixel colour at an x, y coordinate.
type PixelFunc func(int, int, color.Color) (color.Color, error)

// MakeMultiPixelFunc will return a new `PixelFunc` instance that will invoke apply
// the functions defined in 'funcs'.
func MakeMultiPixelFunc(funcs ...PixelFunc) (PixelFunc, error) {

	f := func(x int, y int, c color.Color) (color.Color, error) {

		var err error

		for _, this_f := range funcs {

			c, err = this_f(x, y, c)

			if err != nil {
				return nil, err
			}
		}

		return c, nil
	}

	return f, nil
}

// MakeReplacePixelFunc returns a new `PixelFunc` instance that will match and replace colours
// defined by 'matches'.
func MakeReplacePixelFunc(matches ...ReplacePixelKey) (PixelFunc, error) {

	f := func(x int, y int, c color.Color) (color.Color, error) {

		cr, cg, cb, ca := c.RGBA()

		for _, key := range matches {

			replace := false

			for _, match := range key.Candidates {

				mr, mg, mb, ma := match.RGBA()

				if cr == mr && cg == mg && cb == mb && ca == ma {
					replace = true
					break
				}
			}

			if replace {
				c = key.Replacement
				break
			}
		}

		return c, nil
	}

	return f, nil
}

// MakeTransparentPixelFunc returns a new `PixelFunc` instance that will replace pixels matching
// colours defined by 'matches' with transparent values.
func MakeTransparentPixelFunc(matches ...color.Color) (PixelFunc, error) {

	f := func(x int, y int, c color.Color) (color.Color, error) {

		cr, cg, cb, _ := c.RGBA()

		for _, m := range matches {

			mr, mg, mb, _ := m.RGBA()

			if cr == mr && cg == mg && cb == mb {

				c = color.NRGBA{
					R: uint8(cr / 257),
					G: uint8(cg / 257),
					B: uint8(cg / 257),
					A: 0,
				}

				break
			}
		}

		return c, nil
	}

	return f, nil
}

// ReplacePixels replaces all the pixels in 'im' according to rules defined by 'cb'.
func ReplacePixels(im image.Image, cb PixelFunc) (image.Image, error) {

	bounds := im.Bounds()
	max := bounds.Max

	width := max.X
	height := max.Y

	pr := image.NewNRGBA(image.Rect(0, 0, width, height))

	wg := new(sync.WaitGroup)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			wg.Add(1)

			go func(x int, y int, c color.Color) {

				defer wg.Done()

				new_c, _ := cb(x, y, c)
				pr.Set(x, y, new_c)

			}(x, y, im.At(x, y))
		}
	}

	wg.Wait()

	return pr, nil
}
