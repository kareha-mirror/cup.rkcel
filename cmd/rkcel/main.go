package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/webp"

	"tea.kareha.org/cup/rkcel"
)

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

	rkcel.Print(img)
	fmt.Print("\n")
}
