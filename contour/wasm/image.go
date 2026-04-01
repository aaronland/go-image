//go:build wasmjs

package wasm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"syscall/js"

	"github.com/aaronland/go-image/v2/contour"
)

func ContourImageFunc() js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		im_b64 := args[0].String()
		n := args[1].Int()
		scale := 1.0

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			ctx := context.Background()

			resolve := args[0]
			reject := args[1]

			im_b, err := base64.StdEncoding.DecodeString(im_b64)

			if err != nil {
				reject.Invoke(fmt.Printf("Failed to decode b64 image, %v\n", err))
				return nil
			}

			im_r := bytes.NewReader(im_b)
			im, _, err := image.Decode(im_r)

			if err != nil {
				reject.Invoke(fmt.Printf("Failed to decode image, %v\n", err))
				return nil
			}

			new_im, err := contour.ContourImage(ctx, im, n, scale)

			if err != nil {
				reject.Invoke(fmt.Printf("Failed to contour feature, %v\n", err))
				return nil
			}

			var buf bytes.Buffer
			wr := bufio.NewWriter(&buf)

			err = png.Encode(wr, new_im)

			if err != nil {
				reject.Invoke(fmt.Printf("Failed to encode contour-ed image, %v\n", err))
				return nil
			}

			wr.Flush()

			new_b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

			resolve.Invoke(new_b64)
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})

}
