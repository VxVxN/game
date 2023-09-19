package entity

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

// Entity base structure for any entity
type Entity struct {
	Position  base.Position
	xp        int
	animation *animation.Animation
}

func (entity *Entity) IsDead() bool {
	return entity.xp <= 0
}

func (entity *Entity) XP() int {
	return entity.xp / 100
}

func (entity *Entity) Image() *ebiten.Image {
	return entity.animation.GetCurrentFrame()
}
