//go:build linux && arm
// +build linux,arm

/*
NOTE:
init rpio on not raspberry-os will result to panic.
this happens when you run unit tests on ohter os like mac(development) or unbutu amd(travis-ci),
I use conditional compilation as workaround before I figure out a perfect solution.
*/
package dev

import (
	"github.com/stianeikeland/go-rpio/v4"
)

func init() {
	if err := rpio.Open(); err != nil {
		panic("failed to init rpio")
	}
}
