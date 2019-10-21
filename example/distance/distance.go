package main

import (
	"fmt"

	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/stianeikeland/go-rpio"
)

const (
	pinTrig = 21
	pinEcho = 26
)

func main() {

	if err := rpio.Open(); err != nil {
		return
	}
	h := dev.NewHCSR04(pinTrig, pinEcho)
	d := h.Dist()
	fmt.Printf("d: %v\n", d)

	rpio.Close()
}
