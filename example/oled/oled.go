package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	oled, err := dev.NewOLED(128, 32)
	if err != nil {
		log.Printf("failed to create an oled, error: %v", err)
		return
	}

	base.WaitQuit(oled.Close)
	for {
		if err := oled.Display("30'C", 35); err != nil {
			log.Printf("failed to display temperature, error: %v", err)
			break
		}
		time.Sleep(2 * time.Second)

		if err := oled.Display("75%", 35); err != nil {
			log.Printf("failed to display humidity, error: %v", err)
			break
		}
		time.Sleep(2 * time.Second)
	}
}
