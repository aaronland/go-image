# go-image

There are many "wrapper" packages for working with images in Go. This one is mine.

## Important

These are image tools that I wrote by and for myself tailored to the needs of personal projects. It's possible they are not the image tools you need or want.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/aaronland/go-image.svg)](https://pkg.go.dev/github.com/aaronland/go-image)

## Tools

### resize

```
$> ./bin/resize -h
  -max int
    	The maximum dimension of the resized image
  -profile string
    	An optional colour profile to apply to the resized image. Valid options are: adobergb, displayp3.
  -source-uri string
    	A valid gocloud.dev/blob URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	Zero or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

### transform

```
$> ./bin/transform -h
  -apply-suffix string
    	An optional suffix to apply to the final image filename.
  -format string
    	An optional image format used to encode the final image.
  -source-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are read from. (default "file:///")
  -target-uri string
    	A valid gocloud.dev/blob.Bucket URI where images are written to. (default "file:///")
  -transformation-uri transform.Transformation
    	One or more additional transform.Transformation URIs used to further modify an image after resizing (and before any additional colour profile transformations are performed).
```

For example:

```
$> go run cmd/transform/main.go \
	-transformation-uri 'resize://?max=1280' \
	-transformation-uri displayp3:// \
	-apply-suffix -1280 \
	-format png \
	./fixtures/tokyo.jpg
```

Create a new PNG image with the Apple DisplayP3 colour profile and a maximum dimension of 1280 pixels at `/usr/local/big-fish/big-fish-014-1280.png`.

## See also

* https://github.com/aaronland/go-image-halftone