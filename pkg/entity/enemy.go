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
	"strconv"
)

type Enemy struct {
	BaseEntity
	nameFont       font.Face
	playerPosition base.Position
	cfg            *config.Config
	scriptManager  *scriptmanager.ScriptManager
	rewardCallback func()
}

func NewEnemy(name string, position base.Position, speed float64, imagePath string, x0, y0, framesCount int, gameMap *gamemap.Map, rewardCallback func(), cfg *config.Config) (*Enemy, error) {
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

	return &Enemy{
		BaseEntity: BaseEntity{
			position:  position,
			name:      name,
			xp:        10000,
			animation: animation,
			speed:     speed,
		},
		nameFont:       nameFont,
		cfg:            cfg,
		scriptManager:  scriptManager,
		rewardCallback: rewardCallback,
	}, nil
}

func (enemy *Enemy) Draw(screen *ebiten.Image) {
	if enemy.IsDead() {
		return
	}
	tileSize := float64(enemy.cfg.Common.TileSize)
	windowWidth := float64(enemy.cfg.Common.WindowWidth)
	windowHeight := float64(enemy.cfg.Common.WindowHeight)
	x := tileSize*-enemy.playerPosition.X + tileSize*enemy.Position().X + windowWidth/2
	y := tileSize*-enemy.playerPosition.Y + tileSize*enemy.Position().Y + windowHeight/2
	text.Draw(screen, enemy.name, enemy.nameFont, int(x+2), int(y)-15, color.Black)
	text.Draw(screen, "XP "+strconv.Itoa(enemy.XP()), enemy.nameFont, int(x+2), int(y), color.Black)
}

func (enemy *Enemy) Update(playerPosition base.Position) {
	enemy.playerPosition = playerPosition
	if enemy.IsDead() {
		return
	}
	var action scriptmanager.Action
	position, action := enemy.scriptManager.Update(enemy.Position(), enemy.speed)
	enemy.SetPosition(position)
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
	enemy.animation.Update(key)
}

func (enemy *Enemy) GetAward() {
	enemy.rewardCallback()
}
