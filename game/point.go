package game

import (
	"image"
	"math"
)

// Point is X, Y coordinates in 2D space, like image.Point but with floats
type Point struct {
	X, Y float64
}

func (p *Point) Pt() image.Point {
	return image.Pt(
		int(math.Round(p.X)),
		int(math.Round(p.Y)),
	)
}
