package imager

import (
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestSVG(t *testing.T) {
	cases := []string{"acid.svg", "adobe.svg", "alphachannel.svg", "android.svg", "The United States of America.svg"}
	for _, test := range cases {
		//img, err := svgToImage(inputDir + test)
		file, err := ioutil.ReadFile(inputDir + test)
		if err != nil {
			t.Fatal(err)
		}
		img, err := svgToImage(file)
		if err != nil {
			t.Fatal(err)
		}
		out, err := os.Create(outputDir + test + ".png")
		if err != nil {
			t.Fatal(err)
		}
		err = png.Encode(out, img)
		if err != nil {
			t.Fatal(err)
		}
	}
}
