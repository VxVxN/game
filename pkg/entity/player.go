package entity

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Entity
	satiety int
}

func NewPlayer(position base.Position, imagePath string, x0, y0, tileSize, framesCount int) (*Player, error) {
	animation, err := animation.NewAnimation(imagePath, x0, y0, framesCount, tileSize)
	if err != nil {
		return nil, err
	}

	return &Player{
		Entity: Entity{
			Position:  position,
			xp:        10000,
			animation: animation,
		},
		satiety: 10000,
	}, nil
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