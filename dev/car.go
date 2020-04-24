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

	turnAngleCounts = map[int]int{
		-90: 20,
		-75: 17,
		-60: 14,
		-45: 10,
		-30: 7,
		30:  5,
		45:  8,
		60:  10,
		75:  13,
		90:  17,
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

// WithUlt ...
func WithUlt(ult *US100) Option {
	return func(c *Car) {
		c.ult = ult
	}
}

// WithEncoder ...
func WithEncoder(e *Encoder) Option {
	return func(c *Car) {
		c.encoder = e
	}
}

// WithCSwitchs ...
func WithCSwitchs(cswitchs []*CollisionSwitch) Option {
	return func(c *Car) {
		c.cswitchs = cswitchs
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
	ult         *US100
	encoder     *Encoder
	cswitchs    []*CollisionSwitch
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
	c.encoder.Close()
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
			go c.lightOn()
		case lightoff:
			go c.lightOff()
		case selfdrivingon:
			go c.selfDrivingOn()
		case selfdrivingoff:
			go c.selfDrivingOff()
		default:
			log.Printf("car: invalid op")
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
	c.horn.Beep(5, 100)
}

func (c *Car) blink() {
	for {
		c.led.Blink(1, 1000)
	}
}

func (c *Car) lightOn() {
	log.Printf("car: light on")
	if c.light == nil {
		return
	}
	c.light.On()
}

func (c *Car) lightOff() {
	log.Printf("car: light off")
	if c.light == nil {
		return
	}
	c.light.Off()
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
| start |---->| forward |--->|   obstacles   |---->| distance  |---->| backword |----->| stop |----->| scan |---->| min distance |
+-------+     +-----+---+    |   detected?   | Y   |  < 10cm?  | Y   +----------+      +--+---+      +------+     |    < 10cm    |
                    ^        +-------+-------+     +-----+-----+                          |                       +--------------+
                    |                |                   |                                ^                               |N
                    |                |N                 N|                                |                               |
                    |                |                   |                                |                               v
                    |                v                   |           +----------+ Y       |   Y +----------+   Y  +-------+------+
                    |                |                   +---------->| distance +------>--+-<---| retry<4? |-<----| max distance |
                    |                |                               |  < 40cm? |               +----+-----+      |    < 40cm    |
                    ^                |                               +----------+                    | N          +--------------+
                    |                |                                     |N                        v                    |N
                    |                |                                     |                         |                    |
                    |                |                                     v                    +----+-----+              |
                    +-------<--------+------------------<------------------+---------<----------|   turn   |-------<------+
                                                                                                +----------+



*/
func (c *Car) selfDrivingOn() {
	log.Printf("car: self-drving on")
	if c.ult == nil {
		log.Printf("can't self-driving without the distance sensor")
		return
	}

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	// start self-driving
	c.selfdriving = true
	var (
		op                   = forward
		fwd                  bool
		retry                int
		mindAngle, maxdAngle int
		mind, maxd           float64
	)

	// collided := false
	chOp := make(chan CarOp, 1)
	chDetectCollision := make(chan bool)
	chDetectObstacle := make(chan bool)
	chPauseDetecting := make(chan bool)
	go c.detectCollision(chOp, chDetectCollision, chPauseDetecting)
	go c.detectObstacle(chOp, chDetectObstacle, chPauseDetecting)
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
			mind, maxd, mindAngle, maxdAngle = c.scanDist()
			log.Printf("mind=%.0f, maxd=%.0f, mindAngle=%v, maxdAngle=%v", mind, maxd, mindAngle, maxdAngle)
			if mind < 10 && mindAngle != 90 && mindAngle != -90 && retry < 4 {
				chOp <- backward
				retry++
				continue
			}
			chOp <- turn
			retry = 0
		case turn:
			fwd = false
			c.turn(maxdAngle)
			chDetectCollision <- true // resume to detecting collision
			chDetectObstacle <- true  // resume to detecting obstacles
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
	close(chOp)
	close(chPauseDetecting)
	close(chDetectCollision)
	close(chDetectObstacle)
}

func (c *Car) selfDrivingOff() {
	c.selfdriving = false
	log.Printf("car: self-drving off")
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// detect detects obstacles using an ultrasonic distance meter.
func (c *Car) detectObstacle(chOp chan CarOp, chDetect, chPause chan bool) {
	angles := []int{0, -10, -20, -10, 0, 10, 20, 10}
	pause := func() {
		// notify other sensors to pause detecting
		chPause <- true
	}
	for c.selfdriving {
		for _, angle := range angles {
			c.servo.Roll(angle)
			c.delay(50)
			d := c.ult.Dist()

			select {
			case <-chPause:
				log.Printf("ult: pause")
				<-chDetect
				break
			default:
				// do nothing
			}

			if d < 10 {
				chOp <- backward
				go pause()
				<-chDetect // pause detecting until the car finishs the actions
				break
			}
			if d < 50 {
				chOp <- stop
				go pause()
				<-chDetect // pause detecting until the car finishs the actions
				break
			}
		}
	}
	c.servo.Roll(0)
}

// detect detects obstacles using an ultrasonic distance meter.
func (c *Car) detectCollision(chOp chan CarOp, chDetect, chPause chan bool) {
	pause := func() {
		// notify other sensors to pause detecting
		chPause <- true
	}
	for c.selfdriving {
		select {
		case <-chPause:
			log.Printf("cswitch: pause")
			<-chDetect
			break
		default:
			// do nothing
		}
		for _, cswitch := range c.cswitchs {
			if cswitch.Collided() {
				chOp <- backward
				go pause()
				go c.horn.Beep(1, 100)
				log.Printf("cswitch: got a crash")
				<-chDetect // pause detecting collision until the car finishs the actions
			}
		}

		c.delay(10)
	}
}

// scan for geting the min & max distance, and their corresponding angles
// mind: the min distance
// maxd: the max distance
// mindAngle: the angle correspond to the mind
// maxdAngle: the angle correspond to the maxd
func (c *Car) scanDist() (mind, maxd float64, mindAngle, maxdAngle int) {
	mind = 9999
	maxd = -9999
	for _, ang := range scanningAngles {
		c.servo.Roll(ang)
		c.delay(50)
		d := c.ult.Dist()
		for i := 0; d < 0 && i < 3; i++ {
			c.delay(50)
			d = c.ult.Dist()
		}
		if d < 0 {
			continue
		}
		log.Printf("scan: angle %v, dist: %.0f", ang, d)
		if d < mind {
			mind = d
			mindAngle = ang
		}
		if d > maxd {
			maxd = d
			maxdAngle = ang
		}
	}
	c.servo.Roll(0)
	c.delay(50)
	return
}

func (c *Car) turn(angle int) {
	n, ok := turnAngleCounts[angle]
	if !ok {
		log.Printf("invalid angle: %d", angle)
		return
	}
	if angle < 0 {
		c.engine.Left()
	} else {
		c.engine.Right()
	}
	for i := 0; i < n; {
		i += c.encoder.Count1()
	}
	c.stop()
	return
}
