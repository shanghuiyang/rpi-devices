package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSServer(t *testing.T) {
	s := sserver{}
	assert.NotNil(t, s)
}
