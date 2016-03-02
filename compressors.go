package imager

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"os/exec"
	"strings"
)

var pngEncoder = png.Encoder{png.BestCompression}
var jpgOptions = jpeg.Options{jpeg.DefaultQuality}

type speed string

const (
	bruteforce speed = "1"
	standard   speed = "3"
	fast       speed = "10"
	fastest    speed = "11"
)

//Compress PNG using imagequant
func CompressPNG(img image.Image, s speed) (image.Image, error) {
	var w bytes.Buffer
	err := png.Encode(&w, img)
	if err != nil {
		return nil, err
	}

	compressed, err := compressBytes(w.Bytes(), s)
	if err != nil {
		return nil, err
	}

	output, err := png.Decode(bytes.NewReader(compressed))
	return output, err
}

//Add imagequant structures here
func compressBytes(input []byte, speed speed) ([]byte, error) {
	cmd := exec.Command("pngquant", "-", "--speed", string(speed))
	cmd.Stdin = strings.NewReader(string(input))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
