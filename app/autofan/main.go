/*
Auto-Fan let you the fan working with a single-channel relay and a temperature sensor together.
The temperature sensor will trigger the relay to control the fan running or stopping.

temperature sensor:
 - vcc: pin 1 or any 3.3v pin
 - dat: pin 7(gpio 4)
 - gnd: pin 9 or any gnd pin

realy:
 - vcc: any 5v pin
 - in:  gpio 7
 - gnd: any gnd pin
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

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	relayPin           = 7
	intervalTime       = 1 * time.Minute
	triggerTemperature = 27.3
)

func main() {
	ds18b20 := dev.NewDS18B20()
	if ds18b20 == nil {
		log.Printf("[autofan]failed to new ds18b20 sensor")
		return
	}

	r := dev.NewRelayImp(relayPin)
	if r == nil {
		log.Printf("[autofan]failed to new a relay")
		return
	}

	f := &autoFan{
		tmeter: ds18b20,
		relay:  r,
	}
	util.WaitQuit(func() {
		f.off()
	})
	f.start()
}

type autoFan struct {
	tmeter dev.Thermometer
	relay  dev.Relay
}

func (f *autoFan) start() {
	for {
		time.Sleep(intervalTime)
		c, err := f.tmeter.Temperature()
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
