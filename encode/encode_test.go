package encode

import (
	"context"
	"image"
	"io/ioutil"
	"os"
)

func testEncoder(uri string) error {

	path := "../fixtures/tokyo.jpg"
	ctx := context.Background()

	fh, err := os.Open(path)

	if err != nil {
		return err
	}

	defer fh.Close()

	im, _, err := image.Decode(fh)

	if err != nil {
		return err
	}

	enc, err := NewEncoder(ctx, uri)

	if err != nil {
		return err
	}

	wr := ioutil.Discard

	return enc.Encode(ctx, wr, im)
}
