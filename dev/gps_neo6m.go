/*
Neo6mGPS is a GPS module used to get locations(lat/lon).

Config Raspberry Pi:
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

Connect to Raspberry Pi:
 - VCC: any 3.3v or 5v pin
 - GND: any gnd pin
 - RXT: must connect to GPIO-14/TXD
 - TXD: must connect to GPIO-15/RXD
*/
package dev

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tarm/serial"
)

const (
	// Recommended Minimum Specific Data from GPS
	GPRMC = "$GPRMC,"
	// Recommended Minimum Specific Data from GPS & Beidou/China
	GNRMC = "$GNRMC,"
)

var (
	buf = make([]byte, 1024)
)

// Neo6mGPS implements GPS interface
type Neo6mGPS struct {
	port *serial.Port
}

// NewNeo6mGPS ...
func NewNeo6mGPS(dev string, baud int) (*Neo6mGPS, error) {
	cfg := &serial.Config{
		Name: dev,
		Baud: baud,
	}
	port, err := serial.OpenPort(cfg)
	if err != nil {
		return nil, err
	}
	return &Neo6mGPS{port}, nil
}

// Loc ...
func (gps *Neo6mGPS) Loc() (lat, lon float64, err error) {
	if err := gps.port.Flush(); err != nil {
		return 0, 0, fmt.Errorf("flush port error: %w", err)
	}
	a := 0
	for a < 512 {
		n, err := gps.port.Read(buf[a:])
		if err != nil {
			return 0, 0, fmt.Errorf("read port error: %w", err)
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
		if strings.Contains(line, GPRMC) || strings.Contains(line, GNRMC) {
			loc = line
			break
		}
	}

	if loc == "" {
		return 0, 0, fmt.Errorf("%v and %v not found", GPRMC, GNRMC)
	}
	items := strings.Split(loc, ",")
	if len(items) < 7 {
		return 0, 0, fmt.Errorf("invalid data format")
	}

	var available string
	if _, err := fmt.Sscanf(items[2], "%s", &available); err != nil {
		return 0, 0, fmt.Errorf("parse [available] error: %w", err)
	}
	if available != "A" {
		return 0, 0, fmt.Errorf("invalid data")
	}

	if _, err := fmt.Sscanf(items[3], "%f", &lat); err != nil {
		return 0, 0, fmt.Errorf("parse [lat] error: %w", err)
	}
	northOrSouth := ""
	if _, err := fmt.Sscanf(items[4], "%s", &northOrSouth); err != nil {
		return 0, 0, fmt.Errorf("parse [north/south] error: %w", err)
	}
	if _, err := fmt.Sscanf(items[5], "%f", &lon); err != nil {
		return 0, 0, fmt.Errorf("parse [lon] error: %w", err)
	}
	eastOrWest := ""
	if _, err := fmt.Sscanf(items[6], "%s", &eastOrWest); err != nil {
		return 0, 0, fmt.Errorf("parse [east/west] error: %w", err)
	}
	if northOrSouth == "S" {
		lat = lat * (-1)
	}
	if eastOrWest == "W" {
		lon = lon * (-1)
	}
	dd := int(lat / 100)
	mm := lat - float64(dd*100)
	lat = float64(dd) + mm/60

	dd = int(lon / 100)
	mm = lon - float64(dd*100)
	lon = float64(dd) + mm/60

	return lat, lon, nil
}

// Close ...
func (gps *Neo6mGPS) Close() error {
	return gps.port.Close()
}
