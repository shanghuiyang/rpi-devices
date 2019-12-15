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
	air := dev.NewPMS7003()
	led := dev.NewLed(pinLed)
	oled, err := dev.NewOLED(128, 32)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}

	wsnCfg := &base.WsnConfig{
		Token: "47ccbab9769d6ce64fd9d8b03ef63d9e",
		API:   "http://www.wsncloud.com/api/data/v1/numerical/insert",
	}
	cloud := iot.NewCloud(wsnCfg)

	asst := newHomeAsst(dht11, air, oled, led, cloud)
	base.WaitQuit(func() {
		asst.stop()
		rpio.Close()
	})
	asst.start()
}

type value struct {
	temp float64
	humi float64
	pm25 uint16
}

type homeAsst struct {
	dht11     *dev.DHT11
	air       *dev.PMS7003
	oled      *dev.OLED
	led       *dev.Led
	cloud     iot.Cloud
	chDisplay chan *value // for disploying on oled
	chCloud   chan *value // for pushing to iot cloud
	chAlert   chan *value // for alerting
}

func newHomeAsst(dht11 *dev.DHT11, air *dev.PMS7003, oled *dev.OLED, led *dev.Led, cloud iot.Cloud) *homeAsst {
	return &homeAsst{
		dht11:     dht11,
		air:       air,
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
	go h.alert()
	h.getData()
}

func (h *homeAsst) getData() {
	for {
		temp, humi, err := h.dht11.TempHumidity()
		if err != nil {
			log.Printf("temp|humidity: failed to get temp and humidity, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("temp|humidity: temp: %v, humidity: %v", temp, humi)

		pm25, _, err := h.air.Get()
		if err != nil {
			log.Printf("pm25: failed to get pm2.5, error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("pm2.5: %v ug/m3", pm25)

		v := &value{
			temp: temp,
			humi: humi,
			pm25: pm25,
		}
		h.chDisplay <- v
		h.chCloud <- v
		h.chAlert <- v
		time.Sleep(60 * time.Second)
	}
}

func (h *homeAsst) display() {
	var (
		temp, humi float64 = -999, -999
		pm25       uint16
	)
	on := true
	for {
		select {
		case v := <-h.chDisplay:
			temp, humi, pm25 = v.temp, v.humi, v.pm25
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

		hText := " --"
		if humi > 0 {
			hText = fmt.Sprintf("%.0f%%", humi)
		}
		if err := h.oled.Display(hText, 35, 0, 35); err != nil {
			log.Printf("display: failed to display humidity, error: %v", err)
		}
		time.Sleep(3 * time.Second)

		pmText := "  --"
		if pm25 > 0 {
			pmText = fmt.Sprintf("p%3d", pm25)
		}
		if err := h.oled.Display(pmText, 35, 0, 35); err != nil {
			log.Printf("display: failed to display pm2.5, error: %v", err)
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

			pv := &iot.Value{
				Device: "5df507c4e4b04a9a92a64928",
				Value:  v.pm25,
			}
			if err := h.cloud.Push(pv); err != nil {
				log.Printf("push: failed to push pm2.5 to cloud, error: %v", err)
			}
		}(v)
	}
}

func (h *homeAsst) alert() {
	var (
		temp, humi float64 = -999, -999
		pm25       uint16
	)
	for {
		select {
		case v := <-h.chAlert:
			temp, humi, pm25 = v.temp, v.humi, v.pm25
		default:
			// do nothing
		}

		if (temp > 0 && temp < 15) || humi > 70 || pm25 > 100 {
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
