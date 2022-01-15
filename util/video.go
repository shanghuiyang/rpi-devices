package util

import (
	"fmt"
	"os/exec"
)

// StopMotion ...
func StopMotion() error {
	cmd := "sudo killall motion"
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	DelaySec(1)
	return nil
}

// StartMotion start the Motion Service
// if config file doesn't be provided, the default config file will be used
// the default config file locates /etc/motion/motion.conf
func StartMotion(config ...string) error {
	cmd := "sudo motion"
	if len(config) > 0 {
		cmd = fmt.Sprintf("%v -c %v", cmd, config[0])
	}
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	DelaySec(1)
	return nil
}
