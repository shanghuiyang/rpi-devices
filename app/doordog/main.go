/*
Doordog helps you watch your doors.
When somebody entries your room, you will be alerted by a beeping buzzer and a blinking led.
*/

package main

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/shanghuiyang/face"
	"github.com/shanghuiyang/oauth"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	pinTrig = 2
	pinEcho = 3
	pinBtn  = 7
	pinBzr  = 17
	pinLed  = 23

	ifttAPI = "your-iftt-api"
)

const (
	// the time of keeping alert in second
	alertTime = 60
	// the distance of triggering alert in cm
	alertDist = 80

	groupID = "mygroup"

	baiduFaceRecognitionAppKey    = "your_face_app_key"
	baiduFaceRecognitionSecretKey = "your_face_secret_key"
)

var (
	allowlist = []string{
		"p1",
		"p2",
		"p3",
		"p4",
		"p5",
	}
)

func main() {
	cam := dev.NewMotionCamera()
	bzr := dev.NewBuzzerImp(pinBzr, dev.High)
	led := dev.NewLedImp(pinLed)
	btn := dev.NewButtonImp(pinBtn)
	hcsr04 := dev.NewHCSR04(pinTrig, pinEcho)
	if hcsr04 == nil {
		log.Printf("[doordog]failed to new a HCSR04")
		return
	}

	auth := oauth.NewBaiduOauth(baiduFaceRecognitionAppKey, baiduFaceRecognitionSecretKey, oauth.NewCacheImp())
	f := face.NewBaiduFace(auth, groupID)

	dog := newDoordog(cam, hcsr04, bzr, led, btn, f)
	util.WaitQuit(func() {
		dog.stop()
	})
	dog.start()
}

type doordog struct {
	cam      dev.Camera
	dmeter   dev.DistanceMeter
	buzzer   dev.Buzzer
	led      dev.Led
	button   dev.Button
	face     face.Face
	alerting bool
	chAlert  chan bool
}

func newDoordog(cam dev.Camera, d dev.DistanceMeter, buzzer dev.Buzzer, led dev.Led, btn dev.Button, f face.Face) *doordog {
	return &doordog{
		cam:      cam,
		dmeter:   d,
		buzzer:   buzzer,
		led:      led,
		button:   btn,
		face:     f,
		alerting: false,
		chAlert:  make(chan bool, 4),
	}
}

func (d *doordog) start() {
	log.Printf("[doordog]start to service")
	go d.alert()
	go d.stopAlert()
	d.detect()

}

func (d *doordog) detect() {
	// need to warm-up the ultrasonic distance meter first
	d.dmeter.Dist()
	time.Sleep(500 * time.Millisecond)

	t := 300 * time.Millisecond
	for {
		time.Sleep(t)
		dist, err := d.dmeter.Dist()
		for i := 0; err != nil && i < 3; i++ {
			util.DelayMs(100)
			dist, err = d.dmeter.Dist()
		}
		if err != nil {
			continue
		}
		if dist < 10 {
			log.Printf("[doordog]bad data from distant meter, distance = %.2fcm", dist)
			continue
		}

		detected := (dist < alertDist)
		if detected {
			log.Printf("[doordog]detected objects, distance = %.2fcm", dist)
			who, err := d.RecoginzeFace()
			if err != nil {
				continue
			}

			log.Printf("[doordog]it is %v", who)
			if allowed(who) {
				continue
			}

			d.chAlert <- detected
			continue
		}
	}
}

func (d *doordog) alert() {
	trigTime := time.Now()
	go func() {
		for {
			if d.alerting {
				go d.buzzer.Beep(1, 200)
				go d.led.Blink(1, 200)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for detected := range d.chAlert {
		if detected {
			go ifttt()
			d.alerting = true
			trigTime = time.Now()
			continue
		}
		timeout := time.Since(trigTime).Seconds() > alertTime
		if timeout && d.alerting {
			log.Printf("[doordog]timeout, stop alert")
			d.alerting = false
		}
	}
}

func (d *doordog) RecoginzeFace() (string, error) {
	img, err := d.cam.Photo()
	if err != nil {
		log.Printf("[doordog]failed to take phote, error: %v", err)
		return "unknow", err
	}

	name, err := d.face.Recognize(img)
	if err != nil {
		log.Printf("[doordog]failed to recognize the image, error: %v", err)
		return "unknow", err
	}

	return name, nil
}

func (d *doordog) stopAlert() {
	for {
		pressed := d.button.Pressed()
		if pressed {
			log.Printf("[doordog]the button was pressed")
			if d.alerting {
				d.alerting = false
			}
			// make a dalay detecting
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func ifttt() {
	resp, err := http.PostForm(ifttAPI, url.Values{})
	if err != nil {
		log.Printf("failed to request to ifttt, error: %v", err)
		return
	}
	defer resp.Body.Close()
}

func (d *doordog) stop() {
	d.buzzer.Off()
	d.led.Off()
}

func allowed(user string) bool {
	for _, u := range allowlist {
		if u == user {
			return true
		}
	}
	return false
}
