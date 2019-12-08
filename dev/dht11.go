/*
DHT11 is an sensor for getting temperature and humidity.
config:
1. sudo vim /boot/config.txt
2. add following line to the end of config.txt
	--------------------------
	dtoverlay=dht11,gpiopin=4
	--------------------------
3. connect dht11 to raspberry pi:
	SIGNAL: pin 4
	GND:	any gnd pin
	VCC:	3.3v

-----------------------------------------------------------------------

      +-------------+
      |             |
      |    DHT11    |
      |             |
      +-+----+----+-+
        |    |    |
      S |   VCC   | -
        |    |    |
        |    |    |              +-----------+
        |    +----|--------------+ * 1   2 o |
        |         |              | * 3     o |
        |         |              | o       o |
        +---------|--------------+ * 7     o |
                  +--------------+ * 9     o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o       o |
                                 | o 39 40 o |
								 +-----------+

-----------------------------------------------------------------------
*/
package dev

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	tFile = "/sys/bus/iio/devices/iio:device0/in_temp_input"
	hFile = "/sys/bus/iio/devices/iio:device0/in_humidityrelative_input"
)

const maxRetry = 50

// DHT11 ...
type DHT11 struct {
}

// NewDHT11 ...
func NewDHT11() *DHT11 {
	return &DHT11{}
}

// TempHumidity ...
func (d *DHT11) TempHumidity() (float64, float64, error) {
	var (
		t, h                 float64
		gotTemp, gotHumidity bool
	)

	for i := 0; i < maxRetry; i++ {
		if !gotTemp {
			if data, err := ioutil.ReadFile(tFile); err == nil {
				if t, err = d.parseData(data); err == nil {
					gotTemp = true
				}
			}
		}

		if !gotHumidity {
			if data, err := ioutil.ReadFile(hFile); err == nil {
				if h, err = d.parseData(data); err == nil {
					gotHumidity = true
				}
			}
		}
		if gotTemp && gotHumidity {
			return t, h, nil
		}
		time.Sleep(1 * time.Second)
	}

	return 0, 0, errors.New("bad data")
}

func (d *DHT11) parseData(data []byte) (float64, error) {
	s := strings.Trim(string(data), " \t\n")
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v / 1000.0, nil
}
