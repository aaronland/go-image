package colour

import (
	"fmt"
	"io"

	"github.com/mandykoh/prism/meta/autometa"
)

const ICC_DISPLAY_P3 string = "Display P3"

const ICC_EPSON_RGB_G18 string = "EPSON  Standard RGB - Gamma 1.8"

const ICC_ADOBE_RGB_1998 string = "Adobe RGB (1998)"

const ICC_SRGB_21 string = "sRGB IEC61966-2.1"

const ICC_CAMERA_RGB string = "Camera RGB Profile"

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
