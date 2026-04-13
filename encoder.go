package rkcel

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"sort"
)

type Encoder struct {
	w io.Writer

	color   int
	found   bool
	current int
	count   int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,

		color:   -1,
		found:   false,
		current: -1,
		count:   0,
	}
}

func (enc *Encoder) Start() {
	enc.FlushLine()

	fmt.Fprint(enc.w, "\x1bPq")
}

func (enc *Encoder) End() {
	enc.FlushLine()

	fmt.Fprint(enc.w, "\x1b\\")

	enc.color = -1
}

func (enc *Encoder) Return() {
	enc.FlushLine()

	fmt.Fprint(enc.w, "$")
}

func (enc *Encoder) Newline() {
	enc.FlushLine()

	fmt.Fprint(enc.w, "-")
}

func (enc *Encoder) Palette(i, r, g, b int) {
	pr := 100 * r / 255
	pg := 100 * g / 255
	pb := 100 * b / 255
	fmt.Fprintf(enc.w, "#%d;2;%d;%d;%d", i, pr, pg, pb)
}

func (enc *Encoder) Color(i int) {
	enc.color = i
}

func (enc *Encoder) Put(c int) {
	if c == enc.current {
		enc.count++
		return
	} else {
		enc.Flush()

		enc.current = c
		enc.count = 1

		if c != 0 {
			if !enc.found {
				fmt.Fprintf(enc.w, "#%d", enc.color)
			}
			enc.found = true
		}
	}
}

func (enc *Encoder) Flush() {
	c := enc.current + 63

	if enc.count > 3 {
		fmt.Fprintf(enc.w, "!%d%c", enc.count, c)
	} else {
		switch enc.count {
		case 3:
			fmt.Fprintf(enc.w, "%c%c%c", c, c, c)
		case 2:
			fmt.Fprintf(enc.w, "%c%c", c, c)
		case 1:
			fmt.Fprintf(enc.w, "%c", c)
		}
	}

	enc.current = -1
	enc.count = 0
}

func (enc *Encoder) FlushLine() {
	if enc.found {
		enc.Flush()

		enc.found = false
	} else {
		enc.current = -1
		enc.count = 0
	}
}

func Print(img image.Image, dither bool, median bool) {
	Fprint(os.Stdout, img, dither, median)
}

var defPal []color.Color = make([]color.Color, 216)

func init() {
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				i := r*36 + g*6 + b
				defPal[i] = color.RGBA{
					uint8(r * 51), uint8(g * 51), uint8(b * 51), 255,
				}
			}
		}
	}
}

func Fprint(w io.Writer, img image.Image, dither bool, median bool) {
	rect := img.Bounds()
	pal := defPal
	dst := image.NewPaletted(rect, pal)
	if dither {
		draw.FloydSteinberg.Draw(dst, rect, img, rect.Min)
	} else {
		draw.Draw(dst, rect, img, rect.Min, draw.Src)
	}

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	enc := NewEncoder(bw)

	enc.Start()
	for c, p := range pal {
		r, g, b, _ := p.RGBA()
		enc.Palette(c, int(r>>8), int(g>>8), int(b>>8))
	}

	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y

	for y := 0; y < height; y += 6 {
		used := map[int]bool{}

		for dy := 0; dy < 6; dy++ {
			yy := y + dy
			if yy >= height {
				break
			}

			for x := 0; x < width; x++ {
				c := dst.ColorIndexAt(x, yy)
				used[int(c)] = true
			}
		}

		colors := make([]int, 0, len(used))
		for c := range used {
			colors = append(colors, c)
		}
		sort.Ints(colors)

		for _, col := range colors {
			enc.Color(col)

			for x := 0; x < width; x++ {

				var six byte = 0

				for dy := 0; dy < 6; dy++ {
					yy := y + dy
					if yy >= height {
						break
					}

					c := dst.ColorIndexAt(x, yy)
					if int(c) == col {
						six |= 1 << dy
					}
				}

				enc.Put(int(six))
			}
			enc.Return()
		}
		enc.Newline()
	}

	enc.End()
}
