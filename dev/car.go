package dev

import (
	"log"
	"time"
)

const chSize = 8

const (
	forward  CarOp = "forward"
	backward CarOp = "backward"
	left     CarOp = "left"
	right    CarOp = "right"
	stop     CarOp = "stop"
	turn     CarOp = "turn"
	scan     CarOp = "scan"

	beep  CarOp = "beep"
	blink CarOp = "blink"

	servoleft  CarOp = "servoleft"
	servoright CarOp = "servoright"
	servoahead CarOp = "servoahead"

	lighton  CarOp = "lighton"
	lightoff CarOp = "lightoff"

	selfdrivingon  CarOp = "selfdrivingon"
	selfdrivingoff CarOp = "selfdrivingoff"
)

var (
	scanningAngles = []int{-90, -75, -60, -45, -30, 30, 45, 60, 75, 90}

	// the map between angle(degree) and time(millisecond)
	turnAngleTimes = map[int]int{
		-90: 1250,
		-75: 1000,
		-60: 800,
		-45: 600,
		-30: 400,
		30:  400,
		45:  600,
		60:  800,
		75:  1000,
		90:  1250,
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
func WithDist(dist *US100) Option {
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
	dist        *US100
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
			go c.onSelfDriving()
		case selfdrivingoff:
			go c.offSelfDriving()
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

/*

                                                                          +-----------------------------------------------+
                                                                          |                                               |
                                                                          v                                               |Y
+-------+     +---------+    +---------------+     +-----------+     +----+-----+      +------+      +------+     +--------------+
| start |---->| forward |--->| objects ahead |---->| distance  |---->| backword |----->| stop |----->| scan |---->| min distance |
+-------+     +-----+---+    |   detected?   | Y   |  < 10cm?  | Y   +----------+      +--+---+      +------+     |    < 10cm    |
                    ^        +-------+-------+     +-----+-----+                          ^                       +--------------+
                    |                |                   |                                |                               |N
                    |                |N                 N|                                |                               |
                    |                |                   |                                |                               v
                    |                v                   |           +----------+ Y       |   Y +----------+   Y  +-------+------+
                    |                |                   +---------->| distance +------>--+-----| retry<4? |------| max distance |
                    |                |                               |  < 40cm? |               +----+-----+      |    < 40cm    |
                    ^                |                               +----------+                    | N          +--------------+
                    |                |                                     |N                        v                    |N
                    |                |                                     |                         |                    |
                    |                |                                     v                    +----+-----+              |
                    +-------<--------+------------------<------------------+---------<----------|   turn   |-------<------+
                                                                                                +----------+



*/
func (c *Car) onSelfDriving() {
	log.Printf("car: self-drving")
	if c.dist == nil {
		log.Printf("can't self-driving without the distance sensor")
		return
	}

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	// start self-driving
	c.selfdriving = true
	var (
		op       = forward
		fwd      bool
		retry    int
		angle    int
		min, max float64
	)

	chOp := make(chan CarOp, 1)
	chDetecting := make(chan bool)
	go c.detecting(chOp, chDetecting)
	for c.selfdriving {
		select {
		case p := <-chOp:
			op = p
		default:
			op = forward
		}
		log.Printf("op: %v", op)

		switch op {
		case backward:
			fwd = false
			c.stop()
			c.delay(20)
			c.backward()
			c.delay(500)
			chOp <- stop
			continue
		case stop:
			fwd = false
			c.stop()
			c.delay(20)
			chOp <- scan
			continue
		case scan:
			fwd = false
			min, max, angle = c.scanDist()
			log.Printf("mind=%.0f, maxd=%.0f, angle=%v", min, max, angle)
			if min < 10 && retry < 4 {
				chOp <- backward
				retry++
				continue
			}
			chOp <- turn
			retry = 0
		case turn:
			fwd = false
			c.turn(angle)
			chDetecting <- true // resume to detecting objects ahead
			c.delay(150)
			continue
		case forward:
			if !fwd {
				c.forward()
				fwd = true
			}
			c.delay(50)
			continue
		}
	}
	c.stop()
}

func (c *Car) offSelfDriving() {
	c.selfdriving = false
	log.Printf("car: self-drving off")
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// detects the objects ahead(-30 ~ 30 degree).
func (c *Car) detecting(chOp chan CarOp, chDetecting chan bool) {
	angles := []int{-30, -15, 0, 15, 30, 15, 0, -15}
	for c.selfdriving {
		for _, angle := range angles {
			c.servo.Roll(angle)
			c.delay(50)
			d := c.dist.Dist()
			if d < 10 {
				chOp <- backward
				<-chDetecting // pause detecting until the car finishs the actions
				break
			}
			if d < 40 {
				chOp <- stop
				<-chDetecting // pause detecting until the car finishs the actions
				break
			}
		}
	}
}

// scan for geting the min & max distance, and the angle for max distance
func (c *Car) scanDist() (min, max float64, angle int) {
	min = 9999
	max = -9999
	for _, ang := range scanningAngles {
		c.servo.Roll(ang)
		c.delay(50)
		d := c.dist.Dist()
		log.Printf("scan: angle %v, dist: %.0f", ang, d)
		if d < min {
			min = d
		}
		if d > max {
			max = d
			angle = ang
		}
	}
	c.servo.Roll(0)
	return
}

func (c *Car) turn(angle int) {
	ms, ok := turnAngleTimes[angle]
	if !ok {
		return
	}
	if angle < 0 {
		c.engine.Left()
	} else {
		c.engine.Right()
	}
	c.delay(ms)
	c.stop()
	return
}
