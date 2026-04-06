//go:build wasmjs
package main

import (
	"log"
	"syscall/js"

	"github.com/aaronland/go-image/v2/contour/wasm"
)

func main() {

	contour_image_func := wasm.ContourImageFunc()
	defer contour_image_func.Release()

	contour_svg_func := wasm.ContourSVGFunc()
	defer contour_svg_func.Release()
	
	js.Global().Set("contour_image", contour_image_func)
	js.Global().Set("contour_svg", contour_svg_func)	

	c := make(chan struct{}, 0)

	log.Println("contour function initialized")
	<-c
}
