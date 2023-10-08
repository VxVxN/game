package scriptmanager

import (
	"fmt"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/label"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"os"
)

type DialogueManager struct {
	*PieceDialogue

	face  font.Face
	cfg   *config.Config
	isRun bool
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

	manager := &DialogueManager{
		face: face,
		cfg:  cfg,
	}

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
	replicLabel.Text = manager.CurrentReplica()
	//replicLabel.ContainerColor = color.RGBA{R: 100, G: 200, B: 100, A: 160}
	replicLabel.Draw(screen)

	if !manager.NeedAnswer() {
		return
	}

	answerLabel := label.NewLabel(manager.face)
	for index, answer := range manager.Answers {
		answerColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		if manager.IsActiveAnswer(index) {
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
	if manager == nil || manager.IsEndDialogue() {
		return
	}
	manager.isRun = true
}

func (manager *DialogueManager) IsEndDialogue() bool {
	return manager.isEndDialogue
}

func (manager *DialogueManager) AddDialogue(dialogue *PieceDialogue) {
	manager.PieceDialogue = dialogue
}
