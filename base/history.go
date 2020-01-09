package base

import (
	"errors"
)

// ErrEmpty ...
var ErrEmpty = errors.New("empty")

// History ...
type History struct {
	contains []interface{}
	size     int
	index    int
	full     bool
}

// NewHistory ...
func NewHistory(size int) *History {
	h := &History{
		contains: make([]interface{}, size),
		size:     size,
		index:    0,
		full:     false,
	}
	for i := range h.contains {
		h.contains[i] = int(0)
	}
	return h
}

// Add ...
func (h *History) Add(v interface{}) {
	h.contains[h.index] = v
	h.index++
	if h.index == h.size {
		h.index = 0
		h.full = true
	}
}

// Avg ...
func (h *History) Avg() (float64, error) {
	if h.index == 0 && !h.full {
		return 0, ErrEmpty
	}
	var sum float64
	for _, v := range h.contains {
		switch v.(type) {
		case uint8:
			sum += float64(v.(uint8))
		case uint16:
			sum += float64(v.(uint16))
		case uint32:
			sum += float64(v.(uint32))
		case uint64:
			sum += float64(v.(uint64))
		case int:
			sum += float64(v.(int))
		case int8:
			sum += float64(v.(int8))
		case int16:
			sum += float64(v.(int16))
		case int32:
			sum += float64(v.(int32))
		case int64:
			sum += float64(v.(int64))
		case float32:
			sum += float64(v.(float32))
		case float64:
			sum += v.(float64)
		default:
			return 0, errors.New("the element isn't numerical")
		}
	}
	n := h.index
	if h.full {
		n = h.size
	}
	return sum / float64(n), nil
}
