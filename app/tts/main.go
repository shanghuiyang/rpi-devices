package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/speech"
	"github.com/shanghuiyang/rpi-devices/util"
)

const (
	baiduSpeechAppKey    = "your_speech_app_key"
	baiduSpeechSecretKey = "your_speech_secret_key"
	ttsWav               = "tts.wav"
	ipPattern            = "((000.000.000.000))"
)

type ttsServer struct {
	tts         *speech.TTS
	pageContext []byte
}

func main() {
	s := newTTSServer(baiduSpeechAppKey, baiduSpeechSecretKey)
	util.WaitQuit(func() {})
	if err := s.start(); err != nil {
		log.Printf("[tts]failed to start car server, error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func newTTSServer(appKey, secretKey string) *ttsServer {
	auth := oauth.New(appKey, secretKey, oauth.NewCacheMan())
	tts := speech.NewTTS(auth)
	return &ttsServer{
		tts: tts,
	}
}

func (s *ttsServer) start() error {
	if err := s.loadHomePage(); err != nil {
		return err
	}

	log.Printf("[tts]tts server started successfully")
	http.HandleFunc("/", s.handler)
	if err := http.ListenAndServe(":8082", nil); err != nil {
		return err
	}
	return nil
}

func (s *ttsServer) loadHomePage() error {
	data, err := ioutil.ReadFile("tts.html")
	if err != nil {
		return errors.New("internal error: failed to read car.html")
	}

	ip := util.GetIP()
	if ip == "" {
		return errors.New("internal error: failed to get ip")
	}

	rbuf := bytes.NewBuffer(data)
	wbuf := bytes.NewBuffer([]byte{})
	for {
		line, err := rbuf.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		s := string(line)
		if strings.Index(s, ipPattern) >= 0 {
			s = strings.Replace(s, ipPattern, ip, 1)
		}
		wbuf.Write([]byte(s))
	}
	s.pageContext = wbuf.Bytes()
	return nil
}

func (s *ttsServer) handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write(s.pageContext)
	case "POST":
		txt := r.FormValue("text")
		log.Printf("[tts]receive text: %v", txt)
		go s.playText(txt)

	}
}

func (s *ttsServer) playText(text string) error {
	if text == "" {
		return nil
	}
	log.Printf("[tts]converting text to speech...")
	wav, err := s.toSpeech(text)
	if err != nil {
		log.Printf("[tts]failed to convert test to speech, error: %v", err)
		return err
	}
	log.Printf("[tts]converted in success")

	log.Printf("[tts]playing wav...")
	if err := s.play(wav); err != nil {
		log.Printf("[tts]failed to play wav: %v, error: %v", wav, err)
		return err
	}
	log.Printf("[tts]played in success")
	return nil
}

func (s *ttsServer) toSpeech(text string) (string, error) {
	data, err := s.tts.ToSpeech(text)
	if err != nil {
		log.Printf("[tts]failed to convert text to speech, error: %v", err)
		return "", err
	}

	if err := ioutil.WriteFile(ttsWav, data, 0644); err != nil {
		log.Printf("[tts]failed to save %v, error: %v", ttsWav, err)
		return "", err
	}
	return ttsWav, nil
}

func (s *ttsServer) play(wav string) error {
	// return nil
	// aplay test.wav
	cmd := exec.Command("aplay", wav)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[tts]failed to exec aplay, output: %v, error: %v", string(out), err)
		return err
	}
	return nil
}
