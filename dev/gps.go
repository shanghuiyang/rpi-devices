/*
Package dev ...

GPS is the driver of NEO-M6 module.

Config Your Pi:
1. $ sudo raspi-config
	-> [P5 interface] -> P6 Serial: disable -> [no] -> [yes]
2. $ sudo vim /boot/config.txt
	add following two lines:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	enable_uart=1
	dtoverlay=pi3-miniuart-bt
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. $ sudo reboot now
4. $ sudo cat /dev/ttyAMA0
	should see somethings output

Connect NEO-6M GPS Sensor to Raspberry Pi:
 - VCC: any 5v pin
 - GND: any gnd pin
 - RXT: must connect to pin  8(gpio 14) (TXD)
 - TXD: must connect to pin 10(gpio 15) (RXD)

-----------------------------------------------------------------------

		                   +-----------------+
		                   |       GPS       |
		                   |      NEO-M6     |
		                   |                 |
		                   +--+---+---+---+--+
		                      |   |   |   |
		                     GND TXD RXD VCC
		+-----------+         |   |   |   |
		| o 1   2 * +---------|---|---|---+
		| o       o |         |   |   |
		| 8     6 * +---------+   |   |
		| o     8 * |-------------|---+
		| o    10 * +-------------+
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o       o |
		| o 39 40 o |
		+-----------+

-----------------------------------------------------------------------
*/
package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/jakefau/rpi-devices/util/geo"
	"github.com/tarm/serial"
)

const (
	// Recommended Minimum Specific Data from GPS
	gpsRMC = "$GPRMC"
	// Recommended Minimum Specific Data from GPS & Beidou/China
	gpsAndBdRMC = "$GNRMC"
)

var (
	buf = make([]byte, 1024)
)

// GPS ...
type GPS struct {
	port *serial.Port
}

// NewGPS ...
func NewGPS(dev string, baud int) *GPS {
	g := &GPS{}
	if err := g.open(dev, baud); err != nil {
		return nil
	}
	return g
}

// Loc ...
func (g *GPS) Loc() (*geo.Point, error) {
	if err := g.port.Flush(); err != nil {
		return nil, err
	}
	a := 0
	for a < 512 {
		n, err := g.port.Read(buf[a:])
		if err != nil {
			return nil, err
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
		if strings.Contains(line, gpsRMC) || strings.Contains(line, gpsAndBdRMC) {
			loc = line
			break
		}
	}

	if loc == "" {
		return nil, fmt.Errorf("failed to read location from gps device")
	}
	items := strings.Split(loc, ",")
	if len(items) < 7 {
		return nil, fmt.Errorf("bad data from gps device")
	}

	var pt geo.Point
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
	mm := pt.Lat - float64(dd*100)
	pt.Lat = float64(dd) + mm/60

	dd = int(pt.Lon / 100)
	mm = pt.Lon - float64(dd*100)
	pt.Lon = float64(dd) + mm/60

	return &pt, nil
}

// Close ...
func (g *GPS) Close() {
	g.port.Close()
}

func (g *GPS) open(dev string, baud int) error {
	c := &serial.Config{Name: dev, Baud: baud}
	p, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	g.port = p
	return nil
}
