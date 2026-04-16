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

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
	"tea.kareha.org/cup/termi"
)

func fatal(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(1)
}

func usage(name string) {
	fmt.Printf(`Roku-Cell - A Sixel Image Viewer

Usage: %s [OPTIONS] PATH

PATH: Filename of image file
      (BMP, GIF, JPEG, PNG, TIFF, WebP)
OPTIONS:
  -n N: Use N colors (N: max 255)
  -d: Disable dithering
  -m: Disable median cut
  -c: Run calibration
  -f: Disable fitting
  -sb: Approximate bilinear scaling
  -sn: Nearest neighbor scaling
  -cover: Cover fitting
`, name)
}

func main() {
	numColors := flag.Int("n", 255, "number of colors used (max 255)")
	noDither := flag.Bool("d", false, "disable dithering")
	noMedian := flag.Bool("m", false, "disable median cut")
	runCalib := flag.Bool("c", false, "run calibration")
	noFit := flag.Bool("f", false, "disable fitting")
	scaleApproxBilinear := flag.Bool(
		"sb", false, "approximate bilinear scaling",
	)
	scaleNearestNeighbor := flag.Bool(
		"sn", false, "nearest neighbor scaling",
	)
	cover := flag.Bool("cover", false, "cover fitting")

	flag.Parse()
	args := flag.Args()

	if *runCalib {
		err := rkcel.UserCalibrate()
		if err != nil {
			fatal(err)
		}
		return
	}

	if len(args) < 1 {
		usage(os.Args[0])
		return
	}

	var in io.Reader
	if args[0] == "-" {
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
		fatal(err)
	}

	if !*noFit {
		config, err := rkcel.LoadUserConfig()
		if err != nil {
			fatal(err)
		}
		w, h := termi.Size()
		maxW := config.CellWidth * w
		maxH := config.CellHeight * (h - 1)
		size := img.Bounds().Size()
		if size.X > maxW || size.Y > maxH {
			var method = rkcel.CatmullRom
			if *scaleApproxBilinear {
				method = rkcel.ApproxBilinear
			}
			if *scaleNearestNeighbor {
				method = rkcel.NearestNeighbor
			}
			if *cover {
				img = rkcel.FitCover(img, maxW, maxH, method)
			} else {
				img = rkcel.FitContain(img, maxW, maxH, method)
			}
		}
	}

	err = rkcel.Print(img, *numColors, !*noDither, !*noMedian)
	if err != nil {
		fatal(err)
	}
	fmt.Print("\n")
}
