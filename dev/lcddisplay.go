package dev

/*
LcdDisplay is a LCD Dispaly.

connect to raspberry pi:
- VCC: pin 1 or any 3.3v pin
- GND: pin 9 or and GND pin
- SDA: pin 3 (SDA)
- SCL: pin 5 (SCL)

Config Your Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 interface options] -> [p5 i2c] ->[yes] -> [ok]
5. $ sudo reboot now
6. check: $ sudo i2cdetect -y 1
	it works if you saw following message:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- 27 -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- -- --
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/
import (
	"github.com/shanghuiyang/rpi-devices/util"
	"golang.org/x/exp/io/i2c"
)

const (
	lcdDev               = "/dev/i2c-1"
	lcdAddr              = 0x27
	lcdEnable       byte = 0x04
	lcdBacklightOn  byte = 0x08
	lcdBacklightOff byte = 0x00
)

type LcdDisplay struct {
	width       int
	height      int
	backlightOn bool
	dev         *i2c.Device
}

func NewLcdDisplay(width, height int) (*LcdDisplay, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: lcdDev}, lcdAddr)
	if err != nil {
		return nil, err
	}
	lcd := &LcdDisplay{
		width:  width,
		height: height,
		dev:    dev,
	}
	lcd.sendCommand(0x33) // initialise
	util.DelayMs(1)

	lcd.sendCommand(0x32) // initialise
	util.DelayMs(1)

	lcd.sendCommand(0x06) // cursor move direction
	util.DelayMs(1)

	lcd.sendCommand(0x0C) // display on, cursor off, blink off
	util.DelayMs(1)

	lcd.sendCommand(0x28) // data length, number of lines, font size
	util.DelayMs(1)

	lcd.sendCommand(0x01) // Clear display
	util.DelayMs(1)

	return lcd, nil
}

func (lcd *LcdDisplay) Display(x, y int, text string) error {
	if x < 0 {
		x = 0
	}
	if x >= lcd.width {
		x = lcd.width - 1
	}

	if y < 0 {
		y = 0
	}
	if y >= lcd.height {
		y = lcd.height - 1
	}

	startPos := byte(0x80 + 0x40*y + x)
	if err := lcd.sendCommand(startPos); err != nil {
		return err
	}
	btext := []byte(text)
	for _, c := range btext {
		if err := lcd.sendData(c); err != nil {
			return err
		}
	}
	return nil
}

func (lcd *LcdDisplay) BackLightOn() error {
	lcd.backlightOn = true
	return lcd.dev.Write([]byte{lcdBacklightOn})
}

func (lcd *LcdDisplay) BackLightOff() error {
	lcd.backlightOn = false
	return lcd.dev.Write([]byte{lcdBacklightOff})
}

func (lcd *LcdDisplay) Clear() error {
	return lcd.sendCommand(0x01)
}

func (lcd *LcdDisplay) Close() error {
	if err := lcd.Clear(); err != nil {
		return err
	}
	if err := lcd.BackLightOff(); err != nil {
		return err
	}
	if err := lcd.dev.Close(); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) sendCommand(b byte) error {
	backlight := lcdBacklightOff
	if lcd.backlightOn {
		backlight = lcdBacklightOn
	}
	high := (b & 0xF0) | backlight
	low := ((b << 4) & 0xF0) | backlight

	if err := lcd.dev.Write([]byte{high}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(high); err != nil {
		return err
	}

	if err := lcd.dev.Write([]byte{low}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(low); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) sendData(b byte) error {
	backlight := lcdBacklightOff
	if lcd.backlightOn {
		backlight = lcdBacklightOn
	}
	high := 0x1 | (b & 0xF0) | backlight
	low := 0x1 | ((b << 4) & 0xF0) | backlight

	if err := lcd.dev.Write([]byte{high}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(high); err != nil {
		return err
	}

	if err := lcd.dev.Write([]byte{low}); err != nil {
		return err
	}
	if err := lcd.triggeEnable(low); err != nil {
		return err
	}
	return nil
}

func (lcd *LcdDisplay) triggeEnable(b byte) error {
	util.DelayMs(1)
	buf := b | lcdEnable
	if err := lcd.dev.Write([]byte{buf}); err != nil {
		return err
	}

	util.DelayMs(1)
	buf = b & ^lcdEnable
	if err := lcd.dev.Write([]byte{buf}); err != nil {
		return err
	}
	util.DelayMs(1)
	return nil
}
