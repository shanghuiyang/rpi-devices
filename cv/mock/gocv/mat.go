package gocv

import (
	"errors"
)

// Mat ...
type Mat struct{}

// NewMat ...
func NewMat() Mat {
	return Mat{}
}

// Clone ...
func (m *Mat) Clone() Mat {
	return Mat{}
}

// Empty ...
func (m *Mat) Empty() bool {
	return true
}

// Close ...
func (m *Mat) Close() error {
	return errors.New("not implement")
}
