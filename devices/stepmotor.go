package devices

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio"
)

const (
	angleEachStep   = 0.087
	logTagStepMotor = "stepmotor"
)

var (
	// ChStepMotorOp ...
	ChStepMotorOp = make(chan Operator, 8)

	clockwise = [4][4]uint8{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	anticlockwise = [4][4]uint8{
		{0, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
	}
)

// StepMotor ...
type StepMotor struct {
	pins [4]rpio.Pin
}

// NewStepMotor ...
func NewStepMotor(in1, in2, in3, in4 uint8) *StepMotor {
	if err := rpio.Open(); err != nil {
		return nil
	}
	s := &StepMotor{
		pins: [4]rpio.Pin{
			rpio.Pin(in1),
			rpio.Pin(in2),
			rpio.Pin(in3),
			rpio.Pin(in4),
		},
	}
	for i := 0; i < 4; i++ {
		s.pins[i].Output()
		s.pins[i].Low()
	}
	return s
}

// Start ...
func (s *StepMotor) Start() {
	defer s.Close()

	log.Printf("[%v]start working", logTagRelay)
	for {
		op := <-ChStepMotorOp
		angle := float32(op)
		s.Roll(angle)
	}
}

// Roll ...
func (s *StepMotor) Roll(angle float32) {
	var matrix [4][4]uint8
	if angle > 0 {
		matrix = clockwise
	} else {
		matrix = anticlockwise
		angle = angle * (-1)
	}
	n := int(angle / angleEachStep / 8.0)
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				if matrix[j][k] == 1 {
					s.pins[k].High()
				} else {
					s.pins[k].Low()
				}
			}
			time.Sleep(2 * time.Millisecond)
		}
	}
}

// Stop ...
func (s *StepMotor) Stop() {
	for i := 0; i < 4; i++ {
		s.pins[i].Low()
	}
}

// Close ...
func (s *StepMotor) Close() {
	rpio.Close()
}
