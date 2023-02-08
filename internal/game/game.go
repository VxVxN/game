package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/eventmanager"
	"github.com/VxVxN/game/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
	"os"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/internal/data"
)

type Game struct {
	gameMap      *gamemap.Map
	data         data.GameData
	player       *player.Player
	eventManager *eventmanager.EventManager
	globalTime   time.Time
}

func NewGame() (*Game, error) {
	gameMap, err := gamemap.NewMap()
	if err != nil {
		return nil, err
	}
	player, err := player.NewPlayer(base.NewPosition(5, 2), "assets/player.png")
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	eventManager := eventmanager.NewEventManager()
	eventManager.AddEvent(ebiten.KeyUp, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y-1) {
			return
		}
		player.Move(ebiten.KeyUp)
	})
	eventManager.AddEvent(ebiten.KeyDown, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y+1) {
			return
		}
		player.Move(ebiten.KeyDown)
	})
	eventManager.AddEvent(ebiten.KeyRight, func() {
		if !gameMap.IsCanMove(player.Position.X+1, player.Position.Y) {
			return
		}
		player.Move(ebiten.KeyRight)
	})
	eventManager.AddEvent(ebiten.KeyLeft, func() {
		if !gameMap.IsCanMove(player.Position.X-1, player.Position.Y) {
			return
		}
		player.Move(ebiten.KeyLeft)
	})
	eventManager.AddEvent(ebiten.KeyEscape, func() {
		os.Exit(0) // todo add normal game end processing
	})

	game := &Game{
		gameMap:      gameMap,
		data:         data.NewGameData(),
		player:       player,
		eventManager: eventManager,
		globalTime:   time.Now(),
	}
	return game, nil
}

func (game *Game) Update() error {
	if time.Since(game.globalTime) < time.Second/25 {
		return nil
	}
	game.eventManager.Process()
	game.globalTime = time.Now()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	//Draw the Map
	for x := 0; x < game.data.ScreenWidth; x++ {
		for y := 0; y < game.data.ScreenHeight; y++ {
			tile := game.gameMap.GetTile(x, y)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tile.PixelX), float64(tile.PixelY))
			screen.DrawImage(tile.Image, op)
		}
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.data.TileSize*game.player.Position.X), float64(game.data.TileSize*game.player.Position.Y))
	screen.DrawImage(game.player.Image(), op)
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
