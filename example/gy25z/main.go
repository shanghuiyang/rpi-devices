package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 115200
)

func main() {
	g, err := dev.NewGY25Z(devName, baud)
	if err != nil {
		log.Fatalf("new gy25 error: %v", err)
	}
	defer g.Close()

	// if err := g.SetMode(dev.GY25ZAutoMode); err != nil {
	// 	log.Printf("failed to set mode, error: %v", err)
	// 	return
	// }

	for {
		time.Sleep(100 * time.Millisecond)
		yaw, pitch, roll, err := g.Angles()
		if err != nil {
			log.Printf("failed to gent angles, error: %v", err)
			continue
		}
		fmt.Printf("yaw: %.2f, pitch: %.2f, roll: %.2f\n", yaw, pitch, roll)
	}
}
