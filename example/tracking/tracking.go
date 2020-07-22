package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
	"gocv.io/x/gocv"
)

const (
	pinIn1 = 17
	pinIn2 = 23
	pinIn3 = 27
	pinIn4 = 22
	pinENA = 13
	pinENB = 19
)

var eng *dev.L298N

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[tracking]failed to open rpio, error: %v", err)
		os.Exit(1)
	}
	defer rpio.Close()

	eng = dev.NewL298N(pinIn1, pinIn2, pinIn3, pinIn4, pinENA, pinENB)
	if eng == nil {
		log.Fatal("[tracking]failed to new a L298N as engine, a car can't without any engine")
		os.Exit(1)
	}

	base.WaitQuit(func() {
		rpio.Close()
	})

	tracking()
	os.Exit(0)
}

func tracking() {
	rcolor := color.RGBA{G: 255, A: 255}
	// lcolor := color.RGBA{R: 255, A: 255}

	// the red ball
	lhsv := gocv.Scalar{Val1: 33, Val2: 108, Val3: 138}
	hhsv := gocv.Scalar{Val1: 61, Val2: 255, Val3: 255}

	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}

	img := gocv.NewMat()
	mask := gocv.NewMat()
	frame := gocv.NewMat()
	hsv := gocv.NewMat()
	kernel := gocv.NewMat()
	defer img.Close()
	defer mask.Close()
	defer frame.Close()
	defer hsv.Close()
	defer kernel.Close()

	video, _ := gocv.OpenVideoCapture(0)
	defer video.Close()

	n := 0
	i := 0
	for true {

		n++
		i++

		video.Grab(6)
		if !video.Read(&img) {
			log.Printf("[tracking]failed to read image")
			continue
		}
		// imgf := fmt.Sprintf("img%v.jpg", i+100000)
		// log.Printf("save %v", imgf)
		// gocv.IMWrite(imgf, img)

		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)

		// imgf = fmt.Sprintf("img%v.jpg", i+200000)
		// log.Printf("save %v", imgf)
		// gocv.IMWrite(imgf, img)

		cnt := bestContour(mask, 200)
		if len(cnt) == 0 {
			log.Printf("[tracking]len(cnt)==0")
			continue
		}

		rect := gocv.BoundingRect(cnt)
		fmt.Printf("rect w=%v, h=%v\n", rect.Dx(), rect.Dy())
		fmt.Printf("rect max y=%v\n", rect.Max.Y)

		if rect.Max.Y > 560 {
			stop()
			continue
		}

		gocv.Rectangle(&img, rect, rcolor, 2)
		imgf := fmt.Sprintf("img%v.jpg", i+300000)
		log.Printf("save %v", imgf)
		gocv.IMWrite(imgf, img)
		// ---

		x, y := middle(rect)
		log.Printf("[tracking]]ball at: (%v, %v)\n", x, y)
		if x < 200 {
			right()
			log.Printf("car right, sleep 3s")
			continue
		}
		if x > 400 {
			left()
			log.Printf("car left, sleep 3s")
			continue
		}
		forward()
		continue
	}
}

func bestContour(frame gocv.Mat, minArea float64) []image.Point {
	cnts := gocv.FindContours(frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var (
		bestCnt  []image.Point
		bestArea = minArea
	)
	for _, cnt := range cnts {
		if area := gocv.ContourArea(cnt); area > bestArea {
			bestArea = area
			bestCnt = cnt
		}
	}
	return bestCnt
}

// middle calculates the middle x and y of a rectangle.
func middle(rect image.Rectangle) (x int, y int) {
	return (rect.Max.X-rect.Min.X)/2 + rect.Min.X, (rect.Max.Y-rect.Min.Y)/2 + rect.Min.Y
}

func left() {
	eng.Left()
	time.Sleep(150 * time.Millisecond)
	eng.Stop()
}

func right() {
	eng.Right()
	time.Sleep(150 * time.Millisecond)
	eng.Stop()
}

func forward() {
	eng.Forward()
	time.Sleep(200 * time.Millisecond)
	eng.Stop()
	// time.Sleep(1000 * time.Millisecond)
}

func stop() {
	eng.Stop()
}
