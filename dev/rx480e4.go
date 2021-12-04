/*
RX480E4 is a transmitter receiver decoding module with 4 channels.

Connect to Raspberry Pi:
 - +v: any v3.3 or 5v pin
 - gnd: any gnd pin
 - d0: any data pin
 - d1: any data pin
 - d2: any data pin
 - d3: any data pin
 - vt: null

*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

// RX480E4 implements RFReceiver
type RX480E4 struct {
	channels [4]rpio.Pin
}

// NewRX480E4 ...
func NewRX480E4(d0, d1, d2, d3 uint8) *RX480E4 {
	channels := [4]rpio.Pin{rpio.Pin(d0), rpio.Pin(d1), rpio.Pin(d2), rpio.Pin(d3)}
	rx := &RX480E4{
		channels: channels,
	}
	for _, ch := range rx.channels {
		ch.Input()
		ch.PullDown()
		ch.Detect(rpio.RiseEdge)
	}
	return rx
}

// PressA ...
func (rx *RX480E4) Received(ch int) bool {
	if ch < 0 || ch > 3 {
		return false
	}
	return rx.channels[ch].EdgeDetected()
}
