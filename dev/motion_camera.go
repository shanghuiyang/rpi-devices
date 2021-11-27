/*
MotionCamera is the camera who is controlled by Motion.
More details about Motion, please ref to:
https://motion-project.github.io/index.html
*/

package dev

import (
	"io/ioutil"
	"os/exec"
)

const (
	imageFile = "/var/lib/motion/lastsnap.jpg"
)

// MotionCamera implements Camera interface
type MotionCamera struct {
	// steering *SG90
}

// NewMotionCamera ...
func NewMotionCamera() *MotionCamera {
	return &MotionCamera{}
}

// Photo takes a photo using motion service.
// in default, the created photo file will be in folder: /var/lib/motion
// you can change the directory in the config of motion.
// the config file is in /etc/motion/motion.conf
// and you need to change the 'snapshot_filename' to 'lastsnap.jpg'
// you also need to make sure webcontrol_port=8088 for this function working.
// after changing, this will be look like:
// -----------------------------------------
// target_dir /var/lib/motion
// snapshot_filename lastsnap.jpg
// webcontrol_port 8088
// -----------------------------------------
func (c *MotionCamera) Photo() ([]byte, error) {
	cmd := exec.Command("curl", "-s", "-o", "/dev/null", "http://localhost:8088/0/action/snapshot")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	img, err := ioutil.ReadFile(imageFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}
