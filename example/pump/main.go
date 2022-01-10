package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 26
)

func main() {
	pump := dev.NewPumpImp(pin)

	pump.On()
	time.Sleep(5 * time.Second)
	pump.Off()
}
