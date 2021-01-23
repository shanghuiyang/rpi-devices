package main

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const (
	devName = "/dev/ttyAMA0"
	baud    = 9600
	swPin   = 7
	csPin   = 17

	homeX   = 2.43
	homeY   = 2.56
	homeBuf = 0.15
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	l, err := dev.NewLC12S(devName, baud, csPin)
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

	util.WaitQuit(func() {
		rpio.Close()
	})

	l.Wakeup()

	var curOp, curSpeed byte
	for {
		time.Sleep(200 * time.Millisecond)
		op, speed, err := getOpAndSpeed(j)
		if err != nil {
			continue
		}
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

func getOpAndSpeed(j *dev.Joystick) (op, speed byte, err error) {
	op, speed = 0, 0
	err = nil

	x := j.X()
	y := j.Y()
	z := j.Z()
	// log.Printf("x: %.2f, y: %.2f, z: %v", x, y, z)

	if z == 1 {
		op = 5 // self-driving
		speed = 0
		return
	}

	dx := x - homeX
	dy := y - homeY

	absdx, absdy := math.Abs(dx), math.Abs(dy)
	// log.Printf("x: %.2f, y: %.2f, z: %v, dx: %.2f, dy: %.2f", x, y, z, absdx, absdy)
	// return

	if absdx < homeBuf && absdy < homeBuf {
		// home
		return
	}

	if absdx > 1 && absdy > 1 {
		// invalid data
		log.Printf("invalid data, absdx: %v, absdy: %v", absdx, absdy)
		err = errors.New("invalid data")
		return
	}

	speed = 3
	if absdx > absdy {
		op = 1 // forward
		if dx > 0 {
			op = 2 // backward
		}
		if absdx > 2.40 {
			speed = 4
		}
		return
	}

	op = 4 // right
	if dy > 0 {
		op = 3 // left
	}
	return
}
