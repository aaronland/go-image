package decode

import (
	"context"
	"os"
	"testing"
)

func TestDecodeJPEG(t *testing.T) {

	path := "../fixtures/tokyo.jpg"

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

	if format != "jpeg" {
		t.Fatalf("Invalid format. Expected 'jpeg' but got '%s'", format)
	}

}
