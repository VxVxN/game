package entity

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/gamemap"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/VxVxN/game/pkg/scriptmanager"
	"github.com/VxVxN/game/pkg/utils"
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
	player         *Player
	cfg            *config.Config
	scriptManager  *scriptmanager.ScriptManager
	rewardCallback func()
	delayAttack    int
}

func NewEnemy(name string, position base.Position, speed float64, imagePath string, x0, y0, framesCount int, gameMap *gamemap.Map, player *Player, rewardCallback func(), cfg *config.Config) (*Enemy, error) {
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
		player:         player,
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
	x := tileSize*-enemy.player.Position().X + tileSize*enemy.Position().X + windowWidth/2
	y := tileSize*-enemy.player.Position().Y + tileSize*enemy.Position().Y + windowHeight/2
	text.Draw(screen, enemy.name, enemy.nameFont, int(x+2), int(y)-15, color.Black)
	text.Draw(screen, "XP "+strconv.Itoa(enemy.XP()), enemy.nameFont, int(x+2), int(y), color.Black)
}

func (enemy *Enemy) Update(playerPosition base.Position) {
	if enemy.IsDead() {
		return
	}
	oldPosition := enemy.position
	state := enemy.scriptManager.Update(enemy.Position(), enemy.speed)
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
	case *scriptmanager.FollowForEntityState:
		if oldPosition.Y > enemy.position.Y {
			key = ebiten.KeyUp
		}
		if oldPosition.Y < enemy.position.Y {
			key = ebiten.KeyDown
		}
		if oldPosition.X > enemy.position.X {
			key = ebiten.KeyLeft
		}
		if oldPosition.X < enemy.position.X {
			key = ebiten.KeyRight
		}
	}
	enemy.animation.Update(key)

	if utils.CanAction(enemy.player.Position(), enemy.Position()) && enemy.delayAttack > 50 {
		enemy.player.DecreaseXP(10)
		enemy.delayAttack = 0
	}
	enemy.delayAttack++
}

func (enemy *Enemy) GetAward() {
	enemy.rewardCallback()
}

func (enemy *Enemy) SetScripts(scripts []*scriptmanager.Script) {
	enemy.scriptManager.SetScripts(scripts)
}
