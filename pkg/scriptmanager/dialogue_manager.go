package scriptmanager

import (
	"fmt"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/eventmanager"
	"github.com/VxVxN/game/pkg/label"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"os"
	"time"
)

type DialogueManager struct {
	dialogue *PieceDialogue

	globalTime   time.Time
	face         font.Face
	eventManager *eventmanager.EventManager
	cfg          *config.Config
	isRun        bool
}

func NewDialogueManager(cfg *config.Config) (*DialogueManager, error) {
	data, err := os.ReadFile("assets/fonts/Zack and Sarah.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to open font file: %v", err)
	}
	font, err := opentype.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
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
		cfg:          cfg,
	}

	eventManager.AddEvent(ebiten.KeyUp, func() {
		manager.dialogue.NextAnswer()
	})
	eventManager.AddEvent(ebiten.KeyDown, func() {
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

	replicLabel := label.NewLabel(manager.face)
	replicLabel.X = 0
	replicLabel.Y = 0
	replicLabel.Width = float64(manager.cfg.Common.WindowWidth)
	replicLabel.Height = y - 32
	replicLabel.AlignVertical = label.AlignVerticalBottom
	replicLabel.AlignHorizontal = label.AlignHorizontalCenter
	replicLabel.Text = manager.dialogue.CurrentReplica()
	//replicLabel.ContainerColor = color.RGBA{R: 100, G: 200, B: 100, A: 160}
	replicLabel.Draw(screen)

	if !manager.dialogue.NeedAnswer() {
		return
	}

	answerLabel := label.NewLabel(manager.face)
	for index, answer := range manager.dialogue.Answers {
		answerColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		if manager.dialogue.IsActiveAnswer(index) {
			answerColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		}
		answerLabel.X = 0
		answerLabel.Y = float64(int(y) + 32 + index*32)
		answerLabel.Width = float64(manager.cfg.Common.WindowWidth)
		//answerLabel.Height = y - 32
		answerLabel.AlignHorizontal = label.AlignHorizontalCenter
		answerLabel.Text = answer.Text
		answerLabel.Color = answerColor
		answerLabel.Draw(screen)

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
