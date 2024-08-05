package exif

/*

 	(the value of 0x2 is not standard EXIF. Instead, an Adobe RGB image is indicated by "Uncalibrated" with an InteropIndex of "R03". The values 0xfffd and 0xfffe are also non-standard, and are used by some Sony cameras)
0x1 = sRGB
0x2 = Adobe RGB
0xfffd = Wide Gamut RGB
0xfffe = ICC Profile
0xffff = Uncalibrated

// https://www.exiftool.org/TagNames/EXIF.html

// https://exiv2.org/makernote.html

*/

import (
	"fmt"
	"io"
	_ "log"

	go_exif "github.com/dsoprea/go-exif/v3"
	go_exifcommon "github.com/dsoprea/go-exif/v3/common"
)

func Colorspace(r io.Reader) (string, error) {

	rawExif, err := go_exif.SearchAndExtractExifWithReader(r)

	if err != nil {
		return "", err
	}

	im, err := go_exifcommon.NewIfdMappingWithStandard()

	if err != nil {
		return "", err
	}

	ti := go_exif.NewTagIndex()

	_, index, err := go_exif.Collect(im, ti, rawExif)

	if err != nil {
		return "", err
	}

	var ifd *go_exif.Ifd

	for _, i := range index.Ifds {

		ident := i.IfdIdentity()

		if ident.String() != "IFD/Exif" {
			continue
		}

		ifd = i
		break
	}

	if ifd == nil {
		return "", fmt.Errorf("Failed to find IFD/Exif")
	}

	tagName := "ColorSpace"

	results, err := ifd.FindTagWithName(tagName)

	if err != nil {
		return "", err
	}

	// This should never happen.
	if len(results) != 1 {
		return "", fmt.Errorf("Multiple results")
	}

	ite := results[0]

	valueRaw, err := ite.Value()

	if err != nil {
		return "", err
	}

	value := valueRaw
	return fmt.Sprintf("%v", value), nil
}
