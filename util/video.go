package util

import (
	"fmt"
	"os/exec"
	"time"
)

// StopMotion ...
func StopMotion() error {
	cmd := "sudo killall motion"
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}

// StartMotion start the Motion Service
// if config file doesn't be provided, the default config file will be used
// the default config file locates /etc/motion/motion.conf
func StartMotion(config ...string) error {
	cmd := fmt.Sprintf("sudo motion")
	if len(config) > 0 {
		cmd = fmt.Sprintf("sudo motion -c %v", config[0])
	}
	if _, err := exec.Command("bash", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
