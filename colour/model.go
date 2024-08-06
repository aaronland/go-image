package colour

import (
	"io"
)

// UNKNOWN_MODEL defines an unknown or unspecified colour model.
const UNKNOWN_MODEL string = "unknown"

// SRGB_MODEL defines the sRGB colour space/model.
const SRGB_MODEL string = "sRGB"

// DISPLAYP3_MODEL defines the Apple DisplayP3 colour model
const DISPLAYP3_MODEL string = "DisplayP3"

// ARGB_MODEL defines the Adobe RGB colour model.
const ARGB_MODEL string = "Adobe RGB"

const (
	UnknownModel Model = iota
	SRGBModel
	AdobeRGBModel
	AppleDisplayP3Model
)

type Model uint8

func (p Model) String() string {

	switch p {
	case SRGBModel:
		return SRGB_MODEL
	case AdobeRGBModel:
		return ARGB_MODEL
	case AppleDisplayP3Model:
		return DISPLAYP3_MODEL
	default:
		return UNKNOWN_MODEL
	}
}

func StringToModel(str_model string) Model {

	switch str_model {
	case SRGB_MODEL:
		return SRGBModel
	case ARGB_MODEL:
		return AdobeRGBModel
	case DISPLAYP3_MODEL:
		return AppleDisplayP3Model
	default:
		return UnknownModel
	}
}

func DeriveModel(r io.Reader) (Model, error) {

	return UnknownModel, nil
}
