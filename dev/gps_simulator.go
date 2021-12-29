/*
GPSSimulator simulates GPS module.
*/

package dev

import (
	"errors"
)

// GPSSimulator implements GPS interface
type GPSSimulator struct {
	index   int
	latlons [][]float64
}

// NewGPSSimulator ...
func NewGPSSimulator(latlons [][]float64) (*GPSSimulator, error) {
	return &GPSSimulator{
		index:   0,
		latlons: latlons,
	}, nil
}

// Loc ...
func (gps *GPSSimulator) Loc() (lat, lon float64, err error) {
	delayMs(1000)
	n := len(gps.latlons)
	if n == 0 {
		return 0, 0, errors.New("without data")
	}
	if gps.index >= len(gps.latlons) {
		gps.index = 0
	}
	lat = gps.latlons[gps.index][0]
	lon = gps.latlons[gps.index][1]
	gps.index++
	return lat, lon, nil
}

// Close ...
func (gps *GPSSimulator) Close() error {
	return nil
}
