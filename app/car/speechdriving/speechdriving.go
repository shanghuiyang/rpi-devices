package speechdriving

type SpeechDriving interface {
	Start()
	Stop()
	InDriving() bool
}
