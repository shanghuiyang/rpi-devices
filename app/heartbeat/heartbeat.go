package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	heartBeatInterval = 2 * time.Minute
)

func main() {
	oneNetCfg := &base.OneNetConfig{
		Token: "your token",
		API:   "http://api.heclouds.com/devices/540381180/datapoints",
	}
	cloud := iot.NewCloud(oneNetCfg)
	if cloud == nil {
		log.Printf("failed to new OneNet iot cloud")
		return
	}
	h := &heartBeat{
		cloud: cloud,
	}
	h.start()
}

type heartBeat struct {
	cloud iot.Cloud
}

// Start ...
func (h *heartBeat) start() {
	log.Printf("heart beat start working")
	b := 0
	for {
		time.Sleep(heartBeatInterval)
		b = (b*b - 1) * (b*b - 1)
		v := &iot.Value{
			DeviceName: "heartbeat",
			Value:      b,
		}
		h.cloud.Push(v)
	}
}
