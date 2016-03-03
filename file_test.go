package imager

import (
	"io"
	"os"
	"testing"

	_ "image/gif"

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
	input    string
	output   string
	inFormat filetype
}

const inputDir = "inputData/"
const outputDir = "outputData/"

var cases = []testCase{
	{inputDir + "wafel.webm", outputDir + "wafel.webm.png", webmType},
	{inputDir + "yuno.jpg", outputDir + "yuno.jpg.jpg", jpegType},
	{inputDir + "yuno.png", outputDir + "yuno.png.png", pngType},
	{inputDir + "yuno.gif", outputDir + "yuno.gif.png", gifType},
	{inputDir + "PNG_transparency_demonstration_1.png", outputDir + "PNG_transparency_demonstration_1.png.png", pngType},
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
		thumb, err := Thumbnail(file, normal)
		if err != nil {
			t.Fatal(err, test)
		}
		out, err := os.Create(test.output)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := io.Copy(out, thumb); err != nil {
			t.Fatal(err)
		}
	}
}
