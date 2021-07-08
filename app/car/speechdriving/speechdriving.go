package speechdriving

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/shanghuiyang/rpi-devices/app/car/car"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
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
	aheadAngles = []int{0, -15, 0, 15}
)

type operator string

type SpeechDriving struct {
	car       car.Car
	dmeter    dev.DistMeter
	sg90      *dev.SG90
	led       *dev.Led
	camera    *dev.Camera
	asr       ASR
	tts       TTS
	imgr      ImgRecognizer
	lock      sync.Mutex
	ondriving bool
}

func New(c car.Car, d dev.DistMeter, sg90 *dev.SG90, l *dev.Led, cam *dev.Camera, a ASR, t TTS, imr ImgRecognizer) *SpeechDriving {
	sg90.Roll(0)
	return &SpeechDriving{
		car:       c,
		dmeter:    d,
		sg90:      sg90,
		led:       l,
		camera:    cam,
		asr:       a,
		tts:       t,
		imgr:      imr,
		ondriving: false,
	}
}

func (s *SpeechDriving) Start() {
	if s.ondriving {
		return
	}

	var (
		op   = stop
		fwd  = false
		chOp = make(chan operator, 4)
		wg   sync.WaitGroup
	)

	s.ondriving = true
	wg.Add(1)
	go s.detectSpeech(chOp, &wg)
	for s.ondriving {
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
			util.DelayMs(3000)
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

func (s *SpeechDriving) Status() bool {
	return s.ondriving
}

func (s *SpeechDriving) Stop() {
	s.ondriving = false
}

func (s *SpeechDriving) lookingForObs(chOp chan operator) {
	for s.ondriving {
		for _, angle := range aheadAngles {
			s.sg90.Roll(angle)
			util.DelayMs(70)
			d := s.dmeter.Dist()
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

func (s *SpeechDriving) detectSpeech(chOp chan operator, wg *sync.WaitGroup) {
	defer wg.Done()

	for s.ondriving {
		log.Printf("[%v]start recording", logTag)
		go s.led.On()
		wav := "car.wav"
		if err := util.Record(2, wav); err != nil {
			log.Printf("[%v]failed to record the speech: %v", logTag, err)
			continue
		}
		go s.led.Off()
		log.Printf("[%v]stop recording", logTag)

		text, err := s.asr.ToText(wav)
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

func (s *SpeechDriving) recognize() error {
	log.Printf("[%v]take photo", logTag)
	imagef, err := s.camera.TakePhoto()
	if err != nil {
		log.Printf("[%v]failed to take phote, error: %v", logTag, err)
		return err
	}
	util.PlayWav(letMeThinkWav)

	log.Printf("[%v]recognize image", logTag)
	objname, err := s.imgr.Recognize(imagef)
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

func (s *SpeechDriving) playText(text string) error {
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

func (s *SpeechDriving) toSpeech(text string) (string, error) {
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

func (s *SpeechDriving) volumeUp() {
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

func (s *SpeechDriving) volumeDown() {
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

func (s *SpeechDriving) setvolume(v int) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if err := util.SetVolume(v); err != nil {
		log.Printf("[%v]failed to set volume, error: %v", logTag, err)
		return err
	}
	return nil
}
