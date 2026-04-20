package rkcel

import (
	"image"
	"image/color"
	"math/rand"
	"sort"
)

type rgb uint32

func newRGB(r, g, b uint8) rgb {
	return rgb((uint32(r) << 16) | (uint32(g) << 8) | uint32(b))
}

func (c rgb) r() uint8 {
	return uint8(c >> 16)
}

func (c rgb) g() uint8 {
	return uint8(c >> 8)
}

func (c rgb) b() uint8 {
	return uint8(c)
}

type colorCount struct {
	c rgb
	n int
}

type box struct {
	colors []colorCount
	rRange uint8
	gRange uint8
	bRange uint8
}

func newBox(colors []colorCount) box {
	b := box{colors: colors}
	b.rRange = b.rangeR()
	b.gRange = b.rangeG()
	b.bRange = b.rangeB()
	return b
}

func (b box) rangeR() uint8 {
	var min, max uint8 = 255, 0
	for _, cc := range b.colors {
		if cc.c.r() < min {
			min = cc.c.r()
		}
		if cc.c.r() > max {
			max = cc.c.r()
		}
	}
	return max - min
}

func (b box) rangeG() uint8 {
	var min, max uint8 = 255, 0
	for _, cc := range b.colors {
		if cc.c.g() < min {
			min = cc.c.g()
		}
		if cc.c.g() > max {
			max = cc.c.g()
		}
	}
	return max - min
}

func (b box) rangeB() uint8 {
	var min, max uint8 = 255, 0
	for _, cc := range b.colors {
		if cc.c.b() < min {
			min = cc.c.b()
		}
		if cc.c.b() > max {
			max = cc.c.b()
		}
	}
	return max - min
}

func weightedMedian(colors []colorCount, axis int, target int) int {
	if len(colors) <= 2 {
		return len(colors) / 2
	}

	left, right := 0, len(colors)-1
	for left <= right { // guard
		//pivotIndex := (left + right) / 2
		pivotIndex := left + rand.Intn(right-left+1)
		pivotIndex = partition(colors, left, right, pivotIndex, axis)

		sum := 0
		for i := left; i < pivotIndex; i++ {
			sum += colors[i].n
		}

		pivotWeight := colors[pivotIndex].n
		if target < sum {
			right = pivotIndex - 1
		} else if target < sum+pivotWeight {
			return pivotIndex
		} else {
			target -= sum + pivotWeight
			left = pivotIndex + 1
		}
	}
	return left // fallback
}

func lessAxis(a, b rgb, axis int) bool {
	switch axis {
	case 0: // R
		return a.r() < b.r()
	case 1: // G
		return a.g() < b.g()
	default: // B
		return a.b() < b.b()
	}
}

func partition(
	colors []colorCount, left, right, pivotIndex int, axis int,
) int {
	pivot := colors[pivotIndex]

	colors[pivotIndex], colors[right] = colors[right], colors[pivotIndex]

	store := left
	for i := left; i < right; i++ {
		if lessAxis(colors[i].c, pivot.c, axis) {
			colors[store], colors[i] = colors[i], colors[store]
			store++
		}
	}

	colors[right], colors[store] = colors[store], colors[right]
	return store
}

func splitBox(b box) (box, box) {
	rRange := b.rRange
	gRange := b.gRange
	bRange := b.bRange

	var axis int
	switch {
	case rRange >= gRange && rRange >= bRange:
		axis = 0
	case gRange >= rRange && gRange >= bRange:
		axis = 1
	default:
		axis = 2
	}

	total := 0
	for _, cc := range b.colors {
		total += cc.n
	}
	midWeight := total / 2

	mid := weightedMedian(b.colors, axis, midWeight)
	return newBox(b.colors[:mid]), newBox(b.colors[mid:])
}

func average(b box) color.Color {
	var r, g, bl int
	n := 0
	for _, cc := range b.colors {
		r += int(cc.c.r()) * cc.n
		g += int(cc.c.g()) * cc.n
		bl += int(cc.c.b()) * cc.n
		n += cc.n
	}
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
	// for reproducibility
	//rand.Seed(1)

	bounds := img.Bounds()
	var counts map[rgb]int
	if rgba, ok := img.(*image.RGBA); ok {
		counts = make(map[rgb]int, bounds.Dx()*bounds.Dy()/4)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			offset := (y - bounds.Min.Y) * rgba.Stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := offset + (x-bounds.Min.X)*4
				r := rgba.Pix[i]
				g := rgba.Pix[i+1]
				b := rgba.Pix[i+2]
				c := newRGB(r, g, b)
				counts[c]++
			}
		}
	} else if nrgba, ok := img.(*image.NRGBA); ok {
		counts = make(map[rgb]int, bounds.Dx()*bounds.Dy()/4)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			offset := (y - bounds.Min.Y) * nrgba.Stride
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				i := offset + (x-bounds.Min.X)*4
				r := nrgba.Pix[i]
				g := nrgba.Pix[i+1]
				b := nrgba.Pix[i+2]
				c := newRGB(r, g, b)
				counts[c]++
			}
		}
	} else if pal, ok := img.(*image.Paletted); ok {
		counts = make(map[rgb]int, len(pal.Palette))
		for _, idx := range pal.Pix {
			col := pal.Palette[idx]
			r, g, b, _ := col.RGBA()
			c := newRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			counts[c]++
		}
	} else {
		counts = make(map[rgb]int, bounds.Dx()*bounds.Dy()/4)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				c := newRGB(uint8(r>>8), uint8(g>>8), uint8(b>>8))
				counts[c]++
			}
		}
	}

	var colors []colorCount
	for c, n := range counts {
		colors = append(colors, colorCount{c, n})
	}
	/*
		// for reproducibility
		sort.Slice(colors, func(i, j int) bool {
			if colors[i].c.r() != colors[j].c.r() {
				return colors[i].c.r() < colors[j].c.r()
			}
			if colors[i].c.g() != colors[j].c.g() {
				return colors[i].c.g() < colors[j].c.g()
			}
			return colors[i].c.b() < colors[j].c.b()
		})
	*/
	boxes := []box{newBox(colors)}

	for len(boxes) < maxColors {
		sort.Slice(boxes, func(i, j int) bool {
			a := boxes[i]
			b := boxes[j]
			return max3(a.rRange, a.gRange, a.bRange) >
				max3(b.rRange, b.gRange, b.bRange)
		})

		b := boxes[0]
		if len(b.colors) <= 1 {
			break
		}

		b1, b2 := splitBox(b)
		boxes = append(boxes[1:], b1, b2)
	}

	palette := make([]color.Color, 0, len(boxes))
	for _, b := range boxes {
		palette = append(palette, average(b))
	}

	return palette
}

func max3(a, b, c uint8) uint8 {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}
