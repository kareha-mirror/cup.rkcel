package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"time"

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
  -wait: Wait enter key
`, name)
}

type Options struct {
	numColors            *int
	noDither             *bool
	noMedian             *bool
	runCalib             *bool
	noFit                *bool
	scaleApproxBilinear  *bool
	scaleNearestNeighbor *bool
	cover                *bool
}

func main() {
	opt := Options{}

	opt.numColors = flag.Int("n", 255, "number of colors used (max 255)")
	opt.noDither = flag.Bool("d", false, "disable dithering")
	opt.noMedian = flag.Bool("m", false, "disable median cut")
	opt.runCalib = flag.Bool("c", false, "run calibration")
	opt.noFit = flag.Bool("f", false, "disable fitting")
	opt.scaleApproxBilinear = flag.Bool(
		"sb", false, "approximate bilinear scaling",
	)
	opt.scaleNearestNeighbor = flag.Bool(
		"sn", false, "nearest neighbor scaling",
	)
	opt.cover = flag.Bool("cover", false, "cover fitting")
	wait := flag.Bool("wait", false, "wait enter key")

	flag.Parse()
	args := flag.Args()

	if *opt.runCalib {
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

	r := rkcel.NewCacheReader(in)
	_, format, err := image.DecodeConfig(r)
	if err != nil {
		fatal(err)
	}
	r.Rewind()

	if format == "gif" {
		g, err := gif.DecodeAll(r)
		if err != nil {
			fatal(err)
		}

		if len(g.Image) < 2 {
			if len(g.Image) > 0 {
				err = print(&opt, g.Image[0])
				if err != nil {
					fatal(err)
				}
				fmt.Print("\n")
			}
			return
		}

		termi.HideCursor()
		defer termi.ShowCursor()

		next := time.Now()
		for k := 0; k < g.LoopCount || g.LoopCount < 1; k++ {
			for i, img := range g.Image {
				termi.MoveCursor(0, 0)
				err = print(&opt, img)
				if err != nil {
					fatal(err)
				}

				delay := g.Delay[i]
				if delay == 0 {
					delay = 10
				}
				next = next.Add(time.Duration(delay) * 10 * time.Millisecond)
				time.Sleep(time.Until(next))
			}
		}
	} else {
		img, _, err := image.Decode(r)
		if err != nil {
			fatal(err)
		}
		err = print(&opt, img)
		if err != nil {
			fatal(err)
		}
	}
	fmt.Print("\n")

	if *wait {
		fmt.Print("Press Enter to Continue")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func print(opt *Options, img image.Image) error {
	if !*opt.noFit {
		cfg, err := rkcel.LoadUserConfig()
		if err != nil {
			return err
		}
		w, h := termi.Size()
		maxW := cfg.CellWidth * w
		maxH := cfg.CellHeight * (h - 1)
		size := img.Bounds().Size()
		if size.X > maxW || size.Y > maxH {
			var method = rkcel.CatmullRom
			if *opt.scaleApproxBilinear {
				method = rkcel.ApproxBilinear
			}
			if *opt.scaleNearestNeighbor {
				method = rkcel.NearestNeighbor
			}
			if *opt.cover {
				img = rkcel.FitCover(img, maxW, maxH, method)
			} else {
				img = rkcel.FitContain(img, maxW, maxH, method)
			}
		}
	}

	dst := image.NewRGBA(img.Bounds())
	bg := image.NewUniform(color.White)
	draw.Draw(dst, dst.Bounds(), bg, image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), img, image.Point{}, draw.Over)

	return rkcel.Print(dst, *opt.numColors, !*opt.noDither, !*opt.noMedian)
}
