package rkcel

import (
	"fmt"
	"image"
	"image/color"

	"tea.kareha.org/cup/termi"
)

func newBorderImage(
	width, height int,
	border int,
	col, bg color.Color,
	cw, ch int,
	w, h int,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bg)
		}
	}

	if cw > 0 && ch > 0 {
		for y := 0; y < height; y += ch {
			for x := 0; x < width; x += cw {
				img.Set(x, y, col)
			}
		}
	}

	for x := 0; x < width; x++ {
		img.Set(x, height/2, col)
	}
	for y := 0; y < height; y++ {
		img.Set(width/2, y, col)
	}

	for x := 0; x < width; x++ {
		for b := 0; b < border; b++ {
			img.Set(x, b, col)        // top
			img.Set(x, ch*h-1-b, col) // bottom
		}
	}

	for y := 0; y < height; y++ {
		for b := 0; b < border; b++ {
			img.Set(b, y, col)        // left
			img.Set(cw*w-1-b, y, col) // right
		}
	}

	return img
}

func Calibrate(config *Config) {
	termi.Raw()
	termi.HideCursor()

	cw, ch := config.CellWidth, config.CellHeight
	ub := config.UseBottom

	border := 6
	col := color.RGBA{255, 255, 0, 255}
	bg := color.RGBA{0, 128, 0, 255}

	dx := 1
	dy := 1
	accel := 1
	hasPrev := false
	var prevKey termi.KeyKind
	var prevRune rune

	var img *image.RGBA
	prevWidth, prevHeight := -1, -1

	w, h := termi.Size()
	prevW, prevH := -1, -1

	width, height := -1, -1

loop:
	for {
		w, h = termi.Size()
		if w != prevW || h != prevH {
			width = w * cw
			if ub {
				height = h * ch
			} else {
				height = (h - 1) * ch
			}
			prevW = w
			prevH = h
		}
		cw = width / w
		if ub {
			ch = height / h
		} else {
			ch = height / (h - 1)
		}

		if width != prevWidth || height != prevHeight {
			if ub {
				img = newBorderImage(width, height, border, col, bg, cw, ch, w, h)
			} else {
				img = newBorderImage(width, height, border, col, bg, cw, ch, w, h-1)
			}
			prevWidth = width
			prevHeight = height
		}

		termi.Clear()
		termi.HomeCursor()

		Print(img, 27, false, false)

		termi.MoveCursor(6, 4)
		fmt.Printf("* Use Arrow Keys to Fit the Rectangle to Screen *")
		termi.MoveCursor(6, 5)
		fmt.Printf("             * Push Enter to Exit *              ")

		termi.MoveCursor(6, 7)
		fmt.Printf("width = %d, height = %d", width, height)
		termi.MoveCursor(6, 8)
		fmt.Printf("cw = %d, ch = %d", cw, ch)

		if !ub {
			termi.MoveCursor(0, h-1)
			fmt.Printf("[ Bottm Line Reserved: Push B to Toggle ]")
		}

		key := termi.ReadKey()
		switch key.Kind {
		case termi.KeyRune:
			switch key.Rune {
			case termi.RuneEscape:
				break loop
			case termi.RuneEnter:
				break loop
			case 'h':
				width = max(64, width-dx*accel)
			case 'j':
				height += dy * accel
			case 'k':
				height = max(64, height-dy*accel)
			case 'l':
				width += dx * accel
			case 'q':
				break loop
			case 'b':
				ub = !ub

				if ub {
					height = h * ch
				} else {
					height = (h - 1) * ch
				}
			}
		case termi.KeyUp:
			height = max(64, height-dy*accel)
		case termi.KeyDown:
			height += dy * accel
		case termi.KeyRight:
			width += dx * accel
		case termi.KeyLeft:
			width = max(64, width-dx*accel)
		}

		if hasPrev && prevKey == key.Kind && prevRune == key.Rune {
			accel = min(16, accel+1)
		} else {
			accel = 1
		}
		hasPrev = true
		prevKey = key.Kind
		prevRune = key.Rune
	}

	termi.Clear()
	termi.HomeCursor()
	termi.Cooked()
	termi.ShowCursor()

	config.CellWidth = cw
	config.CellHeight = ch
	config.UseBottom = ub
}
