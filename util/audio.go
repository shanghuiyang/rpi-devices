package util

import (
	"fmt"
	"os/exec"
	"strings"
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

// GetVolume gets the volume using amixer
/*
command line: amixer -M get PCM
output:
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
Simple mixer control 'PCM',0
  Capabilities: pvolume pvolume-joined pswitch pswitch-joined
  Playback channels: Mono
  Limits: Playback 0 - 255
  Mono: Playback 115 [45%] [-127.55dB] [on]
// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
*/
func GetVolume() (int, error) {
	// amixer -M set PCM 20%
	cmd := exec.Command("amixer", "-M", "get", "PCM")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, err
	}

	s := string(output)
	ss := strings.Split(s, "\n")
	n := len(ss)
	if n == 0 {
		return 0, fmt.Errorf("unexpected output from amixer, output: %v", s)
	}
	
	lastline := ss[n-2]
	items := strings.Split(lastline, " ")
	if len(items) != 8 {
		return 0, fmt.Errorf("unexpected output from amixer, last line: %v", lastline)
	}
	
	var v int
	if _, err := fmt.Sscanf(items[5], "[%d%%]", &v); err != nil {
		return 0, err
	}
	return v, nil
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
