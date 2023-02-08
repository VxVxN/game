package eventmanager

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EventManager struct {
	events map[ebiten.Key][]func()
}

func NewEventManager() *EventManager {
	return &EventManager{
		events: make(map[ebiten.Key][]func()),
	}
}

func (eventManager EventManager) Process() {
	var key ebiten.Key
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		key = ebiten.KeyUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		key = ebiten.KeyDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		key = ebiten.KeyLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		key = ebiten.KeyRight
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		key = ebiten.KeyEscape
	}
	events, ok := eventManager.events[key]
	if !ok {
		return // we don't have events
	}
	for _, event := range events {
		event()
	}
}

func (eventManager EventManager) AddEvent(key ebiten.Key, event func()) {
	events, _ := eventManager.events[key]
	events = append(events, event)
	eventManager.events[key] = events
}
