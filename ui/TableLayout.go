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
func (this TableLayout) GetBounds() Bounds {
	return this.Bounds
}

// SetBounds sets the bounds of the component
func (this TableLayout) SetBounds(bounds Bounds) {
	this.Bounds = bounds
	this.Layout()
}

// GetMinimumSize determines the minimum size of the component
func (this TableLayout) GetMinimumSize() Bounds {
	// TODO Implement
	return NewBounds(-1, -1, 0, 0)
}

// Render draws this component and all of its child components onto the active GL context
func (this TableLayout) Render() {
	if this.NeedsLayout {
		this.Layout()
	}
	for _, child := range this.Children {
		// TODO Modify the transformation matrix
		child.Component.Render()
	}
}

// Layout recalculates all of the positions and sizes of the child components
func (this TableLayout) Layout() {
	// TODO Implement
}

// Add adds an additional component to the layout with the given constraints
func (this TableLayout) Add(component Component, row int, col int, rowSpan int, colSpan int) {
	if maxRow := row + rowSpan; maxRow < len(this.Rows) {
		newRows := make([]TableLayoutSize, maxRow)
		copy(newRows, this.Rows)
		for i := len(this.Rows); i < maxRow; i++ {
			newRows[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		this.Rows = newRows
	}
	if maxCol := col + colSpan; maxCol < len(this.Cols) {
		newCols := make([]TableLayoutSize, maxCol)
		copy(newCols, this.Cols)
		for i := len(this.Cols); i < maxCol; i++ {
			newCols[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		this.Cols = newCols
	}
	this.Children = append(this.Children, TableLayoutChild{
		Component: component,
		Row:       row,
		Col:       col,
		RowSpan:   rowSpan,
		ColSpan:   colSpan,
	})
	this.NeedsLayout = true
}

// SetRowSize constrains the size of a row
func (this TableLayout) SetRowSize(row int, spacingType SpacingType, size float32) {
	if row < len(this.Rows) {
		newRows := make([]TableLayoutSize, row)
		copy(this.Rows, newRows)
		for i := len(this.Rows); i < row-1; i++ {
			newRows[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		this.Rows = newRows
	}
	this.Rows[row] = TableLayoutSize{
		SpacingType: spacingType,
		Size:        size,
	}
	this.NeedsLayout = true
}

// SetColSize constrains the size of a column
func (this TableLayout) SetColSize(col int, spacingType SpacingType, size float32) {
	if col < len(this.Cols) {
		newCols := make([]TableLayoutSize, col)
		copy(this.Cols, newCols)
		for i := len(this.Cols); i < col-1; i++ {
			newCols[i] = TableLayoutSize{
				SpacingType: Minimum,
			}
		}
		this.Cols = newCols
	}
	this.Cols[col] = TableLayoutSize{
		SpacingType: spacingType,
		Size:        size,
	}
	this.NeedsLayout = true
}
