package selfnav

import "github.com/shanghuiyang/rpi-devices/util/geo"

type SelfNav interface {
	Start(dest *geo.Point)
	Stop()
	InNaving() bool
}
