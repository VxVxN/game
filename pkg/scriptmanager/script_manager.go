package scriptmanager

import (
	"github.com/VxVxN/game/internal/base"
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

func (manager *ScriptManager) Update(position base.Position, speed float64) (base.Position, Action) {
	var ok bool
	action := Pause
	newPosition := position
	for {
		if manager.currentScript >= len(manager.scripts) {
			break
		}
		script := manager.scripts[manager.currentScript]
		newPosition, action, ok = script.Run(manager.gameMap, position, speed)
		if ok {
			break
		}
		manager.currentScript++
	}
	return newPosition, action
}
