package imager

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os/exec"
	"runtime"
	"strings"
)

var pngEncoder = png.Encoder{png.NoCompression}
var jpgOptions = jpeg.Options{jpeg.DefaultQuality}

type speed string

const (
	bruteforce speed = "1"
	standard   speed = "3"
	fast       speed = "10"
	fastest    speed = "11"
)

//Compress PNG using imagequant
func compressPNG(out io.Writer, img image.Image, s speed) error {
	var w bytes.Buffer
	err := png.Encode(&w, img)
	if err != nil {
		return err
	}

	compressed, err := compressBytes(w.Bytes(), s)
	if err != nil {
		return err
	}
	_, err = out.Write(compressed)
	return err
}

//Add imagequant structures here
func compressBytes(input []byte, speed speed) ([]byte, error) {
	var cmd *exec.Cmd
	//Needed because we can't rely on the existing shell on Windows.
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "pngquant", "-", "--speed", string(speed))
	} else {
		cmd = exec.Command("pngquant", "-", "--speed", string(speed))
	}
	cmd.Stdin = strings.NewReader(string(input))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
