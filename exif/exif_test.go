package exif

import (
	"bufio"
	"bytes"
	"image"
	_ "image/jpeg"
	"os"
	"testing"

	"github.com/rwcarlsen/goexif/exif"
	// "github.com/rwcarlsen/goexif/mknote"
)

func TestUpdateExif(t *testing.T) {

	path := "../fixtures/tokyo.jpg"

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
	buf_wr := bufio.NewWriter(&buf)

	jpeg_dt := "2006:01:02 15:04:05"

	exif_props := map[string]interface{}{
		"DateTime": jpeg_dt,
	}

	err = UpdateExif(im, buf_wr, exif_props)

	if err != nil {
		t.Fatal(err)
	}

	buf_wr.Flush()

	buf_r := bytes.NewReader(buf.Bytes())

	x, err := exif.Decode(buf_r)

	if err != nil {
		t.Fatal(err)
	}

	dt, err := x.DateTime()

	if err != nil {
		t.Fatal(err)
	}

	if dt.Format(jpeg_dt) != jpeg_dt {
		t.Fatalf("Unexpected date time. Expected '%s', got '%s'", jpeg_dt, dt)
	}
}
