// bmw280 reads sensor data from Bosh BME280 sensor.
// taken from https://github.com/quhar/bme280
package dev

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// bus interface is used to communicate with bus where sensor is connected. Sensor supports SPI and I2C interfaces.
type bus interface {
	ReadReg(byte, []byte) error
	WriteReg(byte, []byte) error
}

// Compensation registers addresses.
const (
	T1 byte = 0x88 + iota*2
	T2
	T3
	P1
	P2
	P3
	P4
	P5
	P6
	P7
	P8
	P9
	H1 byte = 0xA1
	H2 byte = 0xE1
	H3 byte = 0xE3
	H4 byte = 0xE4
	H5 byte = 0xE5
	H6 byte = 0xE7
)

const (
	TempCompAddr  byte = 0x88
	PressCompAddr byte = 0x8E
	H1CompAddr    byte = 0xA1
	H2CompAddr    byte = 0xE1
)

// Data registers addresses.
const (
	HumAddr   byte = 0xFD
	TempAddr  byte = 0xFA
	PressAddr byte = 0xF7
	DataAddr  byte = 0xF7
)

// Other registers addresses.
const (
	IDAddr       byte = 0xD0
	ResetAddr    byte = 0xE0
	CtrlHumAddr  byte = 0xF2
	StatusAddr   byte = 0xF3
	CtrlMeasAddr byte = 0xF4
	ConfigAddr   byte = 0xF5
)

// General constants.
const (
	// I2CAddr is default BME280 I2C address.
	I2CAddr int = 0x77
	// ResetVal is a value which when written to ResetAddr resets the sensor.
	ResetVal byte = 0xB6
	// IDVal is a ID value of the sensor.
	IDVal byte = 0x60
)

// Oversampling constants.
const (
	_ byte = iota
	OverSmpl1
	OverSmpl2
	OverSmpl4
	OverSmpl8
	OverSmpl16
)

// Modes of operation.
const (
	SleepMode  byte = 0x00
	ForcedMode byte = 0x01
	NormalMode byte = 0x03
)

// Filter settins
const (
	FilterOff byte = iota << 2
	Filter2
	Filter4
	Filter8
	Filter16
)

// Normal mode standby modes.
const (
	// Standby for 0.5 ms.
	Stndby05 byte = iota << 5
	// Standby for 62.5 ms.
	Stndby625
	// Standby for 125 ms.
	Stndby125
	// Standby for 250 ms.
	Stndby250
	// Standby for 500 ms.
	Stndby500
	// Standby for 1000 ms.
	Stbdby1000
	// Standby for 10 ms.
	Stndby10
	// Standby for 20 ms.
	Stndby20
)

// Available temperature units.
const (
	Celsius = iota
	Fahrenheit
	Kelvin
)

// Available pressure units.
const (
	HPa = iota
	Bar
	PSI
)

// Option is an interface used to set various options in BME280 object.
type Option interface {
	set(b *BME280)
}

// option implements Option interface.
type option func(b *BME280)

func (o option) set(b *BME280) {
	o(b)
}

// Standby sets default Tstandby for Normal mode.
func Standby(s byte) Option {
	return option(func(b *BME280) {
		b.stndby = s
	})
}

// Mode sets default sensors Mode of operation (Sleep, Standby, Forced).
func Mode(m byte) Option {
	return option(func(b *BME280) {
		b.mode = m
	})
}

// TempOverSmpl sets default oversampling for temperature.
func TempOverSmpl(o byte) Option {
	return option(func(b *BME280) {
		b.tempOverSmpl = o
	})
}

// PressOverSmpl sets default oversampling for pressure.
func PressOverSmpl(o byte) Option {
	return option(func(b *BME280) {
		b.pressOverSmpl = o
	})
}

// HumOverSmpl sets default oversampling for humidity.
func HumOverSmpl(o byte) Option {
	return option(func(b *BME280) {
		b.humOverSmpl = o
	})
}

// Filter sets default IIR filter value.
func Filter(f byte) Option {
	return option(func(b *BME280) {
		b.filter = f
	})
}

// TempUnits set default temperature units.
func TempUnit(u int) Option {
	return option(func(b *BME280) {
		b.tempUnit = u
	})
}

// PressUnits set default pressure units.
func PressUnit(u int) Option {
	return option(func(b *BME280) {
		b.pressUnit = u
	})
}

// BME280 is an object representing BME280 sensor.
type BME280 struct {
	dev           bus
	mode          byte
	tempOverSmpl  byte
	pressOverSmpl byte
	humOverSmpl   byte
	filter        byte
	stndby        byte
	tempUnit      int
	pressUnit     int
	compensation  struct {
		temp struct {
			T1 uint16
			T2 int16
			T3 int16
		}
		press struct {
			P1 uint16
			P2 int16
			P3 int16
			P4 int16
			P5 int16
			P6 int16
			P7 int16
			P8 int16
			P9 int16
		}
		hum struct {
			H1 uint8
			H2 int16
			H3 uint8
			H4 int16
			H5 int16
			H6 int8
		}
	}
	tFine int32
}

// New returns new BME280 object.
// New object has following default values:
// mode = ForcedMode
// tempOverSmpl = 16
// pressOverSmpl = 16
// humOverSmpl = 16
// filter = off
// standby = 500ms
// Temperature units Celsius
// Presure units hPa
func New(dev bus, opts ...Option) *BME280 {
	b := &BME280{}
	b.dev = dev
	// set defaults
	b.mode = ForcedMode
	b.tempOverSmpl = OverSmpl16 << 5
	b.pressOverSmpl = OverSmpl16 << 3
	b.humOverSmpl = OverSmpl16
	b.filter = FilterOff
	b.stndby = Stndby500
	b.tempUnit = Celsius
	b.pressUnit = HPa

	for _, o := range opts {
		o.set(b)
	}
	return b
}

// init initializes BME280 and loads calibration data
func (b *BME280) Init() error {
	// read and check chip ID
	buf := make([]byte, 1)
	err := b.dev.ReadReg(IDAddr, buf)
	if buf[0] != IDVal {
		return fmt.Errorf("chip ID is different than expected, want(%X), got(%X)", IDVal, buf[0])
	}
	// maybe reset here
	// write config
	err = b.dev.WriteReg(ConfigAddr, []byte{b.stndby | b.filter})
	if err != nil {
		return fmt.Errorf("failed to write init data to sensor: %v", err)
	}
	// write ctrl_hum register
	err = b.dev.WriteReg(CtrlHumAddr, []byte{b.humOverSmpl})
	if err != nil {
		return fmt.Errorf("failed to write ctrl_hum data to sensor: %v", err)
	}
	// write ctrl_meas register
	err = b.dev.WriteReg(CtrlMeasAddr, []byte{b.tempOverSmpl | b.pressOverSmpl | b.mode})
	if err != nil {
		return fmt.Errorf("failed to write ctrl_meas data to sensor: %v", err)
	}
	return b.loadCompensation()
}

// loadCompensation reads compensation data from sensor.
func (b *BME280) loadCompensation() error {
	// Read temp compensation data.
	tData := make([]byte, 6)
	err := b.dev.ReadReg(TempCompAddr, tData)
	if err != nil {
		return fmt.Errorf("failed to read temperature compensation data: %v", err)
	}
	pData := make([]byte, 18)
	err = b.dev.ReadReg(PressCompAddr, pData)
	if err != nil {
		return fmt.Errorf("failed to read pressure compensation data: %v", err)
	}
	h1Data := make([]byte, 1)
	err = b.dev.ReadReg(H1CompAddr, h1Data)
	if err != nil {
		return fmt.Errorf("failed to first part of humidity compensation data: %v", err)
	}
	h2Data := make([]byte, 7)
	err = b.dev.ReadReg(H2CompAddr, h2Data)
	if err != nil {
		return fmt.Errorf("failed to second part of humidity compensation data: %v", err)
	}
	buf := bytes.NewReader(tData)
	err = binary.Read(buf, binary.LittleEndian, &b.compensation.temp)
	if err != nil {
		return fmt.Errorf("failed to convert raw temp compensation data: %v", err)
	}
	buf = bytes.NewReader(pData)
	err = binary.Read(buf, binary.LittleEndian, &b.compensation.press)
	if err != nil {
		return fmt.Errorf("failed to convert raw temp compensation data: %v", err)
	}
	b.compensation.hum.H1 = uint8(h1Data[0])
	b.compensation.hum.H6 = int8(h2Data[6])
	b.compensation.hum.H2 = int16(h2Data[1])<<8 | int16(h2Data[0])
	b.compensation.hum.H3 = uint8(h2Data[2])
	b.compensation.hum.H4 = int16(h2Data[3])<<4 | int16(h2Data[4]&0x0F)
	b.compensation.hum.H5 = int16(h2Data[4]&0xF0)<<4 | int16(h2Data[5])
	return nil
}

// readRaw reads raw data from the sensor.
func (b *BME280) readRaw() (temp, press, hum int32, err error) {
	if b.mode == ForcedMode {
		err = b.dev.WriteReg(CtrlHumAddr, []byte{b.humOverSmpl})
		if err != nil {
			err = fmt.Errorf("failed to write ctrl_hum data to sensor: %v", err)
			return
		}
		err = b.dev.WriteReg(CtrlMeasAddr, []byte{b.tempOverSmpl | b.pressOverSmpl | b.mode})
		if err != nil {
			err = fmt.Errorf("failed to write ctrl_meas data to sensor: %v", err)
			return
		}
		time.Sleep(b.measTime())
	}
	buf := make([]byte, 8)
	err = b.dev.ReadReg(DataAddr, buf)
	if err != nil {
		err = fmt.Errorf("failed to read env data from sensor: %v", err)
		return
	}
	hum = int32(buf[6])<<8 | int32(buf[7])
	temp = int32(buf[3])<<12 | int32(buf[4])<<4 | int32(buf[5])>>4
	press = int32(buf[0])<<12 | int32(buf[1])<<4 | int32(buf[2])>>4
	return
}

// RawData returnes un-compensated data read from the sensor.
func (b *BME280) RawData() (temp, press, hum int32, err error) {
	return b.readRaw()
}

// compensateTemp compensates temperature read from sensor and computes tFine coeficient used in other compensations.
func (b *BME280) compensateTemp(raw int32) (temp float64, tFine int32) {
	tComp := b.compensation.temp
	v1 := (((raw >> 3) - (int32(tComp.T1) << 1)) * int32(tComp.T2)) >> 11
	v2 := (((((raw >> 4) - int32(tComp.T1)) * (raw>>4 - int32(tComp.T1))) >> 12) * int32(tComp.T3)) >> 14
	tFine = v1 + v2
	temp = float64((tFine*5+128)>>8) / 100.0
	switch b.tempUnit {
	case Fahrenheit:
		temp = temp*1.8 + 32
	case Kelvin:
		temp += 273.15
	default:
		temp = temp
	}
	return
}

// compensatePress returns compensated pressure value.
func (b *BME280) compensatePress(raw, tFine int32) float64 {
	pComp := b.compensation.press
	v1 := int64(tFine - 128000)
	v2 := v1 * v1 * int64(pComp.P6)
	v2 += v1 * int64(pComp.P5) << 17
	v2 += int64(pComp.P4) << 35
	v1 = (v1*v1*int64(pComp.P3))>>8 + (v1*int64(pComp.P2))<<12
	v1 = ((int64(1)<<47 + v1) * int64(pComp.P1)) >> 33
	if v1 == 0 {
		return 0
	}
	press := int64(1048576 - raw)
	press = ((press<<31 - v2) * 3125) / v1
	v1 = (int64(pComp.P9) * (press >> 13) * (press >> 13)) >> 25
	v2 = (int64(pComp.P8) * press) >> 19
	press = (press+v1+v2)>>8 + int64(pComp.P7)<<4
	p := float64(press) / 25600.0
	switch b.pressUnit {
	case Bar:
		return p / 1000.0
	case PSI:
		return p * 1.45038
	default:
		return p
	}
}

// compensateHum returns compensated humidity.
func (b *BME280) compensateHum(raw, tFine int32) float64 {
	hComp := b.compensation.hum
	v1 := tFine - 76800
	v1 = (((raw<<14 - int32(hComp.H4)<<20 - (int32(hComp.H5) * v1)) + int32(16384)) >> 15) *
		(((((((v1*int32(hComp.H6))>>10)*(((v1*int32(hComp.H3))>>11)+int32(32768)))>>10)+int32(2097152))*int32(hComp.H2) + 8192) >> 14)
	v1 -= ((((v1 >> 15) * (v1 >> 15)) >> 7) * int32(hComp.H1)) >> 4
	if v1 < 0 {
		return 0
	}
	if v1 > 419430400 {
		return 419430400
	}
	return float64(v1>>12) / 1024.0
}

// measTime returns time.Duration of measurment time, computed based on the formula in datasheet.
func (b *BME280) measTime() time.Duration {
	var overSmlpMap = map[byte]float64{
		OverSmpl1:  1.0,
		OverSmpl2:  2.0,
		OverSmpl4:  4.0,
		OverSmpl8:  8.0,
		OverSmpl16: 16.0,
	}
	t := 1.25 + 2*overSmlpMap[b.tempOverSmpl>>5] + 2.3*overSmlpMap[b.pressOverSmpl>>3] + 0.575 + 2.3*overSmlpMap[b.humOverSmpl] + 0.575
	return time.Microsecond * time.Duration(t*1000)
}

// EnvData returns compensated temperature, pressure and humidity.
func (b *BME280) EnvData() (temp, press, hum float64, err error) {
	t, p, h, err := b.readRaw()
	if err != nil {
		err = fmt.Errorf("failed to read raw data: %v", err)
		return
	}
	var tFine int32
	temp, tFine = b.compensateTemp(t)
	press = b.compensatePress(p, tFine)
	hum = b.compensateHum(h, tFine)
	return
}

// SetTempUnit sets temperature unit.
func (b *BME280) SetTempUnit(unit int) {
	b.tempUnit = unit
}

// SetPressUnit sets pressure unit.
func (b *BME280) SetPressUnit(unit int) {
	b.pressUnit = unit
}

// Temp returns temperature in set units.
func (b *BME280) Temp() (float64, error) {
	t, _, _, err := b.readRaw()
	if err != nil {
		return 0.0, fmt.Errorf("failed to read raw data: %v", err)
	}
	temp, _ := b.compensateTemp(t)
	return temp, nil
}

// Press returns pressure in set units.
func (b *BME280) Press() (float64, error) {
	t, p, _, err := b.readRaw()
	if err != nil {
		return 0.0, fmt.Errorf("failed to read raw data: %v", err)
	}
	_, tFine := b.compensateTemp(t)
	return b.compensatePress(p, tFine), nil
}

// Hum returns humidity
func (b *BME280) Hum() (float64, error) {
	t, _, h, err := b.readRaw()
	if err != nil {
		return 0.0, fmt.Errorf("failed to read raw data: %v", err)
	}
	_, tFine := b.compensateTemp(t)
	return b.compensateHum(h, tFine), nil
}
