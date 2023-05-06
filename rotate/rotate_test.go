package rotate

import (
	"context"
	"image"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestRotateWithOrientation(t *testing.T) {

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

	bounds := im.Bounds()

	w := bounds.Max.X
	h := bounds.Max.Y

	new_im, err := RotateImageWithOrientation(ctx, im, "5")

	if err != nil {
		t.Fatal(err)
	}

	new_bounds := new_im.Bounds()

	new_w := new_bounds.Max.X
	new_h := new_bounds.Max.Y

	if new_w != h {
		t.Fatalf("Unexpected new width. Expected '%d' but got '%d'", h, new_w)
	}

	if new_h != w {
		t.Fatalf("Unexpected new width. Expected '%d' but got '%d'", w, new_h)
	}

}

func TestRotateWithDegrees(t *testing.T) {

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

	bounds := im.Bounds()

	w := bounds.Max.X
	h := bounds.Max.Y

	new_im, err := RotateImageWithDegrees(ctx, im, 90.0)

	if err != nil {
		t.Fatal(err)
	}

	new_bounds := new_im.Bounds()

	new_w := new_bounds.Max.X
	new_h := new_bounds.Max.Y

	if new_w != h {
		t.Fatalf("Unexpected new width. Expected '%d' but got '%d'", h, new_w)
	}

	if new_h != w {
		t.Fatalf("Unexpected new width. Expected '%d' but got '%d'", w, new_h)
	}

}
