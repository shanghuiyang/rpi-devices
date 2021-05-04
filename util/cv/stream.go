package cv

import (
	"image"
	"image/color"
	"log"
	"net/http"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

// Stream ...
type Stream struct {
	host                   string
	tracked                bool
	lh, ls, lv, hh, hs, hv float64
	cam                    *gocv.VideoCapture
	stream                 *mjpeg.Stream
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
		if s.tracked {
			s.track(&img)
		}
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("[stream]failed to push image, err: %v", err)
			continue
		}
		s.stream.UpdateJPEG(buf)
	}
}

func (s *Stream) Tracked(enabled bool) {
	s.tracked = enabled
}

func (s *Stream) SetHSV(lh, ls, lv, hh, hs, hv float64) {
	s.lh = lh
	s.ls = ls
	s.lv = lv
	s.hh = hh
	s.hs = hs
	s.hv = hv
}

func (s *Stream) track(img *gocv.Mat) {
	rcolor := color.RGBA{G: 255, A: 255}
	lhsv := gocv.Scalar{Val1: s.lh, Val2: s.ls, Val3: s.lv}
	hhsv := gocv.Scalar{Val1: s.hh, Val2: s.hs, Val3: s.hv}

	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}

	mask := gocv.NewMat()
	frame := gocv.NewMat()
	hsv := gocv.NewMat()
	kernel := gocv.NewMat()

	defer mask.Close()
	defer frame.Close()
	defer hsv.Close()
	defer kernel.Close()

	gocv.Flip(*img, img, 1)
	gocv.Resize(*img, img, size, 0, 0, gocv.InterpolationLinear)
	gocv.GaussianBlur(*img, &frame, blur, 0, 0, gocv.BorderReflect101)
	gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
	gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
	gocv.Erode(mask, &mask, kernel)
	gocv.Dilate(mask, &mask, kernel)

	cnt := bestContour(&mask, 200)
	if len(cnt) == 0 {
		return
	}

	rect := gocv.BoundingRect(cnt)
	gocv.Rectangle(img, rect, rcolor, 2)
	return

}

func bestContour(frame *gocv.Mat, minArea float64) []image.Point {
	cnts := gocv.FindContours(*frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var (
		bestCnt  []image.Point
		bestArea = minArea
	)
	for _, cnt := range cnts {
		if area := gocv.ContourArea(cnt); area > bestArea {
			bestArea = area
			bestCnt = cnt
		}
	}
	return bestCnt
}
