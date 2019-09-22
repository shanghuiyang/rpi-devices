package devices

import (
	"time"

	"github.com/stianeikeland/go-rpio"
)

// L298N ...
type L298N struct {
	in1 rpio.Pin
	in2 rpio.Pin
	in3 rpio.Pin
	in4 rpio.Pin
}

// NewL298N ...
func NewL298N(in1, in2, in3, in4 uint8) *L298N {
	l := &L298N{
		in1: rpio.Pin(in1),
		in2: rpio.Pin(in2),
		in3: rpio.Pin(in3),
		in4: rpio.Pin(in4),
	}
	l.in1.Output()
	l.in2.Output()
	l.in3.Output()
	l.in4.Output()
	l.in1.Low()
	l.in2.Low()
	l.in3.Low()
	l.in4.Low()
	return l
}

// Forward ...
func (l *L298N) Forward() {
	l.in1.High()
	l.in2.Low()
	time.Sleep(70 * time.Millisecond)
	l.in3.High()
	l.in4.Low()

	l.in1.Low()
	time.Sleep(93 * time.Millisecond)
	l.in1.High()
}

// Backward ...
func (l *L298N) Backward() {
	l.in1.Low()
	l.in2.High()
	time.Sleep(70 * time.Millisecond)
	l.in3.Low()
	l.in4.High()

	l.in2.Low()
	time.Sleep(80 * time.Millisecond)
	l.in2.High()
}

// Left ...
func (l *L298N) Left() {
	l.in1.Low()
	l.in2.Low()
	l.in3.High()
	l.in4.Low()
}

// Right ...
func (l *L298N) Right() {
	l.in1.High()
	l.in2.Low()
	l.in3.Low()
	l.in4.Low()
}

// Stop ...
func (l *L298N) Stop() {
	l.in1.Low()
	l.in2.Low()
	l.in3.Low()
	l.in4.Low()
}
