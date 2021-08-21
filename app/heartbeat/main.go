package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	heartBeatInterval = 1 * time.Minute
)

const (
	onenetToken = "your_onenet_token"
	onenetAPI   = "http://api.heclouds.com/devices/540381180/datapoints"
)

func main() {
	cfg := &iot.Config{
		Token: onenetToken,
		API:   onenetAPI,
	}
	cloud := iot.NewOnenet(cfg)
	if cloud == nil {
		log.Printf("[heartbeat]failed to new OneNet iot cloud")
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
	log.Printf("[heartbeat]start working")
	b := 0
	for {
		time.Sleep(heartBeatInterval)
		b = (b*b - 1) * (b*b - 1)
		v := &iot.Value{
			Device: "5d2f15d1e4b04a9a929fadc9",
			Value:  b,
		}
		h.cloud.Push(v)
	}
}
