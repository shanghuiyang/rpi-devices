package devices

import (
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

const (
	chSize         = 8
	forward  CarOp = "forward"
	backward CarOp = "backward"
	left     CarOp = "left"
	right    CarOp = "right"
	brake    CarOp = "brake"
	honk     CarOp = "honk"
)

// CarOp ...
type CarOp string

// Car ...
type Car struct {
	engine *L298N
	horn   *Buzzer
	led    *Led
	chOp   chan CarOp
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
		chOp:   make(chan CarOp, chSize),
	}
}

// Start ...
func (c *Car) Start() error {
	go c.blink()
	go c.start()
	return nil
}

// Do ...
func (c *Car) Do(op CarOp) {
	c.chOp <- op
}

// Stop ...
func (c *Car) Stop() error {
	close(c.chOp)
	c.engine.Close()
	return nil
}

func (c *Car) start() {
	for op := range c.chOp {
		switch op {
		case forward:
			c.forward()
		case backward:
			c.backward()
		case left:
			c.left()
		case right:
			c.right()
		case brake:
			c.brake()
		case honk:
			c.honk()
		default:
			c.brake()
		}
	}
}

// Forward ...
func (c *Car) forward() error {
	log.Printf("car: forward")
	c.engine.In1.High()
	c.engine.In2.Low()
	time.Sleep(70 * time.Millisecond)
	c.engine.In3.High()
	c.engine.In4.Low()

	c.engine.In1.Low()
	time.Sleep(70 * time.Millisecond)
	c.engine.In1.High()
	return nil
}

// Backward ...
func (c *Car) backward() error {
	log.Printf("car: backward")
	c.engine.In1.Low()
	c.engine.In2.High()
	time.Sleep(70 * time.Millisecond)
	c.engine.In3.Low()
	c.engine.In4.High()

	c.engine.In2.Low()
	time.Sleep(70 * time.Millisecond)
	c.engine.In2.High()
	return nil
}

// Left ...
func (c *Car) left() error {
	log.Printf("car: left")
	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.High()
	c.engine.In4.Low()
	time.Sleep(70 * time.Millisecond)
	c.brake()
	return nil
}

// Right ...
func (c *Car) right() error {
	log.Printf("car: right")
	c.engine.In1.High()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	time.Sleep(70 * time.Millisecond)
	c.brake()
	return nil
}

// Brake ...
func (c *Car) brake() error {
	log.Printf("car: brake")
	c.engine.In1.Low()
	c.engine.In2.Low()
	c.engine.In3.Low()
	c.engine.In4.Low()
	return nil
}

// Honk ...
func (c *Car) honk() error {
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
