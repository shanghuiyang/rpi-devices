/*
Auto-Light let you control a led light working with a infrared detector together
the led light will light up when the infrared detector detects objects.
and the led will turn off after 30 seconds.

infrared detector:
 - vcc: pin 1 or any 3.3v
 - gnd: ping 9 or any gnd pin
 - out: pin 3(gpio 2) or any date pin

 led:
  - positive: pin 36(gpio 16) or any date pin
  - negative: pin 34 or any gnd pin

-----------------------------------------------------------------------

          o---------o
          |         |
          | Infrared|
          | detector|
          |         |
          o-+--+--+-o
            |  |  |
          gnd out vcc
            |  |  |           +-----------+
            |  |  +-----------+ * 1   2 o |
            +--|--------------+ * 3     o |
               |              | o       o |
               |              | o       o |         \ | | /
               +--------------+ * 9     o |           ___
                              | o       o |         /     \
                              | o       o |        |-------|
                              | o       o |        |  led  |
                              | o       o |        |       |
                              | o       o |        o--+-+--o
                              | o       o |           | |
                              | o       o |         gnd vcc
                              | o       o |           | |
                              | o       o |           | |
                              | o       o |           | |
                              | o       o |           | |
                              | o    34 * +-----------+ |
                              | o    36 * +-------------+
                              | o       o |
                              | o 39 40 o |
                              +-----------+

-----------------------------------------------------------------------

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
	pinInfra = 18
	pinLed   = 26
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	infr := dev.NewInfraredDetector(pinInfra)
	led := dev.NewLed(pinLed)
	light := dev.NewLed(pinLight)
	if light == nil {
		log.Printf("failed to new a led light")
		return
	}

	wsnCfg := &base.WsnConfig{
		Token: base.WsnToken,
		API:   base.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(wsnCfg)

	a := &autoLight{
		infra: infr,
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
	light *dev.Led
	led   *dev.Led
	cloud iot.Cloud
	ch    chan bool
}

func (a *autoLight) start() {
	go a.detect()

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

func (a *autoLight) detect() {
	for {
		hour := time.Now().Hour()
		if hour >= 8 && hour < 18 {
			// disable autolight between 8:00-18:00
			time.Sleep(1 * time.Minute)
			continue
		}
		detected := a.infra.Detected()
		a.ch <- detected

		t := 200 * time.Millisecond
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
