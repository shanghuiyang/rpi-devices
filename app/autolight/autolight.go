/*
Auto-Light let you control a led light by hands or any other objects.
It works with HCSR04, an ultrasonic distance meter, together.
The led light will light up when HCSR04 sensor get distance less then 40cm.
And the led will turn off after 45 seconds.
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
	pinLight = 16
	pinLed   = 4
	pinTrig  = 21
	pinEcho  = 26
)

const (
	statePattern    = "((state))"
	ipPattern       = "((000.000.000.000))"
	datetimePattern = "((yyyy-mm-dd hh:mm:ss))"
	datetimeFormat  = "2006-01-02 15:04:05"
)

var (
	alight      *autoLight
	pageContext []byte
)

var bool2int = map[bool]int{
	false: 0,
	true:  1,
}

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[autolight]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	led := dev.NewLed(pinLed)
	light := dev.NewLed(pinLight)
	if light == nil {
		log.Printf("[autolight]failed to new a led light")
		return
	}
	dist := dev.NewHCSR04(pinTrig, pinEcho)
	if dist == nil {
		log.Printf("[autolight]failed to new a HCSR04")
		return
	}

	wsnCfg := &iot.WsnConfig{
		Token: iot.WsnToken,
		API:   iot.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(wsnCfg)

	alight = newAutoLight(dist, light, led, cloud)
	base.WaitQuit(func() {
		alight.off()
		rpio.Close()
	})
	alight.start()

	http.HandleFunc("/", lightServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("[autolight]ListenAndServe: ", err.Error())
	}
}

func lightServer(w http.ResponseWriter, r *http.Request) {
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
		pageContext, err = ioutil.ReadFile("light.html")
		if err != nil {
			log.Printf("[autolight]failed to read light.html")
			fmt.Fprintf(w, "internal error: failed to read home page")
			return
		}
	}

	ip := base.GetIP()
	if ip == "" {
		log.Printf("[autolight]failed to get ip")
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
		case strings.Index(s, datetimePattern) >= 0:
			datetime := time.Now().Format(datetimeFormat)
			s = strings.Replace(s, datetimePattern, datetime, 1)
		case strings.Index(s, statePattern) >= 0:
			state := "unchecked"
			if alight.state {
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
		log.Printf("[autolight]web op: on")
		alight.on()
	case "off":
		log.Printf("[autolight]web op: off")
		alight.off()
	default:
		log.Printf("[autolight]web op: invalid operator")
	}
}

type autoLight struct {
	dist     *dev.HCSR04
	light    *dev.Led
	led      *dev.Led
	cloud    iot.Cloud
	trigTime time.Time
	state    bool // true: turn on, false: turn off
	chLight  chan bool
	chLed    chan bool
}

func newAutoLight(dist *dev.HCSR04, light *dev.Led, led *dev.Led, cloud iot.Cloud) *autoLight {
	return &autoLight{
		dist:     dist,
		light:    light,
		led:      led,
		state:    false,
		trigTime: time.Now(),
		cloud:    cloud,
		chLight:  make(chan bool, 4),
		chLed:    make(chan bool, 4),
	}
}

func (a *autoLight) start() {
	log.Printf("[autolight]start to service")
	go a.detect()
	go a.ctrLight()
	go a.ctrLed()

}

func (a *autoLight) detect() {
	// need to warm-up the ultrasonic distance meter first
	a.dist.Dist()
	time.Sleep(500 * time.Millisecond)
	for {
		d := a.dist.Dist()
		detected := (d < 20)
		a.chLight <- detected
		a.chLed <- detected

		t := 300 * time.Millisecond
		if detected {
			log.Printf("[autolight]detected objects, distance = %.2fcm", d)
			// make a dalay detecting
			t = 2 * time.Second
		}
		time.Sleep(t)
	}
}

func (a *autoLight) ctrLight() {
	go func() {
		for {
			time.Sleep(10 * time.Second)
			v := &iot.Value{
				Device: "5dd29e1be4b074c40dfe87c4",
				Value:  bool2int[a.state],
			}
			if err := a.cloud.Push(v); err != nil {
				log.Printf("[autolight]push: failed to push the state of light to cloud, error: %v", err)
			}
		}
	}()

	for detected := range a.chLight {
		if detected {
			if !a.state {
				a.on()
			}
			a.trigTime = time.Now()
			continue
		}
		timeout := time.Now().Sub(a.trigTime).Seconds() > 45
		if timeout && a.state {
			log.Printf("[autolight]timeout, light off")
			a.off()
		}
	}
}

func (a *autoLight) ctrLed() {
	for detected := range a.chLed {
		if detected {
			a.led.Blink(1, 200)
		}
	}
}

func (a *autoLight) on() {
	a.state = true
	a.trigTime = time.Now()
	a.light.On()
}

func (a *autoLight) off() {
	a.state = false
	a.light.Off()
}
