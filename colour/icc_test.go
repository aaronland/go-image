package colour

import (
	"log/slog"
	"os"
	"testing"
)

func TestICCProfileDescription(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	path := "../fixtures/tokyo.jpg"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	_, err = ICCProfileDescription(r)

	if err == nil {
		t.Fatalf("Did not expect to derive ICC profile for %s", path)
	}

	if err.Error() != "Missing profile" {
		t.Fatalf("Failed to derive ICC profile for %s, %v", path, err)
	}

}
