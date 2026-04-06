package contour

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"os"
	"testing"
)

func TestContourImageSVG(t *testing.T) {

	ctx := context.Background()

	path := "fixtures/tokyo.jpg"

	r, err := os.Open(path)

	if err != nil {
		t.Fatal(err)
	}

	defer r.Close()

	im, _, err := image.Decode(r)

	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err = ContourImageSVG(ctx, wr, im, 3, 1.0)

	if err != nil {
		t.Fatal(err)
	}

	wr.Flush()

	//

	path_svg := "fixtures/tokyo-003.svg"
	body, err := ioutil.ReadFile(path_svg)

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(body, buf.Bytes()) {
		t.Fatalf("Contour output does not match %s", path_svg)
	}

}
