package lbs

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func loadPoints(csv string) ([][]float64, error) {
	file, err := os.Open(csv)
	if err != nil {
		return nil, err
	}

	latlons := [][]float64{}
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
		latlons = append(latlons, []float64{lat, lon})
	}
	file.Close()

	return latlons, nil
}
