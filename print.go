package rkcel

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"io"
	"os"
	"sort"
)

func Print(img image.Image, numColors int, dither bool, median bool) error {
	return Fprint(os.Stdout, img, numColors, dither, median)
}

func cbrt8(n int) int {
	for i := 0; i <= 6; i++ { // 6^3 = 216
		if i*i*i > n {
			return i - 1
		}
	}
	return 6
}

func createPalette(numColors int) []color.Color {
	n := cbrt8(numColors)
	if n < 2 {
		n = 2
	}
	palette := make([]color.Color, n*n*n)
	k := 255 / (n - 1)
	idx := 0
	for r := 0; r < n; r++ {
		for g := 0; g < n; g++ {
			for b := 0; b < n; b++ {
				palette[idx] = color.RGBA{
					uint8(r * k), uint8(g * k), uint8(b * k), 255,
				}
				idx++
			}
		}
	}
	return palette
}

func Fprint(
	w io.Writer,
	img image.Image,
	numColors int,
	dither bool,
	median bool,
) error {
	var palette []color.Color
	if median {
		palette = MedianCut(img, min(numColors, 255))
	} else {
		palette = createPalette(numColors)
	}

	bounds := img.Bounds()
	dst := image.NewPaletted(bounds, palette)
	if dither {
		draw.FloydSteinberg.Draw(dst, bounds, img, bounds.Min)
	} else {
		draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
	}

	bw := bufio.NewWriter(w)
	enc := NewEncoder(bw)

	err := enc.Start()
	if err != nil {
		return err
	}
	for idx, color := range palette {
		r, g, b, _ := color.RGBA()
		err = enc.SetPalette(idx, int(r>>8), int(g>>8), int(b>>8))
		if err != nil {
			return err
		}
	}

	size := bounds.Size()
	width := size.X
	height := size.Y

	for y := 0; y < height; y += 6 {
		used := map[int]bool{}

		for dy := 0; dy < 6; dy++ {
			yy := y + dy
			if yy >= height {
				break
			}

			row := dst.Pix[yy*dst.Stride : yy*dst.Stride+width]
			for _, c := range row {
				used[int(c)] = true
			}
		}

		idxes := make([]int, 0, len(used))
		for idx := range used {
			idxes = append(idxes, idx)
		}
		sort.Ints(idxes)

		for _, idx := range idxes {
			err = enc.SetIndex(idx)
			if err != nil {
				return err
			}

			maxdy := 6
			if y+6 > height {
				maxdy = height - y
			}
			for x := 0; x < width; x++ {
				var six byte = 0
				k := y*dst.Stride + x

				for dy := 0; dy < maxdy; dy++ {
					if dst.Pix[k] == uint8(idx) {
						six |= 1 << dy
					}
					k += dst.Stride
				}

				err = enc.PutSixel(int(six))
				if err != nil {
					return err
				}
			}

			err = enc.CarriageReturn()
			if err != nil {
				return err
			}
		}

		err = enc.LineFeed()
		if err != nil {
			return err
		}
	}

	err = enc.End()
	if err != nil {
		return err
	}

	return bw.Flush()
}
