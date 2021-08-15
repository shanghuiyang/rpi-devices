package car

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const defaultSpeed = 30

// Car ...
type Car interface {
	Forward()
	Backward()
	Left()
	Right()
	Stop()
	Speed(speed uint32)
	Beep(n int, interval int)
	Turn(angle float64)
}

type CarImp struct {
	engine dev.MotorDriver
	acc    dev.Accelerometer
	buzzer dev.Buzzer
}

func NewCarImp(engine dev.MotorDriver, acc dev.Accelerometer, buz dev.Buzzer) *CarImp {
	c := &CarImp{
		engine: engine,
		acc:    acc,
		buzzer: buz,
	}
	c.Speed(defaultSpeed)
	return c
}

func (c *CarImp) Forward() {
	c.engine.Forward()
}

func (c *CarImp) Backward() {
	c.engine.Backward()
}

func (c *CarImp) Left() {
	c.engine.Left()
}

func (c *CarImp) Right() {
	c.engine.Right()
}

func (c *CarImp) Stop() {
	c.engine.Stop()
}

func (c *CarImp) Speed(speed uint32) {
	c.engine.SetSpeed(speed)
}

func (c *CarImp) Beep(n int, interval int) {
	c.buzzer.Beep(n, interval)
}

func (c *CarImp) Turn(angle float64) {
	turnf := c.Right
	if angle < 0 {
		turnf = c.Left
		angle *= (-1)
	}

	yaw, _, _, err := c.acc.Angles()
	if err != nil {
		log.Printf("[car]failed to get angles from gy-25, error: %v", err)
		return
	}

	retry := 0
	for {
		turnf()
		yaw2, _, _, err := c.acc.Angles()
		if err != nil {
			log.Printf("[car]failed to get angles from gy-25, error: %v", err)
			if retry < 3 {
				retry++
				continue
			}
			break
		}
		ang := util.IncludedAngle(yaw, yaw2)
		if ang >= angle {
			break
		}
		util.DelayMs(100)
		c.Stop()
		util.DelayMs(100)
	}
	c.Stop()
}
