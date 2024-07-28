package utils

import (
	"math/rand"

	"github.com/VxVxN/game/internal/base"
)

func RandomIntByRange(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func RandomFloat64ByRange(min, max float64) float64 {
	return float64(rand.Intn(int(max-min+1))) + min
}

func RandomFloatBetween0And1() float64 {
	return 0 + rand.Float64()*(1-0)
}

func CanAction(playerPosition, dstPosition base.Position) bool {
	return int(playerPosition.X+1) == int(dstPosition.X) && int(playerPosition.Y) == int(dstPosition.Y) ||
		int(playerPosition.X-1) == int(dstPosition.X) && int(playerPosition.Y) == int(dstPosition.Y) ||
		int(playerPosition.X) == int(dstPosition.X) && int(playerPosition.Y+1) == int(dstPosition.Y) ||
		int(playerPosition.X) == int(dstPosition.X) && int(playerPosition.Y-1) == int(dstPosition.Y) ||
		int(playerPosition.X) == int(dstPosition.X) && int(playerPosition.Y) == int(dstPosition.Y)
}

func DeleteElemSlice[T comparable](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
