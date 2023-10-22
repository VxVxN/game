package scriptmanager

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/gamemap"
)

type Script struct {
	currentAction  int
	durationAction int
	states         []State
}

func NewScript(states []State) *Script {
	return &Script{states: states}
}

func (script *Script) Run(gameMap *gamemap.Map, position base.Position) (State, bool) {
	if script.currentAction >= len(script.states) {
		script.currentAction = 0
	}
	state := script.states[script.currentAction]
	state.Do()

	if !gameMap.IsCanMove(position.X, position.Y) {
		return state, false
	}
	if script.durationAction > 10 {
		script.durationAction = 0
		script.currentAction++
	}
	script.durationAction++
	return state, true
}

type Action int

const (
	MoveUp Action = iota
	MoveDown
	MoveLeft
	MoveRight
	Pause
)
