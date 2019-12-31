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

	a := newAutoLight(dist, light, led, cloud)
	base.WaitQuit(func() {
		a.off()
		rpio.Close()
	})
	a.start()
}

type autoLight struct {
	dist    *dev.HCSR04
	light   *dev.Led
	led     *dev.Led
	cloud   iot.Cloud
	chLight chan bool
	chLed   chan bool
}

func newAutoLight(dist *dev.HCSR04, light *dev.Led, led *dev.Led, cloud iot.Cloud) *autoLight {
	return &autoLight{
		dist:    dist,
		light:   light,
		led:     led,
		cloud:   cloud,
		chLight: make(chan bool, 4),
		chLed:   make(chan bool, 4),
	}
}

func (a *autoLight) start() {
	log.Printf("service starting")

	go a.ctrLight()
	go a.ctrLed()

	// need to warm-up the distance sensor first
	a.dist.Dist()
	time.Sleep(500 * time.Millisecond)
	for {
		d := a.dist.Dist()
		detected := (d < 40)
		a.chLight <- detected
		a.chLed <- detected

		t := 300 * time.Millisecond
		if detected {
			log.Printf("detected objects")
			// make a dalay detecting
			t = 1 * time.Second
		}
		time.Sleep(t)
	}
}

func (a *autoLight) ctrLight() {
	on := false
	go func() {
		for {
			time.Sleep(5 * time.Second)
			state := 0
			if on {
				state = 1
			}
			v := &iot.Value{
				Device: "5dd29e1be4b074c40dfe87c4",
				Value:  state,
			}
			if err := a.cloud.Push(v); err != nil {
				log.Printf("push: failed to push the state of light to cloud, error: %v", err)
			}
		}
	}()

	lastTrig := time.Now()
	for detected := range a.chLight {
		if detected {
			if !on {
				a.on()
				on = true
			}
			lastTrig = time.Now()
			continue
		}
		timeout := time.Now().Sub(lastTrig).Seconds() > 45
		if timeout && on {
			log.Printf("timeout, light off")
			a.off()
			on = false
		}
	}
}

func (a *autoLight) ctrLed() {
	for detected := range a.chLed {
		if detected {
			a.led.Blink(1, 300)
		}
	}
}

func (a *autoLight) on() {
	a.light.On()
}

func (a *autoLight) off() {
	a.light.Off()
}
