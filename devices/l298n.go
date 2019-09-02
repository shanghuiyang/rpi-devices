package devices

import (
	"github.com/stianeikeland/go-rpio"
)

// L298N ...
type L298N struct {
	In1 rpio.Pin
	In2 rpio.Pin
	In3 rpio.Pin
	In4 rpio.Pin
}

// NewL298N ...
func NewL298N(in1, in2, in3, in4 uint8) *L298N {
	if err := rpio.Open(); err != nil {
		return nil
	}
	l := &L298N{
		In1: rpio.Pin(in1),
		In2: rpio.Pin(in2),
		In3: rpio.Pin(in3),
		In4: rpio.Pin(in4),
	}
	l.In1.Output()
	l.In2.Output()
	l.In3.Output()
	l.In4.Output()
	l.In1.Low()
	l.In2.Low()
	l.In3.Low()
	l.In4.Low()
	return l
}

// Close ...
func (l *L298N) Close() {
	rpio.Close()
}
