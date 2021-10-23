/*
GPSSimulator simulates GPS module.
*/

package dev

import (
	"errors"

	"github.com/shanghuiyang/rpi-devices/util/geo"
)

// GPSSimulator implements GPS interface
type GPSSimulator struct {
	index  int
	points []*geo.Point
}

// NewGPSSimulator ...
func NewGPSSimulator(points []*geo.Point) (*GPSSimulator, error) {
	return &GPSSimulator{
		index:  0,
		points: points,
	}, nil
}

// Loc ...
func (m *GPSSimulator) Loc() (*geo.Point, error) {
	n := len(m.points)
	if n == 0 {
		return nil, errors.New("without data")
	}
	if m.index >= len(m.points) {
		m.index = 0
	}
	pt := m.points[m.index]
	m.index++
	return pt, nil
}

// Close ...
func (m *GPSSimulator) Close() {
}
