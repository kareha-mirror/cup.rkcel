package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
)

func main() {
	colors := flag.Int("c", 216, "colors (8 - 255)")
	dither := flag.Bool("d", true, "dither")
	median := flag.Bool("m", true, "median")
	flag.Parse()

	var in io.Reader
	args := flag.Args()

	if len(args) < 1 || args[0] == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			fatal(err)
		}
		defer f.Close()
		in = f
	}

	if *colors < 8 || *colors > 255 {
		fatal(errors.New("colors must be 8 - 255"))
	}

	img, _, err := image.Decode(in)
	if err != nil {
		fatal(err)
	}

	rkcel.Print(img, *colors, *dither, *median)
	fmt.Print("\n")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
