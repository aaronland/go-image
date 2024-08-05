// Package colour provides methods for working with colour profiles.
package colour

const (
	UnknownProfile Profile = iota
	SRGBProfile
	AdobeRGBProfile
	AppleDisplayP3Profile
)

type Profile uint8

// UNKNOWN defines an unknown or unspecified colour space/profile.
const UNKNOWN_PROFILE string = "unknown"

// SRGB defines the sRGB colour space/profile.
const SRGB_PROFILE string = "sRGB"

// DISPLAYP3 defines the Apple DisplayP3 colour space/profile
const DISPLAYP3_PROFILE string = "DisplayP3"

// ARGB defines the Adobe RGB colour space/profile.
const ARGB_PROFILE string = "Adobe RGB"

func (p Profile) String() string {

	switch p {
	case SRGBProfile:
		return SRGB_PROFILE
	case AdobeRGBProfile:
		return ARGB_PROFILE
	case AppleDisplayP3Profile:
		return DISPLAYP3_PROFILE
	default:
		return UNKNOWN_PROFILE
	}
}

func StringToProfile(str_profile string) Profile {

	switch str_profile {
	case SRGB_PROFILE:
		return SRGBProfile
	case ARGB_PROFILE:
		return AdobeRGBProfile
	case DISPLAYP3_PROFILE:
		return AppleDisplayP3Profile
	default:
		return UnknownProfile
	}
}
