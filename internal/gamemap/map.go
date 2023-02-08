package gamemap

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"
	"github.com/VxVxN/game/internal/data"
)

type Map struct {
	tiles []MapTile
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

func NewMap() (*Map, error) {
	gd := data.NewGameData()
	tiles := make([]MapTile, 0)

	wall, _, err := ebitenutil.NewImageFromFile("assets/wall.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}
	floor, _, err := ebitenutil.NewImageFromFile("assets/grass.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init grass image: %v", err)
	}
	for x := 0; x < gd.ScreenWidth; x++ {
		for y := 0; y < gd.ScreenHeight; y++ {
			if x == 0 || x == gd.ScreenWidth-1 || y == 0 || y == gd.ScreenHeight-1 {
				tile := MapTile{
					PixelX:  x * gd.TileSize,
					PixelY:  y * gd.TileSize,
					Blocked: true,
					Image:   wall,
				}
				tiles = append(tiles, tile)
			} else {
				tile := MapTile{
					PixelX:  x * gd.TileSize,
					PixelY:  y * gd.TileSize,
					Blocked: false,
					Image:   floor,
				}
				tiles = append(tiles, tile)
			}
		}
	}

	return &Map{tiles: tiles}, nil
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	tile := gameMap.tiles[GetIndexFromXY(x, y)]
	return !tile.Blocked
}

func (gameMap *Map) GetTile(x, y int) MapTile {
	return gameMap.tiles[GetIndexFromXY(x, y)]
}
