package ui

type Component interface {
	GetBounds() Bounds
	SetBounds(bounds Bounds)
	GetMinimumSize() Bounds
	GetMaximumSize() Bounds
	Render()
}
