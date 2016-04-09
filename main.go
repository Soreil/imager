//Converts media types in to size optimised thumbnails
package imager

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	"io"

	_ "github.com/Soreil/pdf"
	_ "github.com/Soreil/svg"
	_ "github.com/Soreil/webm"

	"github.com/nfnt/resize"
)

//Regular quality preset
var Normal = image.Point{X: 250, Y: 250}

//High quality preset
var Sharp = image.Point{X: 500, Y: 500}

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}

//Makes a thumbnail out of a decodable media file.
//Sizes are the maximum dimensions of the thumbnail
func Thumbnail(r io.Reader, s image.Point) (io.Reader, string, error) {
	var outputFormat string
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, outputFormat, err
	}
	img = scale(img, s)

	var out bytes.Buffer
	switch imgString {
	case "jpeg":
		outputFormat = "jpeg"
		err = jpeg.Encode(&out, img, &jpgOptions)
	case "png", "webm", "pdf", "gif", "svg":
		outputFormat = "png"
		err = compressPNG(&out, img, fast)
	default:
		err = errors.New("Unsupported file type")
	}

	return &out, outputFormat, err
}
