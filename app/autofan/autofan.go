/*
Auto-Fan let you the fan working with a relay and a temperature sensor together.
The temperature sensor will trigger the relay to control the fan running or stopping.

temperature sensor:
 - vcc: pin 1 or any 3.3v pin
 - dat: pin 7(gpio 4)
 - gnd: pin 9 or any gnd pin

realy:
 - vcc: pin 2 or any 5v pin
 - in:  pin 26(gpio 7)
 - gnd: pin 34 or any gnd pin
 - on:  fan(+)
 - com: bettery(+)

----------------------------------------------------------------------------------------------------------------------

          o---------o
          |         |
          |temperature
          |sensor   |                                           o---------------o
          |         |               +-----------+               |fan            |
          |         |     +---------+ * 1   2 * +---------+     |    \_   _/    |
          o-+--+--+-o     |         | o       o |         |     |      \ /      |
            |  |  |       |  +------+ * 4     o |         |     |   ----*----   |    +----------+    +----------+
          gnd dat vcc     |  |  +---+ * 9     o |         |     |     _/ \_     |    |          |    |          |
            |  |  +-------+  |  |   | o       o |         |     |    /     \    |    |          |    |          |
            |  |             |  |   | o       o |         |     |               |    |       o-----------o      |
            |  +-------------+  |   | o       o |         |     o-----+---+-----o    |       |  -    +   |      |
            |                   |   | o       o |         |           |   |          |       |           |      |
            +-------------------+   | o       o |         |           |   |          |       |   power   |      |
                                    | o       o |         |           |   +----------+       o-----------o      |
                                    | o       o |         |           |                                         |
                                    | o       o |         |           |    +------------------------------------+
                                    | o    26 * +------+  |           NO   COM
                                    | o       o |      |  |           |    |
                                    | o       o |      |  |     o-----+----+---o
                                    | o       o |      |  +-vcc-+        relay |
                                    | o    34 * +--+   |        |    /         |
                                    | o       o |  +---|----gnd-+   /          |
                                    | o       o |      |        |  *   -----*  |
                                    | o 39 40 o |      +-----in-+              |
                                    +-----------+               o--------------o

----------------------------------------------------------------------------------------------------------------------

*/

package main

import (
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	relayPin           = 7
	intervalTime       = 1 * time.Minute
	triggerTemperature = 27.3
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[autofan]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	temp := dev.NewDS18B20()
	if temp == nil {
		log.Printf("[autofan]failed to new a temperature sensor")
		return
	}

	r := dev.NewRelay(relayPin)
	if r == nil {
		log.Printf("[autofan]failed to new a relay")
		return
	}

	f := &autoFan{
		temp:  temp,
		relay: r,
	}
	base.WaitQuit(func() {
		f.off()
		rpio.Close()
	})
	f.start()
}

type autoFan struct {
	temp  *dev.DS18B20
	relay *dev.Relay
}

func (f *autoFan) start() {
	for {
		time.Sleep(intervalTime)
		c, err := f.temp.GetTemperature()
		if err != nil {
			log.Printf("[autofan]failed to get temperature, error: %v", err)
			continue
		}
		if c >= triggerTemperature {
			f.on()
		} else {
			f.off()
		}
	}
}

func (f *autoFan) on() {
	f.relay.On()
}

func (f *autoFan) off() {
	f.relay.Off()
}
