package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	pinSGH = 12
	pinSGV = 13
	pinLed = 21
	pinBzr = 11
	pinBtn = 4
)

const (
	ipPattern        = "((000.000.000.000))"
	pageOutOfService = `
		<!DOCTYPE html>
		<html>
			<title>Video Monitor</title>
			<h1 style="color:red;font-size:50px;">
				<span style="font-size:100px;">
        			&ensp;&#129318;&#129318;&#129318;<br>
    			</span>
				~~~~~~~~~~~~~~~<br>
				&nbsp;&ensp;Out of Service<br>
				&nbsp;&emsp;20:00 ~ 9:00<br>
				~~~~~~~~~~~~~~~<br>
			</h1>
		</html>
	`
)

type mode string

var (
	normalMode mode = "normal"
	babyMode   mode = "bady"
)

var (
	motionConfs = map[mode]string{
		normalMode: "/home/pi/motion_conf/normal_mode.conf",
		babyMode:   "/home/pi/motion_conf/baby_mode.conf",
	}
)

func main() {
	hServo := dev.NewSG90(pinSGH)
	if hServo == nil {
		log.Printf("[vmonitor]failed to new a sg90")
		return
	}

	vServo := dev.NewSG90(pinSGV)
	if vServo == nil {
		log.Printf("[vmonitor]failed to new a sg90")
		return
	}

	led := dev.NewLedImp(pinLed)
	if led == nil {
		log.Printf("[vmonitor]failed to new a led, will run the monitor without led")
	}

	bzr := dev.NewBuzzerImp(pinBzr, dev.High)
	if bzr == nil {
		log.Printf("[vmonitor]failed to new a buzzer, will run the monitor without buzzer")
	}

	btn := dev.NewButtonImp(pinBtn)
	if btn == nil {
		log.Printf("[vmonitor]failed to new a button, will run the monitor without button")
	}

	server := newVideoServer(hServo, vServo, led, bzr, btn)
	if server == nil {
		log.Printf("[vmonitor]failed to new the video server")
		return
	}

	util.WaitQuit(func() {
		server.stop()
	})

	log.Printf("[vmonitor]video server started")
	server.start()
}

type videoServer struct {
	hServo dev.Motor
	vServo dev.Motor
	led    dev.Led
	buzzer dev.Buzzer
	button dev.Button

	mode        mode
	inServing   bool
	hAngle      float64
	vAngle      float64
	chAlert     chan int
	pageContext []byte
}

func newVideoServer(hServo, vServo dev.Motor, led dev.Led, buzzer dev.Buzzer, button dev.Button) *videoServer {
	v := &videoServer{
		hServo: hServo,
		vServo: vServo,
		led:    led,
		buzzer: buzzer,
		button: button,

		mode:      normalMode,
		inServing: true,
		hAngle:    0,
		vAngle:    0,
		chAlert:   make(chan int, 16),
	}

	if err := v.restartMotion(); err != nil {
		return nil
	}
	return v
}

func (v *videoServer) start() {
	go v.hServo.Roll(v.hAngle)
	go v.vServo.Roll(v.vAngle)
	go v.alert()
	go v.detectConnecting()
	go v.detectServing()
	go v.detectingMode()

	if err := v.loadHomePage(); err != nil {
		log.Fatalf("failed to load home page, error: %v", err)
		return
	}

	http.HandleFunc("/", v.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}

func (v *videoServer) loadHomePage() error {
	if v.mode == normalMode {
		data, err := ioutil.ReadFile("vmonitor.html")
		if err != nil {
			return err
		}
		v.pageContext = data
		return nil
	}

	if v.mode == babyMode {
		ip := util.GetIP()
		if ip == "" {
			return errors.New("internal error: failed to get ip")
		}
		data, err := ioutil.ReadFile("vmonitor_baby.html")
		if err != nil {
			return err
		}
		rbuf := bytes.NewBuffer(data)
		wbuf := bytes.NewBuffer([]byte{})
		for {
			line, err := rbuf.ReadBytes('\n')
			if err == io.EOF {
				break
			}
			s := strings.Replace(string(line), ipPattern, ip, 1)
			wbuf.Write([]byte(s))
		}
		v.pageContext = wbuf.Bytes()
		return nil
	}

	return fmt.Errorf("invalid mode: %v", v.mode)
}

func (v *videoServer) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		v.homePageHandler(w, r)
	case "POST":
		op := r.FormValue("op")
		v.do(op)
	}
}

func (v *videoServer) homePageHandler(w http.ResponseWriter, r *http.Request) {
	if v.outOfService() && v.mode == normalMode {
		w.Write([]byte(pageOutOfService))
		return
	}
	w.Write(v.pageContext)
}

func (v *videoServer) do(op string) {
	switch op {
	case "left":
		go v.left()
	case "right":
		go v.right()
	case "up":
		go v.up()
	case "down":
		go v.down()
	case "beep":
		go v.beep(5, 100)
	default:
		log.Printf("[vmonitor]invalid operator: %v", op)
	}
}

func (v *videoServer) stop() {
	v.led.Off()
	close(v.chAlert)
}

func (v *videoServer) left() {
	log.Printf("[vmonitor]op: left")
	angle := v.hAngle - 15
	if angle < -90 {
		return
	}
	v.hAngle = angle
	log.Printf("[vmonitor]servo: %v", angle)
	v.hServo.Roll(angle)
}

func (v *videoServer) right() {
	log.Printf("[vmonitor]op: right")
	angle := v.hAngle + 15
	if angle > 75 {
		return
	}
	v.hAngle = angle
	log.Printf("[vmonitor]servo: %v", angle)
	v.hServo.Roll(angle)
}

func (v *videoServer) up() {
	log.Printf("[vmonitor]op: up")
	angle := v.vAngle + 15
	if angle > 90 {
		return
	}
	v.vAngle = angle
	log.Printf("[vmonitor]servo: %v", angle)
	v.vServo.Roll(angle)
}

func (v *videoServer) down() {
	log.Printf("[vmonitor]op: down")
	angle := v.vAngle - 15
	if angle < -30 {
		return
	}
	v.vAngle = angle
	log.Printf("[vmonitor]servo: %v", angle)
	v.vServo.Roll(angle)
}

func (v *videoServer) beep(n int, interval int) {
	log.Printf("[vmonitor]op: beep")
	if v.buzzer == nil {
		return
	}
	v.buzzer.Beep(n, interval)
}

func (v *videoServer) detectConnecting() {
	for {
		util.DelaySec(5)
		n, err := v.getConCount()
		if err != nil {
			log.Printf("[vmonitor]failed to get users count, err: %v", err)
			continue
		}
		v.chAlert <- n
	}
}

func (v *videoServer) alert() {
	conCount := 0
	for {
		if v.mode == babyMode {
			util.DelaySec(1)
			continue
		}
		select {
		case n := <-v.chAlert:
			if n > conCount {
				// there are new connections, give an alert
				go v.beep(2, 100)
			}
			conCount = n
		default:
			// do nothing
		}
		if conCount > 0 {
			v.led.Blink(1, 1000)
		}
		util.DelaySec(1)
	}
}

// getConCount is get the count of connecting to the server
func (v *videoServer) getConCount() (int, error) {
	cmd := `netstat -nat|grep -i "127.0.0.1:8081"|wc -l`
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return 0, err
	}
	var count int
	str := string(out)
	if n, err := fmt.Sscanf(str, "%d\n", &count); n != 1 || err != nil {
		return 0, fmt.Errorf("failed to parse the output of netstat")
	}
	return count, nil
}

func (v *videoServer) detectServing() {
	for {
		util.DelaySec(15)
		if v.mode == babyMode {
			// keep serving in baby monitor mode
			continue
		}
		if v.outOfService() {
			if v.inServing {
				log.Printf("[vmonitor]out of service, stop motion")
				util.StopMotion()
				v.inServing = false
			}
			continue
		}

		if !v.inServing {
			log.Printf("[vmonitor]in service time, start motion")
			util.StartMotion(motionConfs[v.mode])
			v.inServing = true
		}
	}
}

func (v *videoServer) outOfService() bool {
	hour := time.Now().Hour()
	if hour >= 20 || hour < 9 {
		// out of service at 20:00~09:00
		return true
	}
	return false
}

func (v *videoServer) detectingMode() {
	count := 0
	for {
		if v.button.Pressed() {
			count++
		} else {
			count = 0
		}
		if count < 5 {
			util.DelayMs(500)
			continue
		}

		count = 0
		log.Printf("[vmonitor]the button was pressed")
		go v.led.Blink(2, 100)
		lastMode := v.mode
		if v.mode == normalMode {
			v.mode = babyMode
		} else if v.mode == babyMode {
			v.mode = normalMode
		} else {
			// make a dalay detecting
			util.DelaySec(1)
			continue
		}
		if err := v.loadHomePage(); err != nil {
			log.Printf("[vmonitor]failed to load home page, error: %v", err)
			continue
		}
		if err := v.restartMotion(); err != nil {
			log.Printf("[vmonitor]failed to restart motion, error: %v", err)
			continue
		}
		go v.led.Blink(5, 100)
		log.Printf("[vmonitor]mode changed: %v --> %v", lastMode, v.mode)
	}
}

func (v *videoServer) restartMotion() error {
	if err := util.StopMotion(); err != nil {
		return err
	}
	if err := util.StartMotion(motionConfs[v.mode]); err != nil {
		return err
	}
	if v.mode == normalMode {
		v.inServing = true
	}
	return nil
}
