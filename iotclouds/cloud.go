package iotclouds

import (
	"github.com/shanghuiyang/rpi-devices/base"
)

var (
	// IotCloud ...
	IotCloud   IOTCloud
	chIoTCloud = make(chan *IoTValue, 32)
)

// IOTCloud is the interface of IOT clound
type IOTCloud interface {
	Push(v *IoTValue) error
}

// IoTValue ...
type IoTValue struct {
	DeviceName string
	Value      interface{}
}

// New ...
func New(config interface{}) IOTCloud {
	var cloud IOTCloud
	switch config.(type) {
	case *base.WsnConfig:
		cfg := config.(*base.WsnConfig)
		cloud = NewWsnClound(cfg)
	case *base.OneNetConfig:
		cfg := config.(*base.OneNetConfig)
		cloud = NewOneNetCloud(cfg)
	default:
		cloud = nil
	}
	return cloud
}
