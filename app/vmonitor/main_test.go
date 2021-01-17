package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	v := &videoServer{}
	assert.NotNil(t, v)
}
