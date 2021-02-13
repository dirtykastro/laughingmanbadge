package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	gu "github.com/dirtykastro/graphicutils"
)

var imagePath = flag.String("image", "", "image file")
var x = flag.Int("x", -1, "pixel X coordinate")
var y = flag.Int("y", -1, "pixel Y coordinate")

func main() {
	flag.Parse()

	if *imagePath == "" || *x == -1 || *y == -1 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	// parse args
	filePath := *imagePath

	img, imgErr := gu.DecodeImage(filePath)

	if imgErr != nil {
		log.Println(imgErr)
		os.Exit(1)
	}

	pixel, pixelErr := gu.GetPixelValue(img, image.Pt(*x, *y))
	if pixelErr != nil {
		log.Println(pixelErr)
		os.Exit(1)
	}

	fmt.Println("Pixel Values of coordinate [", *x, ",", *y, "]", "R:", pixel.R, "G:", pixel.G, "B:", pixel.B, "A:", pixel.A)
}
