/*
Doordog helps you watch your doors.
When somebody entries your room, you will be alerted by a beeping buzzer and a blinking led.
*/

package main

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinTrig = 2
	pinEcho = 3
	pinBtn  = 7
	pinBzr  = 17
	pinLed  = 23

	ifttAPI = "your-iftt-api"
)

const (
	// the time of keeping alert in second
	alertTime = 60
	// the distance of triggering alert in cm
	alertDist = 100
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[doordog]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	bzr := dev.NewBuzzer(pinBzr)
	led := dev.NewLed(pinLed)
	btn := dev.NewButton(pinBtn)
	dist := dev.NewHCSR04(pinTrig, pinEcho)
	if dist == nil {
		log.Printf("[doordog]failed to new a HCSR04")
		return
	}

	dog := newDoordog(dist, bzr, led, btn)
	base.WaitQuit(func() {
		dog.stop()
		rpio.Close()
	})
	dog.start()
}

type doordog struct {
	dist     *dev.HCSR04
	buzzer   *dev.Buzzer
	led      *dev.Led
	button   *dev.Button
	alerting bool
	chAlert  chan bool
}

func newDoordog(dist *dev.HCSR04, buzzer *dev.Buzzer, led *dev.Led, btn *dev.Button) *doordog {
	return &doordog{
		dist:     dist,
		buzzer:   buzzer,
		led:      led,
		button:   btn,
		alerting: false,
		chAlert:  make(chan bool, 4),
	}
}

func (d *doordog) start() {
	log.Printf("[doordog]start to service")
	go d.alert()
	go d.stopAlert()
	d.detect()

}

func (d *doordog) detect() {
	// need to warm-up the ultrasonic distance meter first
	d.dist.Dist()
	time.Sleep(500 * time.Millisecond)
	var t time.Duration
	for {
		time.Sleep(t)
		dist := d.dist.Dist()
		if dist < 10 {
			log.Printf("[doordog]bad data from distant meter, distance = %.2fcm", dist)
			continue
		}
		detected := (dist < alertDist)
		d.chAlert <- detected

		t = 300 * time.Millisecond
		if detected {
			log.Printf("[doordog]detected objects, distance = %.2fcm", dist)
			// make a dalay detecting
			t = 3 * time.Second
			continue
		}
	}
}

func (d *doordog) alert() {
	trigTime := time.Now()
	go func() {
		for {
			if d.alerting {
				go d.buzzer.Beep(1, 200)
				go d.led.Blink(1, 200)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for detected := range d.chAlert {
		if detected {
			go ifttt()
			d.alerting = true
			trigTime = time.Now()
			continue
		}
		timeout := time.Now().Sub(trigTime).Seconds() > alertTime
		if timeout && d.alerting {
			log.Printf("[doordog]timeout, stop alert")
			d.alerting = false
		}
	}
}

func (d *doordog) stopAlert() {
	for {
		pressed := d.button.Pressed()
		if pressed {
			log.Printf("[doordog]the button was pressed")
			if d.alerting {
				d.alerting = false
			}
			// make a dalay detecting
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func ifttt() {
	req, err := http.NewRequest("POST", ifttAPI, bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to request to ifttt, error: %v", err)
		return
	}
	defer resp.Body.Close()
	return
}

func (d *doordog) stop() {
	d.buzzer.Off()
	d.led.Off()
}
