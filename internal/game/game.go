package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/player"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"os"
	"time"
)

type Game struct {
	gameMap         *gamemap.Map
	cfg             *config.Config
	player          *player.Player
	eventManager    *eventmanager.EventManager
	globalTime      time.Time
	isShowDebugInfo bool
}

func NewGame(cfg *config.Config) (*Game, error) {
	gameMap, err := gamemap.NewMap(cfg)
	if err != nil {
		return nil, err
	}
	gameMap.Update()

	playerPosition := base.NewPosition(utils.RandomIntByRange(1, cfg.Map.Width-1), utils.RandomIntByRange(1, cfg.Map.Height-1))
	player, err := player.NewPlayer(playerPosition, cfg.Player.ImagePath, cfg.Common.TileSize, cfg.Player.FrameCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	game := &Game{
		gameMap:         gameMap,
		cfg:             cfg,
		player:          player,
		globalTime:      time.Now(),
		isShowDebugInfo: true,
		eventManager:    eventmanager.NewEventManager(),
	}

	game.addEvents(gameMap, player)

	return game, nil
}

func (game *Game) addEvents(gameMap *gamemap.Map, player *player.Player) {
	game.eventManager.AddEvent(ebiten.KeyUp, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y-1) {
			return
		}
		player.Move(ebiten.KeyUp)
	})
	game.eventManager.AddEvent(ebiten.KeyDown, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y+1) {
			return
		}
		player.Move(ebiten.KeyDown)
	})
	game.eventManager.AddEvent(ebiten.KeyRight, func() {
		if !gameMap.IsCanMove(player.Position.X+1, player.Position.Y) {
			return
		}
		player.Move(ebiten.KeyRight)
	})
	game.eventManager.AddEvent(ebiten.KeyLeft, func() {
		if !gameMap.IsCanMove(player.Position.X-1, player.Position.Y) {
			return
		}
		player.Move(ebiten.KeyLeft)
	})
	game.eventManager.AddEvent(ebiten.KeyTab, func() {
		game.isShowDebugInfo = !game.isShowDebugInfo
	})
	game.eventManager.AddEvent(ebiten.KeyEscape, func() {
		os.Exit(0) // todo add normal game end processing
	})
	game.eventManager.SetDefaultEvent(func() {
		player.Stand()
	})
}

func (game *Game) Update() error {
	if time.Since(game.globalTime) < time.Second/time.Duration(game.cfg.Common.RefreshRateFramesPerSecond) {
		return nil
	}
	game.eventManager.Update()
	game.globalTime = time.Now()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.cfg.Common.TileSize*-game.player.Position.X+game.cfg.Common.WindowWidth/2),
		float64(game.cfg.Common.TileSize*-game.player.Position.Y+game.cfg.Common.WindowHeight/2))
	screen.DrawImage(game.gameMap.BackgroundImage(), op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(game.cfg.Common.WindowWidth/2), float64(game.cfg.Common.WindowHeight/2))
	screen.DrawImage(game.player.Image(), op)

	for _, frontImage := range game.gameMap.FrontImages() {
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(game.cfg.Common.TileSize*-game.player.Position.X+game.cfg.Common.WindowWidth/2),
			float64(game.cfg.Common.TileSize*-game.player.Position.Y+game.cfg.Common.WindowHeight/2))
		screen.DrawImage(frontImage, op)
	}

	if game.isShowDebugInfo {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("X = %d, Y = %d\nLayers: %d", game.player.Position.X, game.player.Position.Y,
			len(game.gameMap.FrontImages())+2)) // +2 -> gameMap.BackgroundImage() + player.Image()
	}
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
