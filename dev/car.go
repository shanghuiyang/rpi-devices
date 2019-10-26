package dev

import (
	"log"
	"time"
)

const chSize = 8

const (
	forward      CarOp = "forward"
	backward     CarOp = "backward"
	left         CarOp = "left"
	right        CarOp = "right"
	stop         CarOp = "stop"
	honk         CarOp = "honk"
	blink        CarOp = "blink"
	camleft      CarOp = "camleft"
	camright     CarOp = "camright"
	camahead     CarOp = "camahead"
	lighton      CarOp = "lighton"
	lightoff     CarOp = "lightoff"
	selfdriveon  CarOp = "selfdriveon"
	selfdriveoff CarOp = "selfdriveoff"
)

type (
	// CarOp ...
	CarOp string
	// Option ...
	Option func(c *Car)
)

// WithEngine ...
func WithEngine(engine *L298N) Option {
	return func(c *Car) {
		c.engine = engine
	}
}

// WithDist ...
func WithDist(dist *HCSR04) Option {
	return func(c *Car) {
		c.dist = dist
	}
}

// WithHorn ...
func WithHorn(horn *Buzzer) Option {
	return func(c *Car) {
		c.horn = horn
	}
}

// WithLed ...
func WithLed(led *Led) Option {
	return func(c *Car) {
		c.led = led
	}
}

// WithLight ...
func WithLight(light *Led) Option {
	return func(c *Car) {
		c.light = light
	}
}

// WithCamera ...
func WithCamera(cam *Camera) Option {
	return func(c *Car) {
		c.camera = cam
	}
}

// Car ...
type Car struct {
	engine    *L298N
	dist      *HCSR04
	horn      *Buzzer
	led       *Led
	light     *Led
	camera    *Camera
	camAngle  int
	selfdrive bool
	chOp      chan CarOp
}

// NewCar ...
func NewCar(opts ...Option) *Car {
	car := &Car{
		camAngle:  0,
		selfdrive: false,
		chOp:      make(chan CarOp, chSize),
	}
	for _, opt := range opts {
		opt(car)
	}
	return car
}

// Start ...
func (c *Car) Start() error {
	go c.start()
	go c.camera.Turn(0)
	go c.led.Blink()
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
		case stop:
			c.stop()
		case honk:
			go c.honk()
		case camleft:
			go c.camLeft()
		case camright:
			go c.camRight()
		case camahead:
			go c.camAhead()
		case lighton:
			go c.light.On()
		case lightoff:
			go c.light.Off()
		case selfdriveon:
			go c.selfDriveOn()
		case selfdriveoff:
			go c.selfDriveOff()
		default:
			c.stop()
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
	c.delay(200)
	c.engine.Stop()
}

// right ...
func (c *Car) right() {
	log.Printf("car: right")
	c.engine.Right()
	c.delay(200)
	c.engine.Stop()
}

// stop ...
func (c *Car) stop() {
	log.Printf("car: stop")
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
			c.horn.Sound()
			c.delay(100)
		}
	}()
}

func (c *Car) camLeft() {
	angle := c.camAngle - 15
	if angle < -90 {
		angle = -90
	}
	c.camAngle = angle
	log.Printf("camera: %v", angle)
	c.camera.Turn(angle)
}

func (c *Car) camRight() {
	angle := c.camAngle + 15
	if angle > 90 {
		angle = 90
	}
	c.camAngle = angle
	log.Printf("camera: %v", angle)
	if c.camera == nil {
		return
	}
	c.camera.Turn(angle)
}

func (c *Car) camAhead() {
	c.camAngle = 0
	log.Printf("camera: %v", 0)
	if c.camera == nil {
		return
	}
	c.camera.Turn(0)
}

func (c *Car) selfDriveOn() {
	if c.dist == nil {
		return
	}
	// need to warm-up the distance sensor first
	c.dist.Dist()

	// make a warning before running into self-driving mode
	for i := 4; i >= 0 && c.horn != nil; i-- {
		log.Printf("self-drive: %v", i)
		c.horn.Sound()
		c.delay(1000)
	}
	// start self-drive
	c.selfdrive = true
	fwd := false
	for c.selfdrive {
		d := c.dist.Dist()
		log.Printf("dist: %.0f cm", d)

		// backward
		if d < 20 {
			c.chOp <- backward
			c.delay(200)
			fwd = false
			continue
		}
		// turn left
		if d < 40 {
			for i := 0; i < 10; i++ {
				c.chOp <- left
				c.delay(500)
			}
			fwd = false
			continue
		}
		// forward
		if !fwd {
			c.chOp <- forward
			fwd = true
		}
		c.delay(200)
	}
	c.chOp <- stop
}

func (c *Car) selfDriveOff() {
	c.selfdrive = false
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
