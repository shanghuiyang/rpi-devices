package main

import (
	"fmt"

	"github.com/shanghuiyang/pi/devices"
)

func main() {
	t := devices.NewTemperature()
	c, err := t.GetTemperature()
	if err != nil {
		fmt.Printf("failed to get temperature, error: %v", err)
		return
	}
	fmt.Printf("current temperature: %v", c)
}
