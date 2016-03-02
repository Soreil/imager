package imager

import (
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"testing"

	//	_ "github.com/Soreil/imager/webm"
)

const dataDirectory = "testData/"

//func TestWebmToJPG(t *testing.T) {
//	const filename = dataDirectory + "wafel.webm"
//	if _, err := os.Stat(filename); err != nil {
//		t.Fatal(err)
//	}
//
//	//TODO(sjon):add this stuff in to a pipeline
//	frame := extractVideoFrame(filename)
//	img := avFrameImage(&frame)
//
//	out, err := os.Create(filename + ".jpg")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	i := scale(&img, normal)
//	if err := jpeg.Encode(out, i, nil); err != nil {
//		t.Fatal(err)
//	}
//}

func TestWebmToPNG(t *testing.T) {
	const filename = dataDirectory + "wafel.webm"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}

	//TODO(sjon):add this stuff in to a pipeline
	img := extractVideoFrame(filename)

	out, err := os.Create(filename + ".png")
	if err != nil {
		t.Fatal(err)
	}

	i := scale(img, normal)
	if err := png.Encode(out, i); err != nil {
		t.Fatal(err)
	}
}

func TestJPGToJPG(t *testing.T) {
	const filename = dataDirectory + "yuno.jpg"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}

	//TODO(sjon):add this stuff in to a pipeline
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create(filename + ".jpg")
	if err != nil {
		t.Fatal(err)
	}

	img = scale(img, normal)
	if err := jpeg.Encode(out, img, nil); err != nil {
		t.Fatal(err)
	}
}

func TestPNGToPNG(t *testing.T) {
	const filename = dataDirectory + "yuno.png"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}

	//TODO(sjon):add this stuff in to a pipeline
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create(filename + ".png")
	if err != nil {
		t.Fatal(err)
	}

	img = scale(img, normal)
	if err := png.Encode(out, img); err != nil {
		t.Fatal(err)
	}
}

func TestGIFToPNG(t *testing.T) {
	const filename = dataDirectory + "yuno.gif"
	if _, err := os.Stat(filename); err != nil {
		t.Fatal(err)
	}

	//TODO(sjon):add this stuff in to a pipeline
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	out, err := os.Create(filename + ".png")
	if err != nil {
		t.Fatal(err)
	}

	img = scale(img, normal)
	if err := png.Encode(out, img); err != nil {
		t.Fatal(err)
	}

}

//TODO(sjon):Debug webm driver
//func TestWebmPackage(t *testing.T) {
//	const filename = dataDirectory + "wafel.webm"
//	if _, err := os.Stat(filename); err != nil {
//		t.Fatal(err)
//	}
//
//	//TODO(sjon):add this stuff in to a pipeline
//	file, err := os.Open(filename)
//	if err != nil {
//		t.Fatal(err)
//	}
//	img, _, err := image.Decode(file)
//	if err != nil {
//		t.Fatal(err)
//	}
//	_ = img
//}
