/*
trigger sends a command to the server when it detects a shake.

*/

package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
)

type state string

const (
	sw420Pin = 2
)

const (
	on  state = "on"
	off state = "off"
)

var (
	api      = "http://192.168.31.50:8080"
	curState = off
)

func main() {
	sw420 := dev.NewSW420(sw420Pin)
	if sw420 == nil {
		log.Printf("[autoairout]failed to new a sw420 sensor")
		return
	}

	for {
		shaked := isKeepShaking(sw420)
		if shaked && curState == off {
			curState = on
			log.Printf("[autoairout]state: on")
			go sendcmd(curState)
			time.Sleep(1 * time.Second)
			continue
		}
		if !shaked && curState == on {
			curState = off
			log.Printf("[autoairout]state: off")
			go sendcmd(curState)
			time.Sleep(1 * time.Second)
			continue
		}

	}
}

func sendcmd(s state) {
	formData := url.Values{
		"op": {string(s)},
	}
	resp, err := http.PostForm(api, formData)
	if err != nil {
		log.Printf("[autoairout]failed to send command to server, error: %v", err)
		return
	}
	defer resp.Body.Close()
}

// isKeepShaking returns true if the sensor detects the object keeps shaking in 100 millisecond,
// or returns false
func isKeepShaking(s *dev.SW420) bool {
	states := map[bool]int{
		true:  0,
		false: 0,
	}
	for i := 0; i < 10; i++ {
		shaked := s.Detected()
		states[shaked]++
		time.Sleep(10 * time.Millisecond)
	}
	return states[true] > states[false]
}
