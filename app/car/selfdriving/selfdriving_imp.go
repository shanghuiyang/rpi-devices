package selfdriving

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	forward  operator = "forward"
	backward operator = "backward"
	left     operator = "left"
	right    operator = "right"
	stop     operator = "stop"
	turn     operator = "turn"
	scan     operator = "scan"

	logTag = "selfdriving"
)

var (
	scanningAngles = []float64{-90, -75, -60, -45, -30, -15, 0, 15, 30, 45, 60, 75, 90}
	aheadAngles    = []float64{0, -15, 0, 15}
)

type operator string

type SelfDrivingImp struct {
	car       car.Car
	dmeter    dev.DistanceMeter
	servo     dev.Motor
	indriving bool
}

func NewSelfDrivingImp(c car.Car, d dev.DistanceMeter, servo dev.Motor) *SelfDrivingImp {
	servo.Roll(0)
	return &SelfDrivingImp{
		car:       c,
		dmeter:    d,
		servo:     servo,
		indriving: false,
	}
}

func (s *SelfDrivingImp) Start() {
	if s.indriving {
		return
	}

	s.indriving = true

	var (
		fwd       bool
		retry     int
		mindAngle float64
		maxdAngle float64
		mind      float64
		maxd      float64
		op        = forward
		chOp      = make(chan operator, 4)
	)

	for s.indriving {
		select {
		case p := <-chOp:
			op = p
			for len(chOp) > 0 {
				log.Printf("[%v]skip op: %v", logTag, <-chOp)
			}
		default:
			// 	do nothing
		}
		log.Printf("[%v]op: %v", logTag, op)

		switch op {
		case backward:
			fwd = false
			s.car.Stop()
			util.DelayMs(20)
			s.car.Backward()
			util.DelayMs(500)
			chOp <- stop
			continue
		case stop:
			fwd = false
			s.car.Stop()
			util.DelayMs(20)
			chOp <- scan
			continue
		case scan:
			fwd = false
			mind, maxd, mindAngle, maxdAngle = s.lookingForWay()
			log.Printf("[%v]mind=%.0f, maxd=%.0f, mindAngle=%v, maxdAngle=%v", logTag, mind, maxd, mindAngle, maxdAngle)
			if mind < 10 && mindAngle != 90 && mindAngle != -90 && retry < 4 {
				chOp <- backward
				retry++
				continue
			}
			chOp <- turn
			retry = 0
		case turn:
			fwd = false
			s.car.Turn(maxdAngle)
			util.DelayMs(150)
			chOp <- forward
			continue
		case forward:
			if !fwd {
				s.car.Forward()
				fwd = true
				go s.lookingForObs(chOp)
			}
			util.DelayMs(50)
			continue
		}
	}
	s.car.Stop()
	util.DelaySec(1)
	close(chOp)
}

func (s *SelfDrivingImp) InDrving() bool {
	return s.indriving
}

func (s *SelfDrivingImp) Stop() {
	s.indriving = false
}

// lookingForWay looks for geting the min & max distance, and their corresponding angles
// mind: the min distance
// maxd: the max distance
// mindAngle: the angle correspond to the mind
// maxdAngle: the angle correspond to the maxd
func (s *SelfDrivingImp) lookingForWay() (mind, maxd, mindAngle, maxdAngle float64) {
	mind = 9999
	maxd = -9999
	for _, ang := range scanningAngles {
		s.servo.Roll(ang)
		util.DelayMs(200)
		d, err := s.dmeter.Dist()
		for i := 0; err != nil && i < 3; i++ {
			util.DelayMs(100)
			d, err = s.dmeter.Dist()
		}
		if err != nil {
			continue
		}
		log.Printf("[%v]scan: angle=%v, dist=%.0f", logTag, ang, d)
		if d < mind {
			mind = d
			mindAngle = ang
		}
		if d > maxd {
			maxd = d
			maxdAngle = ang
		}
	}
	s.servo.Roll(0)
	util.DelayMs(50)
	return
}

func (s *SelfDrivingImp) lookingForObs(chOp chan operator) {
	for s.indriving {
		for _, angle := range aheadAngles {
			s.servo.Roll(angle)
			util.DelayMs(100)
			d, err := s.dmeter.Dist()
			for i := 0; err != nil && i < 3; i++ {
				util.DelayMs(100)
				d, err = s.dmeter.Dist()
			}
			if err != nil {
				continue
			}

			if d < 20 {
				chOp <- backward
				return
			}
			if d < 40 {
				chOp <- stop
				return
			}
		}
	}
}
