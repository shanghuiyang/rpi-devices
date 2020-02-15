package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	u := dev.NewUS100()
	for {
		dist := u.Dist()
		if dist < 0 {
			log.Printf("failed to get distance")
			continue
		}
		log.Printf("%.2f cm\n", dist)
		time.Sleep(50 * time.Millisecond)
	}
}
