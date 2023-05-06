package decode

import (
	"context"
	"os"
	"testing"
)

func TestDecodePNG(t *testing.T) {

	path := "../fixtures/tokyo.png"

	ctx := context.Background()

	r, err := os.Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	dec, err := NewDecoder(ctx, path)

	if err != nil {
		t.Fatal(err)
	}

	_, format, err := dec.Decode(ctx, r)

	if err != nil {
		t.Fatal(err)
	}

	if format != "png" {
		t.Fatalf("Invalid format. Expected 'png' but got '%s'", format)
	}

}
