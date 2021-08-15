package selftracking

import "gocv.io/x/gocv"

type SelfTracking interface {
	Start(chImg chan *gocv.Mat)
	Stop()
	InTracking() bool
}
