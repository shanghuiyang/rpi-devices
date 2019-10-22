package main

import (
	"fmt"

	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	t := dev.NewTemperature()
	c, err := t.GetTemperature()
	if err != nil {
		fmt.Printf("failed to get temperature, error: %v", err)
		return
	}
	fmt.Printf("current temperature: %v", c)
}
