// +build gocv

package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/cv"
	"gocv.io/x/gocv"
)

const (
	host = "0.0.0.0:8088"
)

func main() {
	cam, err := gocv.OpenVideoCapture(0)
	if err != nil {
		log.Printf("failed to open video")
		return
	}
	defer cam.Close()

	cam.Set(gocv.VideoCaptureFocus, cam.ToCodec("MJPG"))
	cam.Set(gocv.VideoCaptureFPS, 100)
	cam.Set(gocv.VideoCaptureFrameWidth, 640)
	cam.Set(gocv.VideoCaptureFrameHeight, 480)

	streamer := cv.NewStreamer(host)
	defer streamer.Close()
	go streamer.Start()

	img := gocv.NewMat()
	defer img.Close()
	for {
		cam.Grab(10)
		if !cam.Read(&img) {
			log.Printf("[car]failed to read img from camera")
			continue
		}
		streamer.SetImage(&img)
	}
}
