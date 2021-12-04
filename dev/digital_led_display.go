/*
DigitalLedDisplay is a digital led display module used to primarily display digits driving by 74HC595 dirver.
Please note that I only test it on a 4-bit led digital module.
And Only following chars were supported. Any char which didn't be spported will be displayed as '-'.
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
0 1 2 3 4 5 6 7 8 9
A B C D E F H I J L O P R S U Y Z
a b c d h i j l o p q s t u y
. - _ =
(and blank char ' ')
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Connect to Raspberry Pi:
 - VCC: 	any v3.3 pin
 - GND: 	any gnd pin
 - DIO: 	any data pin
 - SCLK:	any data pin
 - RCLK:	any data pin
*/
package dev

import (
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	refreshFrequency = 3 * time.Millisecond
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
	'J': 0x8F,
	'L': 0xE3,
	'O': 0x03,
	'P': 0x31,
	'R': 0x11,
	'S': 0x49,
	'U': 0x83,
	'Y': 0x89,
	'Z': 0x25,

	'a': 0x05,
	'b': 0xC1,
	'c': 0x79,
	'd': 0x85,
	'h': 0xD1,
	'i': 0xBF,
	'j': 0xBD,
	'l': 0xFB,
	'o': 0x39,
	'p': 0x31,
	'q': 0x19,
	's': 0x49,
	't': 0xE1,
	'u': 0xB9,
	'y': 0x89,

	'.': 0xFE,
	'-': 0xFD,
	'_': 0xEF,
	'=': 0x7D,

	' ': 0xFF,
}

// DigitalLedDisplay implements Display interface
type DigitalLedDisplay struct {
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

// NewDigitalLedDisplay ...
func NewDigitalLedDisplay(dioPin, rclkPin, sclkPin uint8) *DigitalLedDisplay {
	d := &DigitalLedDisplay{
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
func (d *DigitalLedDisplay) flushShcp() {
	d.sclkPin.Write(d.state ^ 0x01)
	d.sclkPin.Write(d.state)
}

// flushStcp Flush the Stcp pin
// call after data writes are done
func (d *DigitalLedDisplay) flushStcp() {
	d.rclkPin.Write(d.state ^ 0x01)
	d.rclkPin.Write(d.state)
}

// setBit sets an individual bit
func (d *DigitalLedDisplay) setBit(bit rpio.State) {
	d.dioPin.Write(bit)
	d.flushShcp()
}

// sendData sends the bytes to the shiftregister
func (d *DigitalLedDisplay) sendData(data uint8) {
	d.data = data
	for i := uint(0); i < 8; i++ {
		d.setBit(rpio.State((d.data >> i) & 0x01))
	}
	d.flushStcp()
}

func (d *DigitalLedDisplay) display() {
	text := "----"
	for {
		select {
		case txt := <-d.chText:
			text = reverse(txt)
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
func (d *DigitalLedDisplay) Display(text string) {
	d.chText <- text
}

// Open ...
func (d *DigitalLedDisplay) Open() {
	if d.opened {
		return
	}
	go d.display()
	d.opened = true
}

// Close ...
func (d *DigitalLedDisplay) Close() {
	if !d.opened {
		return
	}
	d.chText <- closedSignal
	<-d.chDone
	d.sendData(0xFF)
	d.sendData(0xF0)
	d.opened = false
}
