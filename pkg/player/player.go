package player

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Position  base.Position
	xp        int
	satiety   int
	animation *animation.Animation
}

func NewPlayer(position base.Position, imagePath string, tileSize, framesCount int) (*Player, error) {
	animation, err := animation.NewAnimation(imagePath, framesCount, tileSize)
	if err != nil {
		return nil, err
	}

	return &Player{
		Position:  position,
		xp:        10000,
		satiety:   10000,
		animation: animation,
	}, nil
}

func (player *Player) Image() *ebiten.Image {
	return player.animation.GetCurrentFrame()
}

func (player *Player) Move(key ebiten.Key) {
	if player.IsDead() {
		return
	}
	switch key {
	case ebiten.KeyUp:
		player.Position.Y--
	case ebiten.KeyDown:
		player.Position.Y++
	case ebiten.KeyLeft:
		player.Position.X--
	case ebiten.KeyRight:
		player.Position.X++
	default:
	}
	player.animation.Update(key)
}

func (player *Player) Stand() {
	player.animation.SetDefaultFrame()
}

func (player *Player) Satiety() int {
	return player.satiety / 100
}

func (player *Player) XP() int {
	return player.xp / 100
}

func (player *Player) DecreaseSatiety() {
	if player.IsDead() {
		return
	}
	if player.satiety > 0 {
		player.satiety--
	} else {
		player.xp--
	}
}

func (player *Player) IsDead() bool {
	return player.xp <= 0
}
