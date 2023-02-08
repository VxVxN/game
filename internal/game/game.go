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
	// drawers objects that can be drawn
	drawers []Drawer
}

type Drawer interface {
	Draw(screen *ebiten.Image)
}

func NewGame() (*Game, error) {
	gameData := data.NewGameData()

	gameMap, err := gamemap.NewMap(&gameData)
	if err != nil {
		return nil, err
	}
	player, err := player.NewPlayer(base.NewPosition(5, 2), "assets/player.png", &gameData)
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

	drawers := []Drawer{
		gameMap,
		player,
	}

	game := &Game{
		gameMap:      gameMap,
		data:         gameData,
		player:       player,
		eventManager: eventManager,
		globalTime:   time.Now(),
		drawers:      drawers,
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
	for _, drawer := range game.drawers {
		drawer.Draw(screen)
	}
}

func (game *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
