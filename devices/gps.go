package devices

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shanghuiyang/pi/base"
	"github.com/shanghuiyang/pi/iotclouds"
	"github.com/tarm/serial"
)

const (
	logTagGPS = "gps"
)

var (
	points     []*base.Point
	pointCount int
	index      int
	buf        = make([]byte, 1024)
)

// GPS ...
type GPS struct {
	port    *serial.Port
	tracker *base.Tracker
}

// NewGPS ...
func NewGPS() *GPS {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	p, err := serial.OpenPort(c)
	if err != nil {
		return nil
	}
	tr := base.NewTracker()
	if tr == nil {
		return nil
	}
	g := &GPS{
		port:    p,
		tracker: tr,
	}
	return g
}

// Start ...
func (g *GPS) Start() {
	defer g.Close()

	log.Printf("[%v]start working", logTagGPS)
	for {
		time.Sleep(5 * time.Second)
		// pt, err := g.MockLocFromCSV()
		pt, err := g.Loc()
		if err != nil {
			log.Printf("[%v]failed to get gps locations: %v", logTagGPS, err)
			continue
		}
		g.tracker.AddPoint(pt)
		v := &iotclouds.IoTValue{
			DeviceName: GPSDevice,
			Value:      pt,
		}
		iotclouds.IotCloud.Push(v)
	}
}

// Loc ...
func (g *GPS) Loc() (*base.Point, error) {
	if err := g.port.Flush(); err != nil {
		return nil, err
	}
	a := 0
	for a < 512 {
		n, err := g.port.Read(buf[a:])
		if err != nil {
			return nil, fmt.Errorf("error on read from port, %v", err)
		}
		a += n
	}
	r := bufio.NewReader(bytes.NewReader(buf[:a]))
	loc := ""
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.Trim(line, " \t\n")
		if strings.Contains(line, "$GPRMC") {
			loc = line
			break
		}
	}
	// log.Printf("buf: %v\n", string(buf[:a]))

	if loc == "" {
		return nil, fmt.Errorf("failed to read location from gps device")
	}
	items := strings.Split(loc, ",")
	if len(items) < 7 {
		return nil, fmt.Errorf("bad data from gps device")
	}

	var pt base.Point
	if _, err := fmt.Sscanf(items[3], "%f", &pt.Lat); err != nil {
		return nil, fmt.Errorf("failed to parse lat, %v", err)
	}
	northOrSouth := ""
	if _, err := fmt.Sscanf(items[4], "%s", &northOrSouth); err != nil {
		return nil, fmt.Errorf("failed to parse north or south, %v", err)
	}
	if _, err := fmt.Sscanf(items[5], "%f", &pt.Lon); err != nil {
		return nil, fmt.Errorf("failed to parse lon, %v", err)
	}
	eastOrWest := ""
	if _, err := fmt.Sscanf(items[6], "%s", &eastOrWest); err != nil {
		return nil, fmt.Errorf("failed to parse east or west, %v", err)
	}
	if northOrSouth == "S" {
		pt.Lat = pt.Lat * (-1)
	}
	if eastOrWest == "W" {
		pt.Lon = pt.Lon * (-1)
	}
	dd := int(pt.Lat / 100)
	mm := pt.Lat - float32(dd*100)
	pt.Lat = float32(dd) + mm/60

	dd = int(pt.Lon / 100)
	mm = pt.Lon - float32(dd*100)
	pt.Lon = float32(dd) + mm/60

	return &pt, nil
}

// MockLocFromGPX ...
func (g *GPS) MockLocFromGPX() (*base.Point, error) {
	if pointCount == 0 {
		file, err := os.Open("gps.gpx")
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			var lat, lon float32
			line = strings.Trim(line, " \t\n")
			if !strings.Contains(line, "<trkpt") {
				continue
			}
			if n, err := fmt.Sscanf(line, `<trkpt lat="%f" lon="%f">`, &lat, &lon); n != 2 || err != nil {
				log.Printf("[%v]failed to parse lat/lon, error: %v", logTagGPS, err)
				continue
			}
			pt := &base.Point{
				Lat: lat,
				Lon: lon,
			}
			points = append(points, pt)
		}
		file.Close()
		pointCount = len(points)
		index = 0
	}
	if index >= pointCount {
		index = 0
	}
	pt := points[index]
	index++
	return pt, nil
}

// MockLocFromCSV ...
func (g *GPS) MockLocFromCSV() (*base.Point, error) {
	if pointCount == 0 {
		file, err := os.Open("gps.csv")
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			var timestamp string
			var lat, lon float32
			if _, err := fmt.Sscanf(line, "%19s,%f,%f\n", &timestamp, &lat, &lon); err != nil {
				log.Printf("[%v]failed to parse lat/lon, error: %v", logTagGPS, err)
			}
			pt := &base.Point{
				Lat: lat,
				Lon: lon,
			}
			points = append(points, pt)
		}
		file.Close()
		pointCount = len(points)
		index = 0
	}
	if index >= pointCount {
		index = 0
	}
	pt := points[index]
	index++
	return pt, nil
}

// Close ...
func (g *GPS) Close() {
	g.port.Close()
	g.tracker.Close()
}
