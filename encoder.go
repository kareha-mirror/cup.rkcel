package rkcel

import (
	"fmt"
	"io"
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
