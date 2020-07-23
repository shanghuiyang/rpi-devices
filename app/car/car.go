package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/cv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinLed       = 4
	pinLight     = 16
	pinIn1       = 17
	pinIn2       = 23
	pinIn3       = 27
	pinIn4       = 22
	pinENA       = 13
	pinENB       = 19
	pinBzr       = 10
	pinSG        = 18
	pinEncoder   = 6
	pinCSwaitchL = 20 // the collision switch on left
	pinCSwaitchR = 12 // the collision switch on right

	// use this rpio as 3.3v pin
	// if all 3.3v pins were used
	pin33v = 5

	ipPattern = "((000.000.000.000))"
)

var (
	// the hsv of a tennis
	lh = float64(33)
	ls = float64(108)
	lv = float64(138)
	hh = float64(61)
	hs = float64(255)
	hv = float64(255)
)

type carServer struct {
	car         *dev.Car
	pageContext []byte
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[carapp]failed to open rpio, error: %v", err)
		os.Exit(1)
	}
	defer rpio.Close()

	p33v := rpio.Pin(pin33v)
	p33v.Output()
	p33v.High()

	eng := dev.NewL298N(pinIn1, pinIn2, pinIn3, pinIn4, pinENA, pinENB)
	if eng == nil {
		log.Fatal("[carapp]failed to new a L298N as engine, a car can't without any engine")
		os.Exit(1)
	}

	ult := dev.NewUS100()
	if ult == nil {
		log.Printf("[carapp]failed to new a HCSR04, will build a car without ultrasonic distance meter")
	}

	encoder := dev.NewEncoder(pinEncoder)
	if encoder == nil {
		log.Printf("[carapp]failed to new a encoder, will build a car without encoder")
	}

	cswitchL := dev.NewCollisionSwitch(pinCSwaitchL)
	if cswitchL == nil {
		log.Printf("[carapp]failed to new a collision switch, will build a car without collision switchs")
	}

	cswitchR := dev.NewCollisionSwitch(pinCSwaitchR)
	if cswitchL == nil {
		log.Printf("[carapp]failed to new a collision switch, will build a car without collision switchs")
	}
	cswitchs := []*dev.CollisionSwitch{cswitchL, cswitchR}

	horn := dev.NewBuzzer(pinBzr)
	if horn == nil {
		log.Printf("[carapp]failed to new a buzzer, will build a car without horns")
	}

	led := dev.NewLed(pinLed)
	if led == nil {
		log.Printf("[carapp]failed to new a led, will build a car without leds")
	}

	light := dev.NewLed(pinLight)
	if light == nil {
		log.Printf("[carapp]failed to new a light, will build a car without lights")
	}

	servo := dev.NewSG90(pinSG)
	if servo == nil {
		log.Printf("[carapp]failed to new a sg90, will build a car without servo")
	}
	cam := dev.NewCamera()
	if cam == nil {
		log.Printf("[carapp]failed to new a camera, will build a car without cameras")
	}

	t, err := cv.NewTracker(lh, ls, lv, hh, hs, hv)
	if err != nil {
		log.Printf("[carapp]failed to new a tracker, will build a car without trankers")
	}

	car := dev.NewCar(
		dev.WithEngine(eng),
		dev.WithServo(servo),
		dev.WithUlt(ult),
		dev.WithEncoder(encoder),
		dev.WithCSwitchs(cswitchs),
		dev.WithHorn(horn),
		dev.WithLed(led),
		dev.WithLight(light),
		dev.WithCamera(cam),
		dev.WithTracker(t),
	)
	if car == nil {
		log.Fatal("failed to new a car")
		return
	}

	server := newCarServer(car)
	base.WaitQuit(func() {
		server.stop()
		rpio.Close()
	})
	if err := server.start(); err != nil {
		log.Printf("[carapp]failed to start car server, error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func newCarServer(car *dev.Car) *carServer {
	return &carServer{
		car: car,
	}
}

func (s *carServer) start() error {
	if err := s.loadHomePage(); err != nil {
		return err
	}

	if err := s.car.Start(); err != nil {
		return err
	}
	log.Printf("[carapp]car started successfully")

	http.HandleFunc("/", s.handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}
	return nil
}

func (s *carServer) stop() error {
	return s.car.Stop()
}

func (s *carServer) loadHomePage() error {
	data, err := ioutil.ReadFile("car.html")
	if err != nil {
		return errors.New("internal error: failed to read car.html")
	}

	ip := base.GetIP()
	if ip == "" {
		return errors.New("internal error: failed to get ip")
	}

	rbuf := bytes.NewBuffer(data)
	wbuf := bytes.NewBuffer([]byte{})
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		s := string(line)
		if strings.Index(s, ipPattern) >= 0 {
			s = strings.Replace(s, ipPattern, ip, 1)
		}
		wbuf.Write([]byte(s))
	}
	s.pageContext = wbuf.Bytes()
	return nil
}

func (s *carServer) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write(s.pageContext)
	case "POST":
		op := r.FormValue("op")
		s.car.Do(dev.CarOp(op))
	}
}

// tuningTurnAngle tunings the mapping between angle(degree) and time(millisecond)
func tuningTurnAngle(eng *dev.L298N) {
	if eng == nil {
		log.Fatal("eng is nil")
		return
	}
	for {
		var ms int
		fmt.Printf(">>ms: ")
		if n, err := fmt.Scanf("%d", &ms); n != 1 || err != nil {
			log.Printf("[carapp]invalid operator, error: %v", err)
			continue
		}
		if ms < 0 {
			break
		}
		eng.Right()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		eng.Stop()
	}
	return
}

// tuningTurnAngle tunings the mapping between angle(degree) and time(millisecond)
func tuningEncoder(eng *dev.L298N, encoder *dev.Encoder) {
	if eng == nil {
		log.Fatal("engineer is nil")
		return
	}
	if encoder == nil {
		log.Fatal("encoder is nil")
		return
	}
	eng.Speed(30)
	for {
		var count int
		fmt.Printf(">>count: ")
		if n, err := fmt.Scanf("%d", &count); n != 1 || err != nil {
			log.Printf("[carapp]invalid count, error: %v", err)
			continue
		}
		if count == 0 {
			break
		}
		if count < 0 {
			eng.Left()
			count *= -1
		} else {
			eng.Right()
		}
		for i := 0; i < count; {
			i += encoder.Count1()
		}
		eng.Stop()
	}
	eng.Stop()
	encoder.Close()
	return
}
