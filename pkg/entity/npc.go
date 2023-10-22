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
	BaseEntity
	nameFont        font.Face
	playerPosition  base.Position
	cfg             *config.Config
	scriptManager   *scriptmanager.ScriptManager
	DialogueManager *scriptmanager.DialogueManager
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
		return nil, fmt.Errorf("failed to init DialogueManager: %v", err)
	}

	return &NPC{
		BaseEntity: BaseEntity{
			position:  position,
			name:      name,
			xp:        10000,
			animation: animation,
			speed:     speed,
		},
		nameFont:        nameFont,
		cfg:             cfg,
		scriptManager:   scriptManager,
		DialogueManager: dialogueManager,
	}, nil
}

func (npc *NPC) Update(playerPosition base.Position) {
	npc.playerPosition = playerPosition
	if npc.IsDead() {
		return
	}
	state := npc.scriptManager.Update()
	var key ebiten.Key
	switch state.(type) {
	case *scriptmanager.MoveUpState:
		key = ebiten.KeyUp
	case *scriptmanager.MoveDownState:
		key = ebiten.KeyDown
	case *scriptmanager.MoveLeftState:
		key = ebiten.KeyLeft
	case *scriptmanager.MoveRightState:
		key = ebiten.KeyRight
	}
	npc.animation.Update(key)
}

func (npc *NPC) Draw(screen *ebiten.Image) {
	tileSize := float64(npc.cfg.Common.TileSize)
	windowWidth := float64(npc.cfg.Common.WindowWidth)
	windowHeight := float64(npc.cfg.Common.WindowHeight)
	x := tileSize*-npc.playerPosition.X + tileSize*npc.Position().X + windowWidth/2
	y := tileSize*-npc.playerPosition.Y + tileSize*npc.Position().Y + windowHeight/2
	text.Draw(screen, npc.name, npc.nameFont, int(x+2), int(y), color.Black)

	npc.DialogueManager.Draw(screen, x, y)
}

func (npc *NPC) SetScripts(scripts []*scriptmanager.Script) {
	npc.scriptManager.SetScripts(scripts)
}

func (npc *NPC) Trigger() {
	npc.DialogueManager.Trigger()
}

func (npc *NPC) AddDialogue(dialogue *scriptmanager.PieceDialogue) {
	npc.DialogueManager.AddDialogue(dialogue)
}
