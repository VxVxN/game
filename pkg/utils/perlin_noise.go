package utils

import (
	"math"
	"math/rand"
)

func Noise(x, y float32) float32 {
	// Coordinate left and top vertex square
	left := float32(math.Floor(float64(x)))
	top := float32(math.Floor(float64(y)))

	// Local coordinate
	localX := x - left
	localy := y - top

	topLeft := getRandomVector(x, y)
	topRight := getRandomVector(x+1, y)
	bottomLeft := getRandomVector(x, y+1)
	bottomRight := getRandomVector(x+1, y+1)

	// Vectors from vertices to a point inside
	dTopLeft := []float32{localX, localy}
	dTopRight := []float32{localX - 1, localy}
	dBottomLeft := []float32{localX, localy - 1}
	dBottomRight := []float32{localX - 1, localy - 1}

	// Scalar product
	tx1 := dot(dTopLeft, topLeft)
	tx2 := dot(dTopRight, topRight)
	bx1 := dot(dBottomLeft, bottomLeft)
	bx2 := dot(dBottomRight, bottomRight)

	// Parameters for non-linearity
	pointX := curve(localX)
	pointY := curve(localy)

	// Interpolation
	tx := lerp(tx1, tx2, pointX)
	bx := lerp(bx1, bx2, pointX)
	tb := lerp(tx, bx, pointY)
	return tb

}

func getRandomVector(x, y float32) []float32 {
	v := (int(x+y) + rand.Intn(2)) % 4
	switch v {
	case 0:
		return []float32{-1, 0}
	case 1:
		return []float32{1, 0}
	case 2:
		return []float32{0, 1}
	default:
		return []float32{0, -1}
	}
}

func dot(a []float32, b []float32) float32 {
	return a[0]*b[0] + b[1]*a[1]
}

func lerp(a, b, c float32) float32 {
	return a*(1-c) + b*c
}

func curve(t float32) float32 {
	return t * t * t * (t*(t*6-15) + 10)
}
