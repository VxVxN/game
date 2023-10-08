package item

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Item struct {
	position base.Position
	image    *ebiten.Image
	isTook   bool
}

//type ItemType int
//
//const (
//	AxeType ItemType = iota
//	KeyType
//)

func NewItem(position base.Position, imagePath string, x0, y0, titleSize int) (*Item, error) {
	baseImage, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, err
	}

	image := baseImage.SubImage(image.Rect(x0, y0, x0+titleSize, y0+titleSize)).(*ebiten.Image)

	return &Item{
		position: position,
		image:    image,
	}, nil
}

func (item *Item) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	if item.isTook {
		return
	}
	screen.DrawImage(item.image, options)
}

func (item *Item) Trigger() {
	item.isTook = true
}

func (item *Item) Position() base.Position {
	return item.position
}

func (item *Item) Image() *ebiten.Image {
	return item.image
}
