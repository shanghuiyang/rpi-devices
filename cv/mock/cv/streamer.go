package cv

import (
	"github.com/shanghuiyang/rpi-devices/cv/mock/gocv"
)

// Streamer ...
type Streamer struct{}

// NewStream ...
func NewStreamer(host string) *Streamer {
	return &Streamer{}
}

// Push ...
func (s *Streamer) Push(img *gocv.Mat) {}

// Start ...
func (s *Streamer) Start() {}

// Close ...
func (s *Streamer) Close() {}
