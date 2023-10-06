package game

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/camera"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/entity"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/item"
	"github.com/VxVxN/game/pkg/scriptmanager"
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
	gameMap         *gamemap.Map
	cfg             *config.Config
	player          *entity.Player
	npc             *entity.NPC
	eventManager    *eventmanager.EventManager
	camera          *camera.Camera
	globalTime      time.Time
	gameOverFace    font.Face
	items           []*item.Item
	isShowDebugInfo bool
	isStopWorld     bool
}

func NewGame(cfg *config.Config) (*Game, error) {
	gameMap, err := gamemap.NewMap(cfg)
	if err != nil {
		return nil, err
	}
	gameMap.Update()

	playerPosition := findPosition(cfg, gameMap)
	player, err := entity.NewPlayer(playerPosition, 0.8, 0, 0, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	npcPosition := findPosition(cfg, gameMap)
	npc, err := entity.NewNPC("Bob", npcPosition, 0.2, cfg.Player.ImagePath, 96, 128, cfg.Player.FrameCount, gameMap, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}
	playerPosition.X++
	axeItem, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 160, 4192, cfg.Common.TileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create axe item: %v", err)
	}
	playerPosition.X++
	keyItem, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 224, 4192, cfg.Common.TileSize)
	if err != nil {
		return nil, fmt.Errorf("failed to create axe item: %v", err)
	}
	//npc.SetScripts([]*scriptmanager.Script{
	//	scriptmanager.NewScript([]scriptmanager.Action{scriptmanager.MoveRight, scriptmanager.MoveRight, scriptmanager.MoveRight, scriptmanager.Pause, scriptmanager.MoveUp, scriptmanager.MoveUp, scriptmanager.MoveUp, scriptmanager.Pause, scriptmanager.MoveLeft, scriptmanager.MoveLeft, scriptmanager.MoveLeft, scriptmanager.Pause, scriptmanager.MoveDown, scriptmanager.MoveDown, scriptmanager.MoveDown, scriptmanager.Pause}),
	//	scriptmanager.NewScript([]scriptmanager.Action{scriptmanager.MoveRight, scriptmanager.MoveRight, scriptmanager.MoveRight, scriptmanager.Pause, scriptmanager.MoveLeft, scriptmanager.MoveLeft, scriptmanager.MoveLeft, scriptmanager.Pause}),
	//	scriptmanager.NewScript([]scriptmanager.Action{scriptmanager.MoveUp, scriptmanager.MoveUp, scriptmanager.MoveUp, scriptmanager.Pause, scriptmanager.MoveDown, scriptmanager.MoveDown, scriptmanager.MoveDown, scriptmanager.Pause}),
	//})
	byeReplica := &scriptmanager.PieceDialogue{
		Replicas: []string{"Bye stranger"},
	}
	npc.AddDialogue(&scriptmanager.PieceDialogue{
		Replicas: []string{"Hello stranger", "Do you want a coin?"},
		Answers: []scriptmanager.Answer{
			{
				Text: "Yes, of course",
				Action: func() {
					player.AddCoins(1)
				},
				NextPieceDialogue: byeReplica,
			},
			{
				Text: "No",
				NextPieceDialogue: &scriptmanager.PieceDialogue{
					Replicas: []string{"Do you want two coins?"},
					Answers: []scriptmanager.Answer{
						{
							Text: "Yes",
							Action: func() {
								player.AddCoins(2)
							},
							NextPieceDialogue: byeReplica,
						},
						{
							Text:              "No",
							NextPieceDialogue: byeReplica,
						},
					},
				},
			},
		},
	})

	font, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF font: %v", err)
	}

	gameOverFace, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 64,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new face for game ofer face(font): %v", err)
	}

	game := &Game{
		gameMap:         gameMap,
		cfg:             cfg,
		player:          player,
		npc:             npc,
		globalTime:      time.Now(),
		isShowDebugInfo: true,
		eventManager:    eventmanager.NewEventManager(),
		camera:          camera.NewCamera(cfg),
		gameOverFace:    gameOverFace,
		items:           []*item.Item{axeItem, keyItem},
	}

	game.camera.AddPlayerImage(player.Image())
	game.camera.AddEntityImage(npc.Image())
	game.camera.AddBackgroundImage(gameMap.BackgroundImage())
	game.camera.AddFrontImages(gameMap.FrontImages())
	game.camera.SetItems([]*item.Item{axeItem, keyItem})

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

func findPosition(cfg *config.Config, gameMap *gamemap.Map) base.Position {
	position := base.NewPosition(utils.RandomFloat64ByRange(1, float64(cfg.Map.Width-1)), utils.RandomFloat64ByRange(1, float64(cfg.Map.Height-1)))
	for {
		if gameMap.IsCanMove(position.X, position.Y) {
			break
		}
		position = base.NewPosition(utils.RandomFloat64ByRange(1, float64(cfg.Map.Width-1)), utils.RandomFloat64ByRange(1, float64(cfg.Map.Height-1)))
		continue
	}
	return position
}

func (game *Game) addEvents(gameMap *gamemap.Map, player *entity.Player) {
	game.eventManager.AddEvent(ebiten.KeyUp, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y-1) || game.isStopWorld {
			return
		}
		player.Move(ebiten.KeyUp)
	})
	game.eventManager.AddEvent(ebiten.KeyDown, func() {
		if !gameMap.IsCanMove(player.Position.X, player.Position.Y+1) || game.isStopWorld {
			return
		}
		player.Move(ebiten.KeyDown)
	})
	game.eventManager.AddEvent(ebiten.KeyRight, func() {
		if !gameMap.IsCanMove(player.Position.X+1, player.Position.Y) || game.isStopWorld {
			return
		}
		player.Move(ebiten.KeyRight)
	})
	game.eventManager.AddEvent(ebiten.KeyLeft, func() {
		if !gameMap.IsCanMove(player.Position.X-1, player.Position.Y) || game.isStopWorld {
			return
		}
		player.Move(ebiten.KeyLeft)
	})
	game.eventManager.AddEvent(ebiten.KeySpace, func() {
		if utils.CanAction(game.player.Position, game.npc.Position) {
			game.isStopWorld = !game.npc.IsEndDialogue()
			game.npc.Trigger()
		}
		for _, item := range game.items {
			if utils.CanAction(game.player.Position, item.Position()) {
				item.Trigger()
				game.player.TakeItem(item)
			}
		}
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
	game.globalTime = time.Now()

	game.npc.Update(game.player.Position)
	game.eventManager.Update()
	if game.isStopWorld {
		return nil
	}

	game.camera.AddPlayerImage(game.player.Image())
	game.camera.UpdatePlayer(game.player.Position)
	game.camera.UpdateEntity(game.npc.Position)

	game.player.DecreaseSatiety()
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	game.camera.Draw(screen)

	if game.isShowDebugInfo {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("X = %f, Y = %f\nLayers: %d", game.player.Position.X, game.player.Position.Y,
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
	game.player.Draw(screen)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
