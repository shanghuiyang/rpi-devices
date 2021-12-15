/*
DS18B20 is sensor used to meature temperature.

Config Raspberry Pi:
1. $ sudo vim /boot/config.txt
2. add following line at the end of the file
	~~~~~~~~~~~~~~~~~~~~~~~~~~~
	dtoverlay=w1-gpio,gpiopin=4
	~~~~~~~~~~~~~~~~~~~~~~~~~~~
3. $ sudo reboot now
4. check: $ cat /sys/bus/w1/devices/28-d8baf71d64ff/w1_slave
	should saw:
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	ca 01 55 00 7f ff 0c 10 bf : crc=bf YES
	ca 01 55 00 7f ff 0c 10 bf t=28625
	~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Connect to Raspberry Pi:
 - vcc: any 3.3v pin
 - gnd: any gnd pin
 - dat: must connect to GPIO-4

*/
package dev

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	tempFile = "/sys/bus/w1/devices/28-d8baf71d64ff/w1_slave"
)

// DS18B20 implements Thermometer interface
type DS18B20 struct {
}

// NewDS18B20 ...
func NewDS18B20() *DS18B20 {
	return &DS18B20{}
}

// GetTemperature ...
// temperature file:
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
// ca 01 55 00 7f ff 0c 10 bf : crc=bf YES
// ca 01 55 00 7f ff 0c 10 bf t=28625
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~^^^^^^^~~~~~~~~
func (d *DS18B20) Temperature() (float64, error) {
	data, err := ioutil.ReadFile(tempFile)
	if err != nil {
		return 0, err
	}
	raw := string(data)

	idx := strings.LastIndex(raw, "t=")
	if idx < 0 {
		return 0, fmt.Errorf("can't find 't='")
	}
	t, err := strconv.ParseFloat(raw[idx+2:idx+7], 32)
	if err != nil {
		return 0, fmt.Errorf("bad data")
	}
	return float64(t / 1000), nil
}
