package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/dirtykastro/laughingmanbadge/badge"
)

//go:embed laughing_man.png
//go:embed RobotoMono-Medium.ttf

var f embed.FS

func main() {

	text := flag.String("text", "I thought what I'd do was, I'd pretend I was one of those deaf-mutes.", "text to display in the badge")
	rotation := flag.Float64("rotation", 0.0, "rotate text angle")
	size := flag.Int("size", 200, "size of badge")
	outputFile := flag.String("file", "", "destination file name")

	flag.Parse()

	if *outputFile == "" {
		fmt.Println("the destination file is required")
		flag.PrintDefaults()
		os.Exit(0)
	}

	var badgeNoText image.Image

	file, err := f.ReadFile("laughing_man.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	badgeNoText, err = png.Decode(bytes.NewReader(file))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	fontFile, err := f.ReadFile("RobotoMono-Medium.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lmBadge := &badge.Badge{Img: badgeNoText, FontFile: fontFile}

	im, err := lmBadge.Render(*size, *text, *rotation)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer out.Close()

	png.Encode(out, im)
}
