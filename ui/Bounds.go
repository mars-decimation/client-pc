package ui

type Bounds struct {
	x      float32
	y      float32
	width  float32
	height float32
}

func NewBounds(x float32, y float32, width float32, height float32) Bounds {
	return Bounds{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}
