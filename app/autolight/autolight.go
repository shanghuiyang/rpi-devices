/*
Auto-Light let you control a led light by hands or any other objects.
It works with HCSR04, an ultrasonic sensor, together.
The led light will light up when HCSR04 sensor get distance less then 40cm.
And the led will turn off after 45 seconds.
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
	pinLight = 16
	pinLed   = 4
	pinTrig  = 21
	pinEcho  = 26
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	led := dev.NewLed(pinLed)
	light := dev.NewLed(pinLight)
	if light == nil {
		log.Printf("failed to new a led light")
		return
	}
	dist := dev.NewHCSR04(pinTrig, pinEcho)
	if dist == nil {
		log.Printf("failed to new a HCSR04")
		return
	}

	wsnCfg := &base.WsnConfig{
		Token: base.WsnToken,
		API:   base.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(wsnCfg)

	a := &autoLight{
		dist:  dist,
		light: light,
		led:   led,
		cloud: cloud,
		ch:    make(chan bool, 32),
	}
	base.WaitQuit(func() {
		a.off()
		rpio.Close()
	})
	a.start()
}

type autoLight struct {
	infra *dev.InfraredDetector
	dist  *dev.HCSR04
	light *dev.Led
	led   *dev.Led
	cloud iot.Cloud
	ch    chan bool
}

func (a *autoLight) start() {
	go a.Detect()

	a.off()
	isLightOn := false
	lastTrig := time.Now()
	for b := range a.ch {
		if b {
			log.Printf("detected objects")
			if !isLightOn {
				a.on()
				isLightOn = true
			}
			lastTrig = time.Now()
			go func() {
				// draw a chart looks like:
				//
				// ____|___|____
				//
				v := &iot.Value{
					Device: "5dd29e1be4b074c40dfe87c4",
					Value:  0,
				}
				a.cloud.Push(v)
				time.Sleep(5 * time.Second)
				v.Value = 1
				a.cloud.Push(v)
				time.Sleep(5 * time.Second)
				v.Value = 0
				a.cloud.Push(v)
			}()
			continue
		}
		if time.Now().Sub(lastTrig).Seconds() > 45 && isLightOn {
			log.Printf("timeout, light off")
			a.off()
			isLightOn = false
		}
	}
}

func (a *autoLight) Detect() {
	// need to warm-up the distance sensor first
	a.dist.Dist()
	time.Sleep(500 * time.Millisecond)

	for {
		d := a.dist.Dist()
		detected := (d < 40)
		a.ch <- detected

		t := 300 * time.Millisecond
		if detected {
			go a.led.Blink(1, 300)
			// make a dalay detecting
			t = 1 * time.Second
		}
		time.Sleep(t)
	}
}

func (a *autoLight) on() {
	a.light.On()
}

func (a *autoLight) off() {
	a.light.Off()
}
