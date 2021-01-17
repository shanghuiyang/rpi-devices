package base

import (
	"fmt"
	"os/exec"
)

// PlayWav plays wav audio using aplay
func PlayWav(wav string) error {
	cmd := exec.Command("aplay", wav)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// StopWav stops to play wav
func StopWav() error {
	cmd := "sudo killall aplay"
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// PlayMp3 play mp3 audio using mpg123
// you need to install mpg123 first
func PlayMp3(mp3 string) error {
	cmd := "mpg123 -Z -q " + mp3
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// StopMp3 stops to play mp3
func StopMp3() error {
	cmd := "sudo killall mpg123"
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// Record records voice using arecord
// you need to setup a micro-phone and set it as default device
// arecord usage:
// -D:			device, use commond "$arecord -l" for viewing card number and device number like 1,0
// -d 3:		3 seconds
// -t wav:		wav type
// -r 16000:	Rate 16000 Hz
// -c 1:		1 channel
// -f S16_LE:	Signed 16 bit Little Endian
func Record(sec int, saveTo string) error {
	cmd := fmt.Sprintf(`sudo arecord -D "plughw:1,0" -d %v -t wav -r 16000 -c 1 -f S16_LE %v`, sec, saveTo)
	_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

// SetVolume sets the volume using amixer
func SetVolume(v int) error {
	// amixer -M set PCM 20%
	cmd := exec.Command("amixer", "-M", "set", "PCM", fmt.Sprintf("%v%%", v))
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
