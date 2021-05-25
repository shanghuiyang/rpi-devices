package cv

import (
	"log"
	"net/http"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

const (
	chSize = 512
)

// Streamer ...
type Streamer struct {
	host   string
	stream *mjpeg.Stream
	chImg  chan *gocv.Mat
}

// NewStream ...
func NewStreamer(host string) *Streamer {
	return &Streamer{
		host:   host,
		stream: mjpeg.NewStream(),
		chImg:  make(chan *gocv.Mat, chSize),
	}
}

// SetImage ...
func (s *Streamer) SetImage(img *gocv.Mat) {
	im := img.Clone()
	s.chImg <- &im
}

// Start ...
func (s *Streamer) Start() {
	go func() {
		http.Handle("/video", s.stream)
		if err := http.ListenAndServe(s.host, nil); err != nil {
			log.Printf("[stream]failed to listen and serve, err: %v", err)
			return
		}
	}()

	for img := range s.chImg {
		buf, err := gocv.IMEncode(".jpg", *img)
		if err != nil {
			log.Printf("[stream]failed to encode image, err: %v", err)
			continue
		}
		s.stream.UpdateJPEG(buf)
		img.Close()
	}
}

// Close ...
func (s *Streamer) Close() {
	close(s.chImg)
}
