package imager

import (
	"image"
	"image/png"
	"io"

	"github.com/foobaz/lossypng/lossypng"
)

// PNGQuantization defines the lossyness and strength of PNG thumbnail
// compression. Should be a positive number. 0 is lossless.Should not be
// modified concurently with thumbnailing.
var PNGQuantization = 20

//Compress PNG using imagequant
func compressPNG(w io.Writer, img image.Image) error {
	compressed := lossypng.Compress(img, lossypng.NoConversion, PNGQuantization)
	return png.Encode(w, compressed)
}
