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
	termi.HideCursor()

	cellW, cellH := cfg.CellWidth, cfg.CellHeight

	var cols, rows int
	prevCols, prevRows := -1, -1

	var width, height int
	prevWidth, prevHeight := -1, -1

	var img *image.RGBA

	accel := 1

	hasPrev := false
	var prevKey termi.KeyKind
	var prevRune rune

	tmux := false
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "TMUX=") {
			tmux = true
			break
		}
	}

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

		termi.Clear()
		termi.HomeCursor()

		Print(img, 8, false, false)

		if !tmux || count%2 == 0 {
			termi.MoveCursor(2, 1)
			fmt.Printf("* Use Arrow Keys to Fit the Rectangle to Screen *")
			termi.MoveCursor(2, 2)
			fmt.Printf("             * Push Enter to Exit *              ")
			termi.MoveCursor(10, 4)
			fmt.Printf("CellWidth = %d, CellHeight = %d", cellW, cellH)

			termi.MoveCursor(0, rows-1)
			fmt.Printf("[ Bottm Line Reserved ]")
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
		case termi.KeyUp:
			height = max(minSize, height-accel)
		case termi.KeyDown:
			height += accel
		case termi.KeyRight:
			width += accel
		case termi.KeyLeft:
			width = max(minSize, width-accel)
		}

		if hasPrev && prevKey == key.Kind && prevRune == key.Rune {
			accel = min(maxAccel, accel+1)
		} else {
			accel = 1
		}
		hasPrev = true
		prevKey = key.Kind
		prevRune = key.Rune

		count++
	}

	termi.Clear()
	termi.HomeCursor()
	termi.Cooked()
	termi.ShowCursor()

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
