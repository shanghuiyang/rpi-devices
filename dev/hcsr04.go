/*
HC-SR04 is an ultrasonic distance meter used to measure the distance to objects.
min distance: 2cm
max distance: 600cm

Spec:
  - power supply:	+5V DC
  - range:			2 - 450cm
  - resolution:		0.3cm
    ___________________________
    |                           |
    |          HC-SR04          |
    |                           |
    |___________________________|
    |     |     |     |
    vcc  trig   echo  gnd

Connect to Raspberry Pi:
  - vcc:	any 5v pin
  - gnd:	any gnd pin
  - trig:	any data pin
  - echo:	any data pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	hcsr04MaxDistance = 999  // cm
	hcsr04Timeout     = 1000 // Nanosecond, 612m
	hcsr04Rety        = 5    // must >= 3
)

// HCSR04 implements DistanceMeter interface
type HCSR04 struct {
	trig rpio.Pin
	echo rpio.Pin
}

// NewHCSR04 ...
func NewHCSR04(trig int8, echo int8) *HCSR04 {
	hc := &HCSR04{
		trig: rpio.Pin(trig),
		echo: rpio.Pin(echo),
	}
	hc.trig.Output()
	hc.trig.Low()
	hc.echo.Input()
	return hc
}

// Dist returns distance in cm to objects.
// It takes 10 measurements, drops the 2 largest and 2 smallest, and averages the rest.
func (hc *HCSR04) Dist() (float64, error) {
	max, min, sum := -9999.0, 9999.0, 0.0
	for i := 0; i < hcsr04Rety; i++ {
		d, _ := hc.dist()

		if d > max {
			max = d
		}
		if d < min {
			min = d
		}

		sum += d
		time.Sleep(1 * time.Microsecond)
	}
	return (sum - max - min) / (us100Retry - 2), nil
}

func (hc *HCSR04) dist() (float64, error) {
	hc.trig.Low()
	delayUs(1)
	hc.trig.High()
	delayUs(1)

	for i := 0; hc.echo.Read() != rpio.High; i++ {
		if i >= hcsr04Timeout {
			return hcsr04MaxDistance, nil
		}
		delayNs(1)
	}

	start := time.Now()
	for i := 0; hc.echo.Read() != rpio.Low; i++ {
		if i >= hcsr04Timeout {
			return hcsr04MaxDistance, nil
		}
		delayNs(1)
	}
	return time.Since(start).Seconds() * voiceSpeed / 2.0, nil
}

// Close ...
func (hc *HCSR04) Close() error {
	return nil
}
