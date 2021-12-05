package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

func main() {
	hdc, err := dev.NewHDC1080()
	if err != nil {
		log.Printf("failed to create MPU6050 sensor, error: %v", err)
		return
	}

	util.WaitQuit(func() { hdc.Close() })
	for {
		util.DelayMs(1000)
		t, h, err := hdc.TempHumidity()
		if err != nil {
			log.Printf("failed to get temp & humi, error: %v", err)
			continue
		}
		log.Printf("temp=%.2f°C, humi=%.2f%%", t, h)
	}
}
