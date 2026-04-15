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
	"path/filepath"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
)

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

func usage(name string) {
	fmt.Println("Roku-Cell - An Sixel Image Viewer")
	fmt.Println()
	fmt.Printf("Usage: %s [OPTIONS] PATH\n", name)
	fmt.Println()
	fmt.Println("PATH: Filename of image file")
	fmt.Println("      (BMP, GIF, JPEG, PNG, TIFF, WebP)")
	fmt.Println("OPTIONS:")
	fmt.Println("  -n N: Use N colors (N: 8 - 255)")
	fmt.Println("  -d: Disable dithering")
	fmt.Println("  -m: Disable median cut")
	fmt.Println("  -f: Disable fitting")
	fmt.Println("  -c: Run calibration")
}

func loadConfig() (string, *rkcel.Config) {
	config := rkcel.DefaultConfig()

	var path string
	dir, err := os.UserConfigDir()
	if err == nil {
		path = filepath.Join(dir, "rkcel", "config.yaml")
	}

	if path != "" {
		_, err := os.Stat(path)
		if err == nil { // file exists
			config = rkcel.LoadConfig(path)
		}
	}

	return path, config
}

func calibrate() {
	path, config := loadConfig()

	rkcel.Calibrate(config)

	if path != "" {
		rkcel.SaveConfig(path, config)
	}
}

func main() {
	numColors := flag.Int("n", 216, "number of colors used (8 - 255)")
	noDither := flag.Bool("d", false, "disable dithering")
	noMedian := flag.Bool("m", false, "disable median cut")
	noFit := flag.Bool("f", false, "disable fitting")
	runCalib := flag.Bool("c", false, "run calibration")
	flag.Parse()

	var in io.Reader
	args := flag.Args()

	if *runCalib {
		calibrate()
		return
	}

	if len(args) < 1 {
		usage(os.Args[0])
		return
	}

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

	if *numColors < 8 || *numColors > 255 {
		fatal(errors.New("number of colors must be 8 - 255"))
	}

	img, _, err := image.Decode(in)
	if err != nil {
		fatal(err)
	}

	if !*noFit {
		_, config := loadConfig()
		img = rkcel.Fit(config, img)
	}

	rkcel.Print(img, *numColors, !*noDither, !*noMedian)
	fmt.Print("\n")
}
