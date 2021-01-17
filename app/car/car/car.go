package car

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/speech"
	"github.com/shanghuiyang/image-recognizer/recognizer"
	"github.com/shanghuiyang/rpi-devices/cv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/rpi-devices/util/geo"
)

// Car ...
type Car struct {
	engine *dev.L298N
	horn   *dev.Buzzer
	led    *dev.Led
	light  *dev.Led
	camera *dev.Camera
	lc12s  *dev.LC12S
	chOp   chan Op

	// self-driving
	servo       *dev.SG90
	dmeter      dev.DistMeter
	encoder     *dev.Encoder
	gy25        *dev.GY25
	collisions  []*dev.Collision
	selfdriving bool
	servoAngle  int

	// speed-driving
	asr           *speech.ASR
	tts           *speech.TTS
	imgr          *recognizer.Recognizer
	speechdriving bool
	volume        int

	// self-tracking
	tracker      *cv.Tracker
	selftracking bool

	// nav
	gps       *dev.GPS
	dest      *geo.Point
	gpslogger *dev.GPSLogger
	lastLoc   *geo.Point
	selfnav   bool
}

// New ...
func New(cfg *Config) *Car {
	car := &Car{
		engine:     cfg.Engine,
		horn:       cfg.Horn,
		led:        cfg.Led,
		light:      cfg.Led,
		camera:     cfg.Camera,
		lc12s:      cfg.LC12S,
		servo:      cfg.Servo,
		dmeter:     cfg.DistMeter,
		gy25:       cfg.GY25,
		collisions: cfg.Collisions,
		gps:        cfg.GPS,

		servoAngle:    0,
		selfdriving:   false,
		speechdriving: false,
		selftracking:  false,
		selfnav:       false,
		chOp:          make(chan Op, chSize),
	}
	return car
}

// Start ...
func (c *Car) Start() error {
	go c.start()
	go c.servo.Roll(0)
	go c.blink()
	go c.joystick()
	go c.setVolume(40)
	c.speed(30)
	return nil
}

// Do ...
func (c *Car) Do(op Op) {
	c.chOp <- op
}

// Stop ...
func (c *Car) Stop() error {
	close(c.chOp)
	c.engine.Stop()
	return nil
}

// GetState ...
func (c *Car) GetState() (selfDriving, selfTracking, speechDriving bool) {
	return c.selfdriving, c.selftracking, c.speechdriving
}

// SetDest ...
func (c *Car) SetDest(dest *geo.Point) {
	if c.selfnav {
		return
	}
	c.dest = dest
}

func (c *Car) start() {
	for op := range c.chOp {
		switch op {
		case forward:
			c.forward()
		case backward:
			c.backward()
		case left:
			c.left()
		case right:
			c.right()
		case stop:
			c.stop()
		case beep:
			go c.beep()
		case servoleft:
			go c.servoLeft()
		case servoright:
			go c.servoRight()
		case servoahead:
			go c.servoAhead()
		case musicon:
			go c.musicOn()
		case musicoff:
			go c.musicOff()
		case selfdrivingon:
			go c.selfDrivingOn()
		case selfdrivingoff:
			go c.selfDrivingOff()
		case selftrackingon:
			go c.selfTrackingOn()
		case selftrackingoff:
			go c.selfTrackingOff()
		case speechdrivingon:
			go c.speechDrivingOn()
		case speechdrivingoff:
			go c.speechDrivingOff()
		case selfnavon:
			go c.selfNavOn()
		case selfnavoff:
			go c.selfNavOff()
		default:
			log.Printf("[car]invalid op")
		}
	}
}

// forward ...
func (c *Car) forward() {
	log.Printf("[car]forward")
	c.engine.Forward()
}

// backward ...
func (c *Car) backward() {
	log.Printf("[car]backward")
	c.engine.Backward()
}

// left ...
func (c *Car) left() {
	log.Printf("[car]left")
	c.engine.Left()
}

// right ...
func (c *Car) right() {
	log.Printf("[car]right")
	c.engine.Right()
}

// stop ...
func (c *Car) stop() {
	log.Printf("[car]stop")
	c.engine.Stop()
}

func (c *Car) speed(s uint32) {
	log.Printf("[car]speed %v%%", s)
	c.engine.Speed(s)
}

// beep ...
func (c *Car) beep() {
	log.Printf("[car]beep")
	if c.horn == nil {
		return
	}
	c.horn.Beep(5, 100)
}

func (c *Car) blink() {
	for {
		if c.speechdriving {
			util.DelayMs(2000)
			continue
		}
		c.led.Blink(1, 1000)
	}
}

func (c *Car) musicOn() {
	log.Printf("[car]music on")
	util.PlayMp3("./music/*.mp3")
}

func (c *Car) musicOff() {
	log.Printf("[car]music off")
	util.StopMp3()
	time.Sleep(1 * time.Second)
}

func (c *Car) servoLeft() {
	angle := c.servoAngle - 15
	if angle < -90 {
		angle = -90
	}
	c.servoAngle = angle
	log.Printf("[car]servo roll %v", angle)
	if c.servo == nil {
		return
	}
	c.servo.Roll(angle)
}

func (c *Car) servoRight() {
	angle := c.servoAngle + 15
	if angle > 90 {
		angle = 90
	}
	c.servoAngle = angle
	log.Printf("[car]servo roll %v", angle)
	if c.servo == nil {
		return
	}
	c.servo.Roll(angle)
}

func (c *Car) servoAhead() {
	c.servoAngle = 0
	log.Printf("[car]servo roll %v", 0)
	if c.servo == nil {
		return
	}
	c.servo.Roll(0)
}

func (c *Car) selfDriving() {
	if c.dmeter == nil {
		log.Printf("[car]can't self-driving without the distance sensor")
		return
	}

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	var (
		fwd       bool
		retry     int
		mindAngle int
		maxdAngle int
		mind      float64
		maxd      float64
		op        = forward
		chOp      = make(chan Op, 4)
	)

	for c.selfdriving || c.selftracking {
		select {
		case p := <-chOp:
			op = p
			for len(chOp) > 0 {
				log.Printf("[car]skip op: %v", <-chOp)
				// _ = <-chOp
			}
		default:
			// 	do nothing
		}
		log.Printf("[car]op: %v", op)

		switch op {
		case backward:
			fwd = false
			c.stop()
			util.DelayMs(20)
			c.backward()
			util.DelayMs(500)
			chOp <- stop
			continue
		case stop:
			fwd = false
			c.stop()
			util.DelayMs(20)
			chOp <- scan
			continue
		case scan:
			fwd = false
			mind, maxd, mindAngle, maxdAngle = c.scan()
			log.Printf("[car]mind=%.0f, maxd=%.0f, mindAngle=%v, maxdAngle=%v", mind, maxd, mindAngle, maxdAngle)
			if mind < 10 && mindAngle != 90 && mindAngle != -90 && retry < 4 {
				chOp <- backward
				retry++
				continue
			}
			chOp <- turn
			retry = 0
		case turn:
			fwd = false
			c.turn(maxdAngle)
			util.DelayMs(150)
			chOp <- forward
			continue
		case forward:
			if !fwd {
				c.forward()
				fwd = true
				go c.detecting(chOp)
			}
			util.DelayMs(50)
			continue
		case pause:
			fwd = false
			util.DelayMs(500)
			continue
		}
	}
	c.stop()
	util.DelayMs(1000)
	close(chOp)
}

func (c *Car) speechDriving() {
	var (
		op   = stop
		fwd  = false
		chOp = make(chan Op, 4)
		wg   sync.WaitGroup
	)

	wg.Add(1)
	go c.detectSpeech(chOp, &wg)
	for c.speechdriving {
		select {
		case p := <-chOp:
			op = p
			for len(chOp) > 0 {
				// log.Printf("[car]len(chOp)=%v", len(chOp))
				_ = <-chOp
			}
		default:
			// do nothing
		}
		log.Printf("[car]op: %v", op)

		switch op {
		case forward:
			if !fwd {
				c.forward()
				fwd = true
				go c.detecting(chOp)
			}
			util.DelayMs(50)
			continue
		case backward:
			fwd = false
			c.stop()
			util.DelayMs(20)
			c.backward()
			util.DelayMs(600)
			chOp <- stop
			continue
		case left:
			fwd = false
			c.stop()
			util.DelayMs(20)
			c.turn(-90)
			util.DelayMs(20)
			chOp <- forward
			continue
		case right:
			fwd = false
			c.stop()
			util.DelayMs(20)
			c.turn(90)
			util.DelayMs(20)
			chOp <- forward
			continue
		case roll:
			fwd = false
			c.engine.Left()
			util.DelayMs(3000)
			chOp <- stop
			continue
		case stop:
			fwd = false
			c.stop()
			util.DelayMs(500)
			continue
		}
	}
	c.stop()
	wg.Wait()
	close(chOp)
}

func (c *Car) selfDrivingOn() {
	if c.selfdriving {
		return
	}
	c.selftracking = false
	c.speechdriving = false
	util.DelayMs(1000) // wait for self-tracking and speech-driving quit

	c.selfdriving = true
	log.Printf("[car]self-drving on")
	c.speed(30)
	c.selfDriving()
}

func (c *Car) selfDrivingOff() {
	c.selfdriving = false
	c.servo.Roll(0)
	log.Printf("[car]self-drving off")
}

func (c *Car) selfTrackingOn() {
	if c.selftracking {
		return
	}
	util.StopMotion()
	c.selfdriving = false
	c.speechdriving = false
	c.selfnav = false
	util.DelayMs(1000) // wait to quit self-driving & speech-driving

	// start slef-tracking
	t, err := cv.NewTracker(lh, ls, lv, hh, hs, hv)
	if err != nil {
		log.Printf("[carapp]failed to create a tracker, error: %v", err)
		return
	}
	c.tracker = t
	c.selftracking = true
	log.Printf("[car]self-tracking on")
	c.speed(30)
	c.selfDriving()
}

func (c *Car) selfTrackingOff() {
	c.selftracking = false
	c.tracker.Close()
	c.servo.Roll(0)
	util.DelayMs(500)

	if err := util.StartMotion(); err != nil {
		log.Printf("[car]failed to start motion, error: %v", err)
	}
	log.Printf("[car]self-tracking off")
}

func (c *Car) speechDrivingOn() {
	if c.speechdriving {
		return
	}
	c.selfdriving = false
	c.selftracking = false
	util.DelayMs(1000) // wait for self-driving and self-tracking quit

	c.speechdriving = true
	log.Printf("[car]speech-drving on")
	c.speed(30)
	c.speechDriving()
}

func (c *Car) speechDrivingOff() {
	c.speechdriving = false
	c.servo.Roll(0)
	log.Printf("[car]speech-drving off")
}

func (c *Car) detecting(chOp chan Op) {

	chQuit := make(chan bool, 4)
	var wg sync.WaitGroup

	wg.Add(1)
	go c.detectCollision(chOp, chQuit, &wg)

	wg.Add(1)
	go c.detectObstacles(chOp, chQuit, &wg)

	if c.selftracking {
		wg.Add(1)
		go c.trackingObj(chOp, chQuit, &wg)
	}

	wg.Wait()
	close(chQuit)
}

func (c *Car) detectObstacles(chOp chan Op, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for c.selfdriving || c.selftracking || c.speechdriving {
		for _, angle := range aheadAngles {
			select {
			case quit := <-chQuit:
				if quit {
					return
				}
			default:
				// do nothing
			}
			c.servo.Roll(angle)
			util.DelayMs(70)
			d := c.dmeter.Dist()
			if d < 20 {
				chOp <- backward
				chQuit <- true
				chQuit <- true
				return
			}
			if d < 40 {
				chOp <- stop
				chQuit <- true
				chQuit <- true
				return
			}
		}
	}
}

func (c *Car) detectCollision(chOp chan Op, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for c.selfdriving || c.selftracking || c.speechdriving {
		select {
		case quit := <-chQuit:
			if quit {
				return
			}
		default:
			// do nothing
		}
		for _, collision := range c.collisions {
			if collision.Collided() {
				chOp <- backward
				go c.horn.Beep(1, 100)
				log.Printf("[car]crashed")
				chQuit <- true
				chQuit <- true
				return
			}
		}
		util.DelayMs(10)
	}
}

func (c *Car) trackingObj(chOp chan Op, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	angle := 0
	for c.selftracking {
		select {
		case quit := <-chQuit:
			if quit {
				return
			}
		default:
			// do nothing
		}

		ok, _ := c.tracker.Locate()
		if !ok {
			continue
		}

		// found a ball
		log.Printf("[car]found a ball")
		chQuit <- true
		chQuit <- true
		chOp <- pause
		c.stop()

		firstTime := true // see a ball at the first time
		for c.selftracking {
			ok, rect := c.tracker.Locate()
			if !ok {
				// lost the ball, looking for it by turning 360 degree
				log.Printf("[car]lost the ball")
				firstTime = true
				if angle < 360 {
					c.turn(30)
					angle += 30
					util.DelayMs(200)
					continue
				}
				chOp <- scan
				return
			}
			angle = 0
			if rect.Max.Y > 580 {
				c.stop()
				c.horn.Beep(1, 300)
				continue
			}
			if firstTime {
				go c.horn.Beep(2, 100)
			}
			firstTime = false
			x, y := c.tracker.MiddleXY(rect)
			log.Printf("[car]found a ball at: (%v, %v)", x, y)
			if x < 200 {
				log.Printf("[car]turn right to the ball")
				c.engine.Right()
				util.DelayMs(100)
				c.engine.Stop()
				continue
			}
			if x > 400 {
				log.Printf("[car]turn left to the ball")
				c.engine.Left()
				util.DelayMs(100)
				c.engine.Stop()
				continue
			}
			log.Printf("[car]forward to the ball")
			c.engine.Forward()
			util.DelayMs(100)
			c.engine.Stop()
		}

	}
}

func (c *Car) detectSpeech(chOp chan Op, wg *sync.WaitGroup) {
	defer wg.Done()

	speechAuth := oauth.New(baiduSpeechAppKey, baiduSpeechSecretKey, oauth.NewCacheMan())
	c.asr = speech.NewASR(speechAuth)
	c.tts = speech.NewTTS(speechAuth)

	imgAuth := oauth.New(baiduImgRecognitionAppKey, baiduImgRecognitionSecretKey, oauth.NewCacheMan())
	c.imgr = recognizer.New(imgAuth)

	for c.speechdriving {
		log.Printf("[car]start recording")
		go c.led.On()
		wav := "car.wav"
		if err := util.Record(2, wav); err != nil {
			log.Printf("[car]failed to record the speech: %v", err)
			continue
		}
		go c.led.Off()
		log.Printf("[car]stop recording")

		text, err := c.asr.ToText(wav)
		if err != nil {
			log.Printf("[car]failed to recognize the speech, error: %v", err)
			continue
		}
		log.Printf("[car]speech: %v", text)

		switch {
		case strings.Index(text, "前") >= 0:
			chOp <- forward
		case strings.Index(text, "后") >= 0:
			chOp <- backward
		case strings.Index(text, "左") >= 0:
			chOp <- left
		case strings.Index(text, "右") >= 0:
			chOp <- right
		case strings.Index(text, "停") >= 0:
			chOp <- stop
		case strings.Index(text, "转圈") >= 0:
			chOp <- roll
		case strings.Index(text, "是什么") >= 0:
			c.recognize()
		case strings.Index(text, "开灯") >= 0:
			c.light.On()
		case strings.Index(text, "关灯") >= 0:
			c.light.Off()
		case strings.Index(text, "大声") >= 0:
			c.volumeUp()
		case strings.Index(text, "小声") >= 0:
			c.volumeDown()
		case strings.Index(text, "唱歌") >= 0:
			go util.PlayWav("./music/xiaomaolv.wav")
		default:
			// do nothing
		}
	}
}

// scan for geting the min & max distance, and their corresponding angles
// mind: the min distance
// maxd: the max distance
// mindAngle: the angle correspond to the mind
// maxdAngle: the angle correspond to the maxd
func (c *Car) scan() (mind, maxd float64, mindAngle, maxdAngle int) {
	mind = 9999
	maxd = -9999
	for _, ang := range scanningAngles {
		c.servo.Roll(ang)
		util.DelayMs(100)
		d := c.dmeter.Dist()
		for i := 0; d < 0 && i < 3; i++ {
			util.DelayMs(100)
			d = c.dmeter.Dist()
		}
		if d < 0 {
			continue
		}
		log.Printf("[car]scan: angle=%v, dist=%.0f", ang, d)
		if d < mind {
			mind = d
			mindAngle = ang
		}
		if d > maxd {
			maxd = d
			maxdAngle = ang
		}
	}
	c.servo.Roll(0)
	util.DelayMs(50)
	return
}

func (c *Car) turn(angle int) {
	turnf := c.engine.Right
	if angle < 0 {
		turnf = c.engine.Left
		angle *= (-1)
	}

	yaw, _, _, err := c.gy25.Angles()
	if err != nil {
		log.Printf("[car]failed to get angles from gy-25, error: %v", err)
		return
	}

	retry := 0
	for {
		turnf()
		yaw2, _, _, err := c.gy25.Angles()
		if err != nil {
			log.Printf("[car]failed to get angles from gy-25, error: %v", err)
			if retry < 3 {
				retry++
				continue
			}
			break
		}
		ang := c.gy25.IncludedAngle(yaw, yaw2)
		if ang >= float64(angle) {
			break
		}
		time.Sleep(100 * time.Millisecond)
		c.engine.Stop()
		time.Sleep(100 * time.Millisecond)
	}
	c.engine.Stop()
	return
}

func (c *Car) turnLeft(angle int) {
	n := angle/5 - 1
	c.encoder.Start()
	defer c.encoder.Stop()

	c.chOp <- left
	for i := 0; i < n; {
		i += c.encoder.Count1()
	}
	return
}

func (c *Car) turnRight(angle int) {
	n := angle/5 - 1
	c.encoder.Start()
	defer c.encoder.Stop()

	c.chOp <- right
	for i := 0; i < n; {
		i += c.encoder.Count1()
	}
	return
}

func (c *Car) recognize() error {
	log.Printf("[car]take photo")
	imagef, err := c.camera.TakePhoto()
	if err != nil {
		log.Printf("[car]failed to take phote, error: %v", err)
		return err
	}
	util.PlayWav(letMeThinkWav)

	log.Printf("[car]recognize image")
	objname, err := c.recognizeImg(imagef)
	if err != nil {
		log.Printf("[car]failed to recognize image, error: %v", err)
		util.PlayWav(errorWav)
		return err
	}
	log.Printf("[car]object: %v", objname)

	if err := c.playText("这是" + objname); err != nil {
		log.Printf("[car]failed to play text, error: %v", err)
		return err
	}

	return nil
}

func (c *Car) setVolume(v int) error {
	if err := util.SetVolume(v); err != nil {
		log.Printf("[car]failed to set volume, error: %v", err)
		return err
	}
	c.volume = v
	return nil
}

func (c *Car) volumeUp() {
	v := c.volume + 10
	if v > 100 {
		v = 100
	}
	c.setVolume(v)
	go c.playText(fmt.Sprintf("音量%v%%", v))
}

func (c *Car) volumeDown() {
	v := c.volume - 10
	if v < 0 {
		v = 0
	}
	c.setVolume(v)
	go c.playText(fmt.Sprintf("音量%v%%", v))
}

func (c *Car) recognizeImg(imageFile string) (string, error) {
	if c.imgr == nil {
		return "", errors.New("invalid image recognizer")
	}
	name, err := c.imgr.Recognize(imageFile)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (c *Car) toSpeech(text string) (string, error) {
	data, err := c.tts.ToSpeech(text)
	if err != nil {
		log.Printf("[car]failed to convert text to speech, error: %v", err)
		return "", err
	}

	if err := ioutil.WriteFile(thisIsXWav, data, 0644); err != nil {
		log.Printf("[car]failed to save %v, error: %v", thisIsXWav, err)
		return "", err
	}
	return thisIsXWav, nil
}

func (c *Car) playText(text string) error {
	wav, err := c.toSpeech(text)
	if err != nil {
		log.Printf("[car]failed to tts, error: %v", err)
		return err
	}

	if err := util.PlayWav(wav); err != nil {
		log.Printf("[car]failed to play wav: %v, error: %v", wav, err)
		return err
	}
	return nil
}

func (c *Car) joystick() {
	if c.lc12s == nil {
		return
	}

	c.lc12s.Wakeup()
	defer c.lc12s.Sleep()

	for {
		time.Sleep(200 * time.Millisecond)

		if c.selfdriving {
			continue
		}

		data, err := c.lc12s.Receive()
		if err != nil {
			log.Printf("[car]failed to receive data from LC12S, error: %v", err)
			continue
		}
		log.Printf("[car]LC12S received: %v", data)

		if len(data) != 1 {
			log.Printf("[car]invalid data from LC12S, data len: %v", len(data))
			continue
		}

		op := (data[0] >> 4)
		speed := data[0] & 0x0F

		switch op {
		case 0:
			c.chOp <- stop
		case 1:
			c.chOp <- forward
		case 2:
			c.chOp <- backward
		case 3:
			c.chOp <- left
		case 4:
			c.chOp <- right
		case 5:
			if c.selfdriving {
				c.chOp <- selfdrivingoff
				continue
			}
			c.chOp <- selfdrivingon
		default:
			c.chOp <- stop
		}
		c.speed(uint32(speed * 10))
	}
}

func (c *Car) selfNavOn() {
	if c.selfnav {
		return
	}
	if c.gps == nil {
		log.Printf("[car]failed to nav, error: without gps device")
		return
	}

	c.selfdriving = false
	c.selftracking = false
	c.speechdriving = false
	util.DelayMs(1000) // wait for self-tracking and speech-driving quit

	c.selfnav = true
	log.Printf("[car]nav on")
	if err := c.selfNav(); err != nil {
		return
	}
	c.selfnav = false
}

func (c *Car) selfNavOff() {
	c.selfnav = false
	log.Printf("[car]nav off")
}

func (c *Car) selfNav() error {
	if c.dest == nil {
		log.Printf("[car]destination didn't be set, stop nav")
		return errors.New("destination isn't set")
	}

	c.horn.Beep(3, 300)
	if !bbox.IsInside(c.dest) {
		log.Printf("[car]destination isn't in bbox, stop nav")
		return errors.New("destination isn't in bbox")
	}

	c.gpslogger = dev.NewGPSLogger()
	if c.gpslogger == nil {
		log.Printf("[car]failed to new a tracker, stop nav")
		return errors.New("gpslogger is nil")
	}
	defer c.gpslogger.Close()

	var org *geo.Point
	for c.selfnav {
		pt, err := c.gps.Loc()
		if err != nil {
			log.Printf("[car]gps sensor is not ready")
			util.DelayMs(1000)
			continue
		}
		c.gpslogger.AddPoint(org)
		if !bbox.IsInside(pt) {
			log.Printf("current loc(%v) isn't in bbox(%v)", pt, bbox)
			continue
		}
		org = pt
		break
	}
	if !c.selfnav {
		return errors.New("nav abort")
	}
	c.lastLoc = org

	path, err := findPath(org, c.dest)
	if err != nil {
		log.Printf("[car]failed to find a path, error: %v", err)
		return errors.New("failed to find a path")
	}
	turns := turnPoints(path)

	var turnPts []*geo.Point
	var str string
	for _, xy := range turns {
		pt := xy2geo(xy)
		str += fmt.Sprintf("(%v) ", pt)
		turnPts = append(turnPts, pt)
	}
	log.Printf("[car]turn points(lat,lon): %v", str)

	c.chOp <- forward
	util.DelayMs(1000)
	for i, p := range turnPts {
		if err := c.navTo(p); err != nil {
			log.Printf("[car]failed to nav to (%v), error: %v", p, err)
			break
		}
		if i < len(turnPts)-1 {
			// turn point
			go c.horn.Beep(2, 100)
		} else {
			// destination
			go c.horn.Beep(5, 300)
		}
	}
	c.chOp <- stop
	return nil
}

func (c *Car) navTo(dest *geo.Point) error {
	retry := 8
	for c.selfnav {
		loc, err := c.gps.Loc()
		if err != nil {
			c.chOp <- stop
			log.Printf("[car]gps sensor is not ready")
			util.DelayMs(1000)
			continue
		}

		if !bbox.IsInside(loc) {
			c.chOp <- stop
			log.Printf("current loc(%v) isn't in bbox(%v)", loc, bbox)
			util.DelayMs(1000)
			continue
		}

		c.gpslogger.AddPoint(loc)
		log.Printf("[car]current loc: %v", loc)

		d := loc.DistanceWith(c.lastLoc)
		log.Printf("[car]distance to last loc: %.2f m", d)
		if d > 4 && retry < 5 {
			c.chOp <- stop
			log.Printf("[car]bad gps signal, waiting for better gps signal")
			retry++
			util.DelayMs(1000)
			continue
		}

		retry = 0
		d = loc.DistanceWith(dest)
		log.Printf("[car]distance to destination: %.2f m", d)
		if d < 4 {
			c.chOp <- stop
			log.Printf("[car]arrived at the destination, nav done")
			return nil
		}

		side := geo.Side(c.lastLoc, loc, dest)
		angle := int(180 - geo.Angle(c.lastLoc, loc, dest))
		if angle < 10 {
			side = geo.MiddleSide
		}
		log.Printf("[car]nav angle: %v, side: %v", angle, side)

		switch side {
		case geo.LeftSide:
			c.turnLeft(angle)
		case geo.RightSide:
			c.turnRight(angle)
		case geo.MiddleSide:
			// do nothing
		}
		c.chOp <- forward
		util.DelayMs(1000)
		c.lastLoc = loc
	}
	c.chOp <- stop
	return nil
}
