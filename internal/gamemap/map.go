package gamemap

import (
	"fmt"
	"github.com/VxVxN/game/internal/data"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Map struct {
	tiles    [][]MapTile
	world    *ebiten.Image
	gameData *data.GameData
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func NewMap(gameData *data.GameData) (*Map, error) {
	tiles := make([][]MapTile, gameData.MapWidth)

	wall, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}
	floor, _, err := ebitenutil.NewImageFromFile("assets/grass.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init grass image: %v", err)
	}
	for x := 0; x < gameData.MapWidth; x++ {
		tiles[x] = make([]MapTile, gameData.MapHeight)
		for y := 0; y < gameData.MapHeight; y++ {
			tile := MapTile{
				PixelX: x * gameData.TileSize,
				PixelY: y * gameData.TileSize,
			}
			if x == 0 || x == gameData.MapWidth-1 || y == 0 || y == gameData.MapHeight-1 {
				tile.Blocked = true
				tile.Image = wall
			} else {
				tile.Blocked = false
				tile.Image = floor
			}
			tiles[x][y] = tile
		}
	}

	return &Map{
		tiles:    tiles,
		gameData: gameData,
		world:    ebiten.NewImage(gameData.MapWidth*gameData.TileSize, gameData.MapHeight*gameData.TileSize),
	}, nil
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	tile := gameMap.tiles[x][y]
	return !tile.Blocked
}

func (gameMap *Map) Update() {
	for x := 0; x < gameMap.gameData.MapWidth; x++ {
		for y := 0; y < gameMap.gameData.MapHeight; y++ {
			tile := gameMap.tiles[x][y]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			gameMap.world.DrawImage(tile.Image, op)
		}
	}
}

func (gameMap *Map) Image() *ebiten.Image {
	return gameMap.world
}
