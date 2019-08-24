package devices

import (
	"log"

	"github.com/stianeikeland/go-rpio"
)

const (
	logTagRelay = "relay"
)

var (
	// ChRelayOp ...
	ChRelayOp = make(chan Operator)
)

// Relay ...
type Relay struct {
	pin  rpio.Pin
	isOn bool
}

// NewRelay ...
func NewRelay(pin uint8) *Relay {
	if err := rpio.Open(); err != nil {
		return nil
	}
	r := &Relay{
		pin:  rpio.Pin(pin),
		isOn: false,
	}
	r.pin.Output()
	return r
}

// Start ...
func (r *Relay) Start() {
	defer r.Close()

	log.Printf("[%v]start working", logTagRelay)
	for {
		op := <-ChRelayOp
		switch op {
		case Off:
			r.Off()
		case On:
			r.On()
		default:
			// do nothing
		}
	}
}

// On ...
func (r *Relay) On() {
	if !r.isOn {
		r.pin.High()
		r.isOn = true
	}
}

// Off ...
func (r *Relay) Off() {
	if r.isOn {
		r.pin.Low()
		r.isOn = false
	}
}

// Close ...
func (r *Relay) Close() {
	rpio.Close()
}
