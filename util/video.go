package util

import (
	"fmt"
	"os/exec"
	"time"
)

// StopMotion ...
func StopMotion() error {
	cmd := "sudo killall motion"
	exec.Command("bash", "-c", cmd).CombinedOutput()
	time.Sleep(1 * time.Second)
	return nil
}

// StartMotion ...
func StartMotion() error {
	cmd := fmt.Sprintf("sudo motion")
	_, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return nil
}
