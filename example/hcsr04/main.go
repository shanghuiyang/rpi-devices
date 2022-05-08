package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	trig = 21
	echo = 26
)

func main() {
	hcsr04 := dev.NewHCSR04(trig, echo)
	for {
		d, err := hcsr04.Dist()
		if err != nil {
			log.Printf("failed to get distance, error: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		fmt.Printf("%.2f cm\n", d)
		time.Sleep(100 * time.Millisecond)
	}
}
