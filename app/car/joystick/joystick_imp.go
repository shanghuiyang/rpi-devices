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

type JoystickImp struct {
	car      car.Car
	wireless dev.Wireless
}

func NewJoystickImp(c car.Car, w dev.Wireless) *JoystickImp {
	return &JoystickImp{
		car:      c,
		wireless: w,
	}
}

func (j *JoystickImp) Start() {
	if j.wireless == nil {
		return
	}

	j.wireless.Wakeup()
	defer j.wireless.Sleep()

	for {
		time.Sleep(200 * time.Millisecond)

		data, err := j.wireless.Receive()
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
