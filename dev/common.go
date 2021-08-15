package dev

const (
	// voice speed in cm/s
	voiceSpeed = 34000.0
)

// ComMode ...
type ComMode int

const (
	// UartMode ...
	UartMode ComMode = iota
	// TTLMode ...
	TTLMode
)

// US100Config ...
type US100Config struct {
	Mode  ComMode
	Trig  int8
	Echo  int8
	Dev   string
	Baud  int
	Retry int
}
