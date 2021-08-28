package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

func init() {
	if err := rpio.Open(); err != nil {
		panic("failed to init rpio")
	}
}
