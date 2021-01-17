package mock

import (
	"errors"
	"image"
)

// Tracker ...
type Tracker struct{}

// NewTracker ...
func NewTracker(lh, ls, lv, hh, hs, hv float64) (*Tracker, error) {
	return nil, errors.New("not implement")
}

// Locate ...
func (t *Tracker) Locate() (bool, *image.Rectangle) {
	return false, nil
}

// MiddleXY ...
func (t *Tracker) MiddleXY(rect *image.Rectangle) (x int, y int) {
	return 0, 0
}

// Close ...
func (t *Tracker) Close() {
	return
}
