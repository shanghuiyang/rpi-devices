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
