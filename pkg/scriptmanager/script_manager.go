package scriptmanager

import (
	"github.com/VxVxN/game/internal/gamemap"
)

type ScriptManager struct {
	currentScript int
	scripts       []*Script
	gameMap       *gamemap.Map
}

func NewScriptManager(gameMap *gamemap.Map) *ScriptManager {
	return &ScriptManager{gameMap: gameMap}
}

func (manager *ScriptManager) SetScripts(scripts []*Script) {
	manager.scripts = scripts
}

func (manager *ScriptManager) AddScript(script *Script) {
	manager.scripts = append(manager.scripts, script)
}

func (manager *ScriptManager) Update() State {
	var ok bool
	var state State
	state = NewPauseState()
	for {
		if manager.currentScript >= len(manager.scripts) {
			break
		}
		script := manager.scripts[manager.currentScript]
		state, ok = script.Run()
		if ok {
			break
		}
		manager.currentScript++
	}
	return state
}
