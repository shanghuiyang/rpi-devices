package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/iot"
)

const (
	heartBeatInterval = 1 * time.Minute
)

func main() {
	oneNetCfg := &base.WsnConfig{
		Token: base.WsnToken,
		API:   base.WsnNumericalAPI,
	}
	cloud := iot.NewCloud(oneNetCfg)
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
