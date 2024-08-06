package colour

import (
	"fmt"
	"io"

	"github.com/mandykoh/prism/meta/autometa"
)

func ICCProfileDescription(r io.Reader) (string, error) {

	md, _, err := autometa.Load(r)

	if err != nil {
		return "", fmt.Errorf("Failed to load metadata, %w", err)
	}

	pr, err := md.ICCProfile()

	if err != nil {
		return "", fmt.Errorf("Failed to derive ICC profile, %w", err)
	}

	if pr == nil {
		return "", fmt.Errorf("Missing profile")
	}

	return pr.Description()
}
