package decode

import (
	"context"
	"testing"
)

func TestDecodeFromPath(t *testing.T) {

	ctx := context.Background()

	_, format, err := DecodeFromPath(ctx, "../fixtures/tokyo.jpg")

	if err != nil {
		t.Fatal(err)
	}

	if format != "jpeg" {
		t.Fatalf("Expected image to be 'jpeg' but is '%s'", format)
	}

}

func TestDecoder(t *testing.T) {

	ctx := context.Background()

	_, err := NewDecoder(ctx, "example.jpg")

	if err != nil {
		t.Fatalf("Failed to create new JPEG decoder, %v", err)
	}

	_, err = NewDecoder(ctx, "example.tiff")

	if err == nil {
		t.Fatalf("Expected creation of TIFF decoder to fail")
	}
}
