/*
ch2omonitor detects the concentration of CH2O in the air
which works with ZE08-CH2O, a CH2O sensor.
It will give you a warning when the CH2O concentration more than 0.08 mg/m3
via a blinking led light and a beeping buzzer.

The CH2O concentration will be displayed on a led display screen,
and it also be pushed to iot cloud for drawing a line chart.
*/

package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	pinBzr  = 17
	pinLed  = 26
	dioPin  = 11
	rclkPin = 9
	sclkPin = 10
)

const (
	alertCH2O = float64(0.08)
)

const (
	wsnToken        = "your_wsn_token"
	wsnNumericalAPI = "http://www.wsncloud.com/api/data/v1/numerical/insert"
)

func main() {
	ze08 := dev.NewZE08CH2O()
	led := dev.NewLedImp(pinLed)
	bzr := dev.NewBuzzerImp(pinBzr, true)
	dsp := dev.NewLedDisplay(dioPin, rclkPin, sclkPin)

	cfg := &iot.Config{
		Token: wsnToken,
		API:   wsnNumericalAPI,
	}
	wsn := iot.NewWsn(cfg)

	m := newCH2OMonitor(ze08, led, bzr, dsp, wsn)
	util.WaitQuit(func() {
		m.stop()
	})
	m.start()
}

type ch2oMonitor struct {
	ch2ometer dev.CH2OMeter
	led       dev.Led
	buzzer    dev.Buzzer
	dsp       dev.Display
	cloud     iot.Cloud
	mode      util.Mode
	chAlert   chan float64 // for alerting
	chDisplay chan float64
	chCloud   chan float64 // for pushing to iot cloud
}

func newCH2OMonitor(ch2ometer dev.CH2OMeter, led dev.Led, buzzer dev.Buzzer, dsp dev.Display, cloud iot.Cloud) *ch2oMonitor {
	return &ch2oMonitor{
		ch2ometer: ch2ometer,
		led:       led,
		buzzer:    buzzer,
		dsp:       dsp,
		cloud:     cloud,
		mode:      util.PrdMode,
		chAlert:   make(chan float64, 4),
		chDisplay: make(chan float64, 4),
		chCloud:   make(chan float64, 4),
	}
}

func (m *ch2oMonitor) start() {
	log.Printf("[ch2omonitor]service starting")
	log.Printf("[ch2omonitor]mode: %v", m.mode)
	go m.alert()
	go m.push()
	go m.display()
	m.detect()
}

func (m *ch2oMonitor) detect() {
	log.Printf("[ch2omonitor]detecting ch2o")
	for {
		v, err := m.ch2ometer.Value()
		if err != nil {
			log.Printf("[ch2omonitor]failed to get ch2o, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("[ch2omonitor]ch2o: %.4f mg/m3", v)

		m.chAlert <- v
		m.chCloud <- v
		m.chDisplay <- v

		sec := 60 * time.Second
		if m.mode != util.PrdMode {
			sec = 15 * time.Second
		}
		time.Sleep(sec)
	}
}

func (m *ch2oMonitor) push() {
	for ch2o := range m.chCloud {
		if m.mode != util.PrdMode {
			continue
		}
		go func(ch2o float64) {
			v := &iot.Value{
				Device: "5e134f95e4b04a9a92a79665",
				Value:  math.Round(ch2o*10000) / 10000,
			}
			if err := m.cloud.Push(v); err != nil {
				log.Printf("[ch2omonitor]push: failed to push ch2o to cloud, error: %v", err)
			}
		}(ch2o)
	}
}

func (m *ch2oMonitor) alert() {
	var ch2o float64
	for {
		select {
		case v := <-m.chAlert:
			ch2o = v
		default:
			// do nothing
		}

		if ch2o >= alertCH2O {
			go m.buzzer.Beep(1, 200)
			go m.led.Blink(1, 200)
		}
		time.Sleep(1 * time.Second)
	}
}

func (m *ch2oMonitor) display() {
	var ch2o float64
	m.dsp.Open()
	opened := true
	for {
		select {
		case v := <-m.chDisplay:
			ch2o = v
		default:
			// do nothing, just use the latest temp
		}

		if m.dsp == nil {
			time.Sleep(30 * time.Second)
			continue
		}

		hour := time.Now().Hour()
		if ch2o < alertCH2O && (hour >= 20 || hour < 8) {
			// turn off oled at 20:00-08:00
			if opened {
				m.dsp.Close()
				opened = false
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if !opened {
			m.dsp.Open()
			opened = true
		}
		text := "----"
		if ch2o > 0 {
			text = fmt.Sprintf("%.3f", ch2o)
		}
		m.dsp.Display(text)
		time.Sleep(3 * time.Second)
	}
}

func (m *ch2oMonitor) stop() {
	m.ch2ometer.Close()
	m.led.Off()
	m.buzzer.Off()
	m.dsp.Close()
}
