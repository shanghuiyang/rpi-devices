package server

import (
	"io/ioutil"
	"log"

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	// "github.com/shanghuiyang/rpi-devices/cv/mock/cv"

	// "github.com/shanghuiyang/rpi-devices/cv"
	// "github.com/shanghuiyang/rpi-devices/cv/mock/gocv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	forward  Op = "forward"
	backward Op = "backward"
	left     Op = "left"
	right    Op = "right"
	stop     Op = "stop"
	beep     Op = "beep"

	chSize = 8
)

type Op string

func init() {
	var err error
	pageContext, err = ioutil.ReadFile("home.html")
	if err != nil {
		log.Printf("failed to load home page, error: %v", err)
		panic(err)

	}

	ip = util.GetIP()
	if ip == "" {
		log.Printf("failed to get ip address")
		panic(err)
	}
}

type service struct {
	car        car.Car
	led        *dev.Led
	ledBlinked bool
	chOp       chan Op

	enabledSelfDriving   bool
	enabledSelfTracking  bool
	enabledSelfNav       bool
	enabledSpeechDriving bool
}

func newService(car car.Car, led *dev.Led) *service {
	return &service{
		car:        car,
		led:        led,
		ledBlinked: true,
		chOp:       make(chan Op, chSize),
	}
}

func (s *service) EnableSelfDriving(enabled bool) *service {
	s.enabledSelfDriving = enabled
	return s
}

func (s *service) EnabledSelfTracking(enabled bool) *service {
	s.enabledSelfTracking = enabled
	return s
}

func (s *service) EnabledSpeechDriving(enabled bool) *service {
	s.enabledSpeechDriving = enabled
	return s
}

func (s *service) start() error {
	go s.operate()
	go s.blink()
	return nil
}

// Stop ...
func (s *service) Stop() error {
	close(s.chOp)
	s.car.Stop()
	return nil
}

func (s *service) operate() {
	for op := range s.chOp {
		log.Printf("[car]op: %v", op)
		switch op {
		case forward:
			s.car.Forward()
		case backward:
			s.car.Backward()
		case left:
			s.car.Left()
		case right:
			s.car.Right()
		case stop:
			s.car.Stop()
		case beep:
			go s.car.Beep(3, 100)
		default:
			log.Printf("[car]invalid op")
		}
	}
}

func (s *service) blink() {
	for s.ledBlinked {
		s.led.Blink(1, 1000)
	}
}
