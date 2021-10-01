package iot

type Noop struct{}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) Push(v *Value) error {
	return nil
}

func (n *Noop) Get(params map[string]interface{}) ([]byte, error) {
	return []byte{}, nil
}
