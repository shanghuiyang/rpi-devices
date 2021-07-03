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

var (
	mycar car.Car
	lc12s *dev.LC12S
)

func Init(c car.Car, l *dev.LC12S) {
	mycar = c
	lc12s = l
}

func Start() {
	if lc12s == nil {
		return
	}

	lc12s.Wakeup()
	defer lc12s.Sleep()

	for {
		time.Sleep(200 * time.Millisecond)

		data, err := lc12s.Receive()
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
			mycar.Stop()
		case 1:
			mycar.Forward()
		case 2:
			mycar.Backward()
		case 3:
			mycar.Left()
		case 4:
			mycar.Right()
		case 5:
			// reserve
		default:
			mycar.Stop()
		}
		mycar.Speed(uint32(speed * 10))
	}
}
