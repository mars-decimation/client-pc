package ui

// Component represents a component that can be drawn onto the screen
type Component interface {
	GetBounds() Bounds
	SetBounds(bounds Bounds)
	GetMinimumSize() Bounds
	Render()
}
