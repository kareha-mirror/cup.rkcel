package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
)

func main() {
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

	img, _, err := image.Decode(in)
	if err != nil {
		panic(err)
	}

	rkcel.Print(img, *dither, *median)
	fmt.Print("\n")
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
