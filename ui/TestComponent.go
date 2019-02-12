package ui

import (
	"github.com/go-gl/gl/all-core/gl"
)

// RenderableBox is a renderable box
type RenderableBox struct {
	Bounds      Bounds
	MinimumSize Bounds
	Color       [4]float32
}

// GetBounds TODO describe
func (box RenderableBox) GetBounds() Bounds {
	return box.Bounds
}

// SetBounds TODO describe
func (box *RenderableBox) SetBounds(bounds Bounds) {
	box.Bounds = bounds
}

// GetMinimumSize TODO describe
func (box RenderableBox) GetMinimumSize() Bounds {
	return box.MinimumSize
}

// NewRenderableBox does what it says; nothing more, nothing less
func NewRenderableBox(minimumSize Bounds, color [4]float32) RenderableBox {
	return RenderableBox{
		MinimumSize: minimumSize,
		Color:       color,
	}
}

// CreateRenderableBox creates a new renderable box and adds it to the layout manager
func CreateRenderableBox(layout *TableLayout, width float32, height float32, row int, col int, rowSpan int, colSpan int, color [4]float32) *RenderableBox {
	box := NewRenderableBox(NewBounds(0, 0, width, height), color)
	layout.Add(&box, row, col, rowSpan, colSpan)
	return &box
}

// Render TODO describe
func (box RenderableBox) Render() {
	gl.Color4fv(&(box.Color[0]))
	gl.Begin(gl.QUADS)
	gl.Vertex2f(0, 0)
	gl.Vertex2f(0, box.Bounds.Height)
	gl.Vertex2f(box.Bounds.Width, box.Bounds.Height)
	gl.Vertex2f(box.Bounds.Width, 0)
	gl.End()
}
