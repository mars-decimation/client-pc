package ui

import (
	"fmt"

	"github.com/willauld/lpsimplex"
)

// SpacingType represents what kind of spacing a row or column uses
type SpacingType int

const (
	// Absolute means the spacing value is the exact number of pixels the row or column should be
	Absolute SpacingType = 0
	// Minimum means the row or column should shrink to fit the contents
	Minimum SpacingType = 1
	// Percent means the spacing value is a proportion used to fill extra space with when the layout grows or shrinks
	Percent SpacingType = 2
)

// TableLayoutSize describes a row or column's spacing
type TableLayoutSize struct {
	SpacingType SpacingType
	Size        float32
}

// TableLayoutChild describes all of the layout data for a given child component on the layout
type TableLayoutChild struct {
	Component Component
	Row       int
	Col       int
	RowSpan   int
	ColSpan   int
}

// TableLayout is a component that lays out its children components using a table
type TableLayout struct {
	Bounds       Bounds
	MinSize      Bounds
	Rows         []TableLayoutSize
	Cols         []TableLayoutSize
	Children     []TableLayoutChild
	NeedsLayout  bool
	NeedsMinCalc bool
}

// NewTableLayout creates a new table layout component with default values
func NewTableLayout() TableLayout {
	return TableLayout{
		Bounds:   NewBounds(0, 0, 0, 0),
		Rows:     make([]TableLayoutSize, 0),
		Cols:     make([]TableLayoutSize, 0),
		Children: make([]TableLayoutChild, 0),
	}
}

// GetBounds determines the bounds of the component
func (layout TableLayout) GetBounds() Bounds {
	return layout.Bounds
}

// SetBounds sets the bounds of the component
func (layout *TableLayout) SetBounds(bounds Bounds) {
	layout.Bounds = bounds
	layout.Layout()
}

// GetMinimumSize determines the minimum size of the component
func (layout TableLayout) GetMinimumSize() Bounds {
	if layout.NeedsMinCalc {
		rowPos := layout.CalculateSmooshedLayout(layout.Rows, func(child TableLayoutChild) int {
			return child.Row
		}, func(child TableLayoutChild) int {
			return child.RowSpan
		}, func(bounds Bounds) float32 {
			return bounds.Height
		})
		layout.MinSize.X = 0
		for _, size := range rowPos {
			layout.MinSize.X += size
		}
		colPos := layout.CalculateSmooshedLayout(layout.Cols, func(child TableLayoutChild) int {
			return child.Col
		}, func(child TableLayoutChild) int {
			return child.ColSpan
		}, func(bounds Bounds) float32 {
			return bounds.Width
		})
		layout.MinSize.Y = 0
		for _, size := range colPos {
			layout.MinSize.Y += size
		}
		layout.NeedsMinCalc = false
	}
	return layout.MinSize
}

// Render draws this component and all of its child components onto the active GL context
func (layout *TableLayout) Render() {
	if layout.NeedsLayout {
		layout.Layout()
	}
	for _, child := range layout.Children {
		// TODO Modify the transformation matrix
		child.Component.Render()
	}
}

// CalculateSmooshedLayout determines the minimum size for either the rows or columns, depending on the arguments passed in
func (layout *TableLayout) CalculateSmooshedLayout(elements []TableLayoutSize, elementSelector func(TableLayoutChild) int, spanSelector func(TableLayoutChild) int, sizeSelector func(Bounds) float32) []float32 {
	numSizes := len(elements)
	numChildren := len(layout.Children)
	c := make([]float64, numSizes)
	for i := 0; i < numSizes; i++ {
		c[i] = 1
	}
	A := make([][]float64, numChildren)
	b := make([]float64, numChildren)
	An := make([][]float64, numChildren)
	bn := make([]float64, numChildren)
	for i := 0; i < numChildren; i++ {
		A[i] = make([]float64, numSizes)
		An[i] = make([]float64, numSizes)
		for j := 0; j < numSizes; j++ {
			if el := elementSelector(layout.Children[i]); el <= j && j < el+spanSelector(layout.Children[i]) {
				A[i][j] = -1
			} else {
				A[i][j] = 0
			}
			An[i][j] = 0
		}
		b[i] = -float64(sizeSelector(layout.Children[i].Component.GetMinimumSize()))
		bn[i] = 0
	}
	res := lpsimplex.LPSimplex(c, A, b, An, bn, nil, nil, false, 1000, 0.01, false)
	if !res.Success {
		fmt.Printf("Unable to solve layout: %s\n", res.Message)
		return nil
	}
	pos := make([]float32, numSizes+1)
	accum := 0.0
	for i, size := range elements {
		pos[i] = float32(accum)
		if size.SpacingType == Absolute {
			accum += float64(size.Size)
		} else {
			accum += res.X[i]
		}
	}
	pos[numSizes] = float32(accum)
	return pos
}

// CalculateLayout determines either the size of the rows or columns, depending on the arguments passed in
func (layout *TableLayout) CalculateLayout(elements []TableLayoutSize, elementSelector func(TableLayoutChild) int, spanSelector func(TableLayoutChild) int, sizeSelector func(Bounds) float32) []float32 {
	pos := layout.CalculateSmooshedLayout(elements, elementSelector, spanSelector, sizeSelector)
	return pos
}

// Layout recalculates all of the positions and sizes of the child components
func (layout *TableLayout) Layout() {
	rowPos := layout.CalculateLayout(layout.Rows, func(child TableLayoutChild) int {
		return child.Row
	}, func(child TableLayoutChild) int {
		return child.RowSpan
	}, func(bounds Bounds) float32 {
		return bounds.Height
	})
	colPos := layout.CalculateLayout(layout.Cols, func(child TableLayoutChild) int {
		return child.Col
	}, func(child TableLayoutChild) int {
		return child.ColSpan
	}, func(bounds Bounds) float32 {
		return bounds.Width
	})
	for _, child := range layout.Children {
		x := colPos[child.Col]
		y := rowPos[child.Row]
		child.Component.SetBounds(NewBounds(x, y, colPos[child.Col+child.ColSpan]-x, rowPos[child.Row+child.RowSpan]-y))
	}
	layout.NeedsLayout = false
}

// Add adds an additional component to the layout with the given constraints
func (layout *TableLayout) Add(component Component, row int, col int, rowSpan int, colSpan int) {
	if maxRow := row + rowSpan; maxRow > len(layout.Rows) {
		newRows := make([]TableLayoutSize, maxRow)
		copy(newRows, layout.Rows)
		for i := len(layout.Rows); i < maxRow; i++ {
			newRows[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		layout.Rows = newRows
	}
	if maxCol := col + colSpan; maxCol > len(layout.Cols) {
		newCols := make([]TableLayoutSize, maxCol)
		copy(newCols, layout.Cols)
		for i := len(layout.Cols); i < maxCol; i++ {
			newCols[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		layout.Cols = newCols
	}
	layout.Children = append(layout.Children, TableLayoutChild{
		Component: component,
		Row:       row,
		Col:       col,
		RowSpan:   rowSpan,
		ColSpan:   colSpan,
	})
	layout.NeedsLayout = true
	layout.NeedsMinCalc = true
}

// SetRowSize constrains the size of a row
func (layout *TableLayout) SetRowSize(row int, spacingType SpacingType, size float32) {
	if row < len(layout.Rows) {
		newRows := make([]TableLayoutSize, row)
		copy(layout.Rows, newRows)
		for i := len(layout.Rows); i < row-1; i++ {
			newRows[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		layout.Rows = newRows
	}
	layout.Rows[row] = TableLayoutSize{
		SpacingType: spacingType,
		Size:        size,
	}
	layout.NeedsLayout = true
	layout.NeedsMinCalc = true
}

// SetColSize constrains the size of a column
func (layout *TableLayout) SetColSize(col int, spacingType SpacingType, size float32) {
	if col < len(layout.Cols) {
		newCols := make([]TableLayoutSize, col)
		copy(layout.Cols, newCols)
		for i := len(layout.Cols); i < col-1; i++ {
			newCols[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		layout.Cols = newCols
	}
	layout.Cols[col] = TableLayoutSize{
		SpacingType: spacingType,
		Size:        size,
	}
	layout.NeedsLayout = true
	layout.NeedsMinCalc = true
}
