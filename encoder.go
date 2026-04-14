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

func Print(img image.Image, colors int, dither bool, median bool) {
	Fprint(os.Stdout, img, colors, dither, median)
}

func cbrt8(n int) int {
	for i := 0; i <= 6; i++ { // 6^3 = 216
		if i*i*i > n {
			return i - 1
		}
	}
	return 6
}

func flatPal(colors int) []color.Color {
	n := cbrt8(colors)
	var pal []color.Color = make([]color.Color, n*n*n)
	for r := 0; r < n; r++ {
		for g := 0; g < n; g++ {
			for b := 0; b < n; b++ {
				i := r*n*n + g*n + b
				k := 255 / (n - 1)
				pal[i] = color.RGBA{
					uint8(r * k), uint8(g * k), uint8(b * k), 255,
				}
			}
		}
	}
	return pal
}

func Fprint(
	w io.Writer,
	img image.Image,
	colors int,
	dither bool,
	median bool,
) {
	rect := img.Bounds()
	var pal []color.Color
	if median {
		pal = MedianCut(img, colors)
	} else {
		pal = flatPal(colors)
	}
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

		cols := make([]int, 0, len(used))
		for c := range used {
			cols = append(cols, c)
		}
		sort.Ints(cols)

		for _, col := range cols {
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
