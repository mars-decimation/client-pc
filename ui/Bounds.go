package ui

import "fmt"

// Bounds reperesents a rectangle on the screen
type Bounds struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

// NewBounds creates a new Bounds object with the specified values
func NewBounds(x float32, y float32, width float32, height float32) Bounds {
	return Bounds{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// String converts the bounds to a string
func (bounds Bounds) String() string {
	return fmt.Sprintf("(%f, %f) - <%f, %f>", bounds.X, bounds.Y, bounds.Width, bounds.Height)
}
