package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	air := dev.NewPMS7003()
	pm25, err := air.PM25()
	if err != nil {
		log.Printf("failed, error: %v", err)
		return
	}
	log.Printf("pm2.5: %v ug/m3\n", pm25)
		
	air.Close()
}
