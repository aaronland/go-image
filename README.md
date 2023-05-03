# go-image

There are many "wrapper" packages for working with images in Go. This one is mine.

## Important

This is work in progress and will supersede all of the older `aaronland/go-image-*` packages.

## Documentation

Documentation is incomplete at this time.

## Tools

### transform

```
> go run cmd/transform/main.go \
	-transformation-uri 'resize://?max=1280' \
	-transformation-uri displayp3:// \
	-apply-suffix -1280 \
	-format png \
	/usr/local/big-fish/big-fish-014.jpg
```

Create a new PNG image with the Apple DisplayP3 colour profile and a maximum dimension of 1280 pixels at `/usr/local/big-fish/big-fish-014-1280.png`.