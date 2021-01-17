package iot

// Cloud is the interface of IOT clound
type Cloud interface {
	Push(v *Value) error
}

// Value ...
type Value struct {
	Device string
	Value  interface{}
}

// NewCloud ...
func NewCloud(config interface{}) Cloud {
	var cloud Cloud
	switch config.(type) {
	case *WsnConfig:
		cfg := config.(*WsnConfig)
		cloud = NewWsnClound(cfg)
	case *OneNetConfig:
		cfg := config.(*OneNetConfig)
		cloud = NewOneNetCloud(cfg)
	default:
		cloud = nil
	}
	return cloud
}
