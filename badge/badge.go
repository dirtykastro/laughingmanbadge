package badge

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	gu "github.com/dirtykastro/graphicutils"
)

const badgeCircleRatio = 0.75
const lettersPositionRatio = 0.65

const fontXoffset = 500
const fontYoffset = 500

type node struct {
	x, y, degree int
}

type Badge struct {
	Img      image.Image
	FontFile []byte
}

type CircleStep struct {
	Angle float64
	Sin   float64
	Cos   float64
}

func (badge *Badge) Render(badgeSize int, text string, rotation float64) (out image.Image, err error) {
	blue := color.RGBA{50, 121, 153, 255}

	resizedBadge := resize.Thumbnail(uint(badgeSize), uint(badgeSize), badge.Img, resize.Lanczos3)

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: badgeSize, Y: badgeSize}})

	var circle []node

	halfBadgeSize := float64(badgeSize) / 2

	steps := 50

	var stepsData []CircleStep

	for i := 0; i < steps; i++ {
		angle := float64(i) * -2 / float64(steps) * math.Pi
		sin, cos := math.Sincos(angle)

		circleX := halfBadgeSize + (halfBadgeSize * sin * badgeCircleRatio)
		circleY := halfBadgeSize + (halfBadgeSize * cos * badgeCircleRatio)

		stepsData = append(stepsData, CircleStep{Angle: angle, Sin: sin, Cos: cos})

		circle = append(circle, node{int(circleX), int(circleY), 1})
	}

	circle = append(circle, node{circle[0].x, circle[0].y, -1})

	r := raster.NewRasterizer(badgeSize, badgeSize)
	contour(r, circle)
	//contour(r, inside)
	mask := image.NewAlpha(image.Rect(0, 0, badgeSize, badgeSize))
	p := raster.NewAlphaSrcPainter(mask)
	r.Rasterize(p)

	draw.DrawMask(im, im.Bounds(), image.White, image.ZP, mask, image.ZP, draw.Over)

	f, fontErr := truetype.Parse(badge.FontFile)
	if fontErr != nil {
		err = fontErr
		return
	}
	fupe := fixed.Int26_6(f.FUnitsPerEm())

	//printGlyph(g)

	for stepIndex, step := range stepsData {

		if len(text) > stepIndex {
			letter := rune(text[stepIndex])

			i0 := f.Index(letter)
			//hm := f.HMetric(fupe, i0)
			g := &truetype.GlyphBuf{}
			loadErr := g.Load(f, fupe, i0, font.HintingNone)
			if loadErr != nil {
				err = loadErr
				return
			}

			r := raster.NewRasterizer(badgeSize, badgeSize)

			var letterNodes []node

			e := 0
			for i, p := range g.Points {

				pointX := int(halfBadgeSize+(halfBadgeSize*step.Sin*lettersPositionRatio)) + int((p.X-fontXoffset)>>6)
				pointY := int(halfBadgeSize+(halfBadgeSize*step.Cos*lettersPositionRatio)) + int((p.Y-fontYoffset)>>6)

				if p.Flags&0x01 != 0 {
					letterNodes = append(letterNodes, node{pointX, pointY, 1})
				} else {
					letterNodes = append(letterNodes, node{pointX, pointY, 2})
				}
				if i+1 == int(g.Ends[e]) {

					letterNodes = append(letterNodes, node{letterNodes[0].x, letterNodes[0].y, -1})

					contour(r, letterNodes)
					letterNodes = nil
					e++
				}
			}

			mask := image.NewAlpha(image.Rect(0, 0, badgeSize, badgeSize))
			p := raster.NewAlphaSrcPainter(mask)
			r.Rasterize(p)

			draw.DrawMask(im, im.Bounds(), &image.Uniform{blue}, image.ZP, mask, image.ZP, draw.Over)
		}
	}

	for x := 0; x < badgeSize; x++ {

		for y := 0; y < badgeSize; y++ {
			r1, g1, b1, a1 := im.At(x, y).RGBA()

			bgPixel := gu.Pixel{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(a1)}

			r2, g2, b2, a2 := resizedBadge.At(x, y).RGBA()

			fgPixel := gu.Pixel{R: uint8(r2), G: uint8(g2), B: uint8(b2), A: uint8(a2)}

			pixel := gu.BlendPixel(fgPixel, bgPixel)

			im.SetRGBA(x, y, color.RGBA{R: pixel.R, G: pixel.G, B: pixel.B, A: pixel.A})
		}
	}

	out = im

	return
}

func printGlyph(g *truetype.GlyphBuf) {
	//printBounds(g.Bounds)
	fmt.Print("Points:\n---\n")
	e := 0
	for i, p := range g.Points {
		fmt.Printf("%4d, %4d", p.X, p.Y)
		if p.Flags&0x01 != 0 {
			fmt.Print("  on\n")
		} else {
			fmt.Print("  off\n")
		}
		if i+1 == int(g.Ends[e]) {
			fmt.Print("---\n")
			e++
		}
	}
}

func p(n node) fixed.Point26_6 {
	//x, y := rotate(0, n.x, n.y)

	return fixed.Point26_6{
		X: fixed.Int26_6(n.x << 6),
		Y: fixed.Int26_6(n.y << 6),
	}
}

func rotate(angle float64, x int, y int) (nx int, ny int) {
	sin, cos := math.Sincos(angle * math.Pi / 180)

	nx = int(cos*float64(x) - sin*float64(y))
	ny = int(sin*float64(x) + cos*float64(y))

	return
}

func contour(r *raster.Rasterizer, ns []node) {
	if len(ns) == 0 {
		return
	}
	i := 0
	r.Start(p(ns[i]))
	for {
		switch ns[i].degree {
		case -1:
			// -1 signifies end-of-contour.
			return
		case 1:
			i += 1
			r.Add1(p(ns[i]))
		case 2:
			i += 2
			r.Add2(p(ns[i-1]), p(ns[i]))
		default:
			panic("bad degree")
		}
	}
}
