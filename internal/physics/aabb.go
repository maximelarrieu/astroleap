package physics

// Rect represents an axis-aligned bounding box.
type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

// Overlaps returns true if this rect intersects another rect.
func (r Rect) Overlaps(o Rect) bool {
	return r.X < o.X+o.W &&
		r.X+r.W > o.X &&
		r.Y < o.Y+o.H &&
		r.Y+r.H > o.Y
}

// ContainsPoint returns true if the point (px, py) is inside the rect.
func (r Rect) ContainsPoint(px, py float64) bool {
	return px >= r.X && px <= r.X+r.W && py >= r.Y && py <= r.Y+r.H
}

// IsEmpty returns true if the rectangle has zero width or height.
func (r Rect) IsEmpty() bool {
	return r.W <= 0 || r.H <= 0
}
