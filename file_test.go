package imager

import (
	"io"
	"os"
	"testing"
	"sync"
	"time"
)

type filetype int

const (
	pngType filetype = iota
	jpegType
	gifType
	webmType
	mp4Type
	svgType
	pdfType
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
	{inputDir + "The United States of America.svg", outputDir + "The United States of America.svg.png", svgType},
	{inputDir + "pdf.pdf", outputDir + "pdf.pdf.png", pdfType},
}

func TestDecode(t *testing.T) {
	if _,err := os.Stat(outputDir); err != nil {
		os.Mkdir(outputDir,os.ModeDir+0755)
	}
	var wg sync.WaitGroup
	for _, test := range cases {
		wg.Add(1)
		startingTime := time.Now()
		go func(test testCase) {
			if _, err := os.Stat(test.input); err != nil {
				t.Fatal(err)
			}
			file, err := os.Open(test.input)
			if err != nil {
				t.Fatal(err)
			}
			thumb, err := Thumbnail(file, Normal)
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
			t.Log(time.Now().Sub(startingTime),test)
			wg.Done()
		}(test)
	}
	wg.Wait()
}
