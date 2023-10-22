package scriptmanager

import "github.com/VxVxN/game/pkg/entity/interface"

type State interface {
	Do()
}

type FollowForEntityState struct {
	me          _interface.Entity
	destination _interface.Entity
	speed       float64
}

func NewFollowForEntityState(me _interface.Entity, destination _interface.Entity, speed float64) *FollowForEntityState {
	return &FollowForEntityState{
		me:          me,
		destination: destination,
		speed:       speed,
	}
}

func (state *FollowForEntityState) Do() {
	if state.me.Position().X > state.destination.Position().X {
		state.me.SetX(state.me.Position().X - state.speed)
	}
	if state.me.Position().X < state.destination.Position().X {
		state.me.SetX(state.me.Position().X + state.speed)
	}
	if state.me.Position().Y > state.destination.Position().Y {
		state.me.SetY(state.me.Position().Y - state.speed)
	}
	if state.me.Position().Y < state.destination.Position().Y {
		state.me.SetY(state.me.Position().Y + state.speed)
	}
}

type MoveUpState struct {
	entity _interface.Entity
	speed  float64
}

func NewMoveUpState(entity _interface.Entity, speed float64) *MoveUpState {
	return &MoveUpState{
		entity: entity,
		speed:  speed,
	}
}

func (state *MoveUpState) Do() {
	state.entity.SetY(state.entity.Position().Y - state.speed)
}

type MoveDownState struct {
	entity _interface.Entity
	speed  float64
}

func NewMoveDownState(entity _interface.Entity, speed float64) *MoveDownState {
	return &MoveDownState{
		entity: entity,
		speed:  speed,
	}
}

func (state *MoveDownState) Do() {
	state.entity.SetY(state.entity.Position().Y + state.speed)
}

type MoveLeftState struct {
	entity _interface.Entity
	speed  float64
}

func NewMoveLeftState(entity _interface.Entity, speed float64) *MoveLeftState {
	return &MoveLeftState{
		entity: entity,
		speed:  speed,
	}
}

func (state *MoveLeftState) Do() {
	state.entity.SetX(state.entity.Position().X - state.speed)
}

type MoveRightState struct {
	entity _interface.Entity
	speed  float64
}

func NewMoveRightState(entity _interface.Entity, speed float64) *MoveRightState {
	return &MoveRightState{
		entity: entity,
		speed:  speed,
	}
}

func (state *MoveRightState) Do() {
	state.entity.SetX(state.entity.Position().X + state.speed)
}

type PauseState struct {
}

func NewPauseState() *PauseState {
	return &PauseState{}
}

func (state *PauseState) Do() {}
