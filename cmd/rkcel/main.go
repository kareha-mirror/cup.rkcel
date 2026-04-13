package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

func SixelStart() {
	fmt.Print("\x1bPq")
}

func SixelEnd() {
	SixelFlush()
	fmt.Print("\x1b\\")
}

func SixelReturn() {
	SixelFlush()
	fmt.Print("$")
}

func SixelNewline() {
	SixelFlush()
	fmt.Print("-")
}

func SixelPalette(i, r, g, b int) {
	pr := 100 * r / 255
	pg := 100 * g / 255
	pb := 100 * b / 255
	fmt.Printf("#%d;2;%d;%d;%d", i, pr, pg, pb)
}

func SixelColor(i int) {
	fmt.Printf("#%d", i)
}

var current int = -1
var count int = 0

func SixelPut(c int) {
	if c == current {
		count++
		return
	} else {
		SixelFlush()
		current = c
		count = 1
	}
}

func SixelFlush() {
	if count > 3 {
		fmt.Printf("!%d%c", count, current+63)
	} else if count > 0 {
		for i := 0; i < count; i++ {
			fmt.Printf("%c", current+63)
		}
	}

	current = -1
	count = 0
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s FILENAME\n", os.Args[0])
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	rect := img.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, img, rect.Min, draw.Src)

	SixelStart()
	for r := 0; r < 6; r++ {
		for g := 0; g < 6; g++ {
			for b := 0; b < 6; b++ {
				color := r*36 + g*6 + b
				SixelPalette(color, r*51, g*51, b*51)
			}
		}
	}

	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	for y := 0; y < height; y += 6 {

		for color := 0; color < 216; color++ {
			SixelColor(color)

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

				SixelPut(int(six))
			}
			SixelReturn()
		}
		SixelNewline()
	}

	SixelEnd()
	fmt.Print("\n")
}
