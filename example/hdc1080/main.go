package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	hdc, err := dev.NewHDC1080()
	if err != nil {
		log.Printf("failed to create MPU6050 sensor, error: %v", err)
		return
	}

	defer func() {
		hdc.Close()
	}()

	for {
		time.Sleep(1 * time.Second)
		t, h, err := hdc.TempHumidity()
		if err != nil {
			log.Printf("failed to get temp & humi, error: %v", err)
			continue
		}
		log.Printf("temp=%.2f°C, humi=%.2f%%", t, h)
	}
}
