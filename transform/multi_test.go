package transform

import (
	"context"
	"image"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestMultiTransformation(t *testing.T) {

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

	tr, err := NewMultiTransformationWithURIs(ctx, "null://", "null://")

	if err != nil {
		t.Fatal(err)
	}

	_, err = tr.Transform(ctx, im)

	if err != nil {
		t.Fatal(err)
	}

}
