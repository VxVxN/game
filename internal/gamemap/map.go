package gamemap

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"time"
)

type Map struct {
	layers []Layer
	world  *ebiten.Image
	cfg    *config.Config
}

type MapTile struct {
	PixelX    int
	PixelY    int
	Blocked   bool
	Invisible bool
	Image     *ebiten.Image
}

type Layer [][]MapTile

// todo rewrite this trash
func NewMap(cfg *config.Config) (*Map, error) {
	tileSet, _, err := ebitenutil.NewImageFromFile(cfg.Map.TileSetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}

	x0, y0 := 0, 993
	x1, y1 := (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	wall := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 0, 0
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	floor := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 0, 32
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	forestTop1 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 32, 32
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	forestTop2 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 0, 63
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	forestBotton1 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 32, 63
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	forestBotton2 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	coord := base.Position{Y: 1, X: 1}
	chunk := utils.NewChunk(cfg.Map.Width, int(time.Now().UnixNano()), coord)

	gameMap := &Map{
		cfg:   cfg,
		world: ebiten.NewImage(cfg.Map.Width*cfg.Common.TileSize, cfg.Map.Height*cfg.Common.TileSize),
	}

	var layers []Layer

	// first layer
	tiles := make([][]MapTile, cfg.Map.Width)
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
	layers = append(layers, tiles)

	// second layer
	tiles = make([][]MapTile, cfg.Map.Width)
	for x := 0; x < cfg.Map.Width; x++ {
		tiles[x] = make([]MapTile, cfg.Map.Height)
	}
	for x := 0; x < cfg.Map.Width; x++ {
		for y := 0; y < cfg.Map.Height; y++ {
			tile := MapTile{
				PixelX:    x * cfg.Common.TileSize,
				PixelY:    y * cfg.Common.TileSize,
				Invisible: true,
			}

			if existTile(tiles, x, y) {
				continue
			}

			tileType, ok := chunk.Map[x][y]
			if !ok {
				tiles[x][y] = tile
				continue
			}

			switch tileType {
			case utils.Forest:
				if x != 0 && y != 0 && x < cfg.Map.Width-2 && y < cfg.Map.Height-2 {
					tile1 := MapTile{
						PixelX:    x * cfg.Common.TileSize,
						PixelY:    y * cfg.Common.TileSize,
						Invisible: false,
						Blocked:   true,
						Image:     forestTop1,
					}
					tile2 := MapTile{
						PixelX:    (x + 1) * cfg.Common.TileSize,
						PixelY:    y * cfg.Common.TileSize,
						Invisible: false,
						Blocked:   true,
						Image:     forestTop2,
					}
					tile3 := MapTile{
						PixelX:    x * cfg.Common.TileSize,
						PixelY:    (y + 1) * cfg.Common.TileSize,
						Invisible: false,
						Blocked:   true,
						Image:     forestBotton1,
					}
					tile4 := MapTile{
						PixelX:    (x + 1) * cfg.Common.TileSize,
						PixelY:    (y + 1) * cfg.Common.TileSize,
						Invisible: false,
						Blocked:   true,
						Image:     forestBotton2,
					}
					tiles[x][y] = tile1
					tiles[x+1][y] = tile2
					tiles[x][y+1] = tile3
					tiles[x+1][y+1] = tile4
					continue
				}
			}
			tiles[x][y] = tile
		}
	}
	layers = append(layers, tiles)

	gameMap.layers = layers

	return gameMap, nil
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	for _, layerTiles := range gameMap.layers {
		if layerTiles[x][y].Blocked {
			return false
		}
	}
	return true
}

func (gameMap *Map) Update() {
	for _, layerTiles := range gameMap.layers {
		for x := 0; x < gameMap.cfg.Map.Width; x++ {
			for y := 0; y < gameMap.cfg.Map.Height; y++ {
				tile := layerTiles[x][y]
				if tile.Invisible {
					continue
				}

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
				gameMap.world.DrawImage(tile.Image, op)
			}
		}
	}
}

func (gameMap *Map) Image() *ebiten.Image {
	return gameMap.world
}

func existTile(layer Layer, x, y int) bool {
	return layer[x][y].Image != nil
}
