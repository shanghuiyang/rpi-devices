package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/shanghuiyang/go-speech/asr"
	"github.com/shanghuiyang/go-speech/oauth"
	"github.com/shanghuiyang/go-speech/tts"
	imgoauth "github.com/shanghuiyang/image-recognizer/oauth"
	"github.com/shanghuiyang/image-recognizer/recognizer"
	"github.com/shanghuiyang/rpi-devices/dev"
)

const (
	imageFile     = "/var/lib/motion/lastsnap.jpg"
	wavLetMeThink = "let_me_think.wav"
	wavThisIsX    = "this_is_x.wav"
	wavIDontKnow  = "i_dont_know.wav"

	// replace your_app_key and your_secret_key with yours
	baiduSpeechAppKey    = "your_speech_app_key"
	baiduSpeechSecretKey = "your_speech_secret_key"

	baiduImgRecognitionAppKey    = "your_image_recognition_app_key"
	baiduImgRecognitionSecretKey = "your_image_recognition_secrect_key"
)

var (
	asrEng *asr.Engine
	ttsEng *tts.Engine
	imgr   *recognizer.Recognizer
	cam    *dev.Camera
)

func main() {

	speechOauth := oauth.New(baiduSpeechAppKey, baiduSpeechSecretKey, oauth.NewCacheMan())
	asrEng = asr.NewEngine(speechOauth)
	ttsEng = tts.NewEngine(speechOauth)

	imageOauth := imgoauth.New(baiduImgRecognitionAppKey, baiduImgRecognitionSecretKey, imgoauth.NewCacheMan())
	imgr = recognizer.New(imageOauth)

	cam = dev.NewCamera()

	go play(wavLetMeThink)

	log.Printf("[imgr]take photo")
	cam.TakePhoto()

	log.Printf("[imgr]recognize object")
	objname, err := recognize(imageFile)
	if err != nil {
		log.Printf("[imgr]failed to recognize object, error: %v", err)
		play(wavIDontKnow)
		os.Exit(1)
	}
	log.Printf("[imgr]object: %v", objname)

	wav, err := tospeech("这是" + objname)
	if err != nil {
		log.Printf("[imgr]failed to tts, error: %v", err)
		os.Exit(1)
	}

	if err := play(wav); err != nil {
		log.Printf("[imgr]failed to play wav: %v, error: %v", wav, err)
		os.Exit(1)
	}

	os.Exit(0)
}

func play(wav string) error {
	// omxplayer -o local test.wav
	cmd := exec.Command("omxplayer", "-o", "local", wav)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[imgr]failed to exec omxplayer, output: %v, error: %v", string(out), err)
		return err
	}
	return nil
}

func recognize(image string) (string, error) {
	name, err := imgr.Recognize(image)
	if err != nil {
		return "", err
	}
	return name, nil
}

func tospeech(text string) (string, error) {
	data, err := ttsEng.ToSpeech(text)
	if err != nil {
		log.Printf("[imgr]failed to convert text to speech, error: %v", err)
		return "", err
	}

	if err := ioutil.WriteFile(wavThisIsX, data, 0644); err != nil {
		log.Printf("[imgr]failed to save %v, error: %v", wavThisIsX, err)
		return "", err
	}
	return wavThisIsX, nil
}
