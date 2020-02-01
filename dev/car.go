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
	rudderleft   CarOp = "rudderleft"
	rudderright  CarOp = "rudderright"
	rudderahead  CarOp = "rudderahead"
	lighton      CarOp = "lighton"
	lightoff     CarOp = "lightoff"
	selfdriveon  CarOp = "selfdriveon"
	selfdriveoff CarOp = "selfdriveoff"
)

var (
	scanningAngles = []int{-90, -45, 45, 90}

	turnningAngles = map[int]int{
		-90: 6,
		-45: 4,
		45:  4,
		90:  6,
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

// WithRudder ...
func WithRudder(rudder *SG90) Option {
	return func(c *Car) {
		c.rudder = rudder
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
	rudder      *SG90
	dist        *HCSR04
	horn        *Buzzer
	led         *Led
	light       *Led
	camera      *Camera
	rudderAngle int
	selfdrive   bool
	chOp        chan CarOp
}

// NewCar ...
func NewCar(opts ...Option) *Car {
	car := &Car{
		rudderAngle: 0,
		selfdrive:   false,
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
	go c.rudder.Roll(0)
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
		case honk:
			go c.honk()
		case rudderleft:
			go c.rudderLeft()
		case rudderright:
			go c.rudderRight()
		case rudderahead:
			go c.rudderAhead()
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

func (c *Car) speed(s uint32) {
	c.engine.Speed(s)
}

// honk ...
func (c *Car) honk() {
	log.Printf("car: honk")
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

func (c *Car) rudderLeft() {
	angle := c.rudderAngle - 15
	if angle < -90 {
		angle = -90
	}
	c.rudderAngle = angle
	log.Printf("rudder: %v", angle)
	if c.rudder == nil {
		return
	}
	c.rudder.Roll(angle)
}

func (c *Car) rudderRight() {
	angle := c.rudderAngle + 15
	if angle > 90 {
		angle = 90
	}
	c.rudderAngle = angle
	log.Printf("rudder: %v", angle)
	if c.rudder == nil {
		return
	}
	c.rudder.Roll(angle)
}

func (c *Car) rudderAhead() {
	c.rudderAngle = 0
	log.Printf("rudder: %v", 0)
	if c.rudder == nil {
		return
	}
	c.rudder.Roll(0)
}

func (c *Car) selfDriveOn() {
	if c.dist == nil {
		return
	}
	// need to warm-up the distance sensor first
	c.dist.Dist()

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	// start self-drive
	c.selfdrive = true
	fwd := false
	for c.selfdrive {
		d := c.dist.Dist()
		log.Printf("dist: %.0f cm", d)

		// backward
		if d < 10 {
			fwd = false
			c.backward()
			c.delay(300)
			c.stop()
			continue
		}
		// find a way out
		if d < 40 {
			fwd = false
			c.stop()
			c.delay(500)
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
				c.selfdrive = false
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
		c.delay(200)
	}
	c.stop()
}

func (c *Car) selfDriveOff() {
	c.selfdrive = false
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (c *Car) scan() (maxDist float64, angle int) {
	for _, ang := range scanningAngles {
		c.rudder.Roll(ang)
		c.delay(100)
		d := c.dist.Dist()
		log.Printf("scan: angle %v, dist: %.0f", ang, d)
		if d > maxDist {
			maxDist = d
			angle = ang
		}
	}
	c.rudder.Roll(0)
	return
}

func (c *Car) turn(angle int) {
	c.speed(30)
	c.delay(200)
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
		c.delay(100)
	}
	c.stop()
	c.speed(25)
	c.delay(200)
	return
}
