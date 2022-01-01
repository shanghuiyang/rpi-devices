/*
mpu6050 is a 3-axis gyroscope and a 3-axis accelerometer.

Config Your Pi:
1. $ sudo apt-get install -y python-smbus
2. $ sudo apt-get install -y i2c-tools
3. $ sudo raspi-config
4. 	-> [5 Interface Options] -> [P5 i2c] -> [yes] -> [ok]
5. $ sudo reboot now
6. check: $ sudo i2cdetect -y 1
	it works if you saw following message:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
	00:          -- -- -- -- -- -- -- -- -- -- -- -- --
	10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
	60: -- -- -- -- -- -- -- -- 68 -- -- -- -- -- -- --
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
	"golang.org/x/exp/io/i2c"
)

const (
	mpu6050Dev   = "/dev/i2c-1"
	mpu6050Addr  = 0x68
	accRegister  = 0x3B
	gyroRegister = 0x43
)

// MPU6050 ...
type MPU6050 struct {
	dev *i2c.Device
}

// NewMPU6050 ...
func NewMPU6050() (*MPU6050, error) {
	dev, err := i2c.Open(&i2c.Devfs{Dev: mpu6050Dev}, mpu6050Addr)
	if err != nil {
		return nil, err
	}
	dev.WriteReg(0x6B, []uint8{0}) // power on
	return &MPU6050{
		dev: dev,
	}, nil
}

// GetAcc ...
func (m *MPU6050) GetAcc() (x int32, y int32, z int32) {
	data := make([]byte, 6)
	m.dev.ReadReg(accRegister, data)
	x = int32(int16((uint16(data[0])<<8)|uint16(data[1]))) * 15625 / 256
	y = int32(int16((uint16(data[2])<<8)|uint16(data[3]))) * 15625 / 256
	z = int32(int16((uint16(data[4])<<8)|uint16(data[5]))) * 15625 / 256
	return
}

// GetRot ...
func (m *MPU6050) GetRot() (x int32, y int32, z int32) {
	data := make([]byte, 6)
	m.dev.ReadReg(gyroRegister, data)
	x = int32(int16((uint16(data[0])<<8)|uint16(data[1]))) * 15625 / 2048 * 1000
	y = int32(int16((uint16(data[2])<<8)|uint16(data[3]))) * 15625 / 2048 * 1000
	z = int32(int16((uint16(data[4])<<8)|uint16(data[5]))) * 15625 / 2048 * 1000
	return
}

// Close ...
func (m *MPU6050) Close() error {
	return m.dev.Close()
}
