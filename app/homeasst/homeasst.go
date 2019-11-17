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
	pinLed   = 26
	pinLight = 16
	pinInfr  = 18
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	dht11 := dev.NewDHT11()

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

	infr := dev.NewInfraredDetector(pinInfr)
	light := dev.NewLed(pinLight)

	asst := newHomeAsst(dht11, oled, infr, light, cloud)

	base.WaitQuit(func() {
		asst.stop()
		rpio.Close()
	})
	asst.start()
}

type homeAsst struct {
	dht11    *dev.DHT11
	oled     *dev.OLED
	infrared *dev.InfraredDetector
	light    *dev.Led
	cloud    iot.Cloud

	chDspTemp   chan float64 // for disploy on oled
	chDspHumi   chan float64 // for disploy on oled
	chCloudTemp chan float64 // for push to iot cloud
	chCloudHumi chan float64 // for push to iot cloud
	chObj       chan bool
}

func newHomeAsst(dht11 *dev.DHT11, oled *dev.OLED, infr *dev.InfraredDetector, light *dev.Led, cloud iot.Cloud) *homeAsst {
	return &homeAsst{
		dht11:       dht11,
		oled:        oled,
		infrared:    infr,
		light:       light,
		cloud:       cloud,
		chDspTemp:   make(chan float64, 4),
		chDspHumi:   make(chan float64, 4),
		chCloudTemp: make(chan float64, 4),
		chCloudHumi: make(chan float64, 4),
		chObj:       make(chan bool, 32),
	}
}

func (h *homeAsst) start() {
	go h.display()
	go h.push()
	go h.detect()
	go h.alight()

	h.getTempHumidity()
}

func (h *homeAsst) getTempHumidity() {
	for {
		temp, humi, err := h.dht11.TempHumidity()
		if err != nil {
			log.Printf("failed to get temp and humidity, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		h.chDspTemp <- temp
		h.chDspHumi <- humi

		h.chCloudTemp <- temp
		h.chCloudHumi <- humi
		time.Sleep(30 * time.Second)
	}
}

func (h *homeAsst) display() {
	var (
		temp float64 = -999
		humi float64 = -999
	)
	for {
		select {
		case v := <-h.chDspTemp:
			temp = v
		default:
			// do nothing, just use the latest temp
		}

		select {
		case v := <-h.chDspHumi:
			humi = v
		default:
			// do nothing, just use the latest humidity
		}

		tText := "N/A"
		if temp > -273 {
			tText = fmt.Sprintf("%.0f'C", temp)
		}
		if err := h.oled.Display(tText, 35, 0, 35); err != nil {
			log.Printf("failed to display temperature, error: %v", err)
		}
		time.Sleep(3 * time.Second)

		hText := "N/A"
		if humi > 0 {
			hText = fmt.Sprintf("%.0f%%", humi)
		}
		if err := h.oled.Display(hText, 35, 0, 35); err != nil {
			log.Printf("failed to display humidity, error: %v", err)
		}
		time.Sleep(3 * time.Second)
	}
}

func (h *homeAsst) push() {
	for {
		select {
		case v := <-h.chCloudTemp:
			tv := &iot.Value{
				Device: "5d3c467ce4b04a9a92a02343",
				Value:  v,
			}
			go func() {
				if err := h.cloud.Push(tv); err != nil {
					log.Printf("failed to push temperature to cloud, error: %v", err)
				}
			}()

		case v := <-h.chCloudHumi:
			hv := &iot.Value{
				Device: "5d3c4627e4b04a9a92a02342",
				Value:  v,
			}
			go func() {
				if err := h.cloud.Push(hv); err != nil {
					log.Printf("failed to push humidity to cloud, error: %v", err)
				}
			}()
		}
	}
}

func (h *homeAsst) detect() {
	for {
		h.chObj <- h.infrared.Detected()
		time.Sleep(100 * time.Millisecond)
	}
}

func (h *homeAsst) alight() {
	h.light.Off()
	isLightOn := false
	lastTrig := time.Now()
	for b := range h.chObj {
		if b {
			log.Printf("alight: detected an object")
			if !isLightOn {
				h.light.On()
				isLightOn = true
			}
			lastTrig = time.Now()
			continue
		}
		if time.Now().Sub(lastTrig).Seconds() > 30 && isLightOn {
			log.Printf("alight: timeout, light off")
			h.light.Off()
			isLightOn = false
		}
	}
}

func (h *homeAsst) stop() {
	h.oled.Close()
	h.light.Off()
}
