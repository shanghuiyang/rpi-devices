package geo

import (
	"fmt"
)

// Bbox ...
type Bbox struct {
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
	Top    float64 `json:"top"`
	Bottom float64 `json:"buttom"`
}

// IsInside ...
func (b *Bbox) IsInside(pt *Point) bool {
	if pt.Lat >= b.Bottom && pt.Lat <= b.Top && pt.Lon >= b.Left && pt.Lon <= b.Right {
		return true
	}
	return false
}

// String ...
func (b *Bbox) String() string {
	return fmt.Sprintf("bottom: %.6f, top: %.6f, left: %.6f, right: %.6f", b.Bottom, b.Top, b.Left, b.Right)
}
