package ui

type SpacingType int

const (
	Absolute SpacingType = 0
	Minimum  SpacingType = 1
	Percent  SpacingType = 2
)

type TableLayoutSize struct {
	SpacingType SpacingType
	Size        float32
}

type TableLayoutChild struct {
	Component Component
	Row       int
	Col       int
	RowSpan   int
	ColSpan   int
}

type TableLayout struct {
	Bounds      Bounds
	Rows        []TableLayoutSize
	Cols        []TableLayoutSize
	Children    []TableLayoutChild
	NeedsLayout bool
}

func NewTableLayout() TableLayout {
	return TableLayout{
		Bounds:   NewBounds(0, 0, 0, 0),
		Rows:     make([]TableLayoutSize, 0),
		Cols:     make([]TableLayoutSize, 0),
		Children: make([]TableLayoutChild, 0),
	}
}

func (this TableLayout) GetBounds() Bounds {
	return this.Bounds
}

func (this TableLayout) SetBounds(bounds Bounds) {
	this.Bounds = bounds
	this.Layout()
}

func (this TableLayout) GetMinimumSize() Bounds {
	// TODO Implement
	return NewBounds(-1, -1, 0, 0)
}

func (this TableLayout) GetMaximumSize() Bounds {
	// TODO Implement
	return NewBounds(-1, -1, 0, 0)
}

func (this TableLayout) Render() {
	if this.NeedsLayout {
		this.Layout()
	}
	for _, child := range this.Children {
		// TODO Modify the transformation matrix
		child.Component.Render()
	}
}

func (this TableLayout) Layout() {
	// TODO Implement
}

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
