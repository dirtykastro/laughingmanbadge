package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

func main() {
	fmt.Println("vim-go")

	width := 100
	height := 100

	// Read the font data.
	fontBytes, err := ioutil.ReadFile("/usr/share/fonts/truetype/freefont/FreeMono.ttf")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: width, Y: height}})

	for x := 0; x < width; x++ {

		for y := 0; y < height; y++ {

			im.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})

		}

	}

	fg := image.Black

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(15)
	c.SetClip(im.Bounds())
	c.SetDst(im)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(15)>>6))
	_, err = c.DrawString("this is a test", pt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := os.Create("/home/kastro/Desktop/lm.png")
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer out.Close()

	png.Encode(out, im)
}
