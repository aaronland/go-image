// Package pixels provides methods for altering images at a per-pixel level.
package pixel

import (
	"image"
	"image/color"
	"sync"
)

type ReplacePixelKey struct {
	Candidates  []color.Color
	Replacement color.Color
}

type PixelFunc func(int, int, color.Color) (color.Color, error)

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
