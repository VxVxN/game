package gamemap

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"
	"github.com/VxVxN/game/internal/data"
)

type Map struct {
	tiles    []MapTile
	gameData *data.GameData
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func GetIndexFromXY(x int, y int) int {
	gd := data.NewGameData()
	return (y * gd.ScreenWidth) + x
}

func NewMap(gameData *data.GameData) (*Map, error) {
	tiles := make([]MapTile, 0)

	wall, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}
	floor, _, err := ebitenutil.NewImageFromFile("assets/grass.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init grass image: %v", err)
	}
	for x := 0; x < gameData.ScreenWidth; x++ {
		for y := 0; y < gameData.ScreenHeight; y++ {
			if x == 0 || x == gameData.ScreenWidth-1 || y == 0 || y == gameData.ScreenHeight-1 {
				tile := MapTile{
					PixelX:  x * gameData.TileSize,
					PixelY:  y * gameData.TileSize,
					Blocked: true,
					Image:   wall,
				}
				tiles = append(tiles, tile)
			} else {
				tile := MapTile{
					PixelX:  x * gameData.TileSize,
					PixelY:  y * gameData.TileSize,
					Blocked: false,
					Image:   floor,
				}
				tiles = append(tiles, tile)
			}
		}
	}

	return &Map{
		tiles:    tiles,
		gameData: gameData,
	}, nil
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	tile := gameMap.tiles[GetIndexFromXY(x, y)]
	return !tile.Blocked
}

func (gameMap *Map) getTile(x, y int) MapTile {
	return gameMap.tiles[GetIndexFromXY(x, y)]
}

func (gameMap *Map) Draw(screen *ebiten.Image) {
	for x := 0; x < gameMap.gameData.ScreenWidth; x++ {
		for y := 0; y < gameMap.gameData.ScreenHeight; y++ {
			tile := gameMap.getTile(x, y)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}
}
