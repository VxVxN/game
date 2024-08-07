package utils

import (
	"math"

	"github.com/VxVxN/game/internal/base"
)

type TileType int

const (
	Tree TileType = iota
	Water
	Grass
)

type Chunk struct {
	ChunkID [2]int
	Map     map[int]map[int]TileType
}

func NewChunk(chunkSize, treePerlinSeed, waterPerlinSeed int, idChunk base.Position) Chunk {
	chunk := Chunk{ChunkID: [2]int{int(idChunk.X), int(idChunk.Y)}}

	var chunkXMax, chunkYMax int
	chunkMap := make(map[int]map[int]TileType)
	chunkXMax = int(idChunk.X * float64(chunkSize))
	chunkYMax = int(idChunk.Y * float64(chunkSize))

	for x := 0; x < chunkXMax; x++ {
		chunkMap[x] = make(map[int]TileType)
		for y := 0; y < chunkYMax; y++ {
			chunkMap[x][y] = Grass
		}
	}

	for x := 0; x < chunkXMax; x++ {
		for y := 0; y < chunkYMax; y++ {
			setTile(x, chunkMap, y, waterPerlinSeed, Water)
		}
	}
	fillWaterGaps(chunkMap, chunkXMax)
	fillWaterGaps(chunkMap, chunkXMax)

	for x := 0; x < chunkXMax; x++ {
		for y := 0; y < chunkYMax; y++ {
			setTile(x, chunkMap, y, treePerlinSeed, Tree)
		}
	}

	chunk.Map = chunkMap
	return chunk
}

func fillWaterGaps(chunkMap map[int]map[int]TileType, chunkMax int) {
	for x := 0; x < chunkMax; x++ {
		var firstWaterIndex int
		for y := 0; y < chunkMax; y++ {
			if chunkMap[x][y] != Water {
				continue
			}
			if firstWaterIndex == 0 {
				firstWaterIndex = y
				continue
			}
			if math.Abs(float64(firstWaterIndex-y)) > 20 {
				firstWaterIndex = y
				continue
			}
			for i := firstWaterIndex; i < y; i++ {
				chunkMap[x][i] = Water
			}
			firstWaterIndex = y
		}
	}
	for y := 0; y < chunkMax; y++ {
		var firstWaterIndex int
		for x := 0; x < chunkMax; x++ {
			if chunkMap[x][y] != Water {
				continue
			}
			if firstWaterIndex == 0 {
				firstWaterIndex = x
				continue
			}
			if math.Abs(float64(firstWaterIndex-x)) > 20 {
				firstWaterIndex = x
				continue
			}
			for i := firstWaterIndex; i < x; i++ {
				chunkMap[i][y] = Water
			}
			firstWaterIndex = x
		}
	}
}

func setTile(x int, chunkMap map[int]map[int]TileType, y int, perlinSeed int, tileType TileType) {
	if chunkMap[x][y] != Grass {
		return
	}
	perlinValue := Noise(float32(x)/float32(perlinSeed), float32(y)/float32(perlinSeed))
	switch {
	case perlinValue < -0.3:
		chunkMap[x][y] = tileType
	}
}

func GetChunkID(tileSize, x, y int) base.Position {
	tileX := float64(x) / float64(tileSize)
	tileY := float64(y) / float64(tileSize)

	var ChunkID base.Position
	if tileX < 0 {
		ChunkID.X = math.Floor(tileX / float64(tileSize))
	} else {
		ChunkID.X = math.Ceil(tileX / float64(tileSize))
	}
	if tileY < 0 {
		ChunkID.Y = math.Floor(tileY / float64(tileSize))
	} else {
		ChunkID.Y = math.Ceil(tileY / float64(tileSize))
	}
	if tileX == 0 {
		ChunkID.X = 1
	}
	if tileY == 0 {
		ChunkID.Y = 1
	}
	return ChunkID

}
