package cv

import (
	"image"

	"gocv.io/x/gocv"
)

// Tracker ...
type Tracker struct {
	cam  *gocv.VideoCapture
	lhsv *gocv.Scalar
	hhsv *gocv.Scalar

	size image.Point
	blur image.Point

	img    gocv.Mat
	mask   gocv.Mat
	frame  gocv.Mat
	hsv    gocv.Mat
	kernel gocv.Mat
}

// NewTracker ...
func NewTracker(lh, ls, lv, hh, hs, hv float64) (*Tracker, error) {
	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		return nil, err
	}
	return &Tracker{
		cam:    cam,
		lhsv:   &gocv.Scalar{Val1: lh, Val2: ls, Val3: lv},
		hhsv:   &gocv.Scalar{Val1: hh, Val2: hs, Val3: hv},
		size:   image.Point{X: 600, Y: 600},
		blur:   image.Point{X: 11, Y: 11},
		img:    gocv.NewMat(),
		mask:   gocv.NewMat(),
		frame:  gocv.NewMat(),
		hsv:    gocv.NewMat(),
		kernel: gocv.NewMat(),
	}, nil
}

// Locate ...
func (t *Tracker) Locate() (bool, *image.Rectangle) {
	t.cam.Grab(6)
	if !t.cam.Read(&t.img) {
		return false, nil
	}
	gocv.Flip(t.img, &t.img, 1)
	gocv.Resize(t.img, &t.img, t.size, 0, 0, gocv.InterpolationLinear)
	gocv.GaussianBlur(t.img, &t.frame, t.blur, 0, 0, gocv.BorderReflect101)
	gocv.CvtColor(t.frame, &t.hsv, gocv.ColorBGRToHSV)
	gocv.InRangeWithScalar(t.hsv, *t.lhsv, *t.hhsv, &t.mask)
	gocv.Erode(t.mask, &t.mask, t.kernel)
	gocv.Dilate(t.mask, &t.mask, t.kernel)
	cnt := t.bestContour(t.mask, 200)
	if len(cnt) == 0 {
		return false, nil
	}
	r := gocv.BoundingRect(cnt)
	return true, &r
}

// MiddleXY ...
func (t *Tracker) MiddleXY(rect *image.Rectangle) (x int, y int) {
	return (rect.Max.X-rect.Min.X)/2 + rect.Min.X, (rect.Max.Y-rect.Min.Y)/2 + rect.Min.Y
}

func (t *Tracker) bestContour(frame gocv.Mat, minArea float64) []image.Point {
	cnts := gocv.FindContours(frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
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

// Close ...
func (t *Tracker) Close() {
	t.cam.Close()
	t.img.Close()
	t.mask.Close()
	t.frame.Close()
	t.hsv.Close()
	t.kernel.Close()
}
