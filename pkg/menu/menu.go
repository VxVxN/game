package menu

import (
	"fmt"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/label"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/opentype"
	"image/color"
	"os"
)

type Menu struct {
	cfg            *config.Config
	activeItemMenu int
	buttons        []*label.Label
	buttonOptions  []ButtonOptions
}

type ButtonOptions struct {
	Text   string
	Action func()
}

func NewMenu(cfg *config.Config, buttonOptions []ButtonOptions) (*Menu, error) {
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
	buttons := make([]*label.Label, len(buttonOptions))
	heightLabel := float64(32)

	for i, option := range buttonOptions {
		y := float64(cfg.Common.WindowHeight/2) + float64(i)*heightLabel + 10
		if i == 0 {
			y = float64(cfg.Common.WindowHeight / 2)
		}
		buttons[i] = label.NewLabel(face, 0, y, option.Text)
		buttons[i].ContainerWidth = float64(cfg.Common.WindowWidth)
		buttons[i].ContainerHeight = heightLabel
		buttons[i].AlignVertical = label.AlignVerticalCenter
		buttons[i].AlignHorizontal = label.AlignHorizontalCenter
	}
	return &Menu{cfg: cfg, buttons: buttons, buttonOptions: buttonOptions}, nil
}

func (menu *Menu) NextMenuItem() {
	menu.activeItemMenu++
	if menu.activeItemMenu > len(menu.buttons)-1 {
		menu.activeItemMenu = len(menu.buttons) - 1
	}
}

func (menu *Menu) BeforeMenuItem() {
	menu.activeItemMenu--
	if menu.activeItemMenu < 0 {
		menu.activeItemMenu = 0
	}
}

func (menu *Menu) ClickActiveButton() {
	menu.buttonOptions[menu.activeItemMenu].Action()
}

func (menu *Menu) Draw(screen *ebiten.Image) {
	deactivatedButtonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
	activatedButtonColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	for i, _ := range menu.buttons {
		buttonColor := deactivatedButtonColor
		if menu.activeItemMenu == i {
			buttonColor = activatedButtonColor
		}
		menu.buttons[i].ContainerColor = buttonColor
		menu.buttons[i].Draw(screen)
	}
}
