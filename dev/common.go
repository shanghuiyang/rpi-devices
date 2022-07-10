package dev

import (
	"os/exec"
	"strings"
	"time"
)

type (
	LogicLevel    int
	InterfaceType int
	rpiModel      string
)

type StepperMode int

const (
	FullMode StepperMode = iota
	HalfMode
	QuarterMode
	EighthMode
	SixteenthMode
)

const (
	// voice speed in cm/s
	voiceSpeed = 34000.0
)

const (
	Low  LogicLevel = 0
	High LogicLevel = 1
)

const (
	GPIO InterfaceType = iota
	I2C
	SPI
	UART
	USB
)

const (
	rpiUnknown rpiModel = "Raspberry Pi X Model"
	rpi0       rpiModel = "Raspberry Pi Zero Model"
	rpiA       rpiModel = "Raspberry Pi A Model"
	rpiB       rpiModel = "Raspberry Pi B Model"
	rpi2       rpiModel = "Raspberry Pi 2 Model"
	rpi3       rpiModel = "Raspberry Pi 3 Model"
	rpi4       rpiModel = "Raspberry Pi 4 Model"
)

func delayNs(d time.Duration) {
	time.Sleep(d * time.Nanosecond)
}

func delayUs(d time.Duration) {
	time.Sleep(d * time.Microsecond)
}

func delayMs(d time.Duration) {
	time.Sleep(d * time.Millisecond)
}

func delaySec(d time.Duration) {
	time.Sleep(d * time.Second)
}

func delayMin(d time.Duration) {
	time.Sleep(d * time.Minute)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func getRpiModel() rpiModel {
	cmd := exec.Command("cat", "/proc/device-tree/model")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	s := string(out)
	if strings.Contains(s, string(rpi0)) {
		return rpi0
	}
	if strings.Contains(s, string(rpiA)) {
		return rpiA
	}
	if strings.Contains(s, string(rpiB)) {
		return rpiB
	}
	if strings.Contains(s, string(rpi2)) {
		return rpi2
	}
	if strings.Contains(s, string(rpi3)) {
		return rpi3
	}
	if strings.Contains(s, string(rpi4)) {
		return rpi4
	}
	return rpiUnknown
}
