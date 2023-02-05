package eventmanager

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/VxVxN/game/internal/player"
	"os"
)

type EventManager struct {
	player *player.Player
}

func NewEventManager(player *player.Player) *EventManager {
	return &EventManager{
		player: player,
	}
}

func (eventManager EventManager) Process() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		eventManager.player.Position.Y--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		eventManager.player.Position.Y++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		eventManager.player.Position.X--
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		eventManager.player.Position.X++
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0) // todo add normal game end processing
	}
}
