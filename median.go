package rkcel

import (
	"image"
	"image/color"
	"sort"
)

type rgb struct {
	r, g, b int
}

type box struct {
	colors []rgb
}

func (b box) rangeR() int {
	min, max := 255, 0
	for _, c := range b.colors {
		if c.r < min {
			min = c.r
		}
		if c.r > max {
			max = c.r
		}
	}
	return max - min
}

func (b box) rangeG() int {
	min, max := 255, 0
	for _, c := range b.colors {
		if c.g < min {
			min = c.g
		}
		if c.g > max {
			max = c.g
		}
	}
	return max - min
}

func (b box) rangeB() int {
	min, max := 255, 0
	for _, c := range b.colors {
		if c.b < min {
			min = c.b
		}
		if c.b > max {
			max = c.b
		}
	}
	return max - min
}

func splitBox(b box) (box, box) {
	rRange := b.rangeR()
	gRange := b.rangeG()
	bRange := b.rangeB()

	switch {
	case rRange >= gRange && rRange >= bRange:
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].r < b.colors[j].r
		})
	case gRange >= rRange && gRange >= bRange:
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].g < b.colors[j].g
		})
	default:
		sort.Slice(b.colors, func(i, j int) bool {
			return b.colors[i].b < b.colors[j].b
		})
	}

	mid := len(b.colors) / 2
	return box{b.colors[:mid]}, box{b.colors[mid:]}
}

func average(b box) color.Color {
	var r, g, bl int
	for _, c := range b.colors {
		r += c.r
		g += c.g
		bl += c.b
	}
	n := len(b.colors)
	if n == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	return color.RGBA{
		uint8(r / n),
		uint8(g / n),
		uint8(bl / n),
		255,
	}
}

func MedianCut(img image.Image, maxColors int) []color.Color {
	bounds := img.Bounds()

	pixels := make([]rgb, 0, bounds.Dx()*bounds.Dy())
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, rgb{
				int(r >> 8),
				int(g >> 8),
				int(b >> 8),
			})
		}
	}

	boxes := []box{{pixels}}

	for len(boxes) < maxColors {
		sort.Slice(boxes, func(i, j int) bool {
			a := boxes[i]
			b := boxes[j]
			return max3(a.rangeR(), a.rangeG(), a.rangeB()) >
				max3(b.rangeR(), b.rangeG(), b.rangeB())
		})

		b := boxes[0]
		if len(b.colors) <= 1 {
			break
		}

		b1, b2 := splitBox(b)
		boxes = append(boxes[1:], b1, b2)
	}

	pal := make([]color.Color, 0, len(boxes))
	for _, b := range boxes {
		pal = append(pal, average(b))
	}

	return pal
}

func max3(a, b, c int) int {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}
