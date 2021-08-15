package selfdriving

type SelfDriving interface {
	Start()
	Stop()
	InDrving() bool
}
