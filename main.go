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
	// Psuedo image decoders
	_ "github.com/Soreil/pdf"
	_ "github.com/Soreil/svg"
	_ "github.com/Soreil/video/mkv"
	_ "github.com/Soreil/video/mp4"
	_ "github.com/Soreil/video/webm"

	"github.com/nfnt/resize"
)

// JPEGOptions specifies the options to use for encoding JPEG format image
// thumbnails. Should not be modified concurently with thumbnailing.
var JPEGOptions = jpeg.Options{Quality: jpeg.DefaultQuality}

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
	format := autoSelectFormat(imgString)
	out, err := Encode(img, format)
	return out, format, err
}

// Automatically select the output thumbnail image format
func autoSelectFormat(source string) string {
	if source == "jpeg" {
		return source
	}
	return "png"
}

// Encode encodes a given image.Image into the desired format. Currently only
// JPEG and PNG are supported. PNGs are lossily compressed, as per the
// PNGQuantization setting.
func Encode(img image.Image, format string) (io.Reader, error) {
	var (
		out bytes.Buffer
		err error
	)
	switch format {
	case "jpeg":
		err = jpeg.Encode(&out, img, &JPEGOptions)
	case "png":
		err = compressPNG(&out, img)
	default:
		err = errors.New("Unsupported file type")
	}
	return &out, err
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

// Thumbnails creates a thumbnail per size provided in sorted order from large
// to small to reduce the amount of computation required.
func Thumbnails(r io.Reader, sizes ...image.Point) ([]io.Reader, string, error) {
	scaled, imgString, err := DecodeMany(r, sizes...)
	if err != nil {
		return nil, "", err
	}

	thumbs := make([]io.Reader, len(sizes))
	format := autoSelectFormat(imgString)
	for i, scale := range scaled {
		thumbs[i], err = Encode(scale, format)
		if err != nil {
			return thumbs, format, err
		}
	}
	return thumbs, format, err
}

// DecodeMany creates an image.Image thumbnail per size provided in sorted
// order from large to small to reduce the amount of computation required.
func DecodeMany(r io.Reader, sizes ...image.Point) (
	[]image.Image, string, error,
) {
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
	return scaled, imgString, err
}
