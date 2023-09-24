package scriptmanager

import (
	"fmt"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"image/color"
	"time"
)

type DialogueManager struct {
	dialogue *PieceDialogue

	globalTime   time.Time
	face         font.Face
	eventManager *eventmanager.EventManager
	isRun        bool
}

func NewDialogueManager() (*DialogueManager, error) {
	font, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF font: %v", err)
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 20,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create face(font): %v", err)
	}

	eventManager := eventmanager.NewEventManager()

	manager := &DialogueManager{
		globalTime:   time.Now(),
		face:         face,
		eventManager: eventManager,
	}

	eventManager.AddEvent(ebiten.KeyRight, func() {
		manager.dialogue.NextAnswer()
	})
	eventManager.AddEvent(ebiten.KeyLeft, func() {
		manager.dialogue.BeforeAnswer()
	})
	eventManager.AddEvent(ebiten.KeySpace, func() {
		if manager.dialogue.NeedAnswer() {
			manager.dialogue.DoAnswer()
			if manager.IsEndDialogue() {
				return
			}
			manager.dialogue = manager.dialogue.NextPieceDialogue()
			return
		}
		manager.dialogue.NextReplica()
	})

	return manager, nil
}

func (manager *DialogueManager) Draw(screen *ebiten.Image, x, y float64) {
	if !manager.isRun || manager.IsEndDialogue() {
		return
	}

	alignmentTextByX := int(x) - int(float64(len(manager.dialogue.CurrentReplica()))*2)
	text.Draw(screen, manager.dialogue.CurrentReplica(), manager.face, alignmentTextByX, int(y-20), color.Black)

	if !manager.dialogue.NeedAnswer() {
		return
	}

	for index, answer := range manager.dialogue.Answers {
		activeAnswer := color.Black
		if manager.dialogue.IsActiveAnswer(index) {
			activeAnswer = color.White
		}
		alignmentTextByX = int(x) - int(float64(len(answer.Text))*2)
		text.Draw(screen, answer.Text, manager.face, alignmentTextByX+40*index, int(y+50), activeAnswer)
	}
}

func (manager *DialogueManager) Trigger() {
	if manager.dialogue == nil || manager.IsEndDialogue() {
		return
	}
	if !manager.isRun {
		manager.globalTime = time.Now()
	}
	manager.isRun = true
}

func (manager *DialogueManager) Update() {
	if time.Since(manager.globalTime) < time.Second/time.Duration(5) || !manager.isRun {
		return
	}
	manager.globalTime = time.Now()
	manager.eventManager.Update()
}

func (manager *DialogueManager) IsEndDialogue() bool {
	return manager.dialogue.isEndDialogue
}

func (manager *DialogueManager) AddDialogue(dialogue *PieceDialogue) {
	manager.dialogue = dialogue
}
