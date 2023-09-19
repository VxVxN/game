package base

type Position struct {
	X float64
	Y float64
}

func NewPosition(x, y float64) Position {
	return Position{
		X: x,
		Y: y,
	}
}
