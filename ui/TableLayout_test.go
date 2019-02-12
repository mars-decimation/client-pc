package ui

import (
	"testing"
)

type Box struct {
	Bounds      Bounds
	MinimumSize Bounds
}

func (box Box) GetBounds() Bounds {
	return box.Bounds
}

func (box *Box) SetBounds(bounds Bounds) {
	box.Bounds = bounds
}

func (box Box) GetMinimumSize() Bounds {
	return box.MinimumSize
}

func (box Box) Render() {
}

func (box Box) String() string {
	return box.Bounds.String()
}

func NewBox(width float32, height float32) Box {
	return Box{
		Bounds:      NewBounds(-1, -1, -1, -1),
		MinimumSize: NewBounds(-1, -1, width, height),
	}
}

func CreateBox(layout *TableLayout, width float32, height float32, row int, col int, rowSpan int, colSpan int) *Box {
	box := NewBox(width, height)
	layout.Add(&box, row, col, rowSpan, colSpan)
	return &box
}

func CheckBox(t *testing.T, box *Box, x float32, y float32, width float32, height float32) {
	bounds := box.GetBounds()
	if bounds.X != x {
		t.Error("Invalid x position")
	}
	if bounds.Y != y {
		t.Error("Invalid y position")
	}
	if bounds.Width != width {
		t.Error("Invalid width")
	}
	if bounds.Height != height {
		t.Error("Invalid height")
	}
}

func TestSimpleLayout(t *testing.T) {
	layout := NewTableLayout()
	a := CreateBox(&layout, 20, 10, 0, 0, 1, 1)
	b := CreateBox(&layout, 15, 10, 0, 1, 1, 1)
	c := CreateBox(&layout, 10, 10, 0, 2, 1, 1)
	d := CreateBox(&layout, 60, 10, 1, 0, 1, 2)
	e := CreateBox(&layout, 20, 10, 1, 2, 1, 1)
	f := CreateBox(&layout, 10, 10, 2, 0, 1, 1)
	g := CreateBox(&layout, 80, 10, 2, 1, 1, 2)
	layout.Layout()
	CheckBox(t, a, 0, 0, 20, 10)
	CheckBox(t, b, 20, 0, 60, 10)
	CheckBox(t, c, 80, 0, 20, 10)
	CheckBox(t, d, 0, 10, 80, 10)
	CheckBox(t, e, 80, 10, 20, 10)
	CheckBox(t, f, 0, 20, 20, 10)
	CheckBox(t, g, 20, 20, 80, 10)
}
