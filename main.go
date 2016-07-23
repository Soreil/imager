// Package imager converts media types in to size optimised thumbnails
package imager

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"sort"

	// Import gif decoder
	_ "image/gif"
<<<<<<< HEAD
	// Psuedo image decoders
=======

	// And our own decoders
>>>>>>> fc0710d537f5d41686bc2c9b1f962e5c8e1abd34
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

// Thumb contains an io.Reader of the generated thumbnail and its width
// and height
type Thumb struct {
	image.Rectangle
	bytes.Buffer
}

// Scale resizes the image to max dimensions
// TODO(sjon): evaluate best resizing algorithm
func Scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}

// Thumbnail makes a thumbnail out of a decodable media file. Sizes are the
// maximum dimensions of the thumbnail. Returns a Thumb of the resulting
// thumbnail, the format of the thumbnail, the dimensions of the source image
// and error, if any.
func Thumbnail(r io.Reader, s image.Point) (
	*Thumb, string, image.Rectangle, error,
) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, "", image.Rectangle{}, err
	}

	srcDims := img.Bounds()
	img = Scale(img, s)
	format := autoSelectFormat(imgString)
	out, err := Encode(img, format)
	thumb := &Thumb{
		Rectangle: img.Bounds(),
		Buffer:    *out,
	}
	return thumb, format, srcDims, err
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
func Encode(img image.Image, format string) (*bytes.Buffer, error) {
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
// to small to reduce the amount of computation required. Returns the sorted
// []Thumb of thumbnails, the format string of the thumbnails, dimensions of the
// source image and error, if any.
func Thumbnails(r io.Reader, sizes ...image.Point) (
	[]*Thumb, string, image.Rectangle, error,
) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, "", image.Rectangle{}, err
	}

	srcDims := img.Bounds()

	//Make it so we have them in decreasing sized order
	sort.Sort(points(sizes))
	scaled := make([]image.Image, len(sizes))
	for i, size := range sizes {
		scaled[i] = Scale(img, size)
		img = scaled[i]
	}

	thumbs := make([]*Thumb, len(sizes))
	format := autoSelectFormat(imgString)
	for i, scale := range scaled {
		buf, err := Encode(scale, format)
		if err != nil {
			return thumbs, format, srcDims, err
		}

		thumbs[i] = &Thumb{
			Rectangle: scale.Bounds(),
			Buffer:    *buf,
		}
	}

	return thumbs, format, srcDims, err
}
