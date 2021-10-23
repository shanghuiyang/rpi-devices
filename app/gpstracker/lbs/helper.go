package lbs

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/shanghuiyang/rpi-devices/util/geo"
)

func loadPoints(csv string) ([]*geo.Point, error) {
	file, err := os.Open(csv)
	if err != nil {
		return nil, err
	}

	points := []*geo.Point{}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		var timestamp string
		var lat, lon float64
		if _, err := fmt.Sscanf(line, "%19s,%f,%f\n", &timestamp, &lat, &lon); err != nil {
			return nil, err
		}
		pt := &geo.Point{
			Lat: lat,
			Lon: lon,
		}
		points = append(points, pt)
	}
	file.Close()

	return points, nil
}
