package dev

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
	Beep(n int, intervalInMS int)
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
	Display(text string)
	Open()
	Close()
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
	Blink(n int, intervalInMs int)
}

// Motor ...
type Motor interface {
	// Rolls roll angle dregee in clockwise direction if angle > 0,
	// or roll counter-clockwise if angle < 0
	Roll(angle float64)
	SetSpeed(speed int)
}

// MotorDriver ...
type MotorDriver interface {
	Forward()
	Backward()
	Left()
	Right()
	Stop()
	SetSpeed(s uint32)
}

// RFReciver is the interface of radio-frequency receiver
type RFReceiver interface {
	Received(ch int) bool
}

// Relay ...
type Relay interface {
	On(ch int)
	Off(ch int)
}

// StepperMotor ...
type StepperMotor interface {
	// Step gets the stepper motor rolls n steps in clockwise direction if angle > 0,
	// or roll in counter-clockwise direction if n < 0,
	// or motionless if n = 0.
	Step(n int)
	// Roll gets the stepper motor rolls angle dregee in clockwise direction if angle > 0,
	// or roll in counter-clockwise direction if angle < 0,
	// or motionless if angle = 0.
	Roll(angle float64)
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
