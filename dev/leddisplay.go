/*
Package dev ...

LedDisplay is based on the 74HC595 shiftregister hardware.
*/
package dev

import (
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/stianeikeland/go-rpio"
)

const (
	refreshFrequency = 5 * time.Millisecond
	closedSignal     = "*.*.*.*"
)

var ledchars = map[byte]uint8{
	'0': 0x03,
	'1': 0x9F,
	'2': 0x25,
	'3': 0x0D,
	'4': 0x99,
	'5': 0x49,
	'6': 0x41,
	'7': 0x1B,
	'8': 0x01,
	'9': 0x09,

	'A': 0x11,
	'B': 0x01,
	'C': 0x63,
	'D': 0x03,
	'E': 0x61,
	'F': 0x71,
	'H': 0x91,
	'I': 0x9F,
	'L': 0xE3,
	'O': 0x03,
	'P': 0x31,
	'R': 0x11,
	'S': 0x49,
	'U': 0x83,
	'Y': 0x89,

	'c': 0xE5,
	'.': 0xFE,
	'-': 0xFD,
	' ': 0xFF,
}

// LedDisplay ...
type LedDisplay struct {
	dioPin  rpio.Pin
	rclkPin rpio.Pin
	sclkPin rpio.Pin

	// on    bool
	state rpio.State
	data  uint8

	chText chan string
	chDone chan bool
	opened bool
}

// NewLedDisplay ...
func NewLedDisplay(dioPin, rclkPin, sclkPin uint8) *LedDisplay {
	d := &LedDisplay{
		dioPin:  rpio.Pin(dioPin),
		rclkPin: rpio.Pin(rclkPin),
		sclkPin: rpio.Pin(sclkPin),
		chText:  make(chan string, 4),
		chDone:  make(chan bool),
		opened:  false,
	}

	d.dioPin.Output()
	d.dioPin.Low()

	d.rclkPin.Output()
	d.rclkPin.Low()

	d.sclkPin.Output()
	d.sclkPin.Low()

	return d
}

// flushShcp Flush the Shcp pin
// call after each individual data write
func (d *LedDisplay) flushShcp() {
	d.sclkPin.Write(d.state ^ 0x01)
	d.sclkPin.Write(d.state)
}

// flushStcp Flush the Stcp pin
// call after data writes are done
func (d *LedDisplay) flushStcp() {
	d.rclkPin.Write(d.state ^ 0x01)
	d.rclkPin.Write(d.state)
}

// setBit sets an individual bit
func (d *LedDisplay) setBit(bit rpio.State) {
	d.dioPin.Write(bit)
	d.flushShcp()
}

// SendData sends the bytes to the shiftregister
func (d *LedDisplay) sendData(data uint8) {
	d.data = data
	for i := uint(0); i < 8; i++ {
		d.setBit(rpio.State((d.data >> i) & 0x01))
	}
	d.flushStcp()
}

func (d *LedDisplay) display() {
	text := "----"
	for {
		select {
		case txt := <-d.chText:
			text = base.Reverse(txt)
		default:
			// do nothing, just use the latest text for displaying
		}

		if text == closedSignal {
			break
		}

		dot := 0
		for i, c := range text {
			ledc, ok := ledchars[byte(c)]
			if !ok {
				ledc = ledchars['-']
			}

			if i > 0 && text[i-1] == '.' {
				dot++
			}
			pos := uint8(0x80) >> uint(i-dot)
			d.sendData(ledc)
			d.sendData(pos)
			time.Sleep(refreshFrequency)
		}
	}
	d.chDone <- true
}

// Display ...
func (d *LedDisplay) Display(text string) {
	d.chText <- text
}

// Open ...
func (d *LedDisplay) Open() {
	if d.opened {
		return
	}
	go d.display()
	d.opened = true
}

// Close ...
func (d *LedDisplay) Close() {
	if !d.opened {
		return
	}
	d.chText <- closedSignal
	<-d.chDone
	d.sendData(0xFF)
	d.sendData(0xF0)
	d.opened = false
}
