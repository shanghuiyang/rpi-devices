package devices

// Camera ...
type Camera struct {
	steering *SG90
}

// NewCamera ...
func NewCamera(steering *SG90) *Camera {
	return &Camera{
		steering: steering,
	}
}

// Turn ...
func (c *Camera) Turn(angle int) {
	c.steering.Roll(angle)
}
