//Converts media types in to thumbnails
package imager

import (
	"image"

	"github.com/nfnt/resize"
)

//Regular quality preset
var normal image.Point = image.Point{X: 250, Y: 250}

//High quality preset
var sharp image.Point = image.Point{X: 500, Y: 500}

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}
