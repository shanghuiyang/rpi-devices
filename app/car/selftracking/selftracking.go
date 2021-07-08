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

type SelfTracking struct {
	car        car.Car
	tracker    *cv.Tracker
	streamer   *cv.Streamer
	ontracking bool
}

func New(c car.Car, t *cv.Tracker, s *cv.Streamer) *SelfTracking {
	go s.Start()
	return &SelfTracking{
		car:        c,
		tracker:    t,
		streamer:   s,
		ontracking: false,
	}
}

func (s *SelfTracking) Start(chImg chan *gocv.Mat) {
	if s.ontracking {
		return
	}

	s.ontracking = true

	rcolor := color.RGBA{G: 255, A: 255}
	firstTime := true // saw the ball at the first time
	for s.ontracking {
		util.DelayMs(200)

		img, ok := <-chImg
		if !ok {
			s.ontracking = false
			return
		}

		ok, rect := s.tracker.Locate(img)
		if ok {
			gocv.Rectangle(img, *rect, rcolor, 2)
		}
		s.streamer.Push(img)

		if !ok {
			// looking for the ball by turning 360 degree
			log.Printf("[%v]ball not found", logTag)
			firstTime = true
			continue
		}

		// found the ball, move to it
		if rect.Max.Y > 580 {
			s.car.Stop()
			s.car.Beep(1, 300)
			continue
		}
		if firstTime {
			go s.car.Beep(1, 100)
		}
		firstTime = false
		x, y := s.tracker.MiddleXY(rect)
		log.Printf("[%v]found a ball at: (%v, %v)", logTag, x, y)
		if x < 200 {
			log.Printf("[%v]turn left to the ball", logTag)
			s.car.Left()
			util.DelayMs(100)
			s.car.Stop()
			continue
		}
		if x > 400 {
			log.Printf("[%v]turn right to the ball", logTag)
			s.car.Right()
			util.DelayMs(100)
			s.car.Stop()
			continue
		}
		log.Printf("[%v]forward to the ball", logTag)
		s.car.Forward()
		util.DelayMs(100)
		s.car.Stop()

	}
}

func (s *SelfTracking) Status() bool {
	return s.ontracking
}

func (s *SelfTracking) Stop() {
	s.ontracking = false
}
