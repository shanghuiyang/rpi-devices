package dev

// ComMode ...
type ComMode int

const (
	// UartMode ...
	UartMode ComMode = iota
	// TTLMode ...
	TTLMode
)

// DistMeter ...
type DistMeter interface {
	Dist() float64
	Close()
}

// US100Config ...
type US100Config struct {
	Mode  ComMode
	Trig  int8
	Echo  int8
	Dev   string
	Baud  int
	Retry int
}
