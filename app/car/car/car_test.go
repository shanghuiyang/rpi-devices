package car

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCar(t *testing.T) {
	car := New(&Config{})
	assert.NotNil(t, car)
}
