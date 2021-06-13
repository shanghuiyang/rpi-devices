package car

import (
	"github.com/jakefau/rpi-devices/dev"
)

// Config ...
type Config struct {
	Engine     *dev.L298N
	Servo      *dev.SG90
	GY25       *dev.GY25
	Horn       *dev.Buzzer
	Led        *dev.Led
	Light      *dev.Led
	Camera     *dev.Camera
	GPS        *dev.GPS
	LC12S      *dev.LC12S
	Collisions []*dev.Collision
	DistMeter  dev.DistMeter
}
