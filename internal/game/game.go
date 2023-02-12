package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/player"
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
	"github.com/VxVxN/game/pkg/utils"
)

type Game struct {
	gameMap      *gamemap.Map
	cfg          *config.Config
	player       *player.Player
	eventManager *eventmanager.EventManager
	globalTime   time.Time
}

func NewGame(cfg *config.Config) (*Game, error) {
	gameMap, err := gamemap.NewMap(cfg)
	if err != nil {
		return nil, err
	}

	playerPosition := base.NewPosition(utils.RandomIntByRange(1, cfg.Map.Width-1), utils.RandomIntByRange(1, cfg.Map.Height-1))
	player, err := player.NewPlayer(playerPosition, cfg.Player.ImagePath, cfg.Common.TileSize, cfg.Player.FrameCount)
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
	eventManager.SetDefaultEvent(func() {
		player.Stand()
	})

	game := &Game{
		gameMap:      gameMap,
		cfg:          cfg,
		player:       player,
		eventManager: eventManager,
		globalTime:   time.Now(),
	}
	return game, nil
}

func (game *Game) Update() error {
	if time.Since(game.globalTime) < time.Second/time.Duration(game.cfg.Common.RefreshRateFramesPerSecond) {
		return nil
	}
	game.gameMap.Update()
	game.eventManager.Update()
	game.globalTime = time.Now()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.cfg.Common.TileSize*-game.player.Position.X+game.cfg.Common.WindowWidth/2),
		float64(game.cfg.Common.TileSize*-game.player.Position.Y+game.cfg.Common.WindowHeight/2))
	screen.DrawImage(game.gameMap.Image(), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.cfg.Common.WindowWidth/2), float64(game.cfg.Common.WindowHeight/2))
	screen.DrawImage(game.player.Image(), op)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
