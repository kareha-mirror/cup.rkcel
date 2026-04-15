package rkcel

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"

	"tea.kareha.org/cup/termi"
)

func resize(src image.Image, w, h int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	xdraw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

func Fit(config *Config, img image.Image) image.Image {
	w, h := termi.Size()
	maxW := config.CellWidth * w
	maxH := config.CellHeight * (h - 1)
	bounds := img.Bounds()
	imgW := bounds.Max.X - bounds.Min.X
	imgH := bounds.Max.Y - bounds.Min.Y
	scale := min(
		min(float32(maxW)/float32(imgW), float32(maxH)/float32(imgH)),
		1.0,
	)
	return resize(img, int(float32(imgW)*scale), int(float32(imgH)*scale))
}
