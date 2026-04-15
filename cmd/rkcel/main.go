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
	"path/filepath"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
)

func fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
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
	fmt.Println("  -c: Run calibration")
	fmt.Println("  -f: Disable fitting")
	fmt.Println("  -sb: Approximate bilinear scaling")
	fmt.Println("  -sn: Nearest neighbor scaling")
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
			config, err = rkcel.LoadConfig(path)
			if err != nil {
				fatalf("%v", err)
			}
		}
	}

	return path, config
}

func calibrate() {
	path, config := loadConfig()

	rkcel.Calibrate(config)

	if path != "" {
		err := rkcel.SaveConfig(path, config)
		if err != nil {
			fatalf("%v", err)
		}
	}
}

func main() {
	numColors := flag.Int("n", 216, "number of colors used (8 - 255)")
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
	flag.Parse()
	args := flag.Args()

	if *runCalib {
		calibrate()
		return
	}

	if len(args) < 1 {
		usage(os.Args[0])
		return
	}

	if *numColors < 8 || *numColors > 255 {
		fatalf("number of colors must be 8 - 255")
	}

	var in io.Reader
	if args[0] == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(args[0])
		if err != nil {
			fatalf("%v", err)
		}
		defer f.Close()
		in = f
	}

	img, _, err := image.Decode(in)
	if err != nil {
		fatalf("%v", err)
	}

	if !*noFit {
		_, config := loadConfig()
		var method rkcel.ScaleMethod = rkcel.ScaleCatmullRom
		if *scaleApproxBilinear {
			method = rkcel.ScaleApproxBilinear
		}
		if *scaleNearestNeighbor {
			method = rkcel.ScaleNearestNeighbor
		}
		img = rkcel.Fit(config, img, method)
	}

	rkcel.Print(img, *numColors, !*noDither, !*noMedian)
	fmt.Print("\n")
}
