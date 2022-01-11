package main

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	pin = 26
)

func main() {
	r := dev.NewRelayImp(pin)
	r.On()
	time.Sleep(5 * time.Second)
	r.Off()
}
