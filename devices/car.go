package devices

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
	autodriveon  CarOp = "autodriveon"
	autodriveoff CarOp = "autodriveoff"
)

// CarOp ...
type CarOp string

// IEngine ...
type IEngine interface {
	Forward()
	Backward()
	Left()
	Right()
	Stop()
}

// IHorn ...
type IHorn interface {
	Whistle()
}

// ILed ...
type ILed interface {
	Blink()
}

// ILight ...
type ILight interface {
	On()
	Off()
}

// ICamera ...
type ICamera interface {
	Turn(angle int)
}

// IDistance ...
type IDistance interface {
	Dist() float64
}

// CarBuilder ...
type CarBuilder struct {
	engine IEngine
	dist   IDistance
	camera ICamera
	horn   IHorn
	led    ILed
	light  ILight
}

// NewCarBuilder ...
func NewCarBuilder() *CarBuilder {
	return &CarBuilder{}
}

// Engine ...
func (b *CarBuilder) Engine(eng IEngine) *CarBuilder {
	b.engine = eng
	return b
}

// Distance ...
func (b *CarBuilder) Distance(dist IDistance) *CarBuilder {
	b.dist = dist
	return b
}

// Horn ...
func (b *CarBuilder) Horn(horn IHorn) *CarBuilder {
	b.horn = horn
	return b
}

// Led ...
func (b *CarBuilder) Led(led ILed) *CarBuilder {
	b.led = led
	return b
}

// Light ...
func (b *CarBuilder) Light(light ILight) *CarBuilder {
	b.light = light
	return b
}

// Camera ...
func (b *CarBuilder) Camera(camera ICamera) *CarBuilder {
	b.camera = camera
	return b
}

// Build ...
func (b *CarBuilder) Build() *Car {
	return &Car{
		engine:      b.engine,
		dist:        b.dist,
		horn:        b.horn,
		led:         b.led,
		light:       b.light,
		camera:      b.camera,
		cameraAngle: 0,
		autodrive:   false,
		chOp:        make(chan CarOp, chSize),
	}
}

// Car ...
type Car struct {
	engine      IEngine
	dist        IDistance
	horn        IHorn
	led         ILed
	light       ILight
	camera      ICamera
	cameraAngle int
	autodrive   bool
	chOp        chan CarOp
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
		case autodriveon:
			go c.autoDriveOn()
		case autodriveoff:
			go c.autoDriveOff()
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
	c.delay(150)
	c.engine.Stop()
}

// right ...
func (c *Car) right() {
	log.Printf("car: right")
	c.engine.Right()
	c.delay(150)
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
			c.horn.Whistle()
			c.delay(100)
		}
	}()
}

func (c *Car) camLeft() {
	angle := c.cameraAngle - 15
	if angle < -90 {
		angle = -90
	}
	c.cameraAngle = angle
	log.Printf("camera %v", angle)
	c.camera.Turn(angle)
}

func (c *Car) camRight() {
	angle := c.cameraAngle + 15
	if angle > 90 {
		angle = 90
	}
	c.cameraAngle = angle
	log.Printf("camera %v", angle)
	if c.camera == nil {
		return
	}
	c.camera.Turn(angle)
}

func (c *Car) camAhead() {
	c.cameraAngle = 0
	log.Printf("camera %v", 0)
	if c.camera == nil {
		return
	}
	c.camera.Turn(0)
}

func (c *Car) autoDriveOn() {
	if c.dist == nil {
		return
	}
	// make a warning before into auto-drive mode
	for i := 0; i < 5 && c.horn != nil; i++ {
		c.horn.Whistle()
		c.delay(1000)
	}
	// start auto-drive
	c.autodrive = true
	for c.autodrive {
		d := c.dist.Dist()
		switch {
		case d > 100:
			c.chOp <- forward
			c.delay(2000)
		case d > 50:
			c.chOp <- forward
			c.delay(1000)
		case d > 20:
			c.chOp <- forward
			c.delay(500)
		default:
			c.chOp <- left
			c.delay(100)
		}
	}
}

func (c *Car) autoDriveOff() {
	c.autodrive = false
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
