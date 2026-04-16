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
	src image.Image,
	srcBounds image.Rectangle,
	dstBounds image.Rectangle,
	method Method,
) *image.RGBA {
	dst := image.NewRGBA(dstBounds)
	switch method {
	case CatmullRom:
		xdraw.CatmullRom.Scale(
			dst, dstBounds,
			src, srcBounds,
			draw.Over, nil,
		)
	case ApproxBilinear:
		xdraw.ApproxBiLinear.Scale(
			dst, dstBounds,
			src, srcBounds,
			draw.Over, nil,
		)
	case NearestNeighbor:
		xdraw.NearestNeighbor.Scale(
			dst, dstBounds,
			src, srcBounds,
			draw.Over, nil,
		)
	default:
		panic("unknown method")
	}
	return dst
}

func FitContain(img image.Image, w, h int, method Method) image.Image {
	bounds := img.Bounds()
	size := bounds.Size()

	scale := min(float32(w)/float32(size.X), float32(h)/float32(size.Y))
	dstBounds := image.Rect(
		0, 0, int(float32(size.X)*scale), int(float32(size.Y)*scale),
	)

	return resize(img, bounds, dstBounds, method)
}

func FitCover(img image.Image, w, h int, method Method) image.Image {
	bounds := img.Bounds()
	size := bounds.Size()

	srcRatio := float32(size.X) / float32(size.Y)
	dstRatio := float32(w) / float32(h)

	var srcBounds image.Rectangle
	if srcRatio >= dstRatio {
		width := int(float32(size.X) * dstRatio / srcRatio)
		dx := (size.X - width) / 2
		srcBounds = image.Rect(
			bounds.Min.X+dx,
			bounds.Min.Y,
			bounds.Min.X+dx+width,
			bounds.Min.Y+size.Y,
		)
	} else {
		height := int(float32(size.Y) * srcRatio / dstRatio)
		dy := (size.Y - height) / 2
		srcBounds = image.Rect(
			bounds.Min.X,
			bounds.Min.Y+dy,
			bounds.Min.X+size.X,
			bounds.Min.Y+dy+height,
		)
	}

	return resize(img, srcBounds, image.Rect(0, 0, w, h), method)
}
