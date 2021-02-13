package badge

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype/raster"
	"github.com/nfnt/resize"
	"golang.org/x/image/math/fixed"

	gu "github.com/dirtykastro/graphicutils"
)

const badgeCircleRatio = 0.75

type node struct {
	x, y, degree int
}

type Badge struct {
	Img image.Image
}

func (badge *Badge) Render(badgeSize int, text string, rotation float64) (out image.Image) {
	resizedBadge := resize.Thumbnail(uint(badgeSize), uint(badgeSize), badge.Img, resize.Lanczos3)

	im := image.NewRGBA(image.Rectangle{Max: image.Point{X: badgeSize, Y: badgeSize}})

	var circle []node

	halfBadgeSize := float64(badgeSize) / 2

	steps := 50

	for i := 0; i < steps; i++ {
		sin, cos := math.Sincos(float64(i) * 360 / float64(steps) * math.Pi / 180)

		circleX := halfBadgeSize + (halfBadgeSize * sin * badgeCircleRatio)
		circleY := halfBadgeSize + (halfBadgeSize * cos * badgeCircleRatio)

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
