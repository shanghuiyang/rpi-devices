package iotclouds

import (
	"errors"

	"github.com/shanghuiyang/pi/base"
)

var (
	// IotCloud ...
	IotCloud   IOTCloud
	chIoTCloud = make(chan *IoTValue, 32)
)

// IOTCloud is the interface of IOT clound
type IOTCloud interface {
	Start()
	Push(v *IoTValue)
}

// IoTValue ...
type IoTValue struct {
	DeviceName string
	Value      interface{}
}

// Init ...
func Init(config interface{}) {
	switch config.(type) {
	case *base.WsnConfig:
		cfg := config.(*base.WsnConfig)
		IotCloud = NewWsnClound(cfg)
	case *base.OneNetConfig:
		cfg := config.(*base.OneNetConfig)
		IotCloud = NewOneNetClound(cfg)
	default:
		panic(errors.New("invaild config"))
	}
	go IotCloud.Start()
}
