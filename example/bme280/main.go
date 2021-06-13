// i2c is example usage of BME280 library with I2C bus.
package main

import (
	"fmt"

	"github.com/jakefau/rpi-devices/dev"

	"golang.org/x/exp/io/i2c"
)

func main() {

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, 0x77)
	if err != nil {
		panic(err)
	}
	b := dev.New(d)
	err = b.Init()

	t, p, h, err := b.EnvData()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Temp: %fF, Press: %f, Hum: %f%%\n", toFahrenheit(t), toMercury(p), h)
}

func toFahrenheit(c float64) float64 {
	return (c * 9 / 5.0) + 32
}

func toMercury(m float64) float64 {
	return m * 0.02953
}
