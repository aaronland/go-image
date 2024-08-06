package exif

import (
	"log/slog"
	"os"
	"testing"
)

func TestTagIndex(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	path := "../fixtures/tokyo.jpg"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	_, err = TagIndex(r)

	if err != nil {
		t.Fatalf("Failed to create tag index for %s, %v", path, err)
	}
}

func TestTagValue(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	path := "../fixtures/tokyo.jpg"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	idx, err := TagIndex(r)

	if err != nil {
		t.Fatalf("Failed to create tag index for %s, %v", path, err)
	}

	tag, err := TagValueWithIndex(idx, "ColorSpace", "IFD/Exif")

	if err != nil {
		t.Fatalf("Failed to derive tag value for ColorSpace, %v", err)
	}

	v, err := tag.Value()

	if err != nil {
		t.Fatalf("Failed to derive value for tag, %v", err)
	}

	colorspace := v.([]uint16)

	if len(colorspace) != 1 || colorspace[0] != 1 {
		t.Fatalf("Unexpected value for, %v", colorspace)
	}

	tag2, err := TagValueWithIndex(idx, "Orientation", "IFD")

	if err != nil {
		t.Fatalf("Failed to derive tag value for Orientation, %v", err)
	}

	v2, err := tag2.Value()

	if err != nil {
		t.Fatalf("Failed to derive value for tag, %v", err)
	}

	orientation := v2.([]uint16)

	if len(orientation) != 1 || orientation[0] != 1 {
		t.Fatalf("Unexpected value for, %v", orientation)
	}
}
