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

	// And our own decoders
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

// Dimensions contains the width and height of an image
type Dimensions struct {
	width, height int
}

// Thumb contains an io.Reader of the generated thumbnail and its width
// and height
type Thumb struct {
	Dimensions
	io.Reader
}

//TODO(sjon): evaluate best resizing algorithm
//Resizes the image to max dimensions
func scale(img image.Image, p image.Point) image.Image {
	return resize.Thumbnail(uint(p.X), uint(p.Y), img, resize.Bilinear)
}

// Thumbnail makes a thumbnail out of a decodable media file. Sizes are the
// maximum dimensions of the thumbnail. Returns a Thumb of the resulting
// thumbnail, the format of the thumbnail, the dimensions of the source image
// and error, if any.
func Thumbnail(r io.Reader, s image.Point) (Thumb, string, Dimensions, error) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return Thumb{}, "", Dimensions{}, err
	}

	srcDims := getDims(img)
	img = scale(img, s)
	format := autoSelectFormat(imgString)
	out, err := Encode(img, format)
	thumb := Thumb{
		Dimensions: getDims(img),
		Reader:     out,
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

// Compute the width and height of an image
func getDims(img image.Image) Dimensions {
	rect := img.Bounds()
	return Dimensions{
		width:  rect.Max.X - rect.Min.X,
		height: rect.Max.Y - rect.Min.Y,
	}
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
// to small to reduce the amount of computation required. Returns the sorted
// []Thumb of thumbnails, the format string of the thumbnails, dimensions of the
// source image and error, if any.
func Thumbnails(r io.Reader, sizes ...image.Point) (
	[]Thumb, string, Dimensions, error,
) {
	img, imgString, err := image.Decode(r)
	if err != nil {
		return nil, "", Dimensions{}, err
	}

	srcDims := getDims(img)

	//Make it so we have them in decreasing sized order
	sort.Sort(points(sizes))
	scaled := make([]image.Image, len(sizes))
	for i, size := range sizes {
		scaled[i] = scale(img, size)
		img = scaled[i]
	}

	thumbs := make([]Thumb, len(sizes))
	format := autoSelectFormat(imgString)
	for i, scale := range scaled {
		reader, err := Encode(scale, format)
		if err != nil {
			return thumbs, format, srcDims, err
		}

		thumbs[i] = Thumb{
			Dimensions: getDims(scale),
			Reader:     reader,
		}
	}

	return thumbs, format, srcDims, err
}
