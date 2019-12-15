package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	air := dev.NewPMS7003()
	pm25, pm10, err := air.Get()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("pm2.5: %vug/m3, pm10: %vug/m3\n", pm25, pm10)

	air.Close()
}
