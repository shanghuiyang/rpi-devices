/*
Auto-Fan let you the fan working with a relay and a temperature sensor together.
The temperature sensor will trigger the relay to control the fan running or stopping.

temperature sensor:
 - vcc: phys.1/3.3v
 - dat: phys.7/BCM.4
 - gnd: phys.9/GND

realy:
 - vcc: phys.2/5v
 - in:  phys.26/BCM.7
 - gnd: phys.34/GND
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
                                    | o       o |  +--------gnd-+   /          |
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
	"github.com/stianeikeland/go-rpio"
)

const (
	relayPin           = 7
	intervalTime       = 1 * time.Minute
	triggerTemperature = 27.3
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	t := dev.NewTemperature()
	if t == nil {
		log.Printf("failed to new a temperature device")
		return
	}

	r := dev.NewRelay(relayPin)
	if r == nil {
		log.Printf("failed to new a relay device")
		return
	}

	f := &autoFan{
		temperature: t,
		relay:       r,
	}
	f.start()
}

type autoFan struct {
	temperature *dev.Temperature
	relay       *dev.Relay
}

func (f *autoFan) start() {
	for {
		time.Sleep(intervalTime)
		c, err := f.temperature.GetTemperature()
		if err != nil {
			log.Printf("failed to get temperature, error: %v", err)
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
