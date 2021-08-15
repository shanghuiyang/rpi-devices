package speechdriving

type SpeechDriving interface {
	Start()
	Stop()
	InDriving() bool
}

type ASR interface {
	ToText(wavFile string) (string, error)
}

type TTS interface {
	ToSpeech(text string) ([]byte, error)
}

type ImgRecognizer interface {
	Recognize(image []byte) (string, error)
}
