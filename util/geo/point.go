package geo

import (
	"fmt"
)

// Point is GPS point
type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// String ...
func (p *Point) String() string {
	return fmt.Sprintf("lat: %.6f, lon: %.6f", p.Lat, p.Lon)
}

// DistanceWith ...
func (p *Point) DistanceWith(pt *Point) float64 {
	return Distance(p, pt)
}
