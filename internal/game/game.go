package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/eventmanager"
	"github.com/VxVxN/game/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Game struct {
	tiles        []MapTile
	data         GameData
	player       *player.Player
	eventManager *eventmanager.EventManager
	globalTime   time.Time
}

func NewGame() (*Game, error) {
	tiles, err := NewGameTiles()
	if err != nil {
		return nil, err
	}
	player, err := player.NewPlayer(base.NewPosition(5, 2), "assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	eventManager := eventmanager.NewEventManager(player)

	game := &Game{
		tiles:        tiles,
		data:         NewGameData(),
		player:       player,
		eventManager: eventManager,
		globalTime:   time.Now(),
	}
	return game, nil
}

func (game *Game) Update() error {
	if time.Since(game.globalTime) >= time.Second/25 {
		game.eventManager.Process()
		game.globalTime = time.Now()
	}
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.data.TileSize*game.player.Position.X), float64(game.data.TileSize*game.player.Position.Y))
	screen.DrawImage(game.player.Image(), op)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 800
}

type GameData struct {
	ScreenWidth  int
	ScreenHeight int
	TileSize     int
}

func NewGameData() GameData {
	g := GameData{
		ScreenWidth:  30,
		ScreenHeight: 30,
		TileSize:     32,
	}
	return g
}
