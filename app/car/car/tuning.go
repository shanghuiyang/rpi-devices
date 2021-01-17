package car

import (
	"fmt"
	"log"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

// TuningTurnAngle tunings the mapping between angle(degree) and time(millisecond)
func TuningTurnAngle(eng *dev.L298N) {
	if eng == nil {
		log.Fatal("eng is nil")
		return
	}
	for {
		var ms int
		fmt.Printf(">>ms: ")
		if n, err := fmt.Scanf("%d", &ms); n != 1 || err != nil {
			log.Printf("[carapp]invalid operator, error: %v", err)
			continue
		}
		if ms < 0 {
			break
		}
		eng.Right()
		time.Sleep(time.Duration(ms) * time.Millisecond)
		eng.Stop()
	}
	return
}

// TuningEncoder tunings the mapping between angle(degree) and count
func TuningEncoder(eng *dev.L298N, encoder *dev.Encoder) {
	if eng == nil {
		log.Fatal("engineer is nil")
		return
	}
	if encoder == nil {
		log.Fatal("encoder is nil")
		return
	}
	eng.Speed(30)
	for {
		var count int
		fmt.Printf(">>count: ")
		if n, err := fmt.Scanf("%d", &count); n != 1 || err != nil {
			log.Printf("[carapp]invalid count, error: %v", err)
			continue
		}
		if count == 0 {
			break
		}
		if count < 0 {
			eng.Left()
			count *= -1
		} else {
			eng.Right()
		}

		encoder.Start()
		for i := 0; i < count; {
			i += encoder.Count1()
		}
		eng.Stop()
		encoder.Stop()
	}
	eng.Stop()
	return
}
