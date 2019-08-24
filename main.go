package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/shanghuiyang/rpi-devices/base"
	s "github.com/shanghuiyang/rpi-devices/devices"
	"github.com/shanghuiyang/rpi-devices/iotclouds"
)

var (
	devices []s.Device
)

func main() {
	cfg, err := base.LoadConfig()
	if err != nil {
		panic(err)
	}

	base.Init(cfg)
	iotclouds.Init(cfg.OneNet)

	m := s.NewMemory()
	c := s.NewCPU()
	h := s.NewHeartBeat()

	l := s.NewLed(cfg.Led.Pin)
	r := s.NewRelay(cfg.Relay.Pin)
	t := s.NewTemperature()
	g := s.NewGPS()
	devices = append(devices, m, c, h, l, r, t, g)

	for _, d := range devices {
		go d.Start()
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}
