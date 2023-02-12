package gamemap

import (
	"fmt"
	"github.com/VxVxN/game/internal/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Map struct {
	tiles [][]MapTile
	world *ebiten.Image
	cfg   *config.Config
}

type MapTile struct {
	PixelX  int
	PixelY  int
	Blocked bool
	Image   *ebiten.Image
}

func NewMap(cfg *config.Config) (*Map, error) {
	tiles := make([][]MapTile, cfg.Map.Width)

	tileSet, _, err := ebitenutil.NewImageFromFile(cfg.Map.TileSetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}

	x0, y0 := 0, 993
	x1, y1 := (x0+1)*cfg.Common.TileSize, (y0+1)*cfg.Common.TileSize
	wall := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 0, 0
	x1, y1 = (x0+1)*cfg.Common.TileSize, (y0+1)*cfg.Common.TileSize
	floor := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	for x := 0; x < cfg.Map.Width; x++ {
		tiles[x] = make([]MapTile, cfg.Map.Height)
		for y := 0; y < cfg.Map.Height; y++ {
			tile := MapTile{
				PixelX: x * cfg.Common.TileSize,
				PixelY: y * cfg.Common.TileSize,
			}
			if x == 0 || x == cfg.Map.Width-1 || y == 0 || y == cfg.Map.Height-1 {
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
		tiles: tiles,
		cfg:   cfg,
		world: ebiten.NewImage(cfg.Map.Width*cfg.Common.TileSize, cfg.Map.Height*cfg.Common.TileSize),
	}, nil
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	tile := gameMap.tiles[x][y]
	return !tile.Blocked
}

func (gameMap *Map) Update() {
	for x := 0; x < gameMap.cfg.Map.Width; x++ {
		for y := 0; y < gameMap.cfg.Map.Height; y++ {
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
