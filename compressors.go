package imager

import (
	"image/jpeg"
	"image/png"
	"os/exec"
)

var pngEncoder = png.Encoder{png.BestCompression}
var jpgOptions = jpeg.Options{jpeg.DefaultQuality}

//Compresses a PNG file to a more compressed PNG
func compressPNG(fileName string) (string, error) {
	const prefix = "small_"
	if err := exec.Command("pngquant", fileName, prefix+fileName).Run(); err != nil {
		return "", err
	}
	return prefix + fileName, nil
}

//Compresses a JPG file to a more compressed JPG
func compressJPG(fileName string) (string, error) {
	const prefix = "small_"
	if err := exec.Command("command", fileName, prefix+fileName).Run(); err != nil {
		return "", err
	}
	return prefix + fileName, nil
}
