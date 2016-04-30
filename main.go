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

// Thumbnail makes a thumbnail out of a decodable media file.
// Sizes are the maximum dimensions of the thumbnail
func Thumbnail(r io.Reader, s image.Point) (io.Reader, string, error) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	img = scale(img, s)
	return encode(imgString, img)
}

func encode(imgString string, img image.Image) (io.Reader, string, error) {
	var (
		out    bytes.Buffer
		format string
		err    error
	)
	switch imgString {
	case "jpeg":
		format = "jpeg"
		err = jpeg.Encode(&out, img, &jpgOptions)
	case "png", "webm", "pdf", "gif", "svg":
		format = "png"
		err = compressPNG(&out, img, fast)
	default:
		err = errors.New("Unsupported file type")
	}
	return &out, format, err
}

// TwoThumbnails creates takes a supported format file and outputs two
// thumbnails of desired frame sizes. It is more efficient than calling
// Thumbnail twice, because it only decodes the file once. Large should be
// larger than small.
func TwoThumbnails(r io.Reader, large image.Point, small image.Point) (
	largeThumb io.Reader, smallThumb io.Reader, format string, err error,
) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return
	}
	scaledLarge := scale(img, large)
	scaledSmall := scale(scaledLarge, small)

	largeThumb, format, err = encode(imgString, scaledLarge)
	if err != nil {
		return
	}
	smallThumb, _, err = encode(imgString, scaledSmall)
	return
}
