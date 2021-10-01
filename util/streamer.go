package util

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/hybridgroup/mjpeg"
)

const (
	chSize = 512
)

// Streamer ...
type Streamer struct {
	host   string
	path   string
	stream *mjpeg.Stream
	chImg  chan []byte
}

// NewStream ...
func NewStreamer(url string) (*Streamer, error) {
	pos := strings.Index(url, "/")
	if pos < 0 {
		return nil, errors.New("invalid url, the url would looks like 0.0.0.0:8080/stream")
	}
	host, path := url[:pos], url[pos:]
	s := &Streamer{
		host:   host,
		path:   path,
		stream: mjpeg.NewStream(),
		chImg:  make(chan []byte, chSize),
	}
	go s.start()
	return s, nil
}

// SetImage ...
func (s *Streamer) Push(img []byte) {
	s.chImg <- img
}

// Start ...
func (s *Streamer) start() {
	go func() {
		http.Handle(s.path, s.stream)
		if err := http.ListenAndServe(s.host, nil); err != nil {
			log.Printf("[streamer]failed to listen and serve, err: %v", err)
			return
		}
	}()

	for img := range s.chImg {
		s.stream.UpdateJPEG(img)

	}
}
