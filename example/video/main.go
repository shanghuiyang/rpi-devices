// +build gocv

package main

import (
	"log"
	"os"

	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/cv"
	"github.com/stianeikeland/go-rpio"
	"gocv.io/x/gocv"
)

const (
	host = "0.0.0.0:8088"
)

var stream *cv.Stream

func main() {
	if err := rpio.Open(); err != nil {
		log.Fatalf("[tracking]failed to open rpio, error: %v", err)
		os.Exit(1)
	}
	defer rpio.Close()

	util.WaitQuit(func() {
		rpio.Close()
	})

	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Printf("failed to open video")
		return
	}
	defer cam.Close()

	stream = cv.NewStream(cam, host)
	stream.Start()

	os.Exit(0)
}
