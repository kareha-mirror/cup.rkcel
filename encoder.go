package rkcel

import (
	"fmt"
	"image"
	"image/draw"
	"io"
	"os"
)

type Encoder struct {
	w io.Writer

	current int
	count   int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,

		current: -1,
		count:   0,
	}
}

func (enc *Encoder) Start() {
	enc.Flush()

	fmt.Fprint(enc.w, "\x1bPq")
}

func (enc *Encoder) End() {
	enc.Flush()

	fmt.Fprint(enc.w, "\x1b\\")
}

func (enc *Encoder) Return() {
	enc.Flush()

	fmt.Fprint(enc.w, "$")
}

func (enc *Encoder) Newline() {
	enc.Flush()

	fmt.Fprint(enc.w, "-")
}

func (enc *Encoder) Palette(i, r, g, b int) {
	pr := 100 * r / 255
	pg := 100 * g / 255
	pb := 100 * b / 255
	fmt.Fprintf(enc.w, "#%d;2;%d;%d;%d", i, pr, pg, pb)
}

func (enc *Encoder) Color(i int) {
	fmt.Fprintf(enc.w, "#%d", i)
}

func (enc *Encoder) Put(c int) {
	if c == enc.current {
		enc.count++
		return
	} else {
		enc.Flush()

		enc.current = c
		enc.count = 1
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

func Print(img image.Image) {
	Fprint(os.Stdout, img)
}

func Fprint(w io.Writer, img image.Image) {
	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	enc := NewEncoder(w)

	enc.Start()
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				color := r*36 + g*6 + b
				enc.Palette(color, r*51, g*51, b*51)
			}
		}
	}

	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	for y := 0; y < height; y += 6 {

		for color := 0; color < 216; color++ {
			enc.Color(color)

			for x := 0; x < width; x++ {

				var six byte = 0

				for dy := 0; dy < 6; dy++ {
					yy := y + dy
					if yy >= height {
						break
					}

					i := rgba.PixOffset(x, yy)
					r := rgba.Pix[i]
					g := rgba.Pix[i+1]
					b := rgba.Pix[i+2]

					cr := (int(r) + 25) / 51
					cg := (int(g) + 25) / 51
					cb := (int(b) + 25) / 51

					c := cr*36 + cg*6 + cb
					if int(c) == color {
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
