package world

// taken from https://github.com/StephaneBunel/bresenham

// Line implements the Bresenham algorithm to draw a line
// and returns a Coords array of the points to draw.
// Taken from https://github.com/StephaneBunel/bresenham
func Line(a, b Coords) []Coords {
	var r []Coords
	var dx, dy, e, slope int
	x1, y1, x2, y2 := a.X, a.Y, b.X, b.Y

	// a -> b is equivalent to b -> a, so we simplify
	if x1 > x2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}

	dx, dy = x2-x1, y2-y1
	// dx can't be negative, but dy can, so we fix that
	if dy < 0 {
		dy = -dy
	}

	switch {
	// Is line a point ?
	case x1 == x2 && y1 == y2:
		r = append(r, a)
	// Is line a horizontal?
	case y1 == y2:
		for ; dx != 0; dx-- {
			r = append(r, Coords{X: x1, Y: y1})
			x1++
		}
		r = append(r, Coords{X: x1, Y: y1})
	// Is line a vertical ?
	case x1 == x2:
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		for ; dy != 0; dy-- {
			r = append(r, Coords{X: x1, Y: y1})
			y1++
		}
		r = append(r, Coords{X: x1, Y: y1})

	// Is line a diagonal ?
	case dx == dy:
		if y1 < y2 {
			for ; dx != 0; dx-- {
				r = append(r, Coords{X: x1, Y: y1})
				x1++
				y1++
			}
		} else {
			for ; dx != 0; dx-- {
				r = append(r, Coords{X: x1, Y: y1})
				x1++
				y1--
			}
		}
		r = append(r, Coords{X: x1, Y: y1})

	// wider than high?
	case dx > dy:
		if y1 < y2 {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				r = append(r, Coords{X: x1, Y: y1})
				x1++
				e -= dy
				if e < 0 {
					y1++
					e += slope
				}
			}
		} else {
			dy, e, slope = 2*dy, dx, 2*dx
			for ; dx != 0; dx-- {
				r = append(r, Coords{X: x1, Y: y1})
				x1++
				e -= dy
				if e < 0 {
					y1--
					e += slope
				}
			}
		}
		r = append(r, Coords{X: x1, Y: y1})
	// higher than wide.
	default:
		if y1 < y2 {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				r = append(r, Coords{X: x1, Y: y1})
				y1++
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		} else {
			dx, e, slope = 2*dx, dy, 2*dy
			for ; dy != 0; dy-- {
				r = append(r, Coords{X: x1, Y: y1})
				y1--
				e -= dx
				if e < 0 {
					x1++
					e += slope
				}
			}
		}
		r = append(r, Coords{X: x1, Y: y1})
	}
	return r
}

// Square uses Line to get two horizontal lines and two vertical lines.
// It then removes repeats.
func Square(a, b Coords) []Coords {
	var r []Coords
	r = append(r, Line(a, Coords{X: a.X, Y: b.Y})...)
	r = append(r, Line(a, Coords{X: b.X, Y: a.Y})...)
	r = append(r, Line(b, Coords{X: a.X, Y: b.Y})...)
	r = append(r, Line(b, Coords{X: b.X, Y: a.Y})...)
	r = Reduce(r)
	return r
}
