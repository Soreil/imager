// Package imager converts media types in to size optimised thumbnails
package imager

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"sort"

	// Import decoders
	_ "image/gif"

	_ "github.com/Soreil/pdf"
	_ "github.com/Soreil/svg"
	_ "github.com/Soreil/video/mkv"
	_ "github.com/Soreil/video/mp4"
	_ "github.com/Soreil/video/webm"

	"github.com/nfnt/resize"
)

// JPEGOptions specifies the options to use for encoding JPEG format image
// thumbnails. Should not be modified concurently with thumbnailing.
var JPEGOptions = jpeg.Options{jpeg.DefaultQuality}

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
		err = jpeg.Encode(&out, img, &JPEGOptions)
	case "png", "webm", "pdf", "gif", "svg", "mkv", "mp4":
		format = "png"
		err = compressPNG(&out, img)
	default:
		err = errors.New("Unsupported file type")
	}
	return &out, format, err
}

//Type for sort.Sort
type points []image.Point

func (p points) Len() int {
	return len(p)
}

func (p points) Less(i, j int) bool {
	return !p[i].In(image.Rectangle{Max: p[j]})
}

func (p points) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

//Thumbnails creates a thumbnail per size provided in sorted order from large to small to reduce the amount of computation required.
func Thumbnails(r io.Reader, sizes ...image.Point) ([]io.Reader, string, error) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, "", err
	}
	//Make it so we have them in decreasing sized order
	sort.Sort(points(sizes))

	scaled := make([]image.Image, len(sizes))
	for i, size := range sizes {
		scaled[i] = scale(img, size)
		img = scaled[i]
	}

	thumbs := make([]io.Reader, len(sizes))
	var format string
	for i, scale := range scaled {
		thumbs[i], format, err = encode(imgString, scale)
		if err != nil {
			return thumbs, format, err
		}
	}
	return thumbs, format, err
}
