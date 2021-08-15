package main

import (
	"fmt"
	"log"

	"github.com/shanghuiyang/face-recognizer/face"
	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	groupID = "mygroup"

	// replace your_app_key and your_secret_key with yours
	appKey    = "your_app_key"
	secretKey = "your_secret_key"
)

func main() {
	cam := dev.NewMotionCamera()
	if cam == nil {
		log.Print("failed to new a camera")
		return
	}

	var input string
	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	f := face.New(auth)
	for {
		fmt.Printf(">>press Enter to go: ")
		if _, err := fmt.Scanln(); err != nil {
			log.Print("please press [enter]")
			fmt.Scanln(&input) // discard current inputs
			continue
		}

		img, err := cam.Photo()
		if err != nil {
			log.Printf("failed to take phote, error: %v", err)
			continue
		}

		users, err := f.Recognize(img, groupID)
		if err != nil {
			log.Printf("failed to recognize the image, error: %v", err)
			continue
		}
		for _, u := range users {
			fmt.Println(u)
		}
	}
}
