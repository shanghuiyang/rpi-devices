package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
)

func main() {
	air, err := dev.NewPMS7003(devName, 9600)
	if err != nil {
		log.Printf("failed to new PMS7003, error: %v", err)
		return
	}
	pm25, pm10, err := air.Get()
	if err != nil {
		log.Printf("failed to get data, error: %v", err)
		return
	}
	log.Printf("pm2.5: %vug/m3, pm10: %vug/m3\n", pm25, pm10)

	air.Close()
}
