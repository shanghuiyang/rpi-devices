package main

import (
	"fmt"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const pulsesPerLiter float32 = 450

func main() {
	w := dev.NewWaterFlowSensor(17)
	var numberOfPulsesCounted int = 0

	for {
		if w.Flowing() {
			numberOfPulsesCounted++
			fmt.Printf("MiliLiters Flowed: %f", (float32(1000*numberOfPulsesCounted) / pulsesPerLiter))
			continue
		}

		time.Sleep(100 * time.Millisecond)
	}
}
