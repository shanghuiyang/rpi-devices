package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
)

func main() {
	m, err := dev.NewMPU6050()
	if err != nil {
		log.Printf("failed to create MPU6050 sensor, error: %v", err)
		return
	}

	base.WaitQuit(m.Close)
	for {
		gx, gy, gz := m.GetAcc()
		log.Printf("gx=%v, gy=%v, gz=%v", gx, gy, gz)

		x, y, z := m.GetRot()
		log.Printf("x=%v, y=%v, z=%v", x, y, z)

		time.Sleep(100 * time.Millisecond)
	}
}
