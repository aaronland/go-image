package encode

import (
	"context"
	"image"
	"io/ioutil"
	"os"
	"testing"
)

func TestEncoder(t *testing.T) {

	ctx := context.Background()

	_, err := NewEncoder(ctx, "example.jpg")

	if err != nil {
		t.Fatalf("Failed to create new JPEG encoder, %v", err)
	}

	_, err = NewEncoder(ctx, "example.tiff")

	if err == nil {
		t.Fatalf("Expected creation of TIFF encoder to fail")
	}
}

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
