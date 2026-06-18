package rkcel

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"tea.kareha.org/cup/termi"
)

var fgColor = color.RGBA{255, 255, 0, 255}
var bgColor = color.RGBA{0, 0, 255, 255}
var gridColor = color.RGBA{0, 255, 0, 255}

const borderSize = 6

const maxAccel = 16
const minSize = borderSize * 3

func newBorderImage(
	width, height int,
	cellW, cellH int,
	cols, rows int,
) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// background fill
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// grid dots
	if cellW > 0 && cellH > 0 {
		for y := 0; y < height; y += cellH {
			for x := 0; x < width; x += cellW {
				img.Set(x, y, gridColor)
			}
		}
	}

	// cross
	h := cellH * rows
	for x := 0; x < width; x++ {
		img.Set(x, h/2, fgColor)
	}
	w := cellW * cols
	for y := 0; y < height; y++ {
		img.Set(w/2, y, fgColor)
	}

	// borders
	for x := 0; x < width; x++ {
		for b := 0; b < borderSize; b++ {
			img.Set(x, b, fgColor)              // top
			img.Set(x, cellH*rows-1-b, fgColor) // bottom
		}
	}

	// borders
	for y := 0; y < height; y++ {
		for b := 0; b < borderSize; b++ {
			img.Set(b, y, fgColor)              // left
			img.Set(cellW*cols-1-b, y, fgColor) // right
		}
	}

	return img
}

func Calibrate(cfg *Config) {
	termi.Raw()
	fmt.Print(termi.SetAlternate)
	fmt.Print(termi.HideCursor)
	termi.StartInput()

	cellW, cellH := cfg.CellWidth, cfg.CellHeight

	var cols, rows int
	prevCols, prevRows := -1, -1

	var width, height int
	prevWidth, prevHeight := -1, -1

	var img *image.RGBA

	accel := 1

	hasPrev := false
	var prevSeq termi.SeqKind
	var prevRune rune

	tmux := os.Getenv("TMUX") != "" ||
		strings.Contains(os.Getenv("TERM"), "tmux")

	count := 0
loop:
	for {
		cols, rows = termi.Size()
		if cols != prevCols || rows != prevRows {
			width = cellW * cols
			height = cellH * (rows - 1)
			prevCols = cols
			prevRows = rows
		}
		cellW = width / cols
		if rows > 1 {
			cellH = height / (rows - 1)
		}

		if width != prevWidth || height != prevHeight {
			img = newBorderImage(width, height, cellW, cellH, cols, rows-1)
			prevWidth = width
			prevHeight = height
		}

		fmt.Print(termi.Clear)
		fmt.Print(termi.HomeCursor)

		Print(img, 8, false, false)

		if !tmux || count%2 == 0 {
			fmt.Print(termi.MoveCursor(2, 1))
			fmt.Printf("* Use Arrow Keys to Fit the Rectangle to Screen *")
			fmt.Print(termi.MoveCursor(2, 2))
			fmt.Printf("             * Push Enter to Exit *              ")
			fmt.Print(termi.MoveCursor(10, 4))
			fmt.Printf("CellWidth = %d, CellHeight = %d", cellW, cellH)

			fmt.Print(termi.MoveCursor(0, rows-1))
			fmt.Printf("[ Bottm Line Reserved ]")
		}

		seq := termi.ReadSeq()
		switch seq.Kind {
		case termi.SeqRune:
			switch seq.Rune {
			case termi.RuneEscape:
				break loop
			case termi.RuneEnter:
				break loop
			case 'h':
				width = max(minSize, width-accel)
			case 'j':
				height += accel
			case 'k':
				height = max(minSize, height-accel)
			case 'l':
				width += accel
			case 'q':
				break loop
			}
		case termi.SeqUp:
			height = max(minSize, height-accel)
		case termi.SeqDown:
			height += accel
		case termi.SeqRight:
			width += accel
		case termi.SeqLeft:
			width = max(minSize, width-accel)
		}

		if hasPrev && prevSeq == seq.Kind && prevRune == seq.Rune {
			accel = min(maxAccel, accel+1)
		} else {
			accel = 1
		}
		hasPrev = true
		prevSeq = seq.Kind
		prevRune = seq.Rune

		count++
	}

	termi.StopInput()
	fmt.Print(termi.Clear)
	fmt.Print(termi.HomeCursor)
	fmt.Print(termi.ResetAlternate)
	termi.Cooked()
	fmt.Print(termi.ShowCursor)

	cfg.CellWidth = cellW
	cfg.CellHeight = cellH
}

func UserCalibrate() error {
	cfg, err := LoadUserConfig()
	if err != nil {
		return err
	}

	Calibrate(cfg)

	return SaveUserConfig(cfg)
}
