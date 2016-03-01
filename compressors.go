package imager

import "os/exec"

//Compresses a PNG file to a more compressed PNG
func compressPNG(fileName string) (string, error) {
	if err := exec.Command("pngquant", fileName, "small_"+fileName).Run(); err != nil {
		return "", err
	}
	return "small_" + fileName, nil
}
