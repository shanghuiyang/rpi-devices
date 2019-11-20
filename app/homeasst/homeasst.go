package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinLed = 26
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	dht11 := dev.NewDHT11()
	led := dev.NewLed(pinLed)
	oled, err := dev.NewOLED(128, 32)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}

	wsnCfg := &base.WsnConfig{
		Token: "your token",
		API:   "http://www.wsncloud.com/api/data/v1/numerical/insert",
	}
	cloud := iot.NewCloud(wsnCfg)

	asst := newHomeAsst(dht11, oled, led, cloud)
	base.WaitQuit(func() {
		asst.stop()
		rpio.Close()
	})
	asst.start()
}

type value struct {
	temp float64
	humi float64
}

type homeAsst struct {
	dht11     *dev.DHT11
	oled      *dev.OLED
	led       *dev.Led
	cloud     iot.Cloud
	chDisplay chan *value // for disploying on oled
	chCloud   chan *value // for pushing to iot cloud
	chAlert   chan *value // for alerting
}

func newHomeAsst(dht11 *dev.DHT11, oled *dev.OLED, led *dev.Led, cloud iot.Cloud) *homeAsst {
	return &homeAsst{
		dht11:     dht11,
		oled:      oled,
		led:       led,
		cloud:     cloud,
		chDisplay: make(chan *value, 4),
		chCloud:   make(chan *value, 4),
		chAlert:   make(chan *value, 4),
	}
}

func (h *homeAsst) start() {
	go h.display()
	go h.push()
	go h.alert(false)
	h.getTempHumidity()
}

func (h *homeAsst) getTempHumidity() {
	for {
		temp, humi, err := h.dht11.TempHumidity()
		if err != nil {
			log.Printf("temp|humidity: failed to get temp and humidity, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("temp|humidity: temp: %v, humidity: %v", temp, humi)

		v := &value{
			temp: temp,
			humi: humi,
		}
		h.chDisplay <- v
		h.chCloud <- v
		h.chAlert <- v
		time.Sleep(60 * time.Second)
	}
}

func (h *homeAsst) display() {
	var temp, humi float64 = -999, -999
	on := true
	for {
		select {
		case v := <-h.chDisplay:
			temp, humi = v.temp, v.humi
		default:
			// do nothing, just use the latest temp
		}

		hour := time.Now().Hour()
		if hour >= 20 || hour < 8 {
			// turn off oled at 20:00-08:00
			if on {
				h.oled.Off()
				on = false
			}
			time.Sleep(10 * time.Second)
			continue
		}

		on = true
		tText := "--"
		if temp > -273 {
			tText = fmt.Sprintf("%.0f'C", temp)
		}
		if err := h.oled.Display(tText, 35, 0, 35); err != nil {
			log.Printf("display: failed to display temperature, error: %v", err)
		}
		time.Sleep(3 * time.Second)

		hText := "  --"
		if humi > 0 {
			hText = fmt.Sprintf("%.0f%%", humi)
		}
		if err := h.oled.Display(hText, 35, 0, 35); err != nil {
			log.Printf("display: failed to display humidity, error: %v", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func (h *homeAsst) push() {
	for v := range h.chCloud {
		go func(v *value) {
			tv := &iot.Value{
				Device: "5d3c467ce4b04a9a92a02343",
				Value:  v.temp,
			}
			if err := h.cloud.Push(tv); err != nil {
				log.Printf("push: failed to push temperature to cloud, error: %v", err)
			}

			hv := &iot.Value{
				Device: "5d3c4627e4b04a9a92a02342",
				Value:  v.humi,
			}
			if err := h.cloud.Push(hv); err != nil {
				log.Printf("push: failed to push humidity to cloud, error: %v", err)
			}
		}(v)
	}
}

func (h *homeAsst) alert(enable bool) {
	var temp, humi float64 = -999, -999
	for {
		select {
		case v := <-h.chAlert:
			temp, humi = v.temp, v.humi
		default:
			// do nothing
		}

		if enable && ((temp > -273 && temp < 18) || temp > 32 || (humi > 0 && humi < 50) || humi > 60) {
			h.led.On()
			time.Sleep(1 * time.Second)
			h.led.Off()
			time.Sleep(1 * time.Second)
			continue
		}
		h.led.Off()
		time.Sleep(1 * time.Second)
	}
}

func (h *homeAsst) stop() {
	h.oled.Close()
	h.led.Off()
}
