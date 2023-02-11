package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Animation struct {
	framesCount                        int
	titleSize                          int
	image                              *ebiten.Image
	currentFrameImage                  *ebiten.Image
	startCurrentFrame, endCurrentFrame int
}

func NewAnimation(imagePath string, framesCount, titleSize int) (*Animation, error) {
	charactersImage, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, err
	}

	currentFrameImage := charactersImage.SubImage(image.Rect(0, 0, titleSize, titleSize)).(*ebiten.Image)

	return &Animation{
		framesCount:       framesCount,
		image:             charactersImage,
		currentFrameImage: currentFrameImage,
		titleSize:         titleSize,
	}, nil
}

func (animation *Animation) GetCurrentFrame() *ebiten.Image {
	return animation.currentFrameImage
}

func (animation *Animation) Update(key ebiten.Key) {
	tileSize := animation.titleSize

	var y0, y1 int

	x0, x1 := animation.NextFrame()

	switch key {
	case ebiten.KeyUp:
		y0, y1 = tileSize*3, tileSize*4
	case ebiten.KeyDown:
		y0, y1 = 0, tileSize
	case ebiten.KeyLeft:
		y0, y1 = tileSize, tileSize*2
	case ebiten.KeyRight:
		y0, y1 = tileSize*2, tileSize*3
	default:
		x0, y0, x1, y1 = 0, 0, tileSize, tileSize
	}
	animation.currentFrameImage = animation.image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
}

func (animation *Animation) SetDefaultFrame() {
	x0, y0, x1, y1 := 0, 0, animation.titleSize, animation.titleSize
	animation.currentFrameImage = animation.image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
}

func (animation *Animation) NextFrame() (int, int) {
	animation.startCurrentFrame++
	animation.endCurrentFrame++
	if animation.endCurrentFrame == animation.framesCount {
		animation.startCurrentFrame = 0
		animation.endCurrentFrame = 1
	}
	return animation.startCurrentFrame * animation.titleSize, animation.endCurrentFrame * animation.titleSize
}
