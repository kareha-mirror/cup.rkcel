package rkcel

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"

	"tea.kareha.org/cup/termi"
)

type ScaleMethod int

const (
	ScaleCatmullRom ScaleMethod = iota
	ScaleApproxBilinear
	ScaleNearestNeighbor
)

func resize(src image.Image, w, h int, method ScaleMethod) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	switch method {
	case ScaleCatmullRom:
		xdraw.CatmullRom.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	case ScaleApproxBilinear:
		xdraw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	case ScaleNearestNeighbor:
		xdraw.NearestNeighbor.Scale(
			dst, dst.Bounds(),
			src, src.Bounds(),
			draw.Over, nil,
		)
	default:
		panic("unknown method")
	}
	return dst
}

func Fit(config *Config, img image.Image, method ScaleMethod) image.Image {
	w, h := termi.Size()
	maxW := config.CellWidth * w
	maxH := config.CellHeight * (h - 1)
	size := img.Bounds().Size()
	scale := min(float32(maxW)/float32(size.X), float32(maxH)/float32(size.Y))
	if scale >= 1.0 {
		return img
	}
	return resize(
		img,
		int(float32(size.X)*scale), int(float32(size.Y)*scale),
		method,
	)
}
