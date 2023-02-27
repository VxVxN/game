package utils

import (
	"github.com/VxVxN/game/internal/base"
	"math"
)

type TileType int

const (
	Tree  TileType = iota
	Water TileType = iota
	Grass
)

type Chunk struct {
	ChunkID [2]int
	Map     map[int]map[int]TileType
}

func NewChunk(chunkSize, perlinSeed int, idChunk base.Position) Chunk {
	chunk := Chunk{ChunkID: [2]int{idChunk.X, idChunk.Y}}

	var chunkXMax, chunkYMax int
	chunkMap := make(map[int]map[int]TileType)
	chunkXMax = idChunk.X * chunkSize
	chunkYMax = idChunk.Y * chunkSize

	for x := 0; x < chunkXMax; x++ {
		chunkMap[x] = make(map[int]TileType)
		for y := 0; y < chunkYMax; y++ {
			SetTile(x, chunkMap, y, perlinSeed)
		}
	}

	chunk.Map = chunkMap
	return chunk
}

func SetTile(x int, chunkMap map[int]map[int]TileType, y int, perlinSeed int) {
	var tileType TileType
	perlinValue := Noise(float32(x)/float32(perlinSeed), float32(y)/float32(perlinSeed))
	switch {
	case perlinValue < -0.3:
		tileType = Tree
	case perlinValue >= -0.3:
		tileType = Grass
	case perlinValue > 0.5:
		//tileType = Water
	}
	chunkMap[x][y] = tileType
}

func GetChunkID(tileSize, x, y int) base.Position {
	tileX := float64(x) / float64(tileSize)
	tileY := float64(y) / float64(tileSize)

	var ChunkID base.Position
	if tileX < 0 {
		ChunkID.X = int(math.Floor(tileX / float64(tileSize)))
	} else {
		ChunkID.X = int(math.Ceil(tileX / float64(tileSize)))
	}
	if tileY < 0 {
		ChunkID.Y = int(math.Floor(tileY / float64(tileSize)))
	} else {
		ChunkID.Y = int(math.Ceil(tileY / float64(tileSize)))
	}
	if tileX == 0 {
		ChunkID.X = 1
	}
	if tileY == 0 {
		ChunkID.Y = 1
	}
	return ChunkID

}
