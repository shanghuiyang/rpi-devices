package devices

import (
	"log"
	"time"
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
	camleft  CarOp = "camleft"
	camright CarOp = "camright"
	camahead CarOp = "camahead"
	light    CarOp = "light"
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
	IsOn() bool
}

// ICamera ...
type ICamera interface {
	Turn(angle int)
}

// CarBuilder ...
type CarBuilder struct {
	engine IEngine
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
		engine: b.engine,
		horn:   b.horn,
		led:    b.led,
		light:  b.light,
		camera: b.camera,
		chOp:   make(chan CarOp, chSize),
	}
}

// Car ...
type Car struct {
	engine      IEngine
	horn        IHorn
	led         ILed
	light       ILight
	camera      ICamera
	cameraAngle int
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
		case brake:
			c.brake()
		case honk:
			go c.honk()
		case camleft:
			go c.camLeft()
		case camright:
			go c.camRight()
		case camahead:
			go c.camAhead()
		case light:
			go c.lightOnOrOff()
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
	time.Sleep(150 * time.Millisecond)
	c.engine.Stop()
}

// right ...
func (c *Car) right() {
	log.Printf("car: right")
	c.engine.Right()
	time.Sleep(150 * time.Millisecond)
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
	c.camera.Turn(angle)
}

func (c *Car) camAhead() {
	c.cameraAngle = 0
	log.Printf("camera %v", 0)
	c.camera.Turn(0)
}

func (c *Car) lightOnOrOff() {
	if c.light.IsOn() {
		c.light.Off()
	} else {
		c.light.On()
	}
}
