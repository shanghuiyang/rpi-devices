/*
trigger detect the shaking.

*/

package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
)

const (
	sw420Pin = 7
)

var (
	api = "http://192.168.31.64:8080"
	op  = map[bool]string{
		true:  "on",
		false: "off",
	}
)

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	sw420 := dev.NewSW420(sw420Pin)
	if sw420 == nil {
		log.Printf("failed to new a sw420 sensor")
		return
	}

	base.WaitQuit(func() {
		rpio.Close()
	})

	for {
		shaked := sw420.Shaked()
		go sendcmd(shaked)

		sleepTime := 5 * time.Second
		if shaked {
			sleepTime = 60 * time.Second
		}
		time.Sleep(sleepTime)
	}
}

func sendcmd(shaked bool) {
	formData := url.Values{
		"op": {op[shaked]},
	}
	resp, err := http.PostForm(api, formData)
	if err != nil {
		log.Printf("failed to send command to server, error: %v", err)
		return
	}
	defer resp.Body.Close()
	return
}
