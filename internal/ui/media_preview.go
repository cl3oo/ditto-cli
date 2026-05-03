package ui

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

var asciiRamp = []rune("@%#*+=-:.")

func renderASCIIImage(data []byte, maxWidth int) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	if w == 0 || h == 0 {
		return "", nil
	}
	if maxWidth <= 0 {
		maxWidth = 32
	}
	outW := min(maxWidth, w)
	outH := max(1, h*outW/max(1, w*2))

	var b strings.Builder
	for y := 0; y < outH; y++ {
		if y > 0 {
			b.WriteByte('\n')
		}
		sy := bounds.Min.Y + y*h/outH
		for x := 0; x < outW; x++ {
			sx := bounds.Min.X + x*w/outW
			r, g, bl, _ := img.At(sx, sy).RGBA()
			lum := (299*int(r/257) + 587*int(g/257) + 114*int(bl/257)) / 1000
			idx := lum * (len(asciiRamp) - 1) / 255
			b.WriteRune(asciiRamp[idx])
		}
	}
	return b.String(), nil
}
