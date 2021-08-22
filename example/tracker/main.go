// +build gocv

package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/shanghuiyang/rpi-devices/cv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"gocv.io/x/gocv"
)

const (
	pinIn1 = 17
	pinIn2 = 23
	pinIn3 = 27
	pinIn4 = 22
	pinENA = 13
	pinENB = 19

	// the hsv of a tennis
	lh float64 = 33
	ls float64 = 108
	lv float64 = 138
	hh float64 = 61
	hs float64 = 255
	hv float64 = 255

	host = "0.0.0.0:8088"
)

var eng dev.MotorDriver

func main() {
	eng = dev.NewL298N(pinIn1, pinIn2, pinIn3, pinIn4, pinENA, pinENB)
	if eng == nil {
		log.Fatal("[tracking]failed to new a L298N as engine, a car can't without any engine")
		os.Exit(1)
	}

	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Printf("failed to open video")
		return
	}
	defer cam.Close()

	t, err := cv.NewTracker(lh, ls, lv, hh, hs, hv)
	if err != nil {
		log.Printf("[carapp]failed to create a tracker, error: %v", err)
		return
	}
	defer t.Close()

	streamer := cv.NewStreamer(host)
	defer streamer.Close()
	go streamer.Start()

	img := gocv.NewMat()
	defer img.Close()

	rcolor := color.RGBA{G: 255, A: 255}
	for {
		cam.Grab(10)
		if !cam.Read(&img) {
			log.Printf("[car]failed to read img from camera")
			continue
		}

		ok, rect := t.Locate(&img)
		if ok {
			gocv.Rectangle(&img, *rect, rcolor, 2)
		}
		streamer.Push(&img)

		if !ok {
			continue
		}

		if rect.Max.Y > 400 {
			stop()
			continue
		}

		x, y := middle(rect)
		log.Printf("[tracking]ball at: (%v, %v)\n", x, y)
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
func middle(rect *image.Rectangle) (x int, y int) {
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
}

func stop() {
	eng.Stop()
}
