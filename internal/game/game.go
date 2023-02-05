package game

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	tiles []MapTile
	data  GameData
}

func NewGame() (*Game, error) {
	tiles, err := NewGameTiles()
	if err != nil {
		return nil, err
	}
	game := &Game{
		tiles: tiles,
		data:  NewGameData(),
	}
	return game, nil
}

func (game *Game) Update() error {
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	//Draw the Map
	for x := 0; x < game.data.ScreenWidth; x++ {
		for y := 0; y < game.data.ScreenHeight; y++ {
			tile := game.tiles[GetIndexFromXY(x, y)]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 800
}

type GameData struct {
	ScreenWidth  int
	ScreenHeight int
	TileWidth    int
	TileHeight   int
}

func NewGameData() GameData {
	g := GameData{
		ScreenWidth:  30,
		ScreenHeight: 30,
		TileWidth:    32,
		TileHeight:   32,
	}
	return g
}
