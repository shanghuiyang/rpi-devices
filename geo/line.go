package geo

// Line ...
type Line struct {
	Points []*Point
}

// NewLine ...
func NewLine(pts []*Point) *Line {
	return &Line{
		Points: pts,
	}
}
