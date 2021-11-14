package server

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/shanghuiyang/astar/tilemap"
	"github.com/shanghuiyang/imgr"
	"github.com/shanghuiyang/oauth"
	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/app/car/joystick"
	"github.com/shanghuiyang/rpi-devices/app/car/selfdriving"
	"github.com/shanghuiyang/rpi-devices/app/car/selfnav"
	"github.com/shanghuiyang/rpi-devices/app/car/selftracking"
	"github.com/shanghuiyang/rpi-devices/app/car/speechdriving"
	"github.com/shanghuiyang/rpi-devices/cv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/speech"
)

const (
	forward  Op = "forward"
	backward Op = "backward"
	left     Op = "left"
	right    Op = "right"
	stop     Op = "stop"
	beep     Op = "beep"

	chSize                  = 8
	defaultVolume           = 40
	defaultSpeed            = 30
	defaultHost             = ":8080"
	defaultVideoHost        = ":8081"
	defaultTrackingVideoURL = ":8082/tracking"
)

type Op string

func init() {
	var err error
	pageContext, err = ioutil.ReadFile("home.html")
	if err != nil {
		log.Panicf("failed to load home page, error: %v", err)
	}

	ip = util.GetIP()
	if ip == "" {
		log.Panicf("failed to get ip address")
	}
}

type service struct {
	cfg        *Config
	car        car.Car
	led        dev.Led
	ledBlinked bool
	onMusic    bool
	chOp       chan Op

	selfdriving   selfdriving.SelfDriving
	selftracking  selftracking.SelfTracking
	selfnav       selfnav.SelfNav
	speechdriving speechdriving.SpeechDriving
	joystick      joystick.Joystick
}

func newService(cfg *Config) (*service, error) {
	if cfg.Speed == 0 {
		cfg.Speed = defaultSpeed
	}
	if cfg.Volume == 0 {
		cfg.Volume = defaultVolume
	}
	if cfg.Host == "" {
		cfg.Host = defaultHost
	}
	if cfg.VideoHost == "" {
		cfg.Host = defaultVideoHost
	}
	if cfg.SelfTracking.VideoURL == "" {
		cfg.SelfTracking.VideoURL = defaultTrackingVideoURL
	}

	l298n := dev.NewL298N(
		cfg.L298N.IN1Pin,
		cfg.L298N.IN2Pin,
		cfg.L298N.IN3Pin,
		cfg.L298N.IN4Pin,
		cfg.L298N.ENAPin,
		cfg.L298N.ENBPin)
	if l298n == nil {
		return nil, errors.New("failed to new L298N")
	}

	buz := dev.NewBuzzerImp(cfg.BuzzerPin, true)
	if buz == nil {
		return nil, errors.New("failed to new buzzer")
	}

	led := dev.NewLedImp(cfg.LedPin)
	if led == nil {
		log.Panicf("[%v]failed to new led", logTag)
	}

	sg90 := dev.NewSG90(cfg.SG90DataPin)
	if sg90 == nil {
		log.Panicf("[%v]failed to new sg90", logTag)
	}

	us100 := dev.NewUS100(&dev.US100Config{
		Mode: dev.UartMode,
		Dev:  cfg.US100.Dev,
		Baud: cfg.US100.Baud,
	})
	if us100 == nil {
		log.Panicf("[%v]failed to new us100", logTag)
	}

	gy25 := dev.NewGY25(cfg.GY25.Dev, cfg.GY25.Baud)
	if gy25 == nil {
		log.Panicf("[%v]failed to new gy-25", logTag)
	}

	cam := dev.NewMotionCamera()
	if cam == nil {
		log.Panicf("[%v]failed to new a camera", logTag)
	}

	car := car.NewCarImp(l298n, gy25, buz)
	if car == nil {
		log.Panicf("[%v]failed to new car", logTag)
	}
	car.Speed(cfg.Speed)

	s := &service{
		cfg:        cfg,
		car:        car,
		led:        led,
		ledBlinked: true,
		onMusic:    false,
		chOp:       make(chan Op, chSize),
	}

	if cfg.Joystick.Enabled {
		lc12s, err := dev.NewLC12S(cfg.Joystick.LC12SConfig.Dev, cfg.Joystick.LC12SConfig.Baud, cfg.Joystick.LC12SConfig.CS)
		if err != nil {
			log.Panicf("[%v]failed to new lc12s, error: %v", logTag, err)
		}
		s.joystick = joystick.NewJoystickImp(car, lc12s)
		go s.joystick.Start()
	}

	if cfg.SelfDriving.Enabled {
		s.selfdriving = selfdriving.NewSelfDrivingImp(car, us100, sg90)
	}

	if cfg.SelfTracking.Enabled {
		t, err := cv.NewTracker(cfg.SelfTracking.LH, cfg.SelfTracking.LS, cfg.SelfTracking.LV, cfg.SelfTracking.HH, cfg.SelfTracking.HS, cfg.SelfTracking.HV)
		if err != nil {
			log.Panicf("[%v]failed to create tracker, error: %v", logTag, err)
		}
		st, err := util.NewStreamer(cfg.SelfTracking.VideoURL)
		if err != nil {
			log.Panicf("[%v]failed to create streamer, error: %v", logTag, err)
		}
		s.selftracking = selftracking.NewSelfTrackingImp(car, t, st)
	}

	if cfg.SpeechDriving.Enabled {
		speechAuth := oauth.NewBaiduOauth(cfg.BaiduAPIConfig.Speech.APIKey, cfg.BaiduAPIConfig.Speech.SecretKey, oauth.NewCacheImp())
		imgAuth := oauth.NewBaiduOauth(cfg.BaiduAPIConfig.Image.APIKey, cfg.BaiduAPIConfig.Image.SecretKey, oauth.NewCacheImp())
		asr := speech.NewBaiduASR(speechAuth)
		tts := speech.NewBaiduTTS(speechAuth)
		imgr := imgr.NewBaiduRecognizer(imgAuth)
		s.speechdriving = speechdriving.NewSpeechDrivingImp(car, us100, sg90, led, cam, asr, tts, imgr)
	}

	if cfg.SelfNav.Enabled {
		gps, err := dev.NewNeo6mGPS(cfg.SelfNav.GPSConfig.Dev, cfg.SelfNav.GPSConfig.Baud)
		if err != nil {
			log.Panicf("[%v]failed to create gps, error: %v", logTag, err)
		}
		data, err := ioutil.ReadFile(cfg.SelfNav.TileMapConfig.MapFile)
		if err != nil {
			log.Panicf("[%v]failed to read map file: %v, errror: %v", logTag, cfg.SelfNav.TileMapConfig.MapFile, err)
		}
		m := tilemap.BuildFromStr(string(data))
		s.selfnav = selfnav.NewSelfNavImp(car, gps, m, cfg.SelfNav.TileMapConfig.Box, cfg.SelfNav.TileMapConfig.GridSize)
	}

	if err := util.SetVolume(cfg.Volume); err != nil {
		log.Panicf("[%v]failed to create tracker, error: %v", logTag, err)
	}

	return s, nil
}

func (s *service) start() error {
	go s.operate()
	go s.ledBlink()
	return nil
}

// Stop ...
func (s *service) Shutdown() error {
	s.ledBlinked = false
	close(s.chOp)
	s.car.Stop()
	s.led.Off()
	return nil
}

func (s *service) operate() {
	for op := range s.chOp {
		log.Printf("[car]op: %v", op)
		s.car.Speed(s.cfg.Speed)
		switch op {
		case forward:
			s.car.Forward()
		case backward:
			s.car.Backward()
		case left:
			s.car.Left()
		case right:
			s.car.Right()
		case stop:
			s.car.Stop()
		case beep:
			go s.car.Beep(3, 100)
		default:
			log.Printf("[car]invalid op")
		}
	}
}

func (s *service) ledBlink() {
	for s.ledBlinked {
		s.led.Blink(1, 1000)
	}
}
