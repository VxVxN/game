package utils

import (
	"math/rand"
)

func RandomIntByRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomFloatBetween0And1() float64 {
	return 0 + rand.Float64()*(1-0)
}
