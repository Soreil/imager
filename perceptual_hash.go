package imager

import (
	"image"
	"image/draw"

	"github.com/nfnt/resize"
)

//type phash uint8
//
//const (
//	empty phash = 0
//	full  phash = 255
//)
//
////It's actually quite simple.
////>grab first frame of image, if gif
////>remove alpha channel NOTE(sjon): Just mask out button 8 bits(?)
////>downsample to 160x160
////>convert to grayscale NOTE(sjon): Oh we blur the grayscale. I guess I don't have to manually drop alpha then!
////>blur with a 2x2 area NOTE(sjon): sum and divide by x*y
////>equalize NOTE(sjon): Eh?
////>scale down to 16x16
////>convert to 1 bit colour depth NOTE(sjon): >127 = 1 else 0
////>output as RAW []byte NOTE(sjon): make that a byte, alternatively bit array
//func perceptualHash(img image.Image) phash {
//	img = resize.Resize(160, 160, img, resize.NearestNeighbor)
//
//	bounds := img.Bounds()
//	gray := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
//	draw.Draw(gray, gray.Bounds(), img, bounds.Min, draw.Src)
//	return empty
//}

type phash uint64

const (
	empty phash = 0
	full  phash = 0xffffffff
)

//Extremely simple hashing algorithm
func perceptualHash(img image.Image) phash {
	const w = 8
	const h = 8
	img = resize.Resize(w, h, img, resize.NearestNeighbor)

	bounds := img.Bounds()
	gray := image.NewGray(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(gray, gray.Bounds(), img, bounds.Min, draw.Src)

	var sum int
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			sum += int(gray.GrayAt(x, y).Y)
		}
	}
	mean := uint8(sum / (w * h))

	var hash phash
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if gray.GrayAt(x, y).Y > mean {
				hash += 1 << uint((x*w)+y)
			}
		}
	}
	return hash
}
