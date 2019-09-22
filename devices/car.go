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
	blink    CarOp = "blink"
)

// CarOp ...
type CarOp string

// Engine ...
type Engine interface {
	Forward()
	Backward()
	Left()
	Right()
	Stop()
}

// Horn ...
type Horn interface {
	Whistle()
}

// Light ...
type Light interface {
	On()
	Off()
}

// CarBuilder ...
type CarBuilder struct {
	engine Engine
	horn   Horn
	light  Light
}

// NewCarBuilder ...
func NewCarBuilder() *CarBuilder {
	return &CarBuilder{}
}

// Engine ...
func (b *CarBuilder) Engine(eng Engine) *CarBuilder {
	b.engine = eng
	return b
}

// Horn ...
func (b *CarBuilder) Horn(horn Horn) *CarBuilder {
	b.horn = horn
	return b
}

// Light ...
func (b *CarBuilder) Light(light Light) *CarBuilder {
	b.light = light
	return b
}

// Build ...
func (b *CarBuilder) Build() *Car {
	return &Car{
		engine: b.engine,
		horn:   b.horn,
		light:  b.light,
		chOp:   make(chan CarOp, chSize),
	}
}

// Car ...
type Car struct {
	engine Engine
	horn   Horn
	light  Light
	chOp   chan CarOp
}

// Start ...
func (c *Car) Start() error {
	go c.start()
	c.chOp <- blink
	return nil
}

// Do ...
func (c *Car) Do(op CarOp) {
	c.chOp <- op
}

// Stop ...
func (c *Car) Stop() error {
	close(c.chOp)
	c.engine.Stop()
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
			go c.honk()
		case blink:
			go c.blink()
		default:
			c.brake()
		}
	}
}

// forward ...
func (c *Car) forward() {
	log.Printf("car: forward")
	c.engine.Forward()
}

// backward ...
func (c *Car) backward() {
	log.Printf("car: backward")
	c.engine.Backward()
}

// left ...
func (c *Car) left() {
	log.Printf("car: left")
	c.engine.Left()
	time.Sleep(70 * time.Millisecond)
	c.engine.Stop()
}

// right ...
func (c *Car) right() {
	log.Printf("car: right")
	c.engine.Right()
	time.Sleep(70 * time.Millisecond)
	c.engine.Stop()
}

// brake ...
func (c *Car) brake() {
	log.Printf("car: brake")
	c.engine.Stop()
}

// honk ...
func (c *Car) honk() {
	log.Printf("car: honk")
	if c.horn == nil {
		return
	}
	go func() {
		for i := 0; i < 5; i++ {
			c.horn.Whistle()
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// blink ...
func (c *Car) blink() {
	if c.light == nil {
		return
	}
	for {
		c.light.On()
		time.Sleep(1 * time.Second)
		c.light.Off()
		time.Sleep(1 * time.Second)
	}
}
