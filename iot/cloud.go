package iot

// Cloud ...
type Cloud interface {
	Push(v *Value) error
	Get(params map[string]interface{}) ([]byte, error)
}

// Value ...
type Value struct {
	Device string
	Value  interface{}
}
