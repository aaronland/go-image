package decode

import (
	"context"
	"os"
	"testing"
)

func TestDecodeGIF(t *testing.T) {

	path := "../fixtures/tokyo.gif"
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

	if format != "gif" {
		t.Fatalf("Invalid format. Expected 'gif' but got '%s'", format)
	}

}
