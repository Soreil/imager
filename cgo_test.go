package imager

import (
	"fmt"
	"os"
	"testing"
)

func TestDecodeVP8(t *testing.T) {
	var fileName = "test.webm"
	if _, err := os.Stat(fileName); err != nil {
		t.Fatal(err)
	}
	if err := NewVP8Decoder().DecodeVP8(fileName); err != nil {
		t.Fatal("Failed to thumbnail " + fileName + ": " + fmt.Sprint(err))
	}
}
