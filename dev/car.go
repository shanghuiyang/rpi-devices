package dev

import (
	"log"
	"time"
)

const chSize = 8

const (
	forward        CarOp = "forward"
	backward       CarOp = "backward"
	left           CarOp = "left"
	right          CarOp = "right"
	stop           CarOp = "stop"
	beep           CarOp = "beep"
	blink          CarOp = "blink"
	servoleft      CarOp = "servoleft"
	servoright     CarOp = "servoright"
	servoahead     CarOp = "servoahead"
	lighton        CarOp = "lighton"
	lightoff       CarOp = "lightoff"
	selfdrivingon  CarOp = "selfdrivingon"
	selfdrivingoff CarOp = "selfdrivingoff"
)

var (
	scanningAngles = []int{-90, -75, -60, -45, -30, -15, 15, 30, 45, 60, 75, 90}

	turnningAngles = map[int]int{
		-90: 7,
		-75: 6,
		-60: 5,
		-45: 4,
		-30: 3,
		-15: 2,
		15:  2,
		30:  3,
		45:  4,
		60:  5,
		75:  6,
		90:  7,
	}
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

// WithServo ...
func WithServo(servo *SG90) Option {
	return func(c *Car) {
		c.servo = servo
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
	engine      *L298N
	servo       *SG90
	dist        *HCSR04
	horn        *Buzzer
	led         *Led
	light       *Led
	camera      *Camera
	servoAngle  int
	selfdriving bool
	chOp        chan CarOp
}

// NewCar ...
func NewCar(opts ...Option) *Car {
	car := &Car{
		servoAngle:  0,
		selfdriving: false,
		chOp:        make(chan CarOp, chSize),
	}
	for _, opt := range opts {
		opt(car)
	}
	return car
}

// Start ...
func (c *Car) Start() error {
	go c.start()
	go c.servo.Roll(0)
	go c.blink()
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
		case beep:
			go c.beep()
		case servoleft:
			go c.servoLeft()
		case servoright:
			go c.servoRight()
		case servoahead:
			go c.servoAhead()
		case lighton:
			go c.light.On()
		case lightoff:
			go c.light.Off()
		case selfdrivingon:
			go c.selfDrivingOn()
		case selfdrivingoff:
			go c.selfDrivingOff()
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

func (c *Car) speed(s uint32) {
	c.engine.Speed(s)
}

// beep ...
func (c *Car) beep() {
	log.Printf("car: beep")
	if c.horn == nil {
		return
	}
	go c.horn.Beep(5, 100)
}

func (c *Car) blink() {
	for {
		c.led.Blink(1, 1000)
	}
}

func (c *Car) servoLeft() {
	angle := c.servoAngle - 15
	if angle < -90 {
		angle = -90
	}
	c.servoAngle = angle
	log.Printf("servo: %v", angle)
	if c.servo == nil {
		return
	}
	c.servo.Roll(angle)
}

func (c *Car) servoRight() {
	angle := c.servoAngle + 15
	if angle > 90 {
		angle = 90
	}
	c.servoAngle = angle
	log.Printf("servo: %v", angle)
	if c.servo == nil {
		return
	}
	c.servo.Roll(angle)
}

func (c *Car) servoAhead() {
	c.servoAngle = 0
	log.Printf("servo: %v", 0)
	if c.servo == nil {
		return
	}
	c.servo.Roll(0)
}

func (c *Car) selfDrivingOn() {
	if c.dist == nil {
		return
	}
	// need to warm-up the distance sensor first
	c.dist.Dist()

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	// start self-driving
	c.selfdriving = true
	fwd := false
	for c.selfdriving {
		d := c.dist.Dist()
		log.Printf("dist: %.0f cm", d)

		// find a way out
		if d < 40 {
			fwd = false
			c.stop()
			c.delay(100)
			// backward
			if d < 10 {
				c.backward()
				c.delay(500)
				c.stop()
			}
			maxd, angle := c.scan()
			log.Printf("maxd=%.0f, angle=%v", maxd, angle)
			retry := 3
			i := 0
			for ; i < retry && maxd < 40; i++ {
				c.backward()
				c.delay(300)
				c.stop()
				maxd, angle = c.scan()
			}
			if i == retry {
				// out of self-driving mode
				go c.horn.Beep(60, 300)
				c.selfdriving = false
				break
			}
			c.turn(angle)
			continue
		}
		// forward
		if !fwd {
			c.forward()
			fwd = true
		}
		c.delay(150)
	}
	c.stop()
}

func (c *Car) selfDrivingOff() {
	c.selfdriving = false
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (c *Car) scan() (maxDist float64, angle int) {
	for _, ang := range scanningAngles {
		c.servo.Roll(ang)
		c.delay(50)
		d := c.dist.Dist()
		log.Printf("scan: angle %v, dist: %.0f", ang, d)
		if d > maxDist {
			maxDist = d
			angle = ang
		}
	}
	c.servo.Roll(0)
	return
}

func (c *Car) turn(angle int) {
	n, ok := turnningAngles[angle]
	if !ok {
		n = angle*2/45 + 2
		if angle < 0 {
			n *= -1
		}
	}
	for i := 0; i < n; i++ {
		if angle < 0 {
			c.left()
		} else {
			c.right()
		}
		c.delay(50)
	}
	c.stop()
	return
}
