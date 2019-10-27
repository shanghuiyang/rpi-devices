/*
Auto-Light let you control a led light working with a voice detector together
the led will light up when the voice detector detect a voice.
and the led will turn off after 35 seconds.

voice detector:
 - vcc: phys.1/3.3v
 - out: phys.3/BCM.2
 - gnd: phys.9/GND

 led:
  - positive: phys.36/BCM.16
  - negative: phys.34/GND

-----------------------------------------------------------------------

          +---------+
          |         |
          | voice   |
          | detector|
          |         |
          +-+--+--+-+
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
                              | o       o |        +--+-+--+
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

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	voicePin = 2
	ledPin   = 16
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	v := dev.NewVoiceDetector(voicePin)
	if v == nil {
		log.Printf("failed to new a voice detector")
		return
	}

	light := dev.NewLed(ledPin)
	if light == nil {
		log.Printf("failed to new a led light")
		return
	}

	l := &autoLight{
		voice: v,
		light: light,
		ch:    make(chan bool, 32),
	}
	l.start()
}

type autoLight struct {
	voice *dev.VoiceDetector
	light *dev.Led
	ch    chan bool
}

func (a *autoLight) start() {
	go a.listen()

	a.light.Off()
	isLightOn := false
	lastTrig := time.Now()
	for b := range a.ch {
		if b {
			log.Printf("detected a voice")
			if !isLightOn {
				a.light.On()
				isLightOn = true
			}
			lastTrig = time.Now()
			continue
		}
		if time.Now().Sub(lastTrig).Seconds() > 35 && isLightOn {
			log.Printf("timeout, light off")
			a.light.Off()
			isLightOn = false
		}
	}
}

func (a *autoLight) listen() {
	for {
		a.ch <- a.voice.Detected()
		time.Sleep(10 * time.Millisecond)
	}
}
