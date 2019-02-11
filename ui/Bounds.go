package ui

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
