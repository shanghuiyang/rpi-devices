package devices

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	timeout = 3600
	repeat  = 1
)

// HCSR04 ...
type HCSR04 struct {
	trig rpio.Pin
	echo rpio.Pin
}

// NewHCSR04 ...
func NewHCSR04(trig int8, echo int8) *HCSR04 {
	h := &HCSR04{
		trig: rpio.Pin(trig),
		echo: rpio.Pin(echo),
	}
	h.trig.Output()
	h.trig.Low()
	h.echo.Input()
	return h
}

// Dist is measure the distance in cm
func (h *HCSR04) Dist() float64 {
	h.trig.Low()
	h.delay(200)
	h.trig.High()
	h.delay(15)
	h.echo.Low()

	for n := 0; n < timeout && h.echo.Read() != rpio.High; n++ {
		h.delay(1)
	}
	start := time.Now()

	for n := 0; n < timeout && h.echo.Read() != rpio.Low; n++ {
		h.delay(1)
	}
	end := time.Now()
	return end.Sub(start).Seconds() * 34300.0 / 2.0
}

// delay is to dalay us microsecond
func (h *HCSR04) delay(us int) {
	time.Sleep(time.Duration(us) * time.Microsecond)
}
