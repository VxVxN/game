package player

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	Position  base.Position
	animation *animation.Animation
}

func NewPlayer(position base.Position, imagePath string, tileSize, framesCount int) (*Player, error) {
	animation, err := animation.NewAnimation(imagePath, framesCount, tileSize)
	if err != nil {
		return nil, err
	}

	return &Player{
		Position:  position,
		animation: animation,
	}, nil
}

func (player *Player) Image() *ebiten.Image {
	return player.animation.GetCurrentFrame()
}

func (player *Player) Move(key ebiten.Key) {
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
