package main

import (
	"fmt"
	"image"
	//"image/color"
	"os"

	"gocv.io/x/gocv"
)

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

		// color for the rect when faces detected
		//blue := color.RGBA{0, 0, 255, 0}

		// load classifier to recognize faces
		classifier := gocv.NewCascadeClassifier()
		defer classifier.Close()

		if !classifier.Load(faceCascadeFile) {
			fmt.Printf("Error reading cascade file: %v\n", faceCascadeFile)
			return
		}

		/*launghingManImg := gocv.IMRead("images/laughing_man.png", gocv.IMReadUnchanged)

		if launghingManImg.Empty() {
			fmt.Println("problem reading launghing man png")
		}
		defer launghingManImg.Close()*/

		fmt.Printf("start reading video: %s\n", filePath)
		for {
			if ok := video.Read(&img); !ok {
				fmt.Printf("cannot read video file %v\n", filePath)
				return
			}

			if img.Empty() {
				continue
			}

			dimensions := img.Size()

			videoWidth := dimensions[0]
			videoHeight := dimensions[1]

			//fmt.Println("Video Size", videoWidth, "x", videoHeight, "channels :", img.Channels())

			// detect faces
			//rects := classifier.DetectMultiScale(img)
			rects := classifier.DetectMultiScaleWithParams(img, 1.1, 10, 0, image.Point{X: 400, Y: 400}, image.Point{X: 2000, Y: 2000})
			fmt.Printf("found %d faces\n", len(rects))

			// draw a rectangle around each face on the original image,
			// along with text identifying as "Human"
			for _, r := range rects {
				/*gocv.Rectangle(&img, r, blue, 3)

				size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
				gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)*/

				lmResized := gocv.NewMatWithSize(videoWidth, videoHeight, gocv.MatTypeCV8UC3)
				defer lmResized.Close()

				//gocv.Resize(launghingManImg, &lmResized, image.Pt(videoWidth, videoHeight), 0, 0, gocv.InterpolationNearestNeighbor)

				if r.Min.X > 0 {
					gocv.AddWeighted(img, 1.0, lmResized, 1.0, 0.0, &img)
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
