package speechdriving

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/shanghuiyang/imgr"
	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/speech"
)

const (
	forward  operator = "forward"
	backward operator = "backward"
	left     operator = "left"
	right    operator = "right"
	stop     operator = "stop"
	turn     operator = "turn"
	roll     operator = "roll"
	scan     operator = "scan"

	letMeThinkWav = "let_me_think.wav"
	thisIsXWav    = "this_is_x.wav"
	iDontKnowWav  = "i_dont_know.wav"
	errorWav      = "error.wav"

	logTag = "speechdriving"
)

var (
	aheadAngles = []float64{0, -15, 0, 15}
)

type operator string

type SpeechDrivingImp struct {
	car       car.Car
	dmeter    dev.DistanceMeter
	servo     dev.Motor
	led       dev.Led
	camera    dev.Camera
	asr       speech.ASR
	tts       speech.TTS
	imgr      imgr.Recognizer
	lock      sync.Mutex
	inDriving bool
}

func NewSpeechDrivingImp(car car.Car, dmeter dev.DistanceMeter, servo dev.Motor, led dev.Led, cam dev.Camera, asr speech.ASR, tts speech.TTS, imgr imgr.Recognizer) *SpeechDrivingImp {
	servo.Roll(0)
	return &SpeechDrivingImp{
		car:       car,
		dmeter:    dmeter,
		servo:     servo,
		led:       led,
		camera:    cam,
		asr:       asr,
		tts:       tts,
		imgr:      imgr,
		inDriving: false,
	}
}

func (s *SpeechDrivingImp) Start() {
	if s.inDriving {
		return
	}

	var (
		op   = stop
		fwd  = false
		chOp = make(chan operator, 4)
		wg   sync.WaitGroup
	)

	s.inDriving = true
	wg.Add(1)
	go s.detectSpeech(chOp, &wg)
	for s.inDriving {
		select {
		case p := <-chOp:
			op = p
			for len(chOp) > 0 {
				<-chOp
			}
		default:
			// do nothing
		}
		log.Printf("[%v]op: %v", logTag, op)

		switch op {
		case forward:
			if !fwd {
				s.car.Forward()
				fwd = true
				go s.lookingForObs(chOp)
			}
			util.DelayMs(50)
			continue
		case backward:
			fwd = false
			s.car.Stop()
			util.DelayMs(20)
			s.car.Backward()
			util.DelayMs(600)
			chOp <- stop
			continue
		case left:
			fwd = false
			s.car.Stop()
			util.DelayMs(20)
			s.car.Turn(-90)
			util.DelayMs(20)
			chOp <- forward
			continue
		case right:
			fwd = false
			s.car.Stop()
			util.DelayMs(20)
			s.car.Turn(90)
			util.DelayMs(20)
			chOp <- forward
			continue
		case roll:
			fwd = false
			s.car.Left()
			util.DelaySec(3)
			chOp <- stop
			continue
		case stop:
			fwd = false
			s.car.Stop()
			util.DelayMs(500)
			continue
		}
	}
	s.car.Stop()
	wg.Wait()
	close(chOp)
}

func (s *SpeechDrivingImp) InDriving() bool {
	return s.inDriving
}

func (s *SpeechDrivingImp) Stop() {
	s.inDriving = false
}

func (s *SpeechDrivingImp) lookingForObs(chOp chan operator) {
	for s.inDriving {
		for _, angle := range aheadAngles {
			s.servo.Roll(angle)
			util.DelayMs(70)
			d, err := s.dmeter.Dist()
			for i := 0; err != nil && i < 3; i++ {
				util.DelayMs(100)
				d, err = s.dmeter.Dist()
			}
			if err != nil {
				continue
			}

			if d < 20 {
				chOp <- backward
				return
			}
			if d < 40 {
				chOp <- stop
				return
			}
		}
	}
}

func (s *SpeechDrivingImp) detectSpeech(chOp chan operator, wg *sync.WaitGroup) {
	defer wg.Done()

	for s.inDriving {
		log.Printf("[%v]start recording", logTag)
		go s.led.On()
		wav := "car.wav"
		if err := util.Record(2, wav); err != nil {
			log.Printf("[%v]failed to record the speech: %v", logTag, err)
			continue
		}
		go s.led.Off()
		log.Printf("[%v]stop recording", logTag)

		wavData, err := ioutil.ReadFile(wav)
		if err != nil {
			log.Printf("[%v]failed to read %v, error: %v", logTag, wav, err)
			continue
		}
		text, err := s.asr.ToText(wavData, speech.Wav)
		if err != nil {
			log.Printf("[%v]failed to recognize the speech, error: %v", logTag, err)
			continue
		}
		log.Printf("[%v]speech: %v", logTag, text)

		switch {
		case strings.Contains(text, "前"):
			chOp <- forward
		case strings.Contains(text, "后"):
			chOp <- backward
		case strings.Contains(text, "左"):
			chOp <- left
		case strings.Contains(text, "右"):
			chOp <- right
		case strings.Contains(text, "停"):
			chOp <- stop
		case strings.Contains(text, "转圈"):
			chOp <- roll
		case strings.Contains(text, "是什么"):
			s.recognize()
		case strings.Contains(text, "大声"):
			s.volumeUp()
		case strings.Contains(text, "小声"):
			s.volumeDown()
		case strings.Contains(text, "唱歌"):
			go util.PlayWav("./music/xiaomaolv.wav")
		default:
			// do nothing
		}
	}
}

func (s *SpeechDrivingImp) recognize() error {
	log.Printf("[%v]take photo", logTag)
	img, err := s.camera.Photo()
	if err != nil {
		log.Printf("[%v]failed to take phote, error: %v", logTag, err)
		return err
	}
	go util.PlayWav(letMeThinkWav)

	log.Printf("[%v]recognize image", logTag)
	util.DelaySec(1)
	objname, err := s.imgr.Recognize(img)
	if err != nil {
		log.Printf("[%v]failed to recognize image, error: %v", logTag, err)
		util.PlayWav(errorWav)
		return err
	}
	log.Printf("[%v]object: %v", logTag, objname)

	if err := s.playText("这是" + objname); err != nil {
		log.Printf("[%v]failed to play text, error: %v", logTag, err)
		return err
	}

	return nil
}

func (s *SpeechDrivingImp) playText(text string) error {
	wav, err := s.toSpeech(text)
	if err != nil {
		log.Printf("[%v]failed to tts, error: %v", logTag, err)
		return err
	}

	if err := util.PlayWav(wav); err != nil {
		log.Printf("[%v]failed to play wav: %v, error: %v", logTag, wav, err)
		return err
	}
	return nil
}

func (s *SpeechDrivingImp) toSpeech(text string) (string, error) {
	data, err := s.tts.ToSpeech(text)
	if err != nil {
		log.Printf("[%v]failed to convert text to speech, error: %v", logTag, err)
		return "", err
	}

	if err := ioutil.WriteFile(thisIsXWav, data, 0644); err != nil {
		log.Printf("[%v]failed to save %v, error: %v", logTag, thisIsXWav, err)
		return "", err
	}
	return thisIsXWav, nil
}

func (s *SpeechDrivingImp) volumeUp() {
	v, err := util.GetVolume()
	if err != nil {
		log.Printf("[%v]failed get current volume, error: %v", logTag, err)
		return
	}
	v += 10
	if v > 100 {
		v = 100
	}
	s.setvolume(v)
	go s.playText(fmt.Sprintf("音量%v%%", v))
}

func (s *SpeechDrivingImp) volumeDown() {
	v, err := util.GetVolume()
	if err != nil {
		log.Printf("[%v]failed get current volume, error: %v", logTag, err)
		return
	}
	v -= 10
	if v < 0 {
		v = 0
	}
	s.setvolume(v)
	go s.playText(fmt.Sprintf("音量%v%%", v))
}

func (s *SpeechDrivingImp) setvolume(v int) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := util.SetVolume(v); err != nil {
		log.Printf("[%v]failed to set volume, error: %v", logTag, err)
		return err
	}
	return nil
}
