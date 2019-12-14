package dev

import (
	// "bufio"
	// "bytes"
	"fmt"
	// "io"
	"log"
	// "os"
	// "strings"

	// "github.com/shanghuiyang/rpi-devices/base"
	"github.com/tarm/serial"
)

const (
	logTagPMS7003 = "PMS7003"
)

var (
	pmBuf = make([]byte, 128)
)

// PMS7003 ...
type PMS7003 struct {
	port *serial.Port
}

// NewPMS7003 ...
func NewPMS7003() *PMS7003 {
	p := &PMS7003{}
	if err := p.open(); err != nil {
		return nil
	}
	return p
}

// PM25 ...
func (p *PMS7003) PM25() (float32, error) {
	if err := p.port.Flush(); err != nil {
		return -1, err
	}
	a := 0
	for a < 24 {
		n, err := p.port.Read(pmBuf[a:])
		if err != nil {
			// try to reopen serial
			p.port.Close()
			if err := p.open(); err != nil {
				log.Printf("[%v]failed open serial, error: %v", logTagPMS7003, err)
			}
			return -1, fmt.Errorf("error on read from port, error: %v. try to open serial again", err)
		}
		a += n
	}
	// r := bufio.NewReader(bytes.NewReader(buf[:a]))
	log.Printf("data: %v\n", string(pmBuf[:a]))
	// loc := ""
	// for {
	// 	line, err := r.ReadString('\n')
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	line = strings.Trim(line, " \t\n")
	// 	if strings.Contains(line, gpsRMC) || strings.Contains(line, gpsAndBdRMC) {
	// 		loc = line
	// 		break
	// 	}
	// }

	// if loc == "" {
	// 	return -1, fmt.Errorf("failed to read location from gps device")
	// }

	return -1, nil
}

// Close ...
func (p *PMS7003) Close() {
	p.port.Close()
}

func (p *PMS7003) open() error {
	c := &serial.Config{Name: "/dev/ttyAMA0", Baud: 9600}
	port, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	p.port = port
	return nil
}
