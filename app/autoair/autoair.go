/*
Auto-Air opens the air-cleaner automatically when the pm2.5 >= 130.
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinBzr = 10
	pinSG  = 18
	pinLed = 26

	dioPin  = 9
	rclkPin = 10
	sclkPin = 11
)

const (
	trigOnPM25  = 120
	trigOffPm25 = 100
)

const (
	statePattern    = "((state))"
	ipPattern       = "((000.000.000.000))"
	pm25Pattern     = "((PM2.5))"
	pm10Pattern     = "((PM10))"
	datetimePattern = "((yyyy-mm-dd hh:mm:ss))"
	datetimeFormat  = "2006-01-02 15:04:05"
)

var (
	autoair     *autoAir
	pageContext []byte
)

var bool2int = map[bool]int{
	false: 0,
	true:  1,
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	air := dev.NewPMS7003()
	sg := dev.NewSG90(pinSG)
	led := dev.NewLed(pinLed)
	dsp := dev.NewLedDisplay(dioPin, rclkPin, sclkPin)

	wsnCfg := &base.WsnConfig{
		Token: base.WsnToken,
		API:   base.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(wsnCfg)

	autoair = newAutoAir(air, sg, led, dsp, cloud)
	// autoair.setMode(base.DevMode)
	base.WaitQuit(func() {
		autoair.stop()
		rpio.Close()
	})
	autoair.start()

	http.HandleFunc("/", airServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func airServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		homePageHandler(w, r)
	case "POST":
		operationHandler(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	if len(pageContext) == 0 {
		var err error
		pageContext, err = ioutil.ReadFile("air.html")
		if err != nil {
			log.Printf("failed to read air.html")
			fmt.Fprintf(w, "internal error: failed to read home page")
			return
		}
	}

	ip := base.GetIP()
	if ip == "" {
		log.Printf("failed to get ip")
		fmt.Fprintf(w, "internal error: failed to get ip")
		return
	}

	wbuf := bytes.NewBuffer([]byte{})
	rbuf := bytes.NewBuffer(pageContext)
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		s := string(line)
		switch {
		case strings.Index(s, ipPattern) >= 0:
			s = strings.Replace(s, ipPattern, ip, 1)
		case strings.Index(s, pm25Pattern) >= 0:
			pm25 := fmt.Sprintf("%v", autoair.pm25)
			s = strings.Replace(s, pm25Pattern, pm25, 1)
		case strings.Index(s, pm10Pattern) >= 0:
			pm10 := fmt.Sprintf("%v", autoair.pm10)
			s = strings.Replace(s, pm10Pattern, pm10, 1)
		case strings.Index(s, datetimePattern) >= 0:
			datetime := time.Now().Format(datetimeFormat)
			s = strings.Replace(s, datetimePattern, datetime, 1)
		case strings.Index(s, statePattern) >= 0:
			state := "unchecked"
			if autoair.state {
				state = "checked"
			}
			s = strings.Replace(s, statePattern, state, 1)
		}
		wbuf.Write([]byte(s))
	}
	w.Write(wbuf.Bytes())
}

func operationHandler(w http.ResponseWriter, r *http.Request) {
	op := r.FormValue("op")
	switch op {
	case "on":
		log.Printf("web op: on")
		autoair.on()
	case "off":
		log.Printf("web op: off")
		autoair.off()
	default:
		log.Printf("web op: invalid operator")
	}
}

type autoAir struct {
	air       *dev.PMS7003
	sg        *dev.SG90
	led       *dev.Led
	dsp       *dev.LedDisplay
	cloud     iot.Cloud
	mode      base.Mode
	state     bool // true: turn on, false: turn off
	pm25      uint16
	pm10      uint16
	chClean   chan uint16 // for turning on/off the air-cleaner
	chAlert   chan uint16 // for alerting
	chDisplay chan uint16
	chCloud   chan uint16 // for pushing to iot cloud
}

func newAutoAir(air *dev.PMS7003, sg *dev.SG90, led *dev.Led, dsp *dev.LedDisplay, cloud iot.Cloud) *autoAir {
	return &autoAir{
		air:       air,
		sg:        sg,
		led:       led,
		dsp:       dsp,
		cloud:     cloud,
		mode:      base.PrdMode,
		state:     false,
		chClean:   make(chan uint16, 4),
		chAlert:   make(chan uint16, 4),
		chDisplay: make(chan uint16, 4),
		chCloud:   make(chan uint16, 4),
	}
}

func (a *autoAir) start() {
	log.Printf("service starting")
	log.Printf("mode: %v", a.mode)
	go a.sg.Roll(45)
	go a.detect()
	go a.clean()
	go a.alert()
	go a.push()
	go a.display()
}

func (a *autoAir) setMode(mode base.Mode) {
	a.mode = mode
}

func (a *autoAir) detect() {
	log.Printf("detecting pm2.5")
	for {
		var err error
		if a.mode == base.PrdMode {
			a.pm25, a.pm10, err = a.air.Get()
		} else {
			a.pm25, a.pm10, err = a.air.Mock()
		}
		if err != nil {
			log.Printf("failed to get pm2.5 and pm10, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("pm2.5: %v ug/m3", a.pm25)
		log.Printf("pm10: %v ug/m3", a.pm10)

		a.chClean <- a.pm25
		a.chAlert <- a.pm25
		a.chCloud <- a.pm25
		a.chDisplay <- a.pm25

		sec := 60 * time.Second
		if a.mode != base.PrdMode {
			sec = 15 * time.Second
		}
		time.Sleep(sec)
	}
}

func (a *autoAir) clean() {
	go func() {
		for {
			time.Sleep(60 * time.Second)
			if a.mode != base.PrdMode {
				continue
			}
			v := &iot.Value{
				Device: "5e00eb8fe4b04a9a92a6b3fc",
				Value:  bool2int[a.state],
			}
			if err := a.cloud.Push(v); err != nil {
				log.Printf("push: failed to push the state of air-cleaner to cloud, error: %v", err)
			}
		}
	}()

	for pm25 := range a.chClean {
		hour := time.Now().Hour()
		if pm25 < 400 && (hour >= 20 || hour < 8) {
			// disable at 20:00-08:00
			log.Printf("auto air-cleaner was disabled at 20:00-08:00")
			if a.state {
				a.off()
			}
			continue
		}

		if !a.state && pm25 >= trigOnPM25 {
			a.on()
			log.Printf("air-cleaner was turned on")
			continue
		}
		if a.state && pm25 < trigOffPm25 {
			a.off()
			log.Printf("air-cleaner was turned off")
			continue
		}
	}
}

func (a *autoAir) push() {
	for pm25 := range a.chCloud {
		if a.mode != base.PrdMode {
			continue
		}
		go func(pm25 uint16) {
			v := &iot.Value{
				Device: "5df507c4e4b04a9a92a64928",
				Value:  pm25,
			}
			if err := a.cloud.Push(v); err != nil {
				log.Printf("push: failed to push pm2.5 to cloud, error: %v", err)
			}
		}(pm25)
	}
}

func (a *autoAir) alert() {
	var pm25 uint16
	for {
		select {
		case v := <-a.chAlert:
			pm25 = v
		default:
			// do nothing
		}

		if pm25 >= trigOnPM25 {
			interval := 1000 - int(pm25)
			if interval < 0 {
				interval = 200
			}
			a.led.Blink(1, interval)
			continue
		}
		time.Sleep(1 * time.Second)
	}
}

func (a *autoAir) display() {
	var pm25 uint16
	a.dsp.Open()
	opened := true
	for {
		select {
		case v := <-a.chDisplay:
			pm25 = v
		default:
			// do nothing, just use the latest temp
		}

		if a.dsp == nil {
			time.Sleep(30 * time.Second)
			continue
		}

		hour := time.Now().Hour()
		if hour >= 20 || hour < 8 {
			// turn off oled at 20:00-08:00
			if opened {
				a.dsp.Close()
				opened = false
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if !opened {
			a.dsp.Open()
			opened = true
		}
		text := "----"
		if pm25 > 0 {
			text = fmt.Sprintf("%d", pm25)
		}
		a.dsp.Display(text)
		time.Sleep(3 * time.Second)
	}
}

func (a *autoAir) on() {
	a.sg.Roll(0)
	time.Sleep(1 * time.Second)
	a.sg.Roll(-45)
	a.state = true
}

func (a *autoAir) off() {
	a.sg.Roll(0)
	time.Sleep(1 * time.Second)
	a.sg.Roll(45)
	a.state = false
}

func (a *autoAir) stop() {
	a.sg.Roll(45)
	a.led.Off()
	a.dsp.Close()
}
