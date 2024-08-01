package dev

import "github.com/stianeikeland/go-rpio/v4"

// WaterFlowSensor implements FlowSensor interface
type WaterFlowSensor struct {
	pin rpio.Pin
}

// NewWaterFlowMeter ...
func NewWaterFlowSensor(pin uint8) *WaterFlowSensor {
	w := &WaterFlowSensor{
		pin: rpio.Pin(pin),
	}
	w.pin.Input()
	return w
}

func (w *WaterFlowSensor) Flowing() bool {
	return w.pin.Read() == rpio.High
}
