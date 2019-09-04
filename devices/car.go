package devices

import (
	"errors"
	"log"
	"time"
)

const (
	pinLed = 26
)

var (
	errNoImplement = errors.New("no implement")
)

// Car ...
type Car struct {
	engine *L298N
	horn   *Buzzer
	led    *Led
}

// NewCar ...
func NewCar() *Car {
	// eng := NewL298N(1, 2, 3, 4)
	// if eng == nil {
	// 	return nil
	// }

	// buzzer := NewBuzzer(1)
	// if buzzer == nil {
	// 	return nil
	// }

	// led := NewLed(pinLed)
	// if led == nil {
	// 	return nil
	// }
	return &Car{
		// engine: eng,
		// horn:   buzzer,
		// led: led,
	}
}

// Start ...
func (c *Car) Start() error {
	// go c.blink()
	return nil
}

// Stop ...
func (c *Car) Stop() error {
	// c.engine.Close()
	// c.horn.Close()
	// c.led.Close()
	return nil
}

// Forward ...
func (c *Car) Forward() error {
	log.Printf("car: forward")
	return nil

	c.engine.In1.High()
	c.engine.In2.Low()
	c.engine.In3.High()
	c.engine.In4.Low()
	return nil
}

// Backward ...
func (c *Car) Backward() error {
	log.Printf("car: backward")
	return nil

	c.engine.In1.Low()
	c.engine.In2.High()
	c.engine.In3.Low()
	c.engine.In4.High()
	return nil
}

// Left ...
func (c *Car) Left() error {
	log.Printf("car: left")
	return nil

	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.High()
	c.engine.In4.Low()
	return nil
}

// Right ...
func (c *Car) Right() error {
	log.Printf("car: right")
	return nil

	c.engine.In1.High()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	return nil
}

// Brake ...
func (c *Car) Brake() error {
	log.Printf("car: brake")
	return nil

	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	return nil
}

// Honk ...
func (c *Car) Honk() error {
	log.Printf("car: honk")
	return nil

	go func() {
		for i := 0; i < 3; i++ {
			c.horn.Whistle()
			time.Sleep(300 * time.Millisecond)
		}
	}()
	return nil
}

func (c *Car) blink() {
	for {
		c.led.On()
		time.Sleep(1 * time.Second)
		c.led.Off()
		time.Sleep(1 * time.Second)
	}
}

// Blink ...
// func (c *Car) Blink() {
// 	for i := 0; i < 3; i++ {
// 		c.led.On()
// 		time.Sleep(1 * time.Second)
// 		c.led.Off()
// 		time.Sleep(1 * time.Second)
// 	}
// }
