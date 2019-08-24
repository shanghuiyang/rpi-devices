package base

import (
	"fmt"
)

// Point is GPS point
type Point struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

func (p *Point) String() string {
	return fmt.Sprintf("lat: %.6f, lon: %.6f", p.Lat, p.Lon)
}

// SendEmail ...
func SendEmail(info *EmailInfo) {
	chEmail <- info
}

// GetEmailList ...
func GetEmailList() []string {
	return emailList
}
