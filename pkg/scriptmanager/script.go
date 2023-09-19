package scriptmanager

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/gamemap"
)

type Script struct {
	currentAction  int
	durationAction int
	actions        []Action
}

func NewScript(actions []Action) *Script {
	return &Script{actions: actions}
}

func (script *Script) Run(gameMap *gamemap.Map, position base.Position, speed float64) (base.Position, Action, bool) {
	if script.currentAction >= len(script.actions) {
		script.currentAction = 0
	}
	action := script.actions[script.currentAction]
	switch action {
	case MoveUp:
		position.Y -= speed
	case MoveDown:
		position.Y += speed
	case MoveLeft:
		position.X -= speed
	case MoveRight:
		position.X += speed
	case Pause:
	}
	if !gameMap.IsCanMove(position.X, position.Y) {
		return position, action, false
	}
	if script.durationAction > 10 {
		script.durationAction = 0
		script.currentAction++
	}
	script.durationAction++
	return position, action, true
}

type Action int

const (
	MoveUp Action = iota
	MoveDown
	MoveLeft
	MoveRight
	Pause
)
