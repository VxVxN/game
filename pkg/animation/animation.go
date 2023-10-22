package animation

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Animation struct {
	framesCount                        int
	titleSize                          int
	x0, y0                             int
	image                              *ebiten.Image
	currentFrameImage                  *ebiten.Image
	startCurrentFrame, endCurrentFrame int
}

func NewAnimation(imagePath string, x0, y0, framesCount, titleSize int) (*Animation, error) {
	charactersImage, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, err
	}

	currentFrameImage := charactersImage.SubImage(image.Rect(x0, y0, x0+titleSize, y0+titleSize)).(*ebiten.Image)

	return &Animation{
		framesCount:       framesCount,
		image:             charactersImage,
		currentFrameImage: currentFrameImage,
		titleSize:         titleSize,
		x0:                x0,
		y0:                y0,
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
		y0, y1 = animation.y0+(tileSize*3), animation.y0+(tileSize*4)
	case ebiten.KeyDown:
		y0, y1 = animation.y0, animation.y0+tileSize
	case ebiten.KeyLeft:
		y0, y1 = animation.y0+tileSize, animation.y0+tileSize*2
	case ebiten.KeyRight:
		y0, y1 = animation.y0+tileSize*2, animation.y0+tileSize*3
	default:
		x0, y0, x1, y1 = animation.x0, animation.y0, animation.x0+tileSize, animation.y0+tileSize
	}
	animation.currentFrameImage = animation.image.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image)
}

func (animation *Animation) SetDefaultFrame() {
	x0, y0, x1, y1 := animation.x0, animation.y0, animation.x0+animation.titleSize, animation.y0+animation.titleSize
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
