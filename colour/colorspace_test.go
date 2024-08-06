package colour

import (
	"log/slog"
	"os"
	"testing"
)

func TestColorSpace(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	path := "../fixtures/tokyo.jpg"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	cs, err := ColorSpace(r)

	if err != nil {
		t.Fatalf("Failed to derive colorspacefor %s, %v", path, err)
	}

	if cs != SRGB_COLORSPACE {
		t.Fatalf("Unexpected colorspace value for %s, %d", path, cs)
	}
}
