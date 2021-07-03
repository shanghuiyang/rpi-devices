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

var (
	onnav     bool
	mycar     car.Car
	aStar     *astar.AStar
	mapBBox   *geo.Bbox
	gridSize  float64
	gps       *dev.GPS
	lastLoc   *geo.Point
	gpslogger *util.GPSLogger
)

func Init(c car.Car, g *dev.GPS, tilemap *tilemap.Tilemap, bbox *geo.Bbox, gridsize float64) {
	mycar = c
	gps = g
	aStar = astar.New(tilemap)
	mapBBox = bbox
	gridSize = gridsize
	mapBBox = bbox
	gpslogger = util.NewGPSLogger()
	onnav = false
}

func Start(dest *geo.Point) {
	if dest == nil {
		log.Printf("[%v]destination didn't be set, stop nav", logTag)
		return
	}
	onnav = true
	defer Stop()

	mycar.Beep(3, 300)
	if !mapBBox.IsInside(dest) {
		log.Printf("[%v]destination isn't in bbox, stop nav", logTag)
		return
	}

	var org *geo.Point
	for onnav {
		pt, err := gps.Loc()
		if err != nil {
			log.Printf("[%v]gps sensor is not ready", logTag)
			util.DelayMs(1000)
			continue
		}
		gpslogger.AddPoint(org)
		if !mapBBox.IsInside(pt) {
			log.Printf("[%v]current loc(%v) isn't in bbox(%v)", logTag, pt, mapBBox)
			continue
		}
		org = pt
		break
	}
	if !onnav {
		return
	}
	lastLoc = org

	path, err := findPath(org, dest)
	if err != nil {
		log.Printf("[%v]failed to find a path, error: %v", logTag, err)
		return
	}
	turns := turnPoints(path)

	var turnPts []*geo.Point
	var str string
	for _, xy := range turns {
		pt := xy2geo(xy)
		str += fmt.Sprintf("(%v) ", pt)
		turnPts = append(turnPts, pt)
	}
	log.Printf("[%v]turn points(lat,lon): %v", logTag, str)

	mycar.Forward()
	util.DelayMs(1000)
	for i, p := range turnPts {
		if err := navTo(p); err != nil {
			log.Printf("[%v]failed to nav to (%v), error: %v", logTag, p, err)
			break
		}
		if i < len(turnPts)-1 {
			// turn point
			go mycar.Beep(2, 100)
		} else {
			// destination
			go mycar.Beep(5, 300)
		}
	}
	mycar.Stop()
}

func Status() bool {
	return onnav
}

func Stop() {
	onnav = false
}

func navTo(dest *geo.Point) error {
	retry := 8
	for onnav {
		loc, err := gps.Loc()
		if err != nil {
			mycar.Stop()
			log.Printf("[%v]gps sensor is not ready", logTag)
			util.DelayMs(1000)
			continue
		}

		if !mapBBox.IsInside(loc) {
			mycar.Stop()
			log.Printf("[%v]current loc(%v) isn't in bbox(%v)", logTag, loc, mapBBox)
			util.DelayMs(1000)
			continue
		}

		gpslogger.AddPoint(loc)
		log.Printf("[%v]current loc: %v", logTag, loc)

		d := loc.DistanceWith(lastLoc)
		log.Printf("[%v]distance to last loc: %.2f m", logTag, d)
		if d > 4 && retry < 5 {
			mycar.Stop()
			log.Printf("[%v]bad gps signal, waiting for better gps signal", logTag)
			retry++
			util.DelayMs(1000)
			continue
		}

		retry = 0
		d = loc.DistanceWith(dest)
		log.Printf("[%v]distance to destination: %.2f m", logTag, d)
		if d < 4 {
			mycar.Stop()
			log.Printf("[%v]arrived at the destination, nav done", logTag)
			return nil
		}

		side := geo.Side(lastLoc, loc, dest)
		angle := int(180 - geo.Angle(lastLoc, loc, dest))
		if angle < 10 {
			side = geo.MiddleSide
		}
		log.Printf("[%v]nav angle: %v, side: %v", logTag, angle, side)

		switch side {
		case geo.LeftSide:
			mycar.Turn(angle * (-1))
		case geo.RightSide:
			mycar.Turn(angle)
		case geo.MiddleSide:
			// do nothing
		}
		mycar.Forward()
		util.DelayMs(1000)
		lastLoc = loc
	}
	mycar.Stop()
	return nil
}
