package gamemap

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Map struct {
	layerContainer *LayerContainer
	cfg            *config.Config
	// backgroundImage first layer
	backgroundImage *ebiten.Image
	// frontImages other layers
	frontImages []*ebiten.Image
}

type MapTile struct {
	PixelX    int
	PixelY    int
	Blocked   bool
	Invisible bool
	Image     *ebiten.Image
}

// todo rewrite this trash
func NewMap(cfg *config.Config) (*Map, error) {
	tileSet, _, err := ebitenutil.NewImageFromFile(cfg.Map.TileSetPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}

	waterTileSet, _, err := ebitenutil.NewImageFromFile("assets/tileset/[A]_type1/[A]Water2_pipo.png")
	if err != nil {
		return nil, fmt.Errorf("failed to init wall image: %v", err)
	}

	x0, y0 := 0, 128
	x1, y1 := (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	water := waterTileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 0, 993
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
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
	forestBottom1 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	x0, y0 = 32, 63
	x1, y1 = (x0+1)+cfg.Common.TileSize, (y0+1)+cfg.Common.TileSize
	forestBottom2 := tileSet.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)

	coord := base.Position{Y: 1, X: 1}
	chunk := utils.NewChunk(cfg.Map.Width, 40, 120, coord)

	gameMap := &Map{
		cfg:             cfg,
		backgroundImage: ebiten.NewImage(cfg.Map.Width*cfg.Common.TileSize, cfg.Map.Height*cfg.Common.TileSize),
	}

	layerContainer := NewLayerContainer(cfg.Map.Width, cfg.Map.Height)

	// first layer: grass and walls around the edges
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
	layerContainer.SetCurrent(tiles)
	tiles = layerContainer.Next()

	// build other layers
	for y := 0; y < cfg.Map.Height; y++ {
		for x := 0; x < cfg.Map.Width; x++ {
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

			if x == 0 || y == 0 || x == cfg.Map.Width-1 || y == cfg.Map.Height-1 { // don't touch the edge of the map
				tiles[x][y] = tile
				continue
			}

			switch tileType {
			case utils.Tree:
				if x < cfg.Map.Width-2 && y < cfg.Map.Height-2 { // don't touch the edge of the map
					currentIndex := layerContainer.GetIndex()
					layer := layerContainer.GetLayerWithoutCollisions([]base.Position{{x, y}, {x + 1, y}, {x, y + 1}, {x + 1, y + 1}})
					if layerContainer.GetIndex() < 8 { // [optimization] don't make trees higher than 8 layers
						makeTree(cfg, x, y, forestTop1, forestTop2, forestBottom1, forestBottom2, layer)
					}
					layerContainer.SetIndex(currentIndex)
				}
				continue
			case utils.Water:
				tile.Image = water
				tile.Blocked = true
			}
			tiles[x][y] = tile
		}
	}
	layerContainer.SetCurrent(tiles)

	gameMap.layerContainer = layerContainer

	for i := 0; i < len(layerContainer.Elements()); i++ {
		gameMap.frontImages = append(gameMap.frontImages, ebiten.NewImage(cfg.Map.Width*cfg.Common.TileSize, cfg.Map.Height*cfg.Common.TileSize))
	}

	return gameMap, nil
}

func makeTree(cfg *config.Config, x int, y int, forestTop1 *ebiten.Image, forestTop2 *ebiten.Image, forestBottom1 *ebiten.Image, forestBottom2 *ebiten.Image, tiles [][]MapTile) {
	tile1 := MapTile{
		PixelX:    x * cfg.Common.TileSize,
		PixelY:    y * cfg.Common.TileSize,
		Invisible: false,
		Blocked:   false,
		Image:     forestTop1,
	}
	tile2 := MapTile{
		PixelX:    (x + 1) * cfg.Common.TileSize,
		PixelY:    y * cfg.Common.TileSize,
		Invisible: false,
		Blocked:   false,
		Image:     forestTop2,
	}
	tile3 := MapTile{
		PixelX:    x * cfg.Common.TileSize,
		PixelY:    (y + 1) * cfg.Common.TileSize,
		Invisible: false,
		Blocked:   true,
		Image:     forestBottom1,
	}
	tile4 := MapTile{
		PixelX:    (x + 1) * cfg.Common.TileSize,
		PixelY:    (y + 1) * cfg.Common.TileSize,
		Invisible: false,
		Blocked:   true,
		Image:     forestBottom2,
	}
	tiles[x][y] = tile1
	tiles[x+1][y] = tile2
	tiles[x][y+1] = tile3
	tiles[x+1][y+1] = tile4
}

func (gameMap *Map) IsCanMove(x, y int) bool {
	for _, layerTiles := range gameMap.layerContainer.Elements() {
		if layerTiles[x][y].Blocked {
			return false
		}
	}
	return true
}

func (gameMap *Map) Update() {
	for i, layerTiles := range gameMap.layerContainer.Elements() {
		if i == 0 {
			gameMap.prepareImage(gameMap.backgroundImage, layerTiles)
			continue
		}
		// -1 layer background
		gameMap.prepareImage(gameMap.frontImages[i-1], layerTiles)
	}
}

func (gameMap *Map) prepareImage(backgroundImage *ebiten.Image, layerTiles Layer) {
	for x := 0; x < gameMap.cfg.Map.Width; x++ {
		for y := 0; y < gameMap.cfg.Map.Height; y++ {
			if !existTile(layerTiles, x, y) {
				continue
			}
			tile := layerTiles[x][y]

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			backgroundImage.DrawImage(tile.Image, op)
		}
	}
}

func (gameMap *Map) BackgroundImage() *ebiten.Image {
	return gameMap.backgroundImage
}

func (gameMap *Map) FrontImages() []*ebiten.Image {
	return gameMap.frontImages
}

func existTile(layer Layer, x, y int) bool {
	return layer[x][y].Image != nil
}
