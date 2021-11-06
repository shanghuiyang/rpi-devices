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

	"github.com/shanghuiyang/rpi-devices/util/geo"
	"github.com/tarm/serial"
)

var (
	ht1818Buf = make([]byte, 2048)
)

// HT1818GPS implements GPS interface
type HT1818GPS struct {
	port *serial.Port
}

// NewHT1818GPS ...
func NewHT1818GPS(dev string, baud int) (*HT1818GPS, error) {
	g := &HT1818GPS{}
	if err := g.open(dev, baud); err != nil {
		return nil, err
	}
	return g, nil
}

// Loc ...
func (g *HT1818GPS) Loc() (*geo.Point, error) {
	if err := g.port.Flush(); err != nil {
		return nil, err
	}
	a := 0
	for a < 1024 {
		n, err := g.port.Read(ht1818Buf[a:])
		if err != nil {
			return nil, err
		}
		a += n
	}
	r := bufio.NewReader(bytes.NewReader(ht1818Buf[:a]))
	loc := ""
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = strings.Trim(line, " \t\n")
		if strings.Contains(line, GNRMC) {
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

	var available string
	if _, err := fmt.Sscanf(items[2], "%s", &available); err != nil {
		return nil, fmt.Errorf("failed to parse available, %v", err)
	}
	if available != "A" {
		return nil, fmt.Errorf("invalid data")
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
func (g *HT1818GPS) Close() {
	g.port.Close()
}

func (g *HT1818GPS) open(dev string, baud int) error {
	c := &serial.Config{
		Name: dev,
		Baud: baud,
	}
	p, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	g.port = p
	return nil
}
