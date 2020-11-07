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
		log.Printf("op: %v, speed: %v, data: %v", op, speed, data)

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
	// log.Printf("x: %.2f, y: %.2f, z: %v", x, y, z)

	if z == 1 {
		op = 5 // self-driving
		speed = 0
		return
	}

	dx := x - 2.45
	dy := y - 2.58

	absdx, absdy := math.Abs(dx), math.Abs(dy)

	if absdx < 0.35 && absdy < 0.35 {
		// home
		return
	}

	if absdx > 1 && absdy > 1 {
		// invalid data
		log.Printf("indvild data, absdx: %v, absdy: %v", absdx, absdy)
		return
	}

	speed = 3
	if absdx > absdy {
		op = 1 // forward
		if dx > 0 {
			op = 2 // backward
		}
		if absdx > 2 {
			speed = 4
		}
		return
	}

	op = 4 // right
	if dy > 0 {
		op = 3 // left
	}
	if absdy > 2 {
		speed = 4
	}
	return
}
