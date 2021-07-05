package selftracking

import (
	"image/color"
	"log"

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/cv"

	"github.com/shanghuiyang/rpi-devices/util"
	"gocv.io/x/gocv"
)

const (
	logTag = "selftracking"
)

var (
	ontracking bool
	mycar      car.Car
	tracker    *cv.Tracker
	streamer   *cv.Streamer
)

func Init(c car.Car, t *cv.Tracker, s *cv.Streamer) {
	mycar = c
	tracker = t
	streamer = s
	go streamer.Start()
}

func Start(chImg chan *gocv.Mat) {
	if ontracking {
		return
	}

	ontracking = true

	rcolor := color.RGBA{G: 255, A: 255}
	firstTime := true // saw the ball at the first time
	for ontracking {
		util.DelayMs(200)

		img, ok := <-chImg
		if !ok {
			ontracking = false
			return
		}

		ok, rect := tracker.Locate(img)
		if ok {
			gocv.Rectangle(img, *rect, rcolor, 2)
		}
		streamer.Push(img)

		if !ok {
			// looking for the ball by turning 360 degree
			log.Printf("[%v]ball not found", logTag)
			firstTime = true
			continue
		}

		// found the ball, move to it
		if rect.Max.Y > 580 {
			mycar.Stop()
			mycar.Beep(1, 300)
			continue
		}
		if firstTime {
			go mycar.Beep(1, 100)
		}
		firstTime = false
		x, y := tracker.MiddleXY(rect)
		log.Printf("[%v]found a ball at: (%v, %v)", logTag, x, y)
		if x < 200 {
			log.Printf("[%v]turn left to the ball", logTag)
			mycar.Left()
			util.DelayMs(100)
			mycar.Stop()
			continue
		}
		if x > 400 {
			log.Printf("[%v]turn right to the ball", logTag)
			mycar.Right()
			util.DelayMs(100)
			mycar.Stop()
			continue
		}
		log.Printf("[%v]forward to the ball", logTag)
		mycar.Forward()
		util.DelayMs(100)
		mycar.Stop()

	}
}

func Status() bool {
	return ontracking
}

func Stop() {
	ontracking = false
}
