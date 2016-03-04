//Converts media types in to thumbnails
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

type size image.Point

//Regular quality preset
var Normal size = size(image.Point{X: 250, Y: 250})

//High quality preset
var Sharp size = size(image.Point{X: 500, Y: 500})

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p size) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}

//Makes a thumbnail out of a decodable media file.
//Sizes are the maximum dimensions of the thumbnail
func Thumbnail(r io.Reader, s size) (io.Reader, error) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img = scale(img, s)

	var out bytes.Buffer
	if imgString == "jpeg" {
		if err := jpeg.Encode(&out, img, &jpgOptions); err != nil {
			return nil, err
		}
	} else if imgString == "png" || imgString == "webm" || imgString == "pdf" || imgString == "gif" || imgString == "svg" {
		err := compressPNG(&out, img, fast)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("I give up, I don't know what this file type is")
	}
	return &out, nil
}
