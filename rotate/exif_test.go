package rotate

import (
	"context"
	"os"
	"testing"
)

func TestGetImageOrientation(t *testing.T) {

	path := "../fixtures/tokyo.jpg"

	ctx := context.Background()

	r, err := os.Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	o, err := GetImageOrientation(ctx, r)

	if err != nil {
		t.Fatal(err)
	}

	expected := "1"

	if o != expected {
		t.Fatalf("Unexpected orientation. Got '%s' but expected '%s'", o, expected)
	}

}
