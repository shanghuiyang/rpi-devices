package dev

type DCMotor struct {
	driver MotorDriver
}

// NewDCMotor ...
func NewDCMotor(driver MotorDriver) *DCMotor {
	return &DCMotor{
		driver: driver,
	}
}

// Forward ...
func (m *DCMotor) Forward() {
	m.driver.Forward()
}

// Backward ...
func (m *DCMotor) Backward() {
	m.driver.Backward()
}

// Stop ...
func (m *DCMotor) Stop() {
	m.driver.Stop()
}

// Speed ...
func (m *DCMotor) SetSpeed(percent uint32) {
	m.driver.SetSpeed(percent)
}
