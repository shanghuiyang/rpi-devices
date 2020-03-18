package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	v := &vmonitor{}
	assert.NotNil(t, v)
}
