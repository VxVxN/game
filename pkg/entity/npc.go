package entity

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/animation"
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
	name           string
	nameFont       font.Face
	playerPosition base.Position
	cfg            *config.Config
}

func NewNPC(name string, position base.Position, imagePath string, x0, y0, framesCount int, cfg *config.Config) (*NPC, error) {
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

	return &NPC{
		Entity: Entity{
			Position:  position,
			xp:        10000,
			animation: animation,
		},
		name:     name,
		nameFont: nameFont,
		cfg:      cfg,
	}, nil
}

func (npc *NPC) Update(playerPosition base.Position) {
	npc.playerPosition = playerPosition
	if npc.IsDead() {
		return
	}
	npc.animation.SetDefaultFrame()
}

func (npc *NPC) Draw(screen *ebiten.Image) {
	tileSize := float64(npc.cfg.Common.TileSize)
	windowWidth := float64(npc.cfg.Common.WindowWidth)
	windowHeight := float64(npc.cfg.Common.WindowHeight)
	x := tileSize*-npc.playerPosition.X + tileSize*npc.Position.X + windowWidth/2
	y := tileSize*-npc.playerPosition.Y + tileSize*npc.Position.Y + windowHeight/2
	//game.cfg.Common.WindowWidth
	text.Draw(screen, npc.name, npc.nameFont, int(x+2), int(y), color.Black)
}
