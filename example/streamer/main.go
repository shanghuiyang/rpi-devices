// +build gocv

package main

import (
	"log"

	"github.com/shanghuiyang/rpi-devices/util"
	"gocv.io/x/gocv"
)

const (
	url = ":8088/video"
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

	streamer, err := util.NewStreamer(url)
	if err != nil {
		log.Printf("failed to create streamer, error: %v", err)
		return
	}
	img := gocv.NewMat()
	defer img.Close()
	for {
		cam.Grab(10)
		if !cam.Read(&img) {
			log.Printf("failed to get img from camera")
			continue
		}
		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("failed to encode image, err: %v", err)
			continue
		}
		streamer.Push(buf)
	}
}
