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
	"github.com/stianeikeland/go-rpio"
)

const (
	pinTrig = 2
	pinEcho = 3
	pinBzr  = 17
	pinLed  = 23
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	bzr := dev.NewBuzzer(pinBzr)
	led := dev.NewLed(pinLed)
	dist := dev.NewHCSR04(pinTrig, pinEcho)
	if dist == nil {
		log.Printf("failed to new a HCSR04")
		return
	}

	dog := newDoordog(dist, bzr, led)
	base.WaitQuit(func() {
		dog.stop()
		rpio.Close()
	})
	dog.start()
}

type doordog struct {
	dist    *dev.HCSR04
	buzzer  *dev.Buzzer
	led     *dev.Led
	chAlert chan bool
}

func newDoordog(dist *dev.HCSR04, buzzer *dev.Buzzer, led *dev.Led) *doordog {
	return &doordog{
		dist:    dist,
		buzzer:  buzzer,
		led:     led,
		chAlert: make(chan bool, 4),
	}
}

func (d *doordog) start() {
	log.Printf("doordog start to service")
	go d.alert()
	d.detect()

}

func (d *doordog) detect() {
	// need to warm-up the ultrasonic sensor first
	d.dist.Dist()
	time.Sleep(500 * time.Millisecond)
	for {
		dist := d.dist.Dist()
		detected := (dist < 60)
		d.chAlert <- detected

		t := 300 * time.Millisecond
		if detected {
			log.Printf("detected objects, distance = %.2fcm", dist)
			// make a dalay detecting
			t = 2 * time.Second
		}
		time.Sleep(t)
	}
}

func (d *doordog) alert() {
	alert := false
	trigTime := time.Now()
	go func() {
		for {
			if alert {
				go d.buzzer.Beep(1, 200)
				go d.led.Blink(1, 200)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for detected := range d.chAlert {
		if detected {
			alert = true
			trigTime = time.Now()
			continue
		}
		timeout := time.Now().Sub(trigTime).Seconds() > 30
		if timeout && alert {
			log.Printf("timeout, stop alert")
			alert = false
		}
	}
}

func (d *doordog) stop() {
	d.buzzer.Off()
	d.led.Off()
}
