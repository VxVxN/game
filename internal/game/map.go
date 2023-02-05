package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"fmt"
)

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func GetIndexFromXY(x int, y int) int {
	gd := NewGameData()
	return (y * gd.ScreenWidth) + x
}

func NewGameTiles() ([]MapTile, error) {
	gd := NewGameData()
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
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: true,
					Image:   wall,
				}
				tiles = append(tiles, tile)
			} else {
				tile := MapTile{
					PixelX:  x * gd.TileWidth,
					PixelY:  y * gd.TileHeight,
					Blocked: false,
					Image:   floor,
				}
				tiles = append(tiles, tile)
			}
		}
	}

	return tiles, nil
}
