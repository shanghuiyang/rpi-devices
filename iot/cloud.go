package iot

import (
	"github.com/shanghuiyang/rpi-devices/base"
)

// Cloud is the interface of IOT clound
type Cloud interface {
	Push(v *Value) error
}

// Value ...
type Value struct {
	Device string
	Value  interface{}
}

// NewCloud ...
func NewCloud(config interface{}) Cloud {
	var cloud Cloud
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
