package dev

import (
	"fmt"
	"testing"

	"github.com/shanghuiyang/a-star/astar"
	"github.com/stretchr/testify/assert"
)

func TestTurnPoints(t *testing.T) {
	path := astar.PList{
		&astar.Point{X: 1, Y: 1},
		&astar.Point{X: 2, Y: 2},
		&astar.Point{X: 3, Y: 3},
		&astar.Point{X: 4, Y: 5},
		&astar.Point{X: 5, Y: 7},
		&astar.Point{X: 7, Y: 9},
	}
	c := Car{}
	actual := c.turnPoints(path)
	expected := astar.PList{
		&astar.Point{
			X: 3,
			Y: 3,
		},
		&astar.Point{
			X: 5,
			Y: 7,
		},
		&astar.Point{
			X: 7,
			Y: 9,
		},
	}
	fmt.Println(actual)
	assert.Equal(t, expected, actual)
}
