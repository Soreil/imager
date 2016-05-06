package imager

import (
	"github.com/bakape/lossypng/lossypng"
	"image"
	"image/png"
	"io"
)

// PNGQuantization defines the lossyness and strength of PNG thumbnail
// compression. Should be a positive number. 0 is lossless.Should not be
// modified concurently with thumbnailing.
var PNGQuantization = 20

//Compress PNG using imagequant
func compressPNG(w io.Writer, img image.Image) error {
	comresssed := lossypng.Compress(img, lossypng.NoConversion, PNGQuantization)
	return png.Encode(w, comresssed)
}
