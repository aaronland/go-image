// Package colour provides methods for working with colour profiles.
package colour

const (
	UnknownProfile Profile = iota
	SRGBProfile
	AdobeRGBProfile
	AppleDisplayP3Profile
)

type Profile uint8

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
