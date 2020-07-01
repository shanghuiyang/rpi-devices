package dev

import (
	"os/exec"
)

const (
	imageFile = "/var/lib/motion/lastsnap.jpg"
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
// you can change the directory in the config of motion.
// the config file is in /etc/motion/motion.conf
// and you need to change the 'snapshot_filename' to 'lastsnap.jpg'
// after changing, this will be look like:
// -----------------------------------------
// target_dir /var/lib/motion
// snapshot_filename lastsnap.jpg
// -----------------------------------------
func (c *Camera) TakePhoto() (string, error) {
	// curl -s -o /dev/null http://localhost:8088/0/action/snapshot
	cmd := exec.Command("curl", "-s", "-o", "/dev/null", "http://localhost:8088/0/action/snapshot")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return imageFile, nil
}
