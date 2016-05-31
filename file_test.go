package imager

import (
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"fmt"
)

//Change inputDir to anything you feel like
var (
	// Regular quality preset
	Normal = image.Point{X: 250, Y: 250}

	// High quality preset
	Sharp     = image.Point{X: 500, Y: 500}
	inputDir  = "testdata/"
	outputDir = os.TempDir() + string(os.PathSeparator) + fmt.Sprint(time.Now().UnixNano()) + string(os.PathSeparator)
)

func TestDecode(t *testing.T) {

	if _, err := os.Stat(outputDir); err != nil {
		err := os.Mkdir(outputDir, os.ModeDir+0755)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Output directory: ", outputDir)
	}

	in, err := os.Open(inputDir)
	if err != nil {
		t.Fatal(err)
	}

	names, err := in.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(inputDir); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, test := range names {
		startingTime := time.Now()
		wg.Add(1)

		go func(test string) {
			defer wg.Done()
			if _, err := os.Stat(test); err != nil {
				t.Fatal(err)
			}
			file := open(test, t)
			defer file.Close()
			thumb, outFormat, err := Thumbnail(file, Normal)
			if err != nil {
				t.Fatal(err, outFormat, test)
			}
			out, err := os.Create(outputDir + test + "." + outFormat)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(out, thumb); err != nil {
				t.Fatal(err)
			}
			t.Log(time.Now().Sub(startingTime), test)
		}(test)

	}
	wg.Wait()
}

func TestThumbnails(t *testing.T) {
	file := open("yuno.gif", t)
	thumbs, format, err := Thumbnails(file, Sharp, Normal)
	if format != "png" {
		t.Fatalf("Wrong format: %s", format)
	}
	assertError(err, t)
	largeThumb, err := ioutil.ReadAll(thumbs[0])
	assertError(err, t)
	smallThumb, err := ioutil.ReadAll(thumbs[1])
	assertError(err, t)
	if len(smallThumb) > len(largeThumb) {
		t.Fatal("Thumbnail sizes don't match")
	}
}

func open(path string, t *testing.T) *os.File {
	file, err := os.Open(filepath.FromSlash(path))
	assertError(err, t)
	return file
}

func assertError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
