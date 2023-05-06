package encode

import (
	"testing"
)

func TestJPEGEncoder(t *testing.T) {

	err := testEncoder("test.jpg")

	if err != nil {
		t.Fatal(err)
	}
}
