package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/shanghuiyang/imgr"
	"github.com/shanghuiyang/oauth"
	"github.com/shanghuiyang/rpi-devices/dev"
	"github.com/shanghuiyang/rpi-devices/util"
	"github.com/shanghuiyang/speech"
)

const (
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
	asr   speech.ASR
	tts   speech.TTS
	imgre imgr.Recognizer
	cam   dev.Camera
)

func main() {

	speechAuth := oauth.NewBaiduOauth(baiduSpeechAppKey, baiduSpeechSecretKey, oauth.NewCacheImp())
	imageAuth := oauth.NewBaiduOauth(baiduImgRecognitionAppKey, baiduImgRecognitionSecretKey, oauth.NewCacheImp())
	asr = speech.NewBaiduASR(speechAuth)
	tts = speech.NewBaiduTTS(speechAuth)
	imgre = imgr.NewBaiduRecognizer(imageAuth)
	cam = dev.NewMotionCamera()

	for {
		log.Printf("[imgr]take photo")
		img, err := cam.Photo()
		if err != nil {
			log.Printf("[imgr]failed to take phote, error: %v", err)
			os.Exit(1)
		}
		util.PlayWav(wavLetMeThink)

		log.Printf("[imgr]recognize image")
		objname, err := recognize(img)
		if err != nil {
			log.Printf("[imgr]failed to recognize image, error: %v", err)
			util.PlayWav(wavIDontKnow)
			os.Exit(1)
		}
		log.Printf("[imgr]object: %v", objname)

		wav, err := tospeech("这是" + objname)
		if err != nil {
			log.Printf("[imgr]failed to tts, error: %v", err)
			os.Exit(1)
		}

		if err := util.PlayWav(wav); err != nil {
			log.Printf("[imgr]failed to play wav: %v, error: %v", wav, err)
			os.Exit(1)
		}

		time.Sleep(10 * time.Second)
	}
}

func recognize(image []byte) (string, error) {
	name, err := imgre.Recognize(image)
	if err != nil {
		return "", err
	}
	return name, nil
}

func tospeech(text string) (string, error) {
	data, err := tts.ToSpeech(text)
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
