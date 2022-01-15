/*
TM1637Display is a dirvier for digital led display module drived by TM1637 chip.
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
	"errors"
	"image"

	"github.com/stianeikeland/go-rpio/v4"
)

const (
	refreshFrequencyMs = 3
	closedSignal       = "*.*.*.*"
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

// TM1637Display is a dirvier for digital led display module drived by TM1637 chip.
// It is an implement of Display interface.
type TM1637Display struct {
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

// NewTM1637Display creates a TM1637Display driver.
// Please NOTE that I only test it on a 4-bit digital led module.
func NewTM1637Display(dioPin, rclkPin, sclkPin uint8) *TM1637Display {
	display := &TM1637Display{
		dioPin:  rpio.Pin(dioPin),
		rclkPin: rpio.Pin(rclkPin),
		sclkPin: rpio.Pin(sclkPin),
		chText:  make(chan string, 4),
		chDone:  make(chan bool),
		opened:  false,
	}

	display.dioPin.Output()
	display.dioPin.Low()

	display.rclkPin.Output()
	display.rclkPin.Low()

	display.sclkPin.Output()
	display.sclkPin.Low()

	return display
}

// Image displays an image on the screen.
// NOTE: Digital led display module can't be used to display an image.
// It is here just for implementing the Display interface.
func (display *TM1637Display) Image(img image.Image) error {
	return errors.New("digital led display module can't be used to display an image")
}

// Text display text on the screen.
// NOTE: (x, y) never be used. They are here just for implementing the Display interface.
func (display *TM1637Display) Text(text string, x, y int) error {
	display.chText <- text
	return nil
}

// On ...
func (display *TM1637Display) On() error {
	if display.opened {
		return nil
	}
	go display.display()
	display.opened = true
	return nil
}

// Off ...
func (display *TM1637Display) Off() error {
	if !display.opened {
		return nil
	}
	display.chText <- closedSignal
	<-display.chDone
	display.send(0xFF)
	display.send(0xF0)
	display.opened = false
	return nil
}

// Clear ...
func (display *TM1637Display) Clear() error {
	return display.Text("", 0, 0)
}

// Close ...
func (display *TM1637Display) Close() error {
	return display.Off()
}

// flushShcp Flush the Shcp pin
// call after each individual data write
func (display *TM1637Display) flushShcp() {
	display.sclkPin.Write(display.state ^ 0x01)
	display.sclkPin.Write(display.state)
}

// flushStcp Flush the Stcp pin
// call after data writes are done
func (display *TM1637Display) flushStcp() {
	display.rclkPin.Write(display.state ^ 0x01)
	display.rclkPin.Write(display.state)
}

// setBit sets an individual bit
func (display *TM1637Display) setBit(bit rpio.State) {
	display.dioPin.Write(bit)
	display.flushShcp()
}

// send sends the bytes to the shiftregister
func (display *TM1637Display) send(data uint8) {
	display.data = data
	for i := uint(0); i < 8; i++ {
		display.setBit(rpio.State((display.data >> i) & 0x01))
	}
	display.flushStcp()
}

func (display *TM1637Display) display() {
	text := "----"
	for {
		select {
		case txt := <-display.chText:
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
			display.send(ledc)
			display.send(pos)
			delayMs(refreshFrequencyMs)
		}
	}
	display.chDone <- true
}
