package imager

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Soreil/mp3"
)

type testCase struct {
	input  string
	output string
}

const inputDir = "inputData/"
const outputDir = "outputData/"

var cases = []testCase{
	{inputDir + "wafel.webm", outputDir + "wafel.webm.png"},
	{inputDir + "yuno.jpg", outputDir + "yuno.jpg.jpg"},
	{inputDir + "yuno.png", outputDir + "yuno.png.png"},
	{inputDir + "yuno.gif", outputDir + "yuno.gif.png"},
	{inputDir + "PNG_transparency_demonstration_1.png", outputDir + "PNG_transparency_demonstration_1.png.png"},
	{inputDir + "The United States of America.svg", outputDir + "The United States of America.svg.png"},
	{inputDir + "pdf.pdf", outputDir + "pdf.pdf.png"},
}

func TestDecode(t *testing.T) {
	if _, err := os.Stat(outputDir); err != nil {
		os.Mkdir(outputDir, os.ModeDir+0755)
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
			thumb, inFormat, outFormat, err := Thumbnail(file, Normal)
			if err != nil {
				t.Fatal(err, inFormat, outFormat, test)
			}
			out, err := os.Create(test.output)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(out, thumb); err != nil {
				t.Fatal(err)
			}
			t.Log(time.Now().Sub(startingTime), test)
			wg.Done()
		}(test)
	}
	wg.Wait()
}

func TestMP3(t *testing.T) {
	t.Log(mp3.IsMP3(inputDir + "example.mp3"))
}
