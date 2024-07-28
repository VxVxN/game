package quest

import (
	"fmt"
	"image/color"
	"os"

	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/label"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type QuestsMenu struct {
	cfg        *config.Config
	quests     []*Quest
	face       font.Face
	isLaunched bool
}

func NewQuestsMenu(cfg *config.Config) (*QuestsMenu, error) {
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

	return &QuestsMenu{cfg: cfg, face: face}, nil
}

func (menu *QuestsMenu) OnOff() {
	menu.isLaunched = !menu.isLaunched
}

func (menu *QuestsMenu) Update(quests []*Quest) {
	menu.quests = quests
}

func (menu *QuestsMenu) Draw(screen *ebiten.Image) {
	if !menu.isLaunched {
		return
	}

	x := float64(50)
	y := float64(50)
	ebitenutil.DrawRect(screen, x, y, float64(menu.cfg.Common.WindowWidth-100), float64(menu.cfg.Common.WindowHeight-100), color.RGBA{R: 100, G: 43, B: 43, A: 150})
	for i, quest := range menu.quests {
		var statusQuest string
		if quest.IsCompleted() {
			statusQuest = "(completed)"
		}
		text := quest.Name + statusQuest + ": \n    " + quest.GoalsDescription()
		questLabel := label.NewLabel(menu.face, x+10, (float64(i))*54.0+10, text)
		questLabel.ContainerWidth = float64(menu.cfg.Common.WindowWidth - 100 - 20)
		questLabel.AlignVertical = label.AlignVerticalCenter
		questLabel.ContainerColor = color.RGBA{R: 100, G: 200, B: 100, A: 160}
		questLabel.Draw(screen)
	}
}
