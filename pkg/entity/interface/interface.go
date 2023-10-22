package _interface

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/hajimehoshi/ebiten/v2"
)

type Entity interface {
	IsDead() bool
	XP() int
	DecreaseXP(int)
	Image() *ebiten.Image
	Position() base.Position
	SetPosition(position base.Position)
	SetX(x float64)
	SetY(y float64)
	Update(playerPosition base.Position)
	Draw(screen *ebiten.Image)
}
