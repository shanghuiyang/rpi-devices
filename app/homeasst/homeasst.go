package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/iot"
)

var (
	temperature float64
	humidity    float64
)

func main() {
	dht := dev.NewDHT11()
	wsnCfg := &base.WsnConfig{
		Token: "your token",
		API:   "http://www.wsncloud.com/api/data/v1/numerical/insert",
	}
	cloud := iot.NewCloud(wsnCfg)

	go display()
	for {
		time.Sleep(30 * time.Second)

		t, h, err := dht.TempHumidity()
		if err != nil {
			log.Printf("failed to get temp and humidity, error: %v", err)
			continue
		}
		tv := &iot.Value{
			Device: "5d3c467ce4b04a9a92a02343",
			Value:  t,
		}
		go func() {
			if err := cloud.Push(tv); err != nil {
				log.Printf("failed to push temperature to cloud, error: %v", err)
			}
		}()

		hv := &iot.Value{
			Device: "5d3c4627e4b04a9a92a02342",
			Value:  h,
		}
		go func() {
			if err := cloud.Push(hv); err != nil {
				log.Printf("failed to push humidity to cloud, error: %v", err)
			}
		}()

		temperature = t
		humidity = h
	}
}

func display() {
	oled, err := dev.NewOLED(128, 32)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}
	base.WaitQuit(oled.Close)
	for {
		t := fmt.Sprintf("%.0f'C", temperature)
		if err := oled.Display(t, 35, 0, 35); err != nil {
			log.Printf("failed to display temperature, error: %v", err)
		}
		time.Sleep(3 * time.Second)

		h := fmt.Sprintf("%.0f%%", humidity)
		if err := oled.Display(h, 35, 0, 35); err != nil {
			log.Printf("failed to display humidity, error: %v", err)
		}
		time.Sleep(3 * time.Second)
	}
}
