package encode

import (
	"testing"
)

func TestPNGEncoder(t *testing.T) {

	err := testEncoder("test.png")

	if err != nil {
		t.Fatal(err)
	}
}
