package joystick

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	logTag = "joystick"
)

type Joystick struct {
	car   car.Car
	lc12s *dev.LC12S
}

func New(c car.Car, l *dev.LC12S) *Joystick {
	return &Joystick{
		car:   c,
		lc12s: l,
	}
}

func (j *Joystick) Start() {
	if j.lc12s == nil {
		return
	}

	j.lc12s.Wakeup()
	defer j.lc12s.Sleep()

	for {
		time.Sleep(200 * time.Millisecond)

		data, err := j.lc12s.Receive()
		if err != nil {
			log.Printf("[%v]failed to receive data from LC12S, error: %v", logTag, err)
			continue
		}
		log.Printf("[%v]LC12S received: %v", logTag, data)

		if len(data) != 1 {
			log.Printf("[%v]invalid data from LC12S, data len: %v", logTag, len(data))
			continue
		}

		op := (data[0] >> 4)
		speed := data[0] & 0x0F

		switch op {
		case 0:
			j.car.Stop()
		case 1:
			j.car.Forward()
		case 2:
			j.car.Backward()
		case 3:
			j.car.Left()
		case 4:
			j.car.Right()
		case 5:
			// reserve
		default:
			j.car.Stop()
		}
		j.car.Speed(uint32(speed * 10))
	}
}
