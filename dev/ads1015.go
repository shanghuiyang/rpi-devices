/*
ADS1015 is a 12-bit analog-digital converter. It is a implement of ADC interface.
More details about ADS1015, please ref to:
https://www.adafruit.com/product/1083
https://cdn-learn.adafruit.com/downloads/pdf/adafruit-4-channel-adc-breakouts.pdf

Config Raspberry Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 Interface Options] -> [P5 I2C] -> [yes] -> [ok]
5. $ sudo reboot now
6. check: $ sudo i2cdetect -y 1
	it works if you saw following message:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- 48 -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	70: -- -- -- -- -- -- -- --
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Connect to Raspberry Pi:
- VCC: any 3.3v pin
- GND: any GND pin
- SDA: GPIO-2 (SDA)
- SCL: GPIO-3 (SCL)

*/
package dev

import (
	"errors"

	"golang.org/x/exp/io/i2c"
)

const (
	ads1015Dev  = "/dev/i2c-1"
	ads1015Addr = 0x48
)

const (
	// ConversionRegiserPointer ...
	ConversionRegiserPointer byte = 0x00
	// ConfigRegiserPointer ...
	ConfigRegiserPointer byte = 0x01
	//LoThreshRegiserPointer ...
	LoThreshRegiserPointer byte = 0x10
	// HiThreshRegiserPointer ...
	HiThreshRegiserPointer byte = 0x11

	// ComparatorQueueAssertAfterOne Assert after one conversion
	ComparatorQueueAssertAfterOne uint16 = 0x0000
	// ComparatorQueueAssertAfterTwo Assert after two conversions
	ComparatorQueueAssertAfterTwo uint16 = 0x0001
	// ComparatorQueueAssertAfterFour Assert after four conversions
	ComparatorQueueAssertAfterFour uint16 = 0x0002
	// ComparatorQueueDisable Disable comparator and set ALERT/RDY pin to high-impedance (default)
	ComparatorQueueDisable uint16 = 0x0003

	// LatchingComparatorLatching The ALERT/RDY pin does not latch when asserted (default)
	LatchingComparatorLatching uint16 = 0x0000
	// LatchingComparatorNonLatching The asserted ALERT/RDY pin remains latched until
	LatchingComparatorNonLatching uint16 = 0x0004

	// ComparatorPolarityActiveLow This bit controls the polarity of the ALERT/RDY pin (default)
	ComparatorPolarityActiveLow uint16 = 0x0000
	// ComparatorPolarityActiveHigh This bit controls the polarity of the ALERT/RDY pin
	ComparatorPolarityActiveHigh uint16 = 0x0008

	// ComparatorModeTraditional this bit configures the comparator operating mode. (default)
	ComparatorModeTraditional uint16 = 0x0000
	// ComparatorModeWindow this bit configures the comparator operating mode.
	ComparatorModeWindow uint16 = 0x0010

	// // OperationalStatus determines the operational status of the device. OS can only be written
	// // when in power-down state and has no effect when a conversion is ongoing
	// OperationalStatus uint16 = 0x8000

	// // RegisterPointerConfig ...
	// RegisterPointerConfig byte = 0x01
	// // RegisterConversionConfig Conversion register contains the result of the last conversion in binary two's complement format.
	// RegisterConversionConfig byte = 0x00

	// DataRate128 control the data rate setting. 128 Sample Per Seconds
	DataRate128 uint16 = 0x0000
	// DataRate250 control the data rate setting. 250 Sample Per Seconds
	DataRate250 uint16 = 0x0020
	// DataRate490 control the data rate setting. 490 Sample Per Seconds
	DataRate490 uint16 = 0x0040
	// DataRate920 control the data rate setting. 64 Sample Per Seconds
	DataRate920 uint16 = 0x0060
	// DataRate1600  control the data rate setting. 128 Sample Per Seconds
	DataRate1600 uint16 = 0x0080
	// DataRate2400 control the data rate setting. 250  Sample Per Seconds
	DataRate2400 uint16 = 0x00A0
	// DataRate3300_0 control the data rate setting. 475 Sample Per Seconds
	DataRate3300_0 uint16 = 0x00C0
	// DataRate3300_1 control the data rate setting. 475 Sample Per Seconds
	DataRate3300_1 uint16 = 0x00E0

	// DeviceOperationModeContinous Continuous-conversion mode
	DeviceOperationModeContinous uint16 = 0x0000
	// DeviceOperationModeSingleShot  Single-shot mode or power-down state
	DeviceOperationModeSingleShot uint16 = 0x0100

	// ProgramableGainAmplifier6144 These bits set the FSR of the programmable gain amplifier. For voltages in the range ±6.144
	ProgramableGainAmplifier6144 uint16 = 0x0000
	// ProgramableGainAmplifier4096 set the FSR of the programmable gain amplifier. For voltages in the range ±4.096
	ProgramableGainAmplifier4096 uint16 = 0x0200
	// ProgramableGainAmplifier2048 set the FSR of the programmable gain amplifier. For voltages in the range ±2.048
	ProgramableGainAmplifier2048 uint16 = 0x0400
	// ProgramableGainAmplifier1024 set the FSR of the programmable gain amplifier. For voltages in the range ±1.024
	ProgramableGainAmplifier1024 uint16 = 0x0600
	// ProgramableGainAmplifier0512 set the FSR of the programmable gain amplifier. For voltages in the range ±0.512
	ProgramableGainAmplifier0512 uint16 = 0x0800
	// ProgramableGainAmplifier0256_0 set the FSR of the programmable gain amplifier. For voltages in the range ±0.256
	ProgramableGainAmplifier0256_0 uint16 = 0x0A00
	// ProgramableGainAmplifier0256_1 set the FSR of the programmable gain amplifier. For voltages in the range ±0.256
	ProgramableGainAmplifier0256_1 uint16 = 0x0C00
	// ProgramableGainAmplifier0256_2 set the FSR of the programmable gain amplifier. For voltages in the range ±0.256
	ProgramableGainAmplifier0256_2 uint16 = 0x0E00

	// MultiplexerConfigurationAIN0 AINP = AIN0 and AINN = GND
	MultiplexerConfigurationAIN0 uint16 = 0x4000
	// MultiplexerConfigurationAIN1 AINP = AIN1 and AINN = GND
	MultiplexerConfigurationAIN1 uint16 = 0x5000
	// MultiplexerConfigurationAIN2 AIN2 = AIN2 and AINN = GND
	MultiplexerConfigurationAIN2 uint16 = 0x6000
	// MultiplexerConfigurationAIN3 AIN3 = AIN3 and AINN = GND
	MultiplexerConfigurationAIN3 uint16 = 0x7000
)

var (
	channelMuxConfig = map[int]uint16{
		0: MultiplexerConfigurationAIN0,
		1: MultiplexerConfigurationAIN1,
		2: MultiplexerConfigurationAIN2,
		3: MultiplexerConfigurationAIN3,
	}

	defaultConfig = ComparatorQueueDisable | LatchingComparatorLatching | ComparatorPolarityActiveLow | ComparatorModeTraditional | DataRate3300_0 | DeviceOperationModeContinous | ProgramableGainAmplifier6144
)

// ADS1015 is a 12-bit analog-digital converter. It implements ADC interface
type ADS1015 struct {
	dev    *i2c.Device
	config uint16
}

// NewADS1015 create a driver for ADS1015 module
func NewADS1015() (*ADS1015, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: ads1015Dev}, ads1015Addr)
	if err != nil {
		return nil, err
	}
	return &ADS1015{
		dev:    dev,
		config: defaultConfig,
	}, nil
}

// SetConfig ...
func (m *ADS1015) SetConfig(config uint16) {
	m.config = config
}

// Read ...
func (m *ADS1015) Read(channel int) (float64, error) {
	mux, ok := channelMuxConfig[channel]
	if !ok {
		return 0, errors.New("invalid channel number, should be 0~3")
	}

	conf := m.config | mux
	hiByte := byte(conf >> 8)
	loByte := byte(conf & 0x00FF)
	if err := m.dev.WriteReg(ConfigRegiserPointer, []byte{hiByte, loByte}); err != nil {
		return 0, err
	}

	delayUs(100)
	data := make([]byte, 2)
	if err := m.dev.ReadReg(ConversionRegiserPointer, data); err != nil {
		return 0, err
	}

	val := (uint32(data[0]) << 8) | uint32(data[1])
	v := float64(val*6144/1000) / 32768.0
	return v, nil

}

// Close ...
func (m *ADS1015) Close() error {
	return m.dev.Close()
}
