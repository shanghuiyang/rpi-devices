package geo

import (
	"math"
)

// EarthRadius is the radius of the Earth in meter, UTM, WGS84
const EarthRadius = 6378137

const (
	// MiddleSide ...
	MiddleSide OnSide = 0
	// LeftSide ...
	LeftSide OnSide = 1
	// RightSide ...
	RightSide OnSide = -1
)

// OnSide ...
type OnSide int

// Rad ...
func Rad(degree float64) float64 {
	return degree * math.Pi / 180.0
}

// Distance ...
func Distance(p1, p2 *Point) float64 {
	rad1 := Rad(p1.Lat)
	rad2 := Rad(p2.Lat)
	a := rad1 - rad2
	b := Rad(p1.Lon) - Rad(p2.Lon)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(rad1)*math.Cos(rad2)*math.Pow(math.Sin(b/2), 2)))
	return s * EarthRadius
}

// Angle calculates the angle <AOB>
// the angle value: [0, 180] in degree
//
//     lat /|\
//          |     a *
//          |      /
//          |     /
//          |  o *------* b
//          |
//          +---------------> lon
//
func Angle(a, o, b *Point) float64 {
	oa := Distance(o, a)
	ob := Distance(o, b)
	ab := Distance(a, b)
	if oa == 0 || ob == 0 {
		return 0
	}
	cos := (ob*ob + oa*oa - ab*ab) / (2 * ob * oa)
	return math.Acos(cos) * 180 / math.Pi
}

// Side is calc point p on which side of (a->b)
//
//     a *
//       | \
//       |  \
//       v   \
//       |    * p
//       |
//     b *
//
// vector a->b: ab = b - a
// vector a->p: ap = p - a
// cross-product: m = ab x ap
// m > 0: p on the left
// m < 0: p on the right
// m = 0: p on the middle
func Side(a, b, p *Point) OnSide {
	ab := &Point{
		Lat: b.Lat - a.Lat,
		Lon: b.Lon - a.Lon,
	}
	ap := &Point{
		Lat: p.Lat - a.Lat,
		Lon: p.Lon - a.Lon,
	}

	m := ab.Lon*ap.Lat - ab.Lat*ap.Lon
	if m > 0 {
		return LeftSide
	}
	if m < 0 {
		return RightSide
	}
	return MiddleSide
}
