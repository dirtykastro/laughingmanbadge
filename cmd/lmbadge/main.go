package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
)

func main() {
	badgeSize := 500

	var badgeNoText image.Image

	file, err := os.Open("laughing_man.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	badgeNoText, err = png.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	file.Close()

	resizedBadge := resize.Thumbnail(uint(badgeSize), uint(badgeSize), badgeNoText, resize.Lanczos3)

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: badgeSize, Y: badgeSize}})

	for x := 0; x < badgeSize; x++ {

		for y := 0; y < badgeSize; y++ {
			r, g, b, a := resizedBadge.At(x, y).RGBA()
			im.SetRGBA(x, y, color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		}
	}

	// Read the font data.
	fontBytes, err := ioutil.ReadFile("./fonts/AWAKE.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fg := image.NewUniform(color.NRGBA{255, 128, 128, 255})

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(30.5)
	c.SetClip(im.Bounds())
	c.SetDst(im)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(30.5)>>6))
	_, err = c.DrawString("KOTOKO", pt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create("lm.png")
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer out.Close()

	png.Encode(out, im)
}
