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
	gameMap      *gamemap.Map
	cfg          *config.Config
	player       *entity.Player
	npc          *entity.NPC
	eventManager *eventmanager.EventManager
	camera       *camera.Camera
	globalTime   time.Time
	gameOverFace font.Face
	items        []*item.Item
	//inventory       *inventory.Inventory
	stage           Stage
	isShowDebugInfo bool
}

type Stage int

const (
	GameStage Stage = iota
	DialogutStage
)

func NewGame(cfg *config.Config) (*Game, error) {
	gameMap, err := gamemap.NewMap(cfg)
	if err != nil {
		return nil, err
	}
	gameMap.Update()

	playerPosition := findPosition(cfg, gameMap)
	player, err := entity.NewPlayer(playerPosition, 0.2, 0, 0, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}

	//npcPosition := findPosition(cfg, gameMap)
	npc, err := entity.NewNPC("Bob", playerPosition, 0.1, cfg.Player.ImagePath, 96, 128, cfg.Player.FrameCount, gameMap, cfg)
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
		//inventory:       inventory.NewInventory(cfg),
		stage: GameStage,
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
	game.eventManager.AddPressEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position.X, player.Position.Y-1) {
				return
			}
			player.Move(ebiten.KeyUp)
		case DialogutStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case GameStage:
		case DialogutStage:
			game.npc.DialogueManager.NextAnswer()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position.X, player.Position.Y+1) {
				return
			}
			player.Move(ebiten.KeyDown)
		case DialogutStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case GameStage:
		case DialogutStage:
			game.npc.DialogueManager.BeforeAnswer()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position.X+1, player.Position.Y) {
				return
			}
			player.Move(ebiten.KeyRight)
		case DialogutStage:
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position.X-1, player.Position.Y) {
				return
			}
			player.Move(ebiten.KeyLeft)
		case DialogutStage:
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeySpace, func() {
		switch game.stage {
		case GameStage:
			if utils.CanAction(game.player.Position, game.npc.Position) && !game.npc.DialogueManager.IsEndDialogue() {
				game.npc.Trigger()
				game.stage = DialogutStage
			}
			for _, item := range game.items {
				if utils.CanAction(game.player.Position, item.Position()) {
					item.Trigger()
					game.player.TakeItem(item)
				}
			}
		case DialogutStage:
			if game.npc.DialogueManager.NeedAnswer() {
				game.npc.DialogueManager.DoAnswer()
				game.npc.DialogueManager.PieceDialogue = game.npc.DialogueManager.NextPieceDialogue()
				if game.npc.DialogueManager.IsEndDialogue() {
					game.stage = GameStage
					return
				}
				return
			}
			game.npc.DialogueManager.NextReplica()
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyTab, func() {
		game.isShowDebugInfo = !game.isShowDebugInfo
	})
	//game.eventManager.AddEvent(ebiten.KeyI, func() {
	//	game.inventory.OnOff()
	//})
	game.eventManager.AddPressedEvent(ebiten.KeyEscape, func() {
		os.Exit(0) // todo add normal game end processing
	})
	game.eventManager.SetDefaultEvent(func() {
		player.Stand()
	})
}

func (game *Game) Update() error {
	game.eventManager.Update()
	game.npc.Update(game.player.Position)
	if time.Since(game.globalTime) < time.Second/time.Duration(game.cfg.Common.RefreshRateFramesPerSecond) {
		return nil
	}
	game.globalTime = time.Now()

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
	//game.inventory.Draw(screen)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
