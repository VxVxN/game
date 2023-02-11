package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/data"
	"github.com/VxVxN/game/internal/eventmanager"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/internal/player"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

type Game struct {
	gameMap      *gamemap.Map
	data         data.GameData
	player       *player.Player
	eventManager *eventmanager.EventManager
	globalTime   time.Time
}

func NewGame() (*Game, error) {
	gameData := data.NewGameData()

	gameMap, err := gamemap.NewMap(&gameData)
	if err != nil {
		return nil, err
	}
	player, err := player.NewPlayer(base.NewPosition(1, 1), "assets/characters.png", gameData.TileSize, 3)
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
	eventManager.AddDefaultEvent(func() {
		player.Stand()
	})

	game := &Game{
		gameMap:      gameMap,
		data:         gameData,
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
	game.gameMap.Update()
	game.eventManager.Update()
	game.globalTime = time.Now()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.data.TileSize*-game.player.Position.X+game.data.ScreenWidthPx/2), float64(game.data.TileSize*-game.player.Position.Y+game.data.ScreenHeightPx/2))
	screen.DrawImage(game.gameMap.Image(), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.data.ScreenWidthPx/2), float64(game.data.ScreenHeightPx/2))
	screen.DrawImage(game.player.Image(), op)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	game.data.ScreenWidthPx = screenWidthPx
	game.data.ScreenHeightPx = screenHeightPx
	return screenWidthPx, screenHeightPx
}
