// build with tracking using open cv:
// $ go build -tags=gocv app/car/car.go

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

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinLed       = 3
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
	pinCS        = 2
	pinTrig      = 21
	pinEcho      = 26

	ipPattern          = "((000.000.000.000))"
	selfDrivingState   = "((selfdriving-state))"
	selfTrackingState  = "((selftracking-state))"
	speechDrivingState = "((speechdriving-state))"

	selfDrivingEnabled   = "((selfdriving-enabled))"
	selfTrackingEnabled  = "((selftracking-enabled))"
	speechDrivingEnabled = "((speechdriving-enabled))"
)

type server struct {
	car         *car.Car
	pageContext []byte
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[carapp]failed to open rpio, error: %v", err)
		os.Exit(1)
	}
	defer rpio.Close()

	eng := dev.NewL298N(pinIn1, pinIn2, pinIn3, pinIn4, pinENA, pinENB)
	if eng == nil {
		log.Fatal("[carapp]failed to new a L298N as engine, a car can't without any engine")
		os.Exit(1)
	}

	ult := dev.NewUS100(&dev.US100Config{
		Mode: dev.UartMode,
		Dev:  "/dev/ttyAMA0",
		Baud: 9600,
	})
	if ult == nil {
		log.Printf("[carapp]failed to new a HCSR04, will build a car without ultrasonic distance meter")
	}

	// ult := dev.NewHCSR04(pinTrig, pinEcho)
	// if ult == nil {
	// 	log.Printf("[carapp]failed to new an ultrasonic distance meter, will build a car without ultrasonic distance meter")
	// }

	gy25 := dev.NewGY25("/dev/ttyUSB0", 115200)
	if gy25 == nil {
		log.Printf("[carapp]failed to new a gy-25, will build a car without gy-25")
	}

	collisionL := dev.NewCollision(pinCSwaitchL)
	if collisionL == nil {
		log.Printf("[carapp]failed to new a collision switch, will build a car without collision switchs")
	}

	collisionR := dev.NewCollision(pinCSwaitchR)
	if collisionR == nil {
		log.Printf("[carapp]failed to new a collision switch, will build a car without collision switchs")
	}
	collisions := []*dev.Collision{collisionL, collisionR}

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

	var gps *dev.GPS = nil
	// gps := dev.NewGPS()
	// if gps == nil {
	// 	log.Printf("[carapp]failed to new a gps sensor")
	// 	return
	// }

	var lc12s *dev.LC12S = nil
	// lc12s, err := dev.NewLC12S(pinCS)
	// if err != nil {
	// 	log.Printf("[carapp]failed to new a LC12S, error: %v", err)
	// }

	car := car.New(&car.Config{
		Engine:     eng,
		Servo:      servo,
		DistMeter:  ult,
		GY25:       gy25,
		Collisions: collisions,
		Horn:       horn,
		Led:        led,
		Camera:     cam,
		GPS:        gps,
		LC12S:      lc12s,
	})
	if car == nil {
		log.Fatal("failed to new a car")
		return
	}

	svr := newServer(car)
	util.WaitQuit(func() {
		svr.stop()
		if ult != nil {
			ult.Close()
		}
		if horn != nil {
			horn.Off()
		}
		if gy25 != nil {
			gy25.Close()
		}
		if led != nil {
			led.Off()
		}
		if light != nil {
			light.Off()
		}
		if lc12s != nil {
			lc12s.Close()
		}
		rpio.Close()
	})
	if err := svr.start(); err != nil {
		log.Printf("[carapp]failed to start car server, error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func newServer(car *car.Car) *server {
	return &server{
		car: car,
	}
}

func (s *server) start() error {
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

func (s *server) stop() error {
	return s.car.Stop()
}

func (s *server) loadHomePage(w http.ResponseWriter, r *http.Request) error {
	if len(s.pageContext) == 0 {
		var err error
		s.pageContext, err = ioutil.ReadFile("car.html")
		if err != nil {
			return errors.New("internal error: failed to read car.html")
		}
	}

	ip := util.GetIP()
	if ip == "" {
		return errors.New("internal error: failed to get ip")
	}

	rbuf := bytes.NewBuffer(s.pageContext)
	wbuf := bytes.NewBuffer([]byte{})
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		sline := string(line)

		disabled := false
		selfDriving, selfTracking, speechDriving := s.car.GetState()
		if selfDriving || selfTracking || speechDriving {
			disabled = true
		}

		if strings.Index(sline, ipPattern) >= 0 {
			sline = strings.Replace(sline, ipPattern, ip, 1)
		}

		if strings.Index(sline, selfDrivingState) >= 0 {
			state := "unchecked"
			if selfDriving {
				state = "checked"
			}
			sline = strings.Replace(sline, selfDrivingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, selfDrivingEnabled, able, 1)
		}

		if strings.Index(sline, selfTrackingState) >= 0 {
			state := "unchecked"
			if selfTracking {
				state = "checked"
			}
			sline = strings.Replace(sline, selfTrackingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, selfTrackingEnabled, able, 1)
		}

		if strings.Index(sline, speechDrivingState) >= 0 {
			state := "unchecked"
			if speechDriving {
				state = "checked"
			}
			sline = strings.Replace(sline, speechDrivingState, state, 1)

			able := "enabled"
			if state == "unchecked" && disabled {
				able = "disabled"
			}
			sline = strings.Replace(sline, speechDrivingEnabled, able, 1)
		}

		wbuf.Write([]byte(sline))
	}
	w.Write(wbuf.Bytes())
	return nil
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.loadHomePage(w, r)
	case "POST":
		if dest := r.FormValue("dest"); dest != "" {
			var lat, lon float64
			n, err := fmt.Sscanf(dest, "%f,%f", &lat, &lon)
			if err != nil || n != 2 {
				log.Printf("invalid destination input: %v", dest)
				return
			}
			destPt := &geo.Point{
				Lat: lat,
				Lon: lon,
			}
			log.Printf("dest: %v", destPt)
			s.car.SetDest(destPt)
		}
		op := r.FormValue("op")
		s.car.Do(car.Op(op))
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

		encoder.Start()
		for i := 0; i < count; {
			i += encoder.Count1()
		}
		eng.Stop()
		encoder.Stop()
	}
	eng.Stop()
	return
}
