package imager

import (
	"os"
	"testing"
)

func TestDecodeVP8(t *testing.T) {
	var fileName = "test.webm"
	if _, err := os.Stat(fileName); err != nil {
		t.Fatal(err)
	}
	thumbnailWebm(fileName, "out.png")
}
