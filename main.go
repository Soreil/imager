//Converts media types in to thumbnails
package imager

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os/exec"
	"strings"

	_ "github.com/Soreil/webm"
	"github.com/nfnt/resize"
)

type size image.Point

//Regular quality preset
var normal size = size(image.Point{X: 250, Y: 250})

//High quality preset
var sharp size = size(image.Point{X: 500, Y: 500})

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}

//Makes a thumbnail out of a decodable media file.
//Sizes are the maximum dimensions of the thumbnail
func Thumbnail(r io.Reader, s size) (io.Reader, error) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	img = scale(img, image.Point(s))

	var out bytes.Buffer
	if imgString == "jpeg" {
		if err := jpeg.Encode(&out, img, &jpgOptions); err != nil {
			return nil, err
		}
	} else if imgString == "png" || imgString == "webm" {
		img, err := CompressPNG(img, fastest)
		if err != nil {
			return nil, err
		}
		if err := pngEncoder.Encode(&out, img); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("I give up, I don't know what this file type is")
	}
	return &out, nil
}

//Encode SVG to PNG as image.Image
func svgToImage(input []byte) (image.Image, error) {
	cmd := exec.Command("rsvg-convert")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = strings.NewReader(string(input))

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	img, err := png.Decode(bytes.NewReader(out.Bytes()))
	if err != nil {
		return nil, err
	}
	return img, nil
}
