package resize

import (
	"context"
	"image"
	_ "image/jpeg"
	"os"
	"testing"
)

func TestResizeTransformation(t *testing.T) {

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

	tr, err := NewResizeTransformation(ctx, "resize://?max=300")

	if err != nil {
		t.Fatal(err)
	}

	new_im, err := tr.Transform(ctx, im)

	if err != nil {
		t.Fatal(err)
	}

	new_bounds := new_im.Bounds()

	new_w := new_bounds.Max.X
	expected_w := 300

	if new_w != expected_w {
		t.Fatalf("Unexpected new width. Expected '%d' but got '%d'", expected_w, new_w)
	}

}
