package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	heart := heartBeat{}
	assert.NotNil(t, heart)
}
