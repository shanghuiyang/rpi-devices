package dev

import (
	"os/exec"
)

// Camera ...
type Camera struct {
	// steering *SG90
}

// NewCamera ...
func NewCamera() *Camera {
	return &Camera{}
}

// TakePhoto take a photo using motion service.
// in default, the created photo file will be in folder: /var/lib/motion
// you can change the directory by changing the config of motion
func (c *Camera) TakePhoto() error {
	// curl -s -o /dev/null http://localhost:8088/0/action/snapshot
	cmd := exec.Command("curl", "-s", "-o", "/dev/null", "http://localhost:8088/0/action/snapshot")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
