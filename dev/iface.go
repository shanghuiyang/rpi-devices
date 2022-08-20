package dev

import "image"

// Accelerometer ...
type Accelerometer interface {
	Angles() (yaw, pitch, roll float64, err error)
	Close() error
}

// ADC is the interface of Analog DigitalC onverter
type ADC interface {
	Read(channel int) (float64, error)
	Close() error
}

// Button ...
type Button interface {
	Pressed() bool
}

// Buzzer ...
type Buzzer interface {
	On()
	Off()
	Beep(n int, intervalMs int)
}

// camera ...
type Camera interface {
	Photo() ([]byte, error)
}

// CH2OMeter ...
type CH2OMeter interface {
	Value() (float64, error)
	Close() error
}

// Detector ...
type Detector interface {
	Detected() bool
}

// Display ...
type Display interface {
	Image(img image.Image) error
	Text(text string, x, y int) error
	On() error
	Off() error
	Clear() error
	Close() error
}

// DistanceMeter ...
type DistanceMeter interface {
	Dist() (float64, error)
	Close() error
}

// Encoder ...
type Encoder interface {
	Count1() int
	Detected() bool
	Start()
	Stop()
}

// GPS ...
type GPS interface {
	Loc() (lat, lon float64, err error)
	Close() error
}

// Hygrometer ...
type Hygrometer interface {
	Humidity() (float32, error)
}

// Joystick ...
type Joystick interface {
	X() float64
	Y() float64
	Z() int
}

// Led ...
type Led interface {
	On()
	Off()
	Blink(n int, intervalMs int)
}

// Motor ...
type Motor interface {
	Forward()
	Backward()
	Stop()
	SetSpeed(percent uint32)
}

// MotorDriver ...
type MotorDriver interface {
	Motor
}

// Pump ...
type Pump interface {
	On()
	Off()
	Run(sec int)
}

// RFReciver is the interface of radio-frequency receiver
type RFReceiver interface {
	Received(ch int) bool
}

// Relay ...
type Relay interface {
	On()
	Off()
}

// ServoMotor ...
type ServoMotor interface {
	// Roll gets the servo motor rolls angle dregee in clockwise direction if angle > 0,
	// or roll in counter-clockwise direction if angle < 0,
	// motionless if angle = 0.
	Roll(angle float64)
}

// Stepper ...
type StepperMotor interface {
	// Step gets the stepper motor rolls n steps in clockwise direction if angle > 0,
	// or roll in counter-clockwise direction if n < 0,
	// or motionless if n = 0.
	Step(n int)
	// Roll gets the stepper motor rolls angle dregee in clockwise direction if angle > 0,
	// or roll in counter-clockwise direction if angle < 0,
	// motionless if angle = 0.
	Roll(angle float64)
	// SetMode sets sets the stepping mode.
	// For example: Full, Half, Quarter, Eighth, Sixteenth.
	// Please NOTE not all steppers support all modes. Some steppers only support one or two modes.
	SetMode(mode StepperMode) error
}

// Thermometer ...
type Thermometer interface {
	Temperature() (float64, error)
}

// Thermohygrometer ...
type Thermohygrometer interface {
	TempHumidity() (temp, humi float64, err error)
}

// Wireless ...
type Wireless interface {
	Send(data []byte) error
	Receive() ([]byte, error)
	Sleep()
	Wakeup()
	Close() error
}
