package rkcel

import (
	"image"
	"image/draw"

	xdraw "golang.org/x/image/draw"
)

type Method int

const (
	CatmullRom Method = iota
	ApproxBilinear
	NearestNeighbor
)

func resize(
	src image.Image, bounds image.Rectangle, w, h int, method Method,
) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	switch method {
	case CatmullRom:
		xdraw.CatmullRom.Scale(
			dst, dst.Bounds(),
			src, bounds,
			draw.Over, nil,
		)
	case ApproxBilinear:
		xdraw.ApproxBiLinear.Scale(
			dst, dst.Bounds(),
			src, bounds,
			draw.Over, nil,
		)
	case NearestNeighbor:
		xdraw.NearestNeighbor.Scale(
			dst, dst.Bounds(),
			src, bounds,
			draw.Over, nil,
		)
	default:
		panic("unknown method")
	}
	return dst
}

func Contain(img image.Image, w, h int, method Method) image.Image {
	size := img.Bounds().Size()
	scale := min(float32(w)/float32(size.X), float32(h)/float32(size.Y))
	return resize(
		img,
		img.Bounds(),
		int(float32(size.X)*scale), int(float32(size.Y)*scale),
		method,
	)
}

func Cover(img image.Image, w, h int, method Method) image.Image {
	var bounds image.Rectangle
	size := img.Bounds().Size()
	srcRatio := float32(size.X) / float32(size.Y)
	dstRatio := float32(w) / float32(h)
	if srcRatio >= dstRatio {
		width := int(float32(size.X) * dstRatio / srcRatio)
		dx := (size.X - width) / 2
		bounds = image.Rect(
			img.Bounds().Min.X+dx,
			img.Bounds().Min.Y,
			img.Bounds().Min.X+dx+width,
			img.Bounds().Min.Y+size.Y,
		)
	} else {
		height := int(float32(size.Y) * srcRatio / dstRatio)
		dy := (size.Y - height) / 2
		bounds = image.Rect(
			img.Bounds().Min.X,
			img.Bounds().Min.Y+dy,
			img.Bounds().Min.X+size.X,
			img.Bounds().Min.Y+dy+height,
		)
	}
	return resize(
		img,
		bounds,
		w, h,
		method,
	)
}
