package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/camera"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/entity"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image/color"
	"os"
	"time"
)

type Game struct {
	gameMap          *gamemap.Map
	cfg              *config.Config
	player           *entity.Player
	npc              *entity.NPC
	eventManager     *eventmanager.EventManager
	camera           *camera.Camera
	globalTime       time.Time
	statusPlayerFace font.Face
	gameOverFace     font.Face
	isShowDebugInfo  bool
}

func NewGame(cfg *config.Config) (*Game, error) {
	gameMap, err := gamemap.NewMap(cfg)
	if err != nil {
		return nil, err
	}
	gameMap.Update()

	// looking for position for player
	playerPosition := base.NewPosition(utils.RandomIntByRange(1, cfg.Map.Width-1), utils.RandomIntByRange(1, cfg.Map.Height-1))
	for {
		if gameMap.IsCanMove(playerPosition.X, playerPosition.Y) {
			break
		}
		playerPosition = base.NewPosition(utils.RandomIntByRange(1, cfg.Map.Width-1), utils.RandomIntByRange(1, cfg.Map.Height-1))
		continue
	}

	player, err := entity.NewPlayer(playerPosition, cfg.Player.ImagePath, 0, 0, cfg.Common.TileSize, cfg.Player.FrameCount)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	npcPosition := base.NewPosition(playerPosition.X+2, playerPosition.Y)
	npc, err := entity.NewNPC("Bob", npcPosition, cfg.Player.ImagePath, 96, 128, cfg.Player.FrameCount, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	font, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF font: %v", err)
	}

	statusPlayerFace, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new status player face(font): %v", err)
	}
	gameOverFace, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 64,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new face for game ofer face(font): %v", err)
	}

	game := &Game{
		gameMap:          gameMap,
		cfg:              cfg,
		player:           player,
		npc:              npc,
		globalTime:       time.Now(),
		isShowDebugInfo:  true,
		eventManager:     eventmanager.NewEventManager(),
		camera:           camera.NewCamera(cfg),
		statusPlayerFace: statusPlayerFace,
		gameOverFace:     gameOverFace,
	}

	game.camera.AddPlayerImage(player.Image())
	game.camera.AddEntityImage(npc.Image())
	game.camera.AddBackgroundImage(gameMap.BackgroundImage())
	game.camera.AddFrontImages(gameMap.FrontImages())

	game.addEvents(gameMap, player)

	switch game.cfg.Common.Mode {
	case config.ViewMode:
		//zoom := 0.3
		//zoom := 0.15
		zoom := 1.0

		game.camera.SetZoom(zoom)
	default:
	}

	return game, nil
}

func (game *Game) addEvents(gameMap *gamemap.Map, player *entity.Player) {
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

	game.npc.Update(game.player.Position)

	game.camera.AddPlayerImage(game.player.Image())
	game.camera.UpdatePlayer(game.player.Position)
	game.camera.UpdateEntity(game.npc.Position)

	game.camera.UpdateEntity(game.npc.Position)
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.camera.Draw(screen)

	if game.isShowDebugInfo {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("X = %d, Y = %d\nLayers: %d", game.player.Position.X, game.player.Position.Y,
			len(game.gameMap.FrontImages())+2)) // +2 -> gameMap.BackgroundImage() + player.Image()
	}

	if game.player.IsDead() {
		text.Draw(screen, "Game over", game.gameOverFace, game.cfg.Common.WindowWidth/2-150, game.cfg.Common.WindowHeight/2, color.NRGBA{
			R: 255,
			A: 255,
		})
		return
	}

	game.npc.Draw(screen)

	text.Draw(screen, fmt.Sprintf("XP: %d%%, Satiety: %d%%", game.player.XP(), game.player.Satiety()), game.statusPlayerFace, game.cfg.Common.WindowWidth/2-50, 80, color.Black)
	game.player.DecreaseSatiety()
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
