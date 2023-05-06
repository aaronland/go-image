package encode

import (
	"testing"
)

func TestGIFEncoder(t *testing.T) {

	err := testEncoder("test.gif")

	if err != nil {
		t.Fatal(err)
	}
}
