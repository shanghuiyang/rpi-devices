/*
Auto-Air opens the air-cleaner automatically when the pm2.5 >= 130.
*/

package main

import (
	"log"
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
)

const (
	trigOnPM25  = 130
	trigOffPm25 = 100
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	air := dev.NewPMS7003()
	sg := dev.NewSG90(pinSG)
	led := dev.NewLed(pinLed)

	wsnCfg := &base.WsnConfig{
		Token: "your token",
		API:   "http://www.wsncloud.com/api/data/v1/numerical/insert",
	}
	cloud := iot.NewCloud(wsnCfg)

	a := newAutoAir(air, sg, led, cloud)
	// a.setMode(base.DevMode)
	base.WaitQuit(func() {
		a.stop()
		rpio.Close()
	})
	a.start()
}

type autoAir struct {
	air     *dev.PMS7003
	sg      *dev.SG90
	led     *dev.Led
	cloud   iot.Cloud
	mode    base.Mode
	chClean chan uint16 // for turning on/off the air-cleaner
	chAlert chan uint16 // for alerting
	chCloud chan uint16 // for pushing to iot cloud
}

func newAutoAir(air *dev.PMS7003, sg *dev.SG90, led *dev.Led, cloud iot.Cloud) *autoAir {
	return &autoAir{
		air:     air,
		sg:      sg,
		led:     led,
		cloud:   cloud,
		mode:    base.PrdMode,
		chClean: make(chan uint16, 4),
		chAlert: make(chan uint16, 4),
		chCloud: make(chan uint16, 4),
	}
}

func (a *autoAir) start() {
	log.Printf("service starting")
	log.Printf("mode: %v", a.mode)
	go a.sg.Roll(45)
	go a.clean()
	go a.alert()
	go a.push()
	a.detect()
}

func (a *autoAir) setMode(mode base.Mode) {
	a.mode = mode
}

func (a *autoAir) detect() {
	log.Printf("detecting pm2.5")
	for {
		var pm25 uint16
		var err error
		if a.mode == base.PrdMode {
			pm25, _, err = a.air.Get()
		} else {
			pm25, _, err = a.air.Mock()
		}
		if err != nil {
			log.Printf("pm25: failed to get pm2.5, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("pm2.5: %v ug/m3", pm25)

		a.chClean <- pm25
		a.chAlert <- pm25
		a.chCloud <- pm25

		sec := 60 * time.Second
		if a.mode != base.PrdMode {
			sec = 15 * time.Second
		}
		time.Sleep(sec)
	}
}

func (a *autoAir) clean() {
	on := false
	for pm25 := range a.chClean {
		if !on && pm25 >= trigOnPM25 {
			on = true
			a.trigger()
			log.Printf("air-cleaner was turned on")
			continue
		}
		if on && pm25 < trigOffPm25 {
			on = false
			a.trigger()
			log.Printf("air-clearn was turned off")
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

func (a *autoAir) trigger() {
	a.sg.Roll(0)
	time.Sleep(1 * time.Second)
	a.sg.Roll(45)

}

func (a *autoAir) stop() {
	a.sg.Roll(45)
	a.led.Off()
}
