package colour

import (
	"fmt"
	"io"

	"github.com/aaronland/go-image/exif"
)

const UNKNOWN_COLORSPACE uint16 = 65535

const SRGB_COLORSPACE uint16 = 1

const AGRB_COLORSPACE uint16 = 2

func ColorSpace(r io.Reader) (uint16, error) {

	tag, err := exif.TagValue(r, "ColorSpace", "IFD/Exif")

	if err != nil {
		return UNKNOWN_COLORSPACE, fmt.Errorf("Failed to determine tag value, %w", err)
	}

	v, err := tag.Value()

	if err != nil {
		return UNKNOWN_COLORSPACE, fmt.Errorf("Failed to derive tag value, %w", err)
	}

	colorspace := v.([]uint16)

	if len(colorspace) != 1 {
		return UNKNOWN_COLORSPACE, fmt.Errorf("Multiple values for colorspace")
	}

	return colorspace[0], nil
}
