package speechdriving

type ASR interface {
	ToText(wavFile string) (string, error)
}

type TTS interface {
	ToSpeech(text string) ([]byte, error)
}

type ImgRecognizer interface {
	Recognize(imageFile string) (string, error)
}
