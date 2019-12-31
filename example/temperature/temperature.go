package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	d := dev.NewDS18B20()
	t, err := d.GetTemperature()
	if err != nil {
		fmt.Printf("failed to get temperature, error: %v", err)
		return
	}
	fmt.Printf("current temperature: %v", t)
}
