package ui

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
	Bounds      Bounds
	Rows        []TableLayoutSize
	Cols        []TableLayoutSize
	Children    []TableLayoutChild
	NeedsLayout bool
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
func (layout TableLayout) SetBounds(bounds Bounds) {
	layout.Bounds = bounds
	layout.Layout()
}

// GetMinimumSize determines the minimum size of the component
func (layout TableLayout) GetMinimumSize() Bounds {
	// TODO Implement
	return NewBounds(-1, -1, 0, 0)
}

// Render draws this component and all of its child components onto the active GL context
func (layout TableLayout) Render() {
	if layout.NeedsLayout {
		layout.Layout()
	}
	for _, child := range layout.Children {
		// TODO Modify the transformation matrix
		child.Component.Render()
	}
}

// Layout recalculates all of the positions and sizes of the child components
func (layout TableLayout) Layout() {
	// TODO Implement
}

// Add adds an additional component to the layout with the given constraints
func (layout TableLayout) Add(component Component, row int, col int, rowSpan int, colSpan int) {
	if maxRow := row + rowSpan; maxRow < len(layout.Rows) {
		newRows := make([]TableLayoutSize, maxRow)
		copy(newRows, layout.Rows)
		for i := len(layout.Rows); i < maxRow; i++ {
			newRows[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		layout.Rows = newRows
	}
	if maxCol := col + colSpan; maxCol < len(layout.Cols) {
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
}

// SetRowSize constrains the size of a row
func (layout TableLayout) SetRowSize(row int, spacingType SpacingType, size float32) {
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
}

// SetColSize constrains the size of a column
func (layout TableLayout) SetColSize(col int, spacingType SpacingType, size float32) {
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
}
