package selfnav

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/a-star/astar"
	"github.com/shanghuiyang/a-star/tilemap"
	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

const (
	logTag = "selfnav"
)

type SelfNavImp struct {
	car       car.Car
	astar     *astar.AStar
	mapBBox   *geo.Bbox
	gridSize  float64
	gps       dev.GPS
	lastLoc   *geo.Point
	gpslogger *util.GPSLogger
	inNaving  bool
}

func NewSelfNavImp(c car.Car, gps dev.GPS, tilemap *tilemap.Tilemap, bbox *geo.Bbox, gridsize float64) *SelfNavImp {
	return &SelfNavImp{
		car:       c,
		gps:       gps,
		astar:     astar.New(tilemap),
		mapBBox:   bbox,
		gridSize:  gridsize,
		gpslogger: util.NewGPSLogger(),
		inNaving:  false,
	}
}

func (s *SelfNavImp) Start(dest *geo.Point) {
	if dest == nil {
		log.Printf("[%v]destination didn't be set, stop nav", logTag)
		return
	}
	s.inNaving = true
	defer s.Stop()

	s.car.Beep(3, 300)
	if !s.mapBBox.IsInside(dest) {
		log.Printf("[%v]destination isn't in bbox, stop nav", logTag)
		return
	}

	var org *geo.Point
	for s.inNaving {
		pt, err := s.gps.Loc()
		if err != nil {
			log.Printf("[%v]gps sensor is not ready", logTag)
			util.DelayMs(1000)
			continue
		}
		s.gpslogger.AddPoint(org)
		if !s.mapBBox.IsInside(pt) {
			log.Printf("[%v]current loc(%v) isn't in bbox(%v)", logTag, pt, s.mapBBox)
			continue
		}
		org = pt
		break
	}
	if !s.inNaving {
		return
	}
	s.lastLoc = org

	path, err := s.findPath(org, dest)
	if err != nil {
		log.Printf("[%v]failed to find a path, error: %v", logTag, err)
		return
	}
	turns := s.turnPoints(path)

	var turnPts []*geo.Point
	var str string
	for _, xy := range turns {
		pt := s.xy2geo(xy)
		str += fmt.Sprintf("(%v) ", pt)
		turnPts = append(turnPts, pt)
	}
	log.Printf("[%v]turn points(lat,lon): %v", logTag, str)

	s.car.Forward()
	util.DelayMs(1000)
	for i, p := range turnPts {
		if err := s.navTo(p); err != nil {
			log.Printf("[%v]failed to nav to (%v), error: %v", logTag, p, err)
			break
		}
		if i < len(turnPts)-1 {
			// turn point
			go s.car.Beep(2, 100)
		} else {
			// destination
			go s.car.Beep(5, 300)
		}
	}
	s.car.Stop()
}

func (s *SelfNavImp) InNaving() bool {
	return s.inNaving
}

func (s *SelfNavImp) Stop() {
	s.inNaving = false
}

func (s *SelfNavImp) navTo(dest *geo.Point) error {
	retry := 8
	for s.inNaving {
		loc, err := s.gps.Loc()
		if err != nil {
			s.car.Stop()
			log.Printf("[%v]gps sensor is not ready", logTag)
			util.DelayMs(1000)
			continue
		}

		if !s.mapBBox.IsInside(loc) {
			s.car.Stop()
			log.Printf("[%v]current loc(%v) isn't in bbox(%v)", logTag, loc, s.mapBBox)
			util.DelayMs(1000)
			continue
		}

		s.gpslogger.AddPoint(loc)
		log.Printf("[%v]current loc: %v", logTag, loc)

		d := loc.DistanceWith(s.lastLoc)
		log.Printf("[%v]distance to last loc: %.2f m", logTag, d)
		if d > 4 && retry < 5 {
			s.car.Stop()
			log.Printf("[%v]bad gps signal, waiting for better gps signal", logTag)
			retry++
			util.DelayMs(1000)
			continue
		}

		retry = 0
		d = loc.DistanceWith(dest)
		log.Printf("[%v]distance to destination: %.2f m", logTag, d)
		if d < 4 {
			s.car.Stop()
			log.Printf("[%v]arrived at the destination, nav done", logTag)
			return nil
		}

		side := geo.Side(s.lastLoc, loc, dest)
		angle := 180 - geo.Angle(s.lastLoc, loc, dest)
		if angle < 10 {
			side = geo.MiddleSide
		}
		log.Printf("[%v]nav angle: %v, side: %v", logTag, angle, side)

		switch side {
		case geo.LeftSide:
			s.car.Turn(angle * (-1))
		case geo.RightSide:
			s.car.Turn(angle)
		case geo.MiddleSide:
			// do nothing
		}
		s.car.Forward()
		util.DelayMs(1000)
		s.lastLoc = loc
	}
	s.car.Stop()
	return nil
}

func (s *SelfNavImp) findPath(org, des *geo.Point) (astar.PList, error) {
	orgXY := s.geo2xy(org)
	desXY := s.geo2xy(des)

	path, err := s.astar.FindPath(orgXY, desXY)
	if err != nil {
		log.Printf("[car]failed to find the path from A(%v) to B(%v)", org, des)
		return nil, err
	}
	log.Printf("[car]path: %v", path)
	s.astar.Draw()
	return path, nil
}

func (s *SelfNavImp) turnPoints(path astar.PList) astar.PList {
	if len(path) <= 2 {
		return path
	}

	var ks []float64
	for i := 0; i < len(path)-1; i++ {
		k := 99999.99
		if path[i].Y != path[i+1].Y {
			k = float64(path[i].X-path[i+1].X) / float64(path[i].Y-path[i+1].Y)
		}
		ks = append(ks, k)
	}
	log.Printf("ks: %v\n", ks)

	var turns astar.PList
	for i := 0; i < len(ks)-1; i++ {
		if ks[i] == ks[i+1] {
			continue
		}
		turns = append(turns, path[i+1])
	}
	turns = append(turns, path[len(path)-1])
	log.Printf("turn points(x,y): %v", turns)
	return turns
}

func (s *SelfNavImp) geo2xy(p *geo.Point) *astar.Point {
	return &astar.Point{
		X: int((s.mapBBox.Top-p.Lat)/s.gridSize + 0.5),
		Y: int((p.Lon-s.mapBBox.Left)/s.gridSize + 0.5),
	}
}

func (s *SelfNavImp) xy2geo(p *astar.Point) *geo.Point {
	return &geo.Point{
		Lat: s.mapBBox.Top - float64(p.X)*s.gridSize,
		Lon: s.mapBBox.Left + float64(p.Y)*s.gridSize,
	}
}
