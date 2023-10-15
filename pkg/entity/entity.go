package entity

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	IsDead() bool
	XP() int
	DecreaseXP(int)
	Image() *ebiten.Image
	Position() base.Position
	SetPosition(position base.Position)
	Update(playerPosition base.Position)
	Draw(screen *ebiten.Image)
}

// BaseEntity base structure for any entity
type BaseEntity struct {
	name      string
	position  base.Position
	xp        int
	animation *animation.Animation
	speed     float64
}

func (entity *BaseEntity) IsDead() bool {
	return entity.xp <= 0
}

func (entity *BaseEntity) XP() int {
	return entity.xp / 100
}

func (entity *BaseEntity) DecreaseXP(value int) {
	entity.xp -= value * 100
}

func (entity *BaseEntity) Image() *ebiten.Image {
	return entity.animation.GetCurrentFrame()
}

func (entity *BaseEntity) Position() base.Position {
	return entity.position
}

func (entity *BaseEntity) SetPosition(position base.Position) {
	entity.position = position
}

func (entity *BaseEntity) SetX(x float64) {
	entity.position.X = x
}

func (entity *BaseEntity) SetY(y float64) {
	entity.position.Y = y
}
