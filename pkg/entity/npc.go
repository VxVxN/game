package entity

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/VxVxN/game/pkg/scriptmanager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image/color"
)

type NPC struct {
	Entity
	name            string
	nameFont        font.Face
	playerPosition  base.Position
	cfg             *config.Config
	scriptManager   *scriptmanager.ScriptManager
	dialogueManager *scriptmanager.DialogueManager
}

func NewNPC(name string, position base.Position, speed float64, imagePath string, x0, y0, framesCount int, gameMap *gamemap.Map, cfg *config.Config) (*NPC, error) {
	animation, err := animation.NewAnimation(imagePath, x0, y0, framesCount, cfg.Common.TileSize)
	if err != nil {
		return nil, err
	}

	font, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF font: %v", err)
	}

	nameFont, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create name face(font): %v", err)
	}

	scriptManager := scriptmanager.NewScriptManager(gameMap)

	dialogueManager, err := scriptmanager.NewDialogueManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init dialogueManager: %v", err)
	}

	return &NPC{
		Entity: Entity{
			Position:  position,
			xp:        10000,
			animation: animation,
			speed:     speed,
		},
		name:            name,
		nameFont:        nameFont,
		cfg:             cfg,
		scriptManager:   scriptManager,
		dialogueManager: dialogueManager,
	}, nil
}

func (npc *NPC) Update(playerPosition base.Position) {
	npc.playerPosition = playerPosition
	if npc.IsDead() {
		return
	}
	var action scriptmanager.Action
	npc.Position, action = npc.scriptManager.Update(npc.Position, npc.speed)
	var key ebiten.Key
	switch action {
	case scriptmanager.MoveUp:
		key = ebiten.KeyUp
	case scriptmanager.MoveDown:
		key = ebiten.KeyDown
	case scriptmanager.MoveLeft:
		key = ebiten.KeyLeft
	case scriptmanager.MoveRight:
		key = ebiten.KeyRight
	}
	npc.animation.Update(key)
	npc.dialogueManager.Update()
}

func (npc *NPC) Draw(screen *ebiten.Image) {
	tileSize := float64(npc.cfg.Common.TileSize)
	windowWidth := float64(npc.cfg.Common.WindowWidth)
	windowHeight := float64(npc.cfg.Common.WindowHeight)
	x := tileSize*-npc.playerPosition.X + tileSize*npc.Position.X + windowWidth/2
	y := tileSize*-npc.playerPosition.Y + tileSize*npc.Position.Y + windowHeight/2
	text.Draw(screen, npc.name, npc.nameFont, int(x+2), int(y), color.Black)

	npc.dialogueManager.Draw(screen, x, y)
}

func (npc *NPC) SetScripts(scripts []*scriptmanager.Script) {
	npc.scriptManager.SetScripts(scripts)
}

func (npc *NPC) Trigger() {
	npc.dialogueManager.Trigger()
}

func (npc *NPC) IsEndDialogue() bool {
	return npc.dialogueManager.IsEndDialogue()
}

func (npc *NPC) AddDialogue(dialogue *scriptmanager.PieceDialogue) {
	npc.dialogueManager.AddDialogue(dialogue)
}
