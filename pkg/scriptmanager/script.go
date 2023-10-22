package scriptmanager

type Script struct {
	currentAction  int
	durationAction int
	states         []State
}

func NewScript(states []State) *Script {
	return &Script{states: states}
}

func (script *Script) Run() (State, bool) {
	if script.currentAction >= len(script.states) {
		script.currentAction = 0
	}
	state := script.states[script.currentAction]
	state.Do()

	if script.durationAction > 10 {
		script.durationAction = 0
		script.currentAction++
	}
	script.durationAction++
	return state, true
}
