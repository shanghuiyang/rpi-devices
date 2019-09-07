package devices

import (
	"errors"
	"log"
	"time"
)

const (
	pinLed    = 26
	pinIn1    = 17
	pinIn2    = 18
	pinIn3    = 27
	pinIn4    = 22
	pinBuzzer = 10
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
	eng := NewL298N(pinIn1, pinIn2, pinIn3, pinIn4)
	if eng == nil {
		return nil
	}

	buzzer := NewBuzzer(pinBuzzer)
	if buzzer == nil {
		return nil
	}

	led := NewLed(pinLed)
	if led == nil {
		return nil
	}
	return &Car{
		engine: eng,
		horn:   buzzer,
		led:    led,
	}
}

// Start ...
func (c *Car) Start() error {
	go c.blink()
	return nil
}

// Stop ...
func (c *Car) Stop() error {
	c.engine.Close()
	// c.horn.Close()
	// c.led.Close()
	return nil
}

// Forward ...
func (c *Car) Forward() error {
	log.Printf("car: forward")
	c.engine.In1.High()
	c.engine.In2.Low()
	time.Sleep(50 * time.Millisecond)
	c.engine.In3.High()
	c.engine.In4.Low()
	return nil
}

// Backward ...
func (c *Car) Backward() error {
	log.Printf("car: backward")
	c.engine.In1.Low()
	c.engine.In2.High()
	time.Sleep(50 * time.Millisecond)
	c.engine.In3.Low()
	c.engine.In4.High()
	return nil
}

// Left ...
func (c *Car) Left() error {
	log.Printf("car: left")
	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.High()
	c.engine.In4.Low()
	return nil
}

// Right ...
func (c *Car) Right() error {
	log.Printf("car: right")
	c.engine.In1.High()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	return nil
}

// Brake ...
func (c *Car) Brake() error {
	log.Printf("car: brake")
	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	return nil
}

// Honk ...
func (c *Car) Honk() error {
	log.Printf("car: honk")
	go func() {
		for i := 0; i < 5; i++ {
			c.horn.Whistle()
			time.Sleep(100 * time.Millisecond)
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
