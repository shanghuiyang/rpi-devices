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

// NewStream ...
func NewStream(cam *gocv.VideoCapture, host string) *Stream {
	return &Stream{
		host:   host,
		cam:    cam,
		stream: mjpeg.NewStream(),
	}
}

// Start ...
func (s *Stream) Start() {
	go func() {
		http.Handle("/video", s.stream)
		if err := http.ListenAndServe(s.host, nil); err != nil {
			log.Printf("[stream]failed to listen and serve, err: %v", err)
			return
		}
	}()

	img := gocv.NewMat()
	defer img.Close()
	for true {
		s.cam.Grab(6)
		if !s.cam.Read(&img) {
			continue
		}
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("[stream]failed to push image, err: %v", err)
			continue
		}
		s.stream.UpdateJPEG(buf)
	}
}
