package inventory

import (
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/item"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type Inventory struct {
	cfg             *config.Config
	items           []*item.Item
	activeItemIndex int
	isLaunched      bool
}

func NewInventory(cfg *config.Config) *Inventory {
	return &Inventory{cfg: cfg}
}

func (inventory *Inventory) OnOff() {
	inventory.isLaunched = !inventory.isLaunched
}

func (inventory *Inventory) Update(items []*item.Item) {
	inventory.items = items
}

func (inventory *Inventory) Draw(screen *ebiten.Image) {
	if !inventory.isLaunched {
		return
	}
	ebitenutil.DrawRect(screen, 50, 50, float64(inventory.cfg.Common.WindowWidth-100), float64(inventory.cfg.Common.WindowHeight-100), color.RGBA{R: 100, G: 43, B: 43, A: 150})
	itemColor := color.RGBA{R: 200, G: 200, B: 200, A: 255}
	activatedItemColor := color.RGBA{R: 255, G: 165, B: 0, A: 255}

	for i, item := range inventory.items {
		x0 := float64(100 + i*32 + 20*i)
		y0 := float64(100)
		if inventory.activeItemIndex == i {
			ebitenutil.DrawRect(screen, x0, y0, 32, 32, activatedItemColor)
		} else {
			ebitenutil.DrawRect(screen, x0, y0, 32, 32, itemColor)
		}

		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(x0, y0)
		screen.DrawImage(item.Image(), options)
	}
}

func (inventory *Inventory) NextItem() {
	inventory.activeItemIndex++
	if inventory.activeItemIndex > len(inventory.items)-1 {
		inventory.activeItemIndex = len(inventory.items) - 1
	}
}

func (inventory *Inventory) BeforeItem() {
	inventory.activeItemIndex--
	if inventory.activeItemIndex < 0 {
		inventory.activeItemIndex = 0
	}
}

func (inventory *Inventory) GetActiveItem() *item.Item {
	return inventory.items[inventory.activeItemIndex]
}
