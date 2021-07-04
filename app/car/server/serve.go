package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/app/car/joystick"
	"github.com/shanghuiyang/rpi-devices/app/car/selfdriving"
	"github.com/shanghuiyang/rpi-devices/app/car/selftracking"
	"github.com/shanghuiyang/rpi-devices/app/car/speechdriving"

	"github.com/shanghuiyang/rpi-devices/cv"
	// "github.com/shanghuiyang/rpi-devices/cv/mock/cv"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/stianeikeland/go-rpio"
)

const logTag = "server"

func Start(cfg *Config) {
	if err := rpio.Open(); err != nil {
		log.Printf("failed to open rpio, error: %v", err)
		return
	}

	defer rpio.Close()

	l298n := dev.NewL298N(
		cfg.L298N.IN1Pin,
		cfg.L298N.IN2Pin,
		cfg.L298N.IN3Pin,
		cfg.L298N.IN4Pin,
		cfg.L298N.ENAPin,
		cfg.L298N.ENBPin)
	if l298n == nil {
		log.Panicf("[%v]failed to new L298N", logTag)
	}

	buz := dev.NewBuzzer(int8(cfg.BuzzerPin))
	if buz == nil {
		log.Panicf("[%v]failed to new buzzer", logTag)
	}

	led := dev.NewLed(cfg.LedPin)
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

	cam := dev.NewCamera()
	if cam == nil {
		log.Panicf("[%v]failed to new a camera", logTag)
	}

	car := car.NewCarImp(l298n, gy25, buz)
	if car == nil {
		log.Panicf("[%v]failed to new car", logTag)
	}

	if cfg.Joystick.Enabled {
		lc12s, err := dev.NewLC12S(cfg.Joystick.LC12SConfig.Dev, cfg.Joystick.LC12SConfig.Baud, cfg.Joystick.LC12SConfig.CS)
		if err != nil {
			log.Panicf("[%v]failed to new lc12s, error: %v", logTag, err)
		}
		joystick.Init(car, lc12s)
		go joystick.Start()
	}

	if cfg.SelfDriving.Enabled {
		selfdriving.Init(car, us100, sg90)
	}

	if cfg.SelfTracking.Enabled {
		t, err := cv.NewTracker(cfg.SelfTracking.LH, cfg.SelfTracking.LS, cfg.SelfTracking.LV, cfg.SelfTracking.HH, cfg.SelfTracking.HS, cfg.SelfTracking.HV)
		if err != nil {
			log.Panicf("[%v]failed to create tracker, error: %v", logTag, err)
		}
		s := cv.NewStreamer(cfg.SelfTracking.VideoHost)
		selftracking.Init(car, t, s)
	}

	if cfg.SpeechDriving.Enabled {
		// TODO
		// create asr, tts, imgr
		speechdriving.Init(car, us100, sg90, led, cam, nil, nil, nil)
	}

	if err := util.SetVolume(cfg.Volume); err != nil {
		log.Panicf("[%v]failed to create tracker, error: %v", logTag, err)
	}

	s := newService(car, led).
		EnableSelfDriving(cfg.SelfDriving.Enabled).
		EnabledSelfTracking(cfg.SelfTracking.Enabled).
		EnabledSpeechDriving(cfg.SpeechDriving.Enabled)

	if err := s.start(); err != nil {
		log.Panicf("[%v]failed to start server, error: %v", logTag, err)
	}
	if err := s.start(); err != nil {
		log.Panicf("[%v]failed to start server, error: %v", logTag, err)
	}
	log.Printf("[%v]service started", logTag)

	util.WaitQuit(func() {
		s.Stop()
		rpio.Close()
	})

	r := mux.NewRouter()

	// home
	r.HandleFunc("/", s.loadHomeHandler).Methods("GET")

	// car operation
	r.HandleFunc("/car/{op:[a-z]+}", s.opHandler).Methods("POST")

	// self-driving
	r.HandleFunc("/selfdriving/on", s.selfDrivingOnHandler).Methods("POST")
	r.HandleFunc("/selfdriving/off", s.selfDrivingOffHandler).Methods("POST")

	// self-tracking
	r.HandleFunc("/selftracking/on", s.selfTrackingOnHandler).Methods("POST")
	r.HandleFunc("/selftracking/off", s.selfTrackingOffHandler).Methods("POST")

	// speech-driving
	r.HandleFunc("/speechdriving/on", s.speechDrivingOnHandler).Methods("POST")
	r.HandleFunc("/speechdriving/off", s.speechDrivingOffHandler).Methods("POST")

	// self-nav
	r.HandleFunc("/selfnav/{lat:[0-9]+}/{lon:[0-9]+", s.selfNavOnHandler).Methods("POST")
	r.HandleFunc("/selfnav/off", s.selfNavOffHandler).Methods("POST")

	// set volume
	r.HandleFunc("/volume/{v:[0-9]+}", s.setVolumeHandler).Methods("POST")

	// music
	r.HandleFunc("/music/on", s.musicOnHandler).Methods("POST")
	r.HandleFunc("/music/off", s.musicOffHandler).Methods("POST")

	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicf("[%v]failed to start http server, error: %v", logTag, err)
	}
	log.Printf("[%v]http server stop", logTag)
}
