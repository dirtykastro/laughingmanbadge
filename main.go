package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"os"

	gu "github.com/dirtykastro/graphicutils"
	"github.com/dirtykastro/laughingmanbadge/badge"

	"gocv.io/x/gocv"
)

//go:embed cmd/lmbadge/laughing_man.png
//go:embed cmd/lmbadge/RobotoMono-Medium.ttf

var f embed.FS

// location of the frontface haarcascade file
const faceCascadeFile = "/usr/local/share/opencv4/haarcascades/haarcascade_frontalface_default.xml"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Cover Faces with Laughing Man Badge")
		fmt.Println("usage:", os.Args[0], " [video file]")
		return
	}

	// parse args
	filePath := os.Args[1]

	if exists(filePath) && !isDirectory(filePath) {

		// open video file
		video, err := gocv.VideoCaptureFile(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer video.Close()

		// open display window
		window := gocv.NewWindow("Laughing Man Badge")
		defer window.Close()

		// prepare image matrix
		img := gocv.NewMat()
		defer img.Close()

		// load classifier to recognize faces
		classifier := gocv.NewCascadeClassifier()
		defer classifier.Close()

		if !classifier.Load(faceCascadeFile) {
			fmt.Printf("Error reading cascade file: %v\n", faceCascadeFile)
			return
		}

		var badgeNoText image.Image

		file, err := f.ReadFile("cmd/lmbadge/laughing_man.png")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)

		}

		badgeNoText, err = png.Decode(bytes.NewReader(file))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)

		}

		fontFile, err := f.ReadFile("cmd/lmbadge/RobotoMono-Medium.ttf")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		lmBadge := &badge.Badge{Img: badgeNoText, FontFile: fontFile}

		fmt.Printf("start reading video: %s\n", filePath)
		for {
			if ok := video.Read(&img); !ok {
				fmt.Printf("cannot read video file %v\n", filePath)
				return
			}

			if img.Empty() {
				continue
			}

			/*dimensions := img.Size()

			videoWidth := dimensions[0]
			videoHeight := dimensions[1]*/

			//fmt.Println("Video Size", videoWidth, "x", videoHeight, "channels :", img.Channels())

			// detect faces
			//rects := classifier.DetectMultiScale(img)
			rects := classifier.DetectMultiScaleWithParams(img, 1.1, 10, 0, image.Point{X: 400, Y: 400}, image.Point{X: 2000, Y: 2000})
			fmt.Printf("found %d faces\n", len(rects))

			// draw a rectangle around each face on the original image,
			// along with text identifying as "Human"

			badgeSize := 0
			//badgePositionX := 0
			//badgePositionY := 0

			for _, r := range rects {
				rectWidth := r.Max.X - r.Min.X
				rectHeight := r.Max.Y - r.Min.Y

				rectSize := rectWidth

				if rectHeight > rectSize {
					rectSize = rectHeight
				}

				if rectSize > badgeSize {
					badgeSize = rectSize
					//badgePositionX = r.Min.X
					//badgePositionY = r.Min.Y
				}
			}

			if badgeSize > 0 {
				im, err := lmBadge.Render(badgeSize, "I thought what I'd do was, I'd pretend I was one of those deaf-mutes.", 0.0)
				if err == nil {

					totalChannels := img.Channels()

					for x := 0; x < badgeSize; x++ {
						for y := 0; y < badgeSize; y++ {
							// get video pixel color values
							b0 := img.GetUCharAt(x, y*totalChannels+0)
							g0 := img.GetUCharAt(x, y*totalChannels+1)
							r0 := img.GetUCharAt(x, y*totalChannels+2)

							bgPixel := gu.Pixel{R: uint8(r0), G: uint8(g0), B: uint8(b0), A: 255}

							// opencv order is inverted
							r1, g1, b1, a1 := im.At(y, x).RGBA()

							fgPixel := gu.Pixel{R: uint8(r1), G: uint8(g1), B: uint8(b1), A: uint8(a1)}

							pixel := gu.BlendPixel(fgPixel, bgPixel)

							img.SetUCharAt(x, y*totalChannels+0, pixel.B)
							img.SetUCharAt(x, y*totalChannels+1, pixel.G)
							img.SetUCharAt(x, y*totalChannels+2, pixel.R)

						}
					}
				} else {
					fmt.Println(err)
				}
			}

			// show the image in the window, and wait 1 millisecond
			window.IMShow(img)
			if window.WaitKey(1) >= 0 {
				break
			}
		}
	} else {
		fmt.Println(filePath, "is not valid")
	}

}

// exists checks if file exists
//
func exists(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// isDirectory checks if path is a directory
//
func isDirectory(file string) bool {
	if stat, err := os.Stat(file); err == nil && stat.IsDir() {
		return true
	}
	return false
}
