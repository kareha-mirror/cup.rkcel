package rkcel

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"

	"tea.kareha.org/cup/termi"
)

func resize(src image.Image, w, h int, reduce int) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	switch reduce {
	case 2:
		xdraw.NearestNeighbor.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	case 1:
		xdraw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	default:
		xdraw.CatmullRom.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	}
	return dst
}

func Fit(config *Config, img image.Image, reduce int) image.Image {
	w, h := termi.Size()
	maxW := config.CellWidth * w
	maxH := config.CellHeight * (h - 1)
	bounds := img.Bounds()
	imgW := bounds.Max.X - bounds.Min.X
	imgH := bounds.Max.Y - bounds.Min.Y
	scale := min(float32(maxW)/float32(imgW), float32(maxH)/float32(imgH))
	if scale >= 1.0 {
		return img
	}
	return resize(
		img,
		int(float32(imgW)*scale), int(float32(imgH)*scale),
		reduce,
	)
}
