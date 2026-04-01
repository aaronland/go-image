package contour

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
)

func TestContourImage(t *testing.T) {

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

	im, err = ContourImage(ctx, im, 3, 1.0)

	if err != nil {
		t.Fatal(err)
	}

	//

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	jpeg_opts := &jpeg.Options{
		Quality: 100,
	}

	err = jpeg.Encode(wr, im, jpeg_opts)

	if err != nil {
		t.Fatal(err)
	}

	wr.Flush()

	//

	path_contour := "fixtures/tokyo-contour-3.jpg"
	body, err := ioutil.ReadFile(path_contour)

	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(body, buf.Bytes()) {
		t.Fatalf("Contour output does not match %s", path_contour)
	}

}
