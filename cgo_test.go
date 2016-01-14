package imager

import (
	"fmt"
	"testing"
)

func TestDecodeVP8(t *testing.T) {
	fmt.Println(NewVP8Decoder().DecodeVP8("test.webm"))
}
