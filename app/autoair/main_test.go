package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	a := autoAir{}
	assert.NotNil(t, a)
}
