/*
MQ-7 is a sensor used to detected co gas.

Connect to Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 -  do: any data pin
 -  ao: (don't support currently)

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// MQ7 implements Detector interface
type MQ7 struct {
	do rpio.Pin
	ao rpio.Pin
}

// NewMQ7 ...
func NewMQ7(do uint8) *MQ7 {
	mq7 := &MQ7{
		do: rpio.Pin(do),
	}
	mq7.do.Input()
	return mq7
}

// Detected ...
func (mq7 *MQ7) Detected() bool {
	return mq7.do.Read() == rpio.Low
}
