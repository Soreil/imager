package imager

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

//Change inputDir to anything you feel like
var inputDir = "inputData/"
var outputDir = os.TempDir() + "/" + time.Now().String() + "/"

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
			file, err := os.Open(test)
			if err != nil {
				t.Fatal(err)
			}
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
