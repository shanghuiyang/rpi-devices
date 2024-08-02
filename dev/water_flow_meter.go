package dev

import "github.com/stianeikeland/go-rpio/v4"

// WaterFlowMeter implements Detector interface
type WaterFlowMeter struct {
	pin rpio.Pin
}

// NewWaterFlowMeter ...
func NewWaterFlowMeter(pin uint8) *WaterFlowMeter {
	w := &WaterFlowMeter{
		pin: rpio.Pin(pin),
	}
	w.pin.Input()
	return w
}

func (w *WaterFlowMeter) Detected() bool {
	return w.pin.Read() == rpio.High
}
