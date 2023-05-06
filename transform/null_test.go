package transform

import (
	"context"
	"image"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestNullTransformation(t *testing.T) {

	path := "../fixtures/tokyo.jpg"

	ctx := context.Background()

	r, err := os.Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	im, _, err := image.Decode(r)

	if err != nil {
		t.Fatal(err)
	}

	tr, err := NewTransformation(ctx, "null://")

	if err != nil {
		t.Fatal(err)
	}

	_, err = tr.Transform(ctx, im)

	if err != nil {
		t.Fatal(err)
	}

}
