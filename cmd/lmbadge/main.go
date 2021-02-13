package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/dirtykastro/laughingmanbadge/badge"
)

func main() {

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

	lmBadge := &badge.Badge{Img: badgeNoText}

	/*// Read the font data.
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

	fg := image.NewUniform(color.NRGBA{255, 128, 128, 128})

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(30.5)
	c.SetClip(im.Bounds())
	c.SetDst(im)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 10+int(c.PointToFixed(30.5)>>6))
	_, err = c.DrawString("A", pt)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/

	im := lmBadge.Render(500, "This is a test", 0)

	out, err := os.Create("lm.png")
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer out.Close()

	png.Encode(out, im)
}
