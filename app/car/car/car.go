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
	Turn(angle int)
}

type CarImp struct {
	l298n  *dev.L298N
	gy25   *dev.GY25
	buzzer *dev.Buzzer
}

func NewCarImp(l298n *dev.L298N, gy25 *dev.GY25, buz *dev.Buzzer) *CarImp {
	c := &CarImp{
		l298n:  l298n,
		gy25:   gy25,
		buzzer: buz,
	}
	c.Speed(defaultSpeed)
	return c
}

func (c *CarImp) Forward() {
	c.l298n.Forward()
}

func (c *CarImp) Backward() {
	c.l298n.Backward()
}

func (c *CarImp) Left() {
	c.l298n.Left()
}

func (c *CarImp) Right() {
	c.l298n.Right()
}

func (c *CarImp) Stop() {
	c.l298n.Stop()
}

func (c *CarImp) Speed(speed uint32) {
	c.l298n.Speed(speed)
}

func (c *CarImp) Beep(n int, interval int) {
	c.buzzer.Beep(n, interval)
}

func (c *CarImp) Turn(angle int) {
	turnf := c.Right
	if angle < 0 {
		turnf = c.Left
		angle *= (-1)
	}

	yaw, _, _, err := c.gy25.Angles()
	if err != nil {
		log.Printf("[car]failed to get angles from gy-25, error: %v", err)
		return
	}

	retry := 0
	for {
		turnf()
		yaw2, _, _, err := c.gy25.Angles()
		if err != nil {
			log.Printf("[car]failed to get angles from gy-25, error: %v", err)
			if retry < 3 {
				retry++
				continue
			}
			break
		}
		ang := c.gy25.IncludedAngle(yaw, yaw2)
		if ang >= float64(angle) {
			break
		}
		util.DelayMs(100)
		c.Stop()
		util.DelayMs(100)
	}
	c.Stop()
}
