// +build gocv

package cv

import (
	"log"
	"net/http"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

// Stream ...
type Stream struct {
	host   string
	cam    *gocv.VideoCapture
	stream *mjpeg.Stream
}

// NewTracker ...
func NewStream(cam *gocv.VideoCapture, host string) *Stream {
	return &Stream{
		host:   host,
		cam:    cam,
		stream: mjpeg.NewStream(),
	}
}

// Locate ...
func (s *Stream) StartService() {
	http.Handle("/video", s.stream)
	if err := http.ListenAndServe(s.host, nil); err != nil {
		log.Printf("[stream]failed to listen and serve, err: %v", err)
		return
	}

}

func (s *Stream) Push(img *gocv.Mat) {
	buf, err := gocv.IMEncode(".jpg", *img)
	if err != nil {
		log.Printf("[stream]failed to push image, err: %v", err)
		return
	}
	s.stream.UpdateJPEG(buf)
}
