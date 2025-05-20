# go-image

There are many "wrapper" packages for working with images in Go. This one is mine.

## Important

These are image tools that I wrote by and for myself tailored to the needs of personal projects. It's possible they are not the image tools you need or want.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-image.svg)](https://pkg.go.dev/github.com/aaronland/go-image/v2)

## Example

```
package main

import (
	"context"
	"flag"
	"os"

	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-image/v2/encode"
	"github.com/aaronland/go-image/v2/exif"
)

func main() {

	flag.Parse()

	ctx := context.Background()

	decode_opts := &decode.DecodeImageOptions{
		Rotate: true,
	}		  
		    
	for _, path := range flag.Args() {

		r, _ := os.Open(path)
		defer r.Close()

		im, _, ifd, _ := decode.DecodeImageWithOptions(ctx, r, decode_opts)
		ib, _ := exif.NewIfdBuilderWithOrientation(ifd, "1")

		new_path := fmt.Sprintf("%s.jpg", path)
		wr, _ := os.OpenFile(new_path, os.O_RDWR|os.O_CREATE, 0644)

		encode.EncodeJPEG(ctx, wr, im, ib, nil)
		wr.Close()
	}
}
```

_Error handling removed for the sake of brevity._

## Tools

```
$> make cli TAGS=libheif
go build -mod vendor -ldflags="-s -w" -tags libheif -o bin/transform cmd/transform/main.go
go build -mod vendor -ldflags="-s -w" -tags libheif -o bin/resize cmd/resize/main.go
```

### resize

Resize one or more images.

```
$> ./bin/resize -h
Resize one or more images.
Usage:
	./bin/resize uri(N) uri(N)
  -max int
    	The maximum dimension of the resized image
  -preserve-exif
    	Copy EXIF data from source image final target image.
  -profile string
    	An optional colour profile to apply to the resized image. Valid options are: adobergb, displayp3.
  -rotate
    	Automatically rotate based on EXIF orientation. This does NOT update any of the original EXIF data with one exception: If the -rotate flag is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1". (default true)
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	Zero or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

#### Example

Create a new JPEG image with a maximum dimension of 640 pixels and no EXIF data at `./fixtures/tokyo-1280.jpg`.

```
$> ./bin/resize -max 640 ./fixtures/tokyo.jpg
```

### transform

Transform one or more images applying one or more transform:// transformation URIs.

```
$> ./bin/transform -h
Transform one or more images applying one or more transform:// transformation URIs.
Usage:
	./bin/transform uri(N) uri(N)
  -apply-suffix string
    	An optional suffix to apply to the final image filename.
  -format string
    	An optional image format used to encode the final image.
  -preserve-exif
    	Copy EXIF data from source image final target image.
  -rotate
    	Automatically rotate based on EXIF orientation. This does NOT update any of the original EXIF data with one exception: If the -rotate flag is true OR the original image of type HEIC then the EXIF "Orientation" tag is re-written to be "1". (default true)
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	One or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

#### Example

Create a new PNG image with the Apple DisplayP3 colour profile and a maximum dimension of 1280 pixels at `./fixtures/tokyo-1280.png`

```
$> go run cmd/transform/main.go \
	-transformation-uri 'resize://?max=1280' \
	-transformation-uri displayp3:// \
	-apply-suffix -1280 \
	-preserve-exif \
	-format png \
	./fixtures/tokyo.jpg
```

## Decoders

The following image decoders are supported by default:

* `image/bmp`
* `image/gif`
* `image/heic` (if built with `libheif` tag)
* `image/jpeg`
* `image/png`
* `image/tiff`
* `image/webp`

### HEIC images

By default this package supports decoding HEIC images using the [strukturag/libheif-go](http://github.com/strukturag/libheif-go) package which, in turn, depends on the presence of the `libheif` library but when you are compiling your code (or the command line tools) you will need to pass in the `-tags libheif` flag.

## Encoders

* `image/bmp`
* `image/heic` (if built with `libheif` tag)
* `image/jpeg`
* `image/png`
* `image/tiff`

### HEIC images

By default this package supports encoding HEIC images using the [strukturag/libheif-go](http://github.com/strukturag/libheif-go) package which, in turn, depends on the presence of the `libheif` library but when you are compiling your code (or the command line tools) you will need to pass in the `-tags libheif` flag.