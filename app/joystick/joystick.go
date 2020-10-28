package main

import (
	"log"
	"math"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	swPin = 7
	csPin = 17
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	l, err := dev.NewLC12S(csPin)
	if err != nil {
		log.Fatalf("failed to new LC12S, error: %v", err)
		return
	}
	defer l.Close()

	j, err := dev.NewJoystick(swPin)
	if err != nil {
		log.Printf("failed to new joystick")
		return
	}

	base.WaitQuit(func() {
		rpio.Close()
	})

	l.Wakeup()

	var curOp, curSpeed byte
	for {
		time.Sleep(300 * time.Millisecond)
		op, speed := getOpAndSpeed(j)
		if op == curOp && speed == curSpeed {
			continue
		}
		data := (op << 4) | speed
		// log.Printf("op: %v, speed: %v, data: %v", op, speed, data)

		if err := l.Send([]byte{data}); err != nil {
			log.Printf("lc12s failed to send data, error: %v", err)
			continue
		}
		curOp, curSpeed = op, speed
	}
}

func getOpAndSpeed(j *dev.Joystick) (op, speed byte) {
	op, speed = 0, 0

	x := j.X()
	y := j.Y()
	z := j.Z()
	log.Printf("x: %.2f, y: %.2f, z: %v", x, y, z)

	if z == 1 {
		op = 5 // self-driving
		speed = 0
		return
	}
	if x >= 2.1 && x <= 2.8 && y >= 2.2 && y <= 2.9 {
		op = 0 // stop
		speed = 0
		return
	}

	dx := math.Abs(x - 2.45)
	dy := math.Abs(y - 2.58)

	if dx > dy {
		// if x < 2.2 {
		// 	op = 1 // forward
		// 	speed = byte(50 * (1 - x/2.3))

		// } else if x > 2.7 {
		// 	op = 2 // backward
		// 	speed = byte(37.5*x - 97.5)
		// }
		if x < 2.1 && x >= 1.4 {
			op = 1 // forward
			speed = 2
		} else if x < 1.4 && x >= 0.7 {
			op = 1
			speed = 3
		} else if x < 0.7 {
			op = 1
			speed = 4
		} else if x > 2.8 && x <= 3.2 {
			op = 2 // backward
			speed = 2
		} else if x > 3.2 && x <= 3.6 {
			op = 2
			speed = 3
		} else if x > 3.6 {
			op = 2
			speed = 4
		}

	} else {
		// if y < 2.2 {
		// 	op = 4 // right
		// 	speed = byte(90 * (1 - y/2.3))
		// } else if y > 2.7 {
		// 	op = 3 // left
		// 	speed = byte(37.5*y - 97.5)
		// }
		if y < 2.2 && y >= 1.46 {
			op = 4 // right
			speed = 2
		} else if y < 1.46 && y >= 0.73 {
			op = 4
			speed = 3
		} else if y < 0.73 {
			op = 4
			speed = 4
		} else if y > 2.9 && y <= 3.3 {
			op = 3 // left
			speed = 2
		} else if y > 3.3 && y <= 3.7 {
			op = 3
			speed = 3
		} else if y > 3.7 {
			op = 3
			speed = 4
		}
	}
	return
}
