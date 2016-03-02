package imager

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"os"
	"testing"

	_ "github.com/Soreil/webm"
)

type filetype int

const (
	pngType filetype = iota
	jpegType
	gifType
	webmType
	mp4Type
)

type testCase struct {
	input     string
	output    string
	inFormat  filetype
	outFormat filetype
}

const inputDir = "inputData/"
const outputDir = "outputData/"

var cases = []testCase{
	{inputDir + "wafel.webm", outputDir + "wafel.webm.png", webmType, pngType},
	{inputDir + "wafel.webm", outputDir + "wafel.webm.jpg", webmType, jpegType},
	{inputDir + "yuno.jpg", outputDir + "yuno.jpg.png", jpegType, pngType},
	{inputDir + "yuno.jpg", outputDir + "yuno.jpg.jpg", jpegType, jpegType},
	{inputDir + "yuno.png", outputDir + "yuno.png.png", pngType, pngType},
	{inputDir + "yuno.png", outputDir + "yuno.png.jpg", pngType, jpegType},
	{inputDir + "yuno.gif", outputDir + "yuno.gif.png", gifType, pngType},
	{inputDir + "yuno.gif", outputDir + "yuno.gif.jpg", gifType, jpegType},
}

func TestDecode(t *testing.T) {
	for _, test := range cases {
		if _, err := os.Stat(test.input); err != nil {
			t.Fatal(err)
		}
		file, err := os.Open(test.input)
		if err != nil {
			t.Fatal(err)
		}
		img, _, err := image.Decode(file)
		if err != nil {
			t.Fatal(err)
		}
		out, err := os.Create(test.output)
		if err != nil {
			t.Fatal(err)
		}
		img = scale(img, normal)

		switch test.outFormat {
		case pngType:
			if err := pngEncoder.Encode(out, img); err != nil {
				t.Fatal(err)
			}
		case webmType:
			t.Fatal(errors.New("Can only decode webm, not encode"))
		case jpegType:
			if err := jpeg.Encode(out, img, nil); err != nil {
				t.Fatal(err)
			}
		case gifType:
			if err := gif.Encode(out, img, nil); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatal(errors.New("Format not implemented"))
		}
	}
}
