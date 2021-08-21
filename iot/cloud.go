package iot

// Cloud ...
type Cloud interface {
	Push(v *Value) error
}

// Value ...
type Value struct {
	Device string
	Value  interface{}
}
