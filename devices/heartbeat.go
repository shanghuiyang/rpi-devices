package devices

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/iotclouds"
)

const (
	logTagHeartBeat   = "heartbeat"
	heartBeatInterval = 2 * time.Minute
)

// HeartBeat ...
type HeartBeat struct {
}

// NewHeartBeat ...
func NewHeartBeat() *HeartBeat {
	return &HeartBeat{}
}

// Start ...
func (h *HeartBeat) Start() {
	log.Printf("[%v]start working", logTagHeartBeat)
	b := 0
	for {
		time.Sleep(heartBeatInterval)
		b = (b*b - 1) * (b*b - 1)
		v := &iotclouds.IoTValue{
			DeviceName: HeartBeatDevice,
			Value:      b,
		}
		iotclouds.IotCloud.Push(v)
		ChLedOp <- Blink
	}
}
