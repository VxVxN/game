package game

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/camera"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/entity"
	_interface "github.com/VxVxN/game/pkg/entity/interface"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/inventory"
	"github.com/VxVxN/game/pkg/item"
	"github.com/VxVxN/game/pkg/menu"
	"github.com/VxVxN/game/pkg/quest"
	"github.com/VxVxN/game/pkg/scriptmanager"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Game struct {
	gameMap         *gamemap.Map
	cfg             *config.Config
	player          *entity.Player
	npc             *entity.NPC
	entities        []_interface.Entity
	eventManager    *eventmanager.EventManager
	camera          *camera.Camera
	globalTime      time.Time
	gameOverFace    font.Face
	items           []*item.Item
	inventory       *inventory.Inventory
	questsMenu      *quest.QuestsMenu
	stage           Stage
	menu            *menu.Menu
	isShowDebugInfo bool
}

type Stage int

const (
	GameStage Stage = iota
	DialogueStage
	InventoryStage
	QuestMenuStage
	MenuStage
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

	baseEntitySpped := 0.1
	playerPosition.Y++
	npc, err := entity.NewNPC("Bob", playerPosition, baseEntitySpped, cfg.Player.ImagePath, 96, 128, cfg.Player.FrameCount, gameMap, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}
	enemy, err := entity.NewEnemy("Monster", findPosition(cfg, gameMap), baseEntitySpped, cfg.Player.ImagePath, 0, 128, cfg.Player.FrameCount, gameMap, player, func() {
		player.AddExperience(10)
		player.AddCoins(5)
	}, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create player: %v", err)
	}
	playerPosition.X++
	keyItem, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 224, 4192, cfg.Common.TileSize, item.KeyType)
	if err != nil {
		return nil, fmt.Errorf("failed to create key item: %v", err)
	}
	//enemy.SetScripts([]*scriptmanager.Script{
	//	scriptmanager.NewScript([]scriptmanager.State{
	//		scriptmanager.NewMoveRightState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveRightState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveRightState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState(),
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState(),
	//		scriptmanager.NewMoveLeftState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveLeftState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveLeftState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState(),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState()}),
	//	scriptmanager.NewScript([]scriptmanager.State{
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveUpState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState(),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewMoveDownState(enemy, baseEntitySpped),
	//		scriptmanager.NewPauseState()}),
	//})
	enemy.SetScripts([]*scriptmanager.Script{scriptmanager.NewScript([]scriptmanager.State{scriptmanager.NewFollowForEntityState(enemy, player, baseEntitySpped, gameMap)})})
	npc.AddDialogue(&scriptmanager.PieceDialogue{
		CanStartDialogue: true,
		Replicas:         []string{"Hello stranger", "Do you want a quest?"},
		Answers: []scriptmanager.Answer{
			{
				Text: "Yes, of course",
				Action: func() {
					player.TakeQuest(quest.NewQuest("First quest", []*quest.Goal{
						{
							NeedItems: []quest.NeedItem{
								{
									Type:       item.KeyType,
									NumberNeed: 1,
								},
							},
						},
					}, func() {
						player.AddCoins(100)
						player.DeleteItem(item.KeyType)

						axeItem, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 160, 4192, cfg.Common.TileSize, item.AxeType)
						if err != nil {
							panic(err)
						}
						player.TakeItem(axeItem)
						axeItem2, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 160, 4192, cfg.Common.TileSize, item.AxeType)
						if err != nil {
							panic(err)
						}
						player.TakeItem(axeItem2)
						axeItem3, err := item.NewItem(playerPosition, cfg.Map.TileSetPath, 160, 4192, cfg.Common.TileSize, item.AxeType)
						if err != nil {
							panic(err)
						}
						player.TakeItem(axeItem3)
					}))
				},
				NextPieceDialogue: &scriptmanager.PieceDialogue{
					Replicas: []string{"Bye stranger"},
				},
			},
			{
				Text: "No",
				NextPieceDialogue: &scriptmanager.PieceDialogue{
					Replicas:         []string{"Bye stranger, see you"},
					CanStartDialogue: true,
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

	questMenu, err := quest.NewQuestsMenu(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init quest menu: %v", err)
	}

	game := &Game{
		gameMap:         gameMap,
		cfg:             cfg,
		player:          player,
		npc:             npc,
		entities:        []_interface.Entity{npc, enemy},
		globalTime:      time.Now(),
		isShowDebugInfo: false,
		eventManager:    eventmanager.NewEventManager(),
		camera:          camera.NewCamera(cfg),
		gameOverFace:    gameOverFace,
		items:           []*item.Item{keyItem},
		inventory:       inventory.NewInventory(cfg),
		questsMenu:      questMenu,
		stage:           MenuStage,
	}

	menu, err := menu.NewMenu(cfg, []menu.ButtonOptions{
		{
			Text: "New game",
			Action: func() {
				game.stage = GameStage
			},
		},
		{
			Text: "Exit",
			Action: func() {
				os.Exit(0)
			},
		}})
	if err != nil {
		return nil, fmt.Errorf("failed new menu: %v", err)
	}

	game.menu = menu

	game.camera.AddPlayerImage(player.Image())
	game.camera.AddBackgroundImage(gameMap.BackgroundImage())
	game.camera.AddFrontImages(gameMap.FrontImages())
	game.camera.SetItems([]*item.Item{keyItem})

	game.addEvents(gameMap, player)

	switch game.cfg.Common.Mode {
	case config.ViewMode:
		zoom := 0.2
		// 2.0833 calculation constant for aligning a map of any size based on zoom
		game.player.SetPosition(base.NewPosition(float64(game.cfg.Map.Width)/2.0833*zoom, float64(game.cfg.Map.Height)/2.0833*zoom))

		game.camera.SetZoom(zoom)
	case config.DeveloperMode:
		game.stage = GameStage
		game.isShowDebugInfo = true
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
			if !gameMap.IsCanMove(player.Position().X, player.Position().Y-1) {
				return
			}
			player.Move(ebiten.KeyUp)
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyUp, func() {
		switch game.stage {
		case DialogueStage:
			game.npc.DialogueManager.NextAnswer()
		case MenuStage:
			game.menu.BeforeMenuItem()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position().X, player.Position().Y+1) {
				return
			}
			player.Move(ebiten.KeyDown)
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyDown, func() {
		switch game.stage {
		case DialogueStage:
			game.npc.DialogueManager.BeforeAnswer()
		case MenuStage:
			game.menu.NextMenuItem()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyRight, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position().X+1, player.Position().Y) {
				return
			}
			player.Move(ebiten.KeyRight)
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyRight, func() {
		switch game.stage {
		case InventoryStage:
			game.inventory.NextItem()
		}
	})
	game.eventManager.AddPressEvent(ebiten.KeyLeft, func() {
		switch game.stage {
		case GameStage:
			if !gameMap.IsCanMove(player.Position().X-1, player.Position().Y) {
				return
			}
			player.Move(ebiten.KeyLeft)
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyLeft, func() {
		switch game.stage {
		case InventoryStage:
			game.inventory.BeforeItem()
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeySpace, func() {
		switch game.stage {
		case GameStage:
			if utils.CanAction(game.player.Position(), game.npc.Position()) && game.npc.DialogueManager.CanStartDialogue {
				game.npc.Trigger()
				game.stage = DialogueStage
			}
			for _, item := range game.items {
				if utils.CanAction(game.player.Position(), item.Position()) {
					item.Trigger()
					game.player.TakeItem(item)
				}
			}
			for _, e := range game.entities {
				if e.IsDead() {
					continue
				}
				if enemy, ok := e.(*entity.Enemy); ok {
					if utils.CanAction(game.player.Position(), e.Position()) {
						e.DecreaseXP(game.player.Attack())
					}
					if e.IsDead() {
						enemy.GetAward()
					}
				}
			}
		case DialogueStage:
			if game.npc.DialogueManager.NeedAnswer() {
				game.npc.DialogueManager.DoAnswer()
				game.npc.DialogueManager.PieceDialogue = game.npc.DialogueManager.NextPieceDialogue()
				if game.npc.DialogueManager.IsEndDialogue {
					game.stage = GameStage
					return
				}
				return
			}
			game.npc.DialogueManager.NextReplica()
			if game.npc.DialogueManager.IsEndDialogue {
				game.stage = GameStage
			}
		case MenuStage:
			game.menu.ClickActiveButton()
		case InventoryStage:
			player.TakeItemInHand(game.inventory.GetActiveItem())
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyTab, func() {
		game.isShowDebugInfo = !game.isShowDebugInfo
	})
	game.eventManager.AddPressedEvent(ebiten.KeyI, func() {
		switch game.stage {
		case GameStage:
			game.inventory.OnOff()
			game.inventory.Update(game.player.Items())
			game.stage = InventoryStage
		case InventoryStage:
			game.inventory.OnOff()
			game.stage = GameStage
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyQ, func() {
		switch game.stage {
		case GameStage:
			game.questsMenu.OnOff()
			game.questsMenu.Update(game.player.Quests())
			game.stage = QuestMenuStage
		case QuestMenuStage:
			game.questsMenu.OnOff()
			game.stage = GameStage
		}
	})
	game.eventManager.AddPressedEvent(ebiten.KeyEscape, func() {
		switch game.stage {
		case MenuStage:
			game.stage = GameStage
		default:
			game.stage = MenuStage
		}
	})
	game.eventManager.SetDefaultEvent(func() {
		player.Stand()
	})
}

func (game *Game) Update() error {
	game.eventManager.Update()
	if game.stage != GameStage {
		return nil
	}
	game.npc.Update(game.player.Position())
	if time.Since(game.globalTime) < time.Second/time.Duration(game.cfg.Common.RefreshRateFramesPerSecond) {
		return nil
	}
	game.globalTime = time.Now()

	for _, entity := range game.entities {
		entity.Update(game.player.Position())
	}

	game.camera.AddPlayerImage(game.player.Image())
	game.camera.UpdatePlayer(game.player.Position())
	game.camera.UpdateEntities(game.entities)

	game.player.Update(base.Position{})
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	switch game.stage {
	case MenuStage:
		game.menu.Draw(screen)
		return
	default:
	}
	game.camera.Draw(screen)

	if game.isShowDebugInfo {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("x = %f, y = %f\nLayers: %d", game.player.Position().X, game.player.Position().Y,
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
	for _, entity := range game.entities {
		entity.Draw(screen)
	}
	game.player.Draw(screen)
	game.inventory.Draw(screen)
	game.questsMenu.Draw(screen)
}

func (game *Game) Layout(screenWidthPx, screenHeightPx int) (int, int) {
	return game.cfg.Common.WindowWidth, game.cfg.Common.WindowHeight
}
