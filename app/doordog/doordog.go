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

	"github.com/shanghuiyang/face-recognizer/face"
	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/rpi-devices/base"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/stianeikeland/go-rpio"
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
	if err := rpio.Open(); err != nil {
		log.Fatalf("[doordog]failed to open rpio, error: %v", err)
		return
	}
	defer rpio.Close()

	cam := dev.NewCamera()
	bzr := dev.NewBuzzer(pinBzr)
	led := dev.NewLed(pinLed)
	btn := dev.NewButton(pinBtn)
	dist := dev.NewHCSR04(pinTrig, pinEcho)
	if dist == nil {
		log.Printf("[doordog]failed to new a HCSR04")
		return
	}

	auth := oauth.New(baiduFaceRecognitionAppKey, baiduFaceRecognitionSecretKey, oauth.NewCacheMan())
	f := face.New(auth)

	dog := newDoordog(cam, dist, bzr, led, btn, f)
	base.WaitQuit(func() {
		dog.stop()
		rpio.Close()
	})
	dog.start()
}

type doordog struct {
	cam      *dev.Camera
	dist     *dev.HCSR04
	buzzer   *dev.Buzzer
	led      *dev.Led
	button   *dev.Button
	face     *face.Face
	alerting bool
	chAlert  chan bool
}

func newDoordog(cam *dev.Camera, dist *dev.HCSR04, buzzer *dev.Buzzer, led *dev.Led, btn *dev.Button, face *face.Face) *doordog {
	return &doordog{
		cam:      cam,
		dist:     dist,
		buzzer:   buzzer,
		led:      led,
		button:   btn,
		face:     face,
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
	d.dist.Dist()
	time.Sleep(500 * time.Millisecond)

	t := 300 * time.Millisecond
	for {
		time.Sleep(t)
		dist := d.dist.Dist()
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
		timeout := time.Now().Sub(trigTime).Seconds() > alertTime
		if timeout && d.alerting {
			log.Printf("[doordog]timeout, stop alert")
			d.alerting = false
		}
	}
}

func (d *doordog) RecoginzeFace() (name string, err error) {
	imgf, e := d.cam.TakePhoto()
	if e != nil {
		log.Printf("[doordog]failed to take phote, error: %v", e)
		name, err = "unknow", e
		return
	}

	users, e := d.face.Recognize(imgf, groupID)
	if e != nil {
		log.Printf("[doordog]failed to recognize the image, error: %v", e)
		name, err = "unknow", e
		return
	}

	if len(users) == 0 {
		name, err = "unknow", nil
		return
	}

	log.Printf("who: %v", *(users[0]))
	if users[0].Score > 50 {
		return users[0].UserID, nil
	}

	name, err = "unknow", nil
	return
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
