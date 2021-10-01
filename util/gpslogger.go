package util

import (
	"fmt"
	"os"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Close()
}

// GPSLogger ...
type GPSLogger struct {
	f      *os.File
	chLine chan string
}

// NewGPSLogger ...
func NewGPSLogger(file string) (*GPSLogger, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	t := &GPSLogger{
		f:      f,
		chLine: make(chan string, 32),
	}
	go t.start()
	return t, nil
}

func (l *GPSLogger) start() {
	for line := range l.chLine {
		l.f.WriteString(line)
	}
}

func (l *GPSLogger) Printf(format string, v ...interface{}) {
	l.chLine <- fmt.Sprintf(format, v...)
}

// Close ...
func (l *GPSLogger) Close() {
	l.f.Close()
	close(l.chLine)
}

type NoopLogger struct{}

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (n *NoopLogger) Printf(format string, v ...interface{}) {}
func (l *NoopLogger) Close()                                 {}
