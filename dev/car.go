package dev

import (
	"errors"
	"image"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/speech"
	"github.com/shanghuiyang/image-recognizer/recognizer"
	"gocv.io/x/gocv"
)

const (
	chSize        = 8
	letMeThinkWav = "let_me_think.wav"
	thisIsXWav    = "this_is_x.wav"
	iDontKnowWav  = "i_dont_know.wav"
	errorWav      = "error.wav"
)

const (
	baiduSpeechAppKey    = "your_speech_app_key"
	baiduSpeechSecretKey = "your_speech_secret_key"

	baiduImgRecognitionAppKey    = "your_image_recognition_app_key"
	baiduImgRecognitionSecretKey = "your_image_recognition_secrect_key"
)

const (
	forward  CarOp = "forward"
	backward CarOp = "backward"
	left     CarOp = "left"
	right    CarOp = "right"
	stop     CarOp = "stop"
	turn     CarOp = "turn"
	scan     CarOp = "scan"

	beep  CarOp = "beep"
	blink CarOp = "blink"

	servoleft  CarOp = "servoleft"
	servoright CarOp = "servoright"
	servoahead CarOp = "servoahead"

	lighton  CarOp = "lighton"
	lightoff CarOp = "lightoff"

	selfdrivingon  CarOp = "selfdrivingon"
	selfdrivingoff CarOp = "selfdrivingoff"

	speechdrivingon  CarOp = "speechdrivingon"
	speechdrivingoff CarOp = "speechdrivingoff"
)

var (
	scanningAngles  = []int{-90, -75, -60, -45, -30, 30, 45, 60, 75, 90}
	turnAngleCounts = map[int]int{
		-90: 20,
		-75: 17,
		-60: 14,
		-45: 10,
		-30: 7,
		30:  5,
		45:  8,
		60:  10,
		75:  13,
		90:  17,
	}
	aheadAngles = []int{0, -15, 0, 15}
)

type (
	// CarOp ...
	CarOp string
	// Option ...
	Option func(c *Car)
)

// WithEngine ...
func WithEngine(engine *L298N) Option {
	return func(c *Car) {
		c.engine = engine
	}
}

// WithServo ...
func WithServo(servo *SG90) Option {
	return func(c *Car) {
		c.servo = servo
	}
}

// WithUlt ...
func WithUlt(ult *US100) Option {
	return func(c *Car) {
		c.ult = ult
	}
}

// WithEncoder ...
func WithEncoder(e *Encoder) Option {
	return func(c *Car) {
		c.encoder = e
	}
}

// WithCSwitchs ...
func WithCSwitchs(cswitchs []*CollisionSwitch) Option {
	return func(c *Car) {
		c.cswitchs = cswitchs
	}
}

// WithHorn ...
func WithHorn(horn *Buzzer) Option {
	return func(c *Car) {
		c.horn = horn
	}
}

// WithLed ...
func WithLed(led *Led) Option {
	return func(c *Car) {
		c.led = led
	}
}

// WithLight ...
func WithLight(light *Led) Option {
	return func(c *Car) {
		c.light = light
	}
}

// WithCamera ...
func WithCamera(cam *Camera) Option {
	return func(c *Car) {
		c.camera = cam
	}
}

// Car ...
type Car struct {
	engine   *L298N
	servo    *SG90
	ult      *US100
	encoder  *Encoder
	cswitchs []*CollisionSwitch
	horn     *Buzzer
	led      *Led
	light    *Led
	camera   *Camera

	asr  *speech.ASR
	tts  *speech.TTS
	imgr *recognizer.Recognizer

	servoAngle    int
	selfdriving   bool
	speechdriving bool
	chOp          chan CarOp
}

// NewCar ...
func NewCar(opts ...Option) *Car {
	car := &Car{
		servoAngle:  0,
		selfdriving: false,
		chOp:        make(chan CarOp, chSize),
	}
	for _, opt := range opts {
		opt(car)
	}
	return car
}

// Start ...
func (c *Car) Start() error {
	go c.start()
	go c.servo.Roll(0)
	go c.blink()
	return nil
}

// Do ...
func (c *Car) Do(op CarOp) {
	c.chOp <- op
}

// Stop ...
func (c *Car) Stop() error {
	close(c.chOp)
	c.engine.Stop()
	c.encoder.Close()
	return nil
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
		case lighton:
			go c.lightOn()
		case lightoff:
			go c.lightOff()
		case selfdrivingon:
			go c.selfDrivingOn()
		case selfdrivingoff:
			go c.selfDrivingOff()
		case speechdrivingon:
			go c.speechDrivingOn()
		case speechdrivingoff:
			go c.speechDrivingOff()
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
	c.delay(250)
	c.engine.Stop()
}

// right ...
func (c *Car) right() {
	log.Printf("[car]right")
	c.engine.Right()
	c.delay(250)
	c.engine.Stop()
}

// stop ...
func (c *Car) stop() {
	log.Printf("[car]stop")
	c.engine.Stop()
}

func (c *Car) speed(s uint32) {
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
			c.delay(2000)
			continue
		}
		c.led.Blink(1, 1000)
	}
}

func (c *Car) lightOn() {
	log.Printf("[car]light on")
	if c.light == nil {
		return
	}
	c.light.On()
}

func (c *Car) lightOff() {
	log.Printf("[car]light off")
	if c.light == nil {
		return
	}
	c.light.Off()
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

/*

                                                                          +-----------------------------------------------+
                                                                          |                                               |
                                                                          v                                               |Y
+-------+     +---------+    +---------------+     +-----------+     +----+-----+      +------+      +------+     +--------------+
| start |---->| forward |--->|   obstacles   |---->| distance  |---->| backword |----->| stop |----->| scan |---->| min distance |
+-------+     +-----+---+    |   detected?   | Y   |  < 10cm?  | Y   +----------+      +--+---+      +------+     |    < 10cm    |
                    ^        +-------+-------+     +-----+-----+                          |                       +--------------+
                    |                |                   |                                ^                               |N
                    |                |N                 N|                                |                               |
                    |                |                   |                                |                               v
                    |                v                   |           +----------+ Y       |   Y +----------+   Y  +-------+------+
                    |                |                   +---------->| distance +------>--+-<---| retry<4? |-<----| max distance |
                    |                |                               |  < 40cm? |               +----+-----+      |    < 40cm    |
                    ^                |                               +----------+                    | N          +--------------+
                    |                |                                     |N                        v                    |N
                    |                |                                     |                         |                    |
                    |                |                                     v                    +----+-----+              |
                    +-------<--------+------------------<------------------+---------<----------|   turn   |-------<------+
                                                                                                +----------+



*/
func (c *Car) selfDrivingOn() {
	if c.selfdriving {
		return
	}

	log.Printf("[car]self-drving on")
	if c.ult == nil {
		log.Printf("[car]can't self-driving without the distance sensor")
		return
	}

	// make a warning before running into self-driving mode
	c.horn.Beep(3, 300)

	// start self-driving
	c.selfdriving = true
	c.speechdriving = false
	var (
		fwd       bool
		retry     int
		mindAngle int
		maxdAngle int
		mind      float64
		maxd      float64
		op        = forward
		chOp      = make(chan CarOp, 4)
	)

	for c.selfdriving {
		select {
		case p := <-chOp:
			op = p
			for len(chOp) > 0 {
				// log.Printf("[car]len(chOp)=%v, op=%v", len(chOp), <-chOp)
				_ = <-chOp
			}
		default:
			// 	do nothing
		}
		log.Printf("[car]op: %v", op)

		switch op {
		case backward:
			fwd = false
			c.stop()
			c.delay(20)
			c.backward()
			c.delay(500)
			chOp <- stop
			continue
		case stop:
			fwd = false
			c.stop()
			c.delay(20)
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
			c.delay(150)
			chOp <- forward
			continue
		case forward:
			if !fwd {
				c.forward()
				fwd = true
				go c.detecting(chOp)
			}
			c.delay(50)
			continue
		}
	}
	c.stop()
	c.delay(1000)
	close(chOp)
}

func (c *Car) speechDrivingOn() {
	if c.speechdriving {
		return
	}
	log.Printf("[car]speech-drving on")
	c.speechdriving = true
	c.selfdriving = false

	var (
		op   = stop
		fwd  = false
		chOp = make(chan CarOp, 4)
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
			c.delay(50)
			continue
		case backward:
			fwd = false
			c.stop()
			c.delay(20)
			c.backward()
			c.delay(600)
			chOp <- stop
			continue
		case left:
			fwd = false
			c.stop()
			c.delay(20)
			c.turn(-90)
			c.delay(20)
			chOp <- forward
			continue
		case right:
			fwd = false
			c.stop()
			c.delay(20)
			c.turn(90)
			c.delay(20)
			chOp <- forward
			continue
		case stop:
			fwd = false
			c.stop()
			c.delay(500)
			continue
		}
	}
	c.stop()
	wg.Wait()
	close(chOp)
}

func (c *Car) selfDrivingOff() {
	c.selfdriving = false
	log.Printf("[car]self-drving off")
}

func (c *Car) speechDrivingOff() {
	c.speechdriving = false
	log.Printf("[car]speech-drving off")
}

func (c *Car) detecting(chOp chan CarOp) {

	chQuit := make(chan bool, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go c.detectCollision(chOp, chQuit, &wg)

	wg.Add(1)
	go c.detectObstacles(chOp, chQuit, &wg)

	wg.Add(1)
	go c.detectBalls(chOp, chQuit, &wg)

	wg.Wait()
	close(chQuit)
}

func (c *Car) detectObstacles(chOp chan CarOp, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for c.selfdriving || c.speechdriving {
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
			c.delay(70)
			d := c.ult.Dist()
			if d < 10 {
				chOp <- backward
				chQuit <- true
				return
			}
			if d < 40 {
				chOp <- stop
				chQuit <- true
				return
			}
		}
	}
}

func (c *Car) detectCollision(chOp chan CarOp, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for c.selfdriving || c.speechdriving {
		select {
		case quit := <-chQuit:
			if quit {
				return
			}
		default:
			// do nothing
		}
		for _, cswitch := range c.cswitchs {
			if cswitch.Collided() {
				chOp <- backward
				go c.horn.Beep(1, 100)
				log.Printf("[car]crashed")
				chQuit <- true
				return
			}
		}
		c.delay(10)
	}
}

func (c *Car) detectBalls(chOp chan CarOp, chQuit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	// rcolor := color.RGBA{G: 255, A: 255}
	// lcolor := color.RGBA{R: 255, A: 255}

	// the red ball
	lhsv := gocv.Scalar{Val1: 0, Val2: 100, Val3: 100}
	hhsv := gocv.Scalar{Val1: 42, Val2: 255, Val3: 255}

	size := image.Point{X: 600, Y: 600}
	blur := image.Point{X: 11, Y: 11}

	img := gocv.NewMat()
	mask := gocv.NewMat()
	frame := gocv.NewMat()
	hsv := gocv.NewMat()
	kernel := gocv.NewMat()
	defer img.Close()
	defer mask.Close()
	defer frame.Close()
	defer hsv.Close()
	defer kernel.Close()

	video, _ := gocv.OpenVideoCapture(0)
	defer video.Close()

	for c.selfdriving {
		select {
		case quit := <-chQuit:
			if quit {
				return
			}
		default:
			// do nothing
		}

		if !video.Read(&img) {
			continue
		}
		gocv.Flip(img, &img, 1)
		gocv.Resize(img, &img, size, 0, 0, gocv.InterpolationLinear)
		gocv.GaussianBlur(img, &frame, blur, 0, 0, gocv.BorderReflect101)
		gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lhsv, hhsv, &mask)
		gocv.Erode(mask, &mask, kernel)
		gocv.Dilate(mask, &mask, kernel)
		cnt := bestContour(mask, 200)
		if len(cnt) == 0 {
			continue
		}
		chQuit <- true
		chOp <- stop

		rect := gocv.BoundingRect(cnt)
		// ---
		// gocv.Rectangle(&img, rect, rcolor, 2)
		// imgf := fmt.Sprintf("test%v.jpg", i+100000)
		// gocv.IMWrite(imgf, img)
		// i++
		// ---
		x, y := middle(rect)
		log.Printf("[car]ball at: (%v, %v)\n", x, y)
		dx := x - 300
		for dx < -50 {
			c.left()
		}
		for dx > 50 {
			c.right()
		}
		chOp <- forward
		return
	}
}

func (c *Car) detectSpeech(chOp chan CarOp, wg *sync.WaitGroup) {
	defer wg.Done()

	speechAuth := oauth.New(baiduSpeechAppKey, baiduSpeechSecretKey, oauth.NewCacheMan())
	c.asr = speech.NewASR(speechAuth)
	c.tts = speech.NewTTS(speechAuth)

	imgAuth := oauth.New(baiduImgRecognitionAppKey, baiduImgRecognitionSecretKey, oauth.NewCacheMan())
	c.imgr = recognizer.New(imgAuth)

	for c.speechdriving {
		// -D:			device
		// -d 3:		3 seconds
		// -t wav:		wav type
		// -r 16000:	Rate 16000 Hz
		// -c 1:		1 channel
		// -f S16_LE:	Signed 16 bit Little Endian
		cmd := `sudo arecord -D "plughw:1,0" -d 2 -t wav -r 16000 -c 1 -f S16_LE car.wav`
		log.Printf("[car]start recording")
		go c.led.On()
		_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
		if err != nil {
			log.Printf("[car]failed to record the speech: %v", err)
			continue
		}
		go c.led.Off()
		log.Printf("[car]stop recording")

		text, err := c.asr.ToText("car.wav")
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
		case strings.Index(text, "是什么") >= 0:
			c.recognize()
		case strings.Index(text, "唱歌") >= 0:
			go c.play("./music/xiaomaolv.wav")
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
		c.delay(120)
		d := c.ult.Dist()
		for i := 0; d < 0 && i < 3; i++ {
			c.delay(120)
			d = c.ult.Dist()
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
	c.delay(50)
	return
}

func (c *Car) turn(angle int) {
	n, ok := turnAngleCounts[angle]
	if !ok {
		log.Printf("[car]invalid angle: %d", angle)
		return
	}
	if angle < 0 {
		c.engine.Left()
	} else {
		c.engine.Right()
	}
	for i := 0; i < n; {
		i += c.encoder.Count1()
	}
	c.stop()
	return
}

func (c *Car) delay(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func (c *Car) recognize() error {
	log.Printf("[car]take photo")
	imagef, err := c.camera.TakePhoto()
	if err != nil {
		log.Printf("[car]failed to take phote, error: %v", err)
		return err
	}
	c.play(letMeThinkWav)

	log.Printf("[car]recognize image")
	objname, err := c.recognizeImg(imagef)
	if err != nil {
		log.Printf("[car]failed to recognize image, error: %v", err)
		c.play(errorWav)
		return err
	}
	log.Printf("[car]object: %v", objname)

	if err := c.playText("这是" + objname); err != nil {
		log.Printf("[car]failed to play text, error: %v", err)
		return err
	}

	return nil
}

func (c *Car) play(wav string) error {
	// aplay test.wav
	cmd := exec.Command("aplay", wav)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[car]failed to exec aplay, output: %v, error: %v", string(out), err)
		return err
	}
	return nil
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

	if err := c.play(wav); err != nil {
		log.Printf("[car]failed to play wav: %v, error: %v", wav, err)
		return err
	}
	return nil
}

func bestContour(frame gocv.Mat, minArea float64) []image.Point {
	cnts := gocv.FindContours(frame, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var (
		bestCnt  []image.Point
		bestArea = minArea
	)
	for _, cnt := range cnts {
		if area := gocv.ContourArea(cnt); area > bestArea {
			bestArea = area
			bestCnt = cnt
		}
	}
	return bestCnt
}

// middle calculates the middle x and y of a rectangle.
func middle(rect image.Rectangle) (x int, y int) {
	return (rect.Max.X-rect.Min.X)/2 + rect.Min.X, (rect.Max.Y-rect.Min.Y)/2 + rect.Min.Y
}
