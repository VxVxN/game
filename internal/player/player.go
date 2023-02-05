package player

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	position base.Position
	image    *ebiten.Image
}

func NewPlayer(position base.Position, imagePath string) (*Player, error) {
	image, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, err
	}

	return &Player{
		position: position,
		image:    image,
	}, nil
}

func (player Player) X() int {
	return player.position.X
}

func (player Player) Y() int {
	return player.position.Y
}

func (player Player) Image() *ebiten.Image {
	return player.image
}
