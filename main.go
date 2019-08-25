package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shanghuiyang/rpi-devices/base"
	dev "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/shanghuiyang/rpi-devices/iotclouds"
)

var (
	devices []dev.Device
)

func main() {
	cfg, err := base.LoadConfig()
	if err != nil {
		panic(err)
	}

	base.Init(cfg)
	iotclouds.Init(cfg.OneNet)

	m := dev.NewMemory()
	c := dev.NewCPU()
	h := dev.NewHeartBeat()

	l := dev.NewLed(cfg.Led.Pin)
	r := dev.NewRelay(cfg.Relay.Pin)
	t := dev.NewTemperature()
	g := dev.NewGPS()
	devices = append(devices, m, c, h, l, r, t, g)

	for _, d := range devices {
		go d.Start()
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}
