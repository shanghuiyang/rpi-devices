package devices

// Operator ...
type Operator int

const (
	// Off will turn off the led
	Off Operator = iota
	// On will turn on the led
	On
	// Blink will make the led blinking
	Blink
)
