package rkcel

import (
	"fmt"
	"io"
)

type Encoder struct {
	w io.Writer

	index int
	found bool
	six   int
	count int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,

		index: -1,
		found: false,
		six:   -1,
		count: 0,
	}
}

func (enc *Encoder) Start() error {
	err := enc.flushLine()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(enc.w, "\x1bPq")
	return err
}

func (enc *Encoder) End() error {
	err := enc.flushLine()
	if err == nil {
		_, err = fmt.Fprint(enc.w, "\x1b\\")
	}
	enc.index = -1
	return err
}

func (enc *Encoder) CarriageReturn() error {
	err := enc.flushLine()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(enc.w, "$")
	return err
}

func (enc *Encoder) LineFeed() error {
	err := enc.flushLine()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(enc.w, "-")
	return err
}

func (enc *Encoder) SetPalette(idx, r, g, b int) error {
	pr := 100 * r / 255
	pg := 100 * g / 255
	pb := 100 * b / 255
	_, err := fmt.Fprintf(enc.w, "#%d;2;%d;%d;%d", idx, pr, pg, pb)
	return err
}

func (enc *Encoder) SetIndex(idx int) error {
	enc.index = idx
	return nil
}

func (enc *Encoder) PutSixel(s int) error {
	if s == enc.six {
		enc.count++
		return nil
	} else {
		err := enc.flush()

		enc.six = s
		enc.count = 1

		if s != 0 {
			if !enc.found {
				if err == nil {
					_, err = fmt.Fprintf(enc.w, "#%d", enc.index)
				}
			}
			enc.found = true
		}

		return err
	}
}

func (enc *Encoder) flush() error {
	c := enc.six + 63

	var err error
	switch enc.count {
	case 0:
	case 1:
		_, err = fmt.Fprintf(enc.w, "%c", c)
	case 2:
		_, err = fmt.Fprintf(enc.w, "%c%c", c, c)
	case 3:
		_, err = fmt.Fprintf(enc.w, "%c%c%c", c, c, c)
	default:
		_, err = fmt.Fprintf(enc.w, "!%d%c", enc.count, c)
	}

	enc.six = -1
	enc.count = 0

	return err
}

func (enc *Encoder) flushLine() error {
	var err error
	if enc.found {
		err = enc.flush()
		enc.found = false
	} else {
		enc.six = -1
		enc.count = 0
	}
	return err
}
