package eventmanager

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EventManager struct {
	events       map[ebiten.Key][]func()
	defaultEvent func()
}

func NewEventManager() *EventManager {
	return &EventManager{
		events: make(map[ebiten.Key][]func()),
	}
}

var supportedKeys = []ebiten.Key{
	ebiten.KeyUp,
	ebiten.KeyDown,
	ebiten.KeyLeft,
	ebiten.KeyRight,
	ebiten.KeyTab,
	ebiten.KeyEscape,
	ebiten.KeySpace,
}

func (eventManager *EventManager) Update() {
	var key ebiten.Key

	for _, supportedKey := range supportedKeys {
		if ebiten.IsKeyPressed(supportedKey) {
			key = supportedKey
		}
	}
	events, ok := eventManager.events[key]
	if !ok && eventManager.defaultEvent != nil {
		eventManager.defaultEvent()
		return // we don't have events
	}
	for _, event := range events {
		event()
	}
}

func (eventManager *EventManager) AddEvent(key ebiten.Key, event func()) {
	events, _ := eventManager.events[key]
	events = append(events, event)
	eventManager.events[key] = events
}

func (eventManager *EventManager) SetDefaultEvent(event func()) {
	eventManager.defaultEvent = event
}
