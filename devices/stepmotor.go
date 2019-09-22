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
	chOp chan Operator
}

// NewStepMotor ...
func NewStepMotor(in1, in2, in3, in4 uint8) *StepMotor {
	s := &StepMotor{
		pins: [4]rpio.Pin{
			rpio.Pin(in1),
			rpio.Pin(in2),
			rpio.Pin(in3),
			rpio.Pin(in4),
		},
		chOp: make(chan Operator, 8),
	}
	for i := 0; i < 4; i++ {
		s.pins[i].Output()
		s.pins[i].Low()
	}
	go s.start()
	return s
}

// Start ...
func (s *StepMotor) start() {
	log.Printf("[%v]start working", logTagStepMotor)
	for op := range s.chOp {
		angle := int(op)
		s.roll(angle)
	}
}

func (s *StepMotor) roll(angle int) {
	var matrix [4][4]uint8
	if angle > 0 {
		matrix = clockwise
	} else {
		matrix = anticlockwise
		angle = angle * (-1)
	}
	n := int(float32(angle) / angleEachStep / 8.0)
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

// Roll ...
func (s *StepMotor) Roll(angle int) {
	s.chOp <- Operator(angle)
}

// Stop ...
func (s *StepMotor) Stop() {
	for i := 0; i < 4; i++ {
		s.pins[i].Low()
	}
}
