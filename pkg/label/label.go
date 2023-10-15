package label

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"math"
	"strings"
)

type AlignVertical uint8

const (
	AlignVerticalTop AlignVertical = iota
	AlignVerticalCenter
	AlignVerticalBottom
)

type AlignHorizontal uint8

const (
	AlignHorizontalLeft AlignHorizontal = iota
	AlignHorizontalCenter
	AlignHorizontalRight
)

type GrowVertical uint8

const (
	GrowVerticalDown GrowVertical = iota
	GrowVerticalUp
	GrowVerticalNone
)

type GrowHorizontal uint8

const (
	GrowHorizontalRight GrowHorizontal = iota
	GrowHorizontalLeft
	GrowHorizontalNone
)

type Label struct {
	x float64
	y float64

	ContainerWidth  float64
	ContainerHeight float64

	text string

	Color          color.RGBA
	ContainerColor color.RGBA

	AlignVertical   AlignVertical
	AlignHorizontal AlignHorizontal
	GrowVertical    GrowVertical
	GrowHorizontal  GrowHorizontal

	Visible bool

	fontFace   font.Face
	capHeight  float64
	lineHeight float64

	// lineSpacing line spacing increase factor
	// for 1.0 we do not over-wrap
	lineSpacing float64
}

func NewLabel(fontFace font.Face, x, y float64, text string) *Label {
	m := fontFace.Metrics()
	capHeight := math.Abs(float64(m.CapHeight.Floor()))
	lineHeight := float64(m.Height.Floor())
	return &Label{
		x:          x,
		y:          y,
		text:       text,
		fontFace:   fontFace,
		capHeight:  capHeight,
		lineHeight: lineHeight,
		Color:      color.RGBA{A: 0xff},
		Visible:    true,
	}
}

func (l *Label) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}
	posX := l.x
	posY := l.y + l.capHeight
	var (
		containerX0 float64
		containerY0 float64
		containerX1 float64
		containerY1 float64
	)
	bounds := text.BoundString(l.fontFace, l.text)
	boundsWidth := float64(bounds.Dx())
	boundsHeight := float64(bounds.Dy())
	if l.ContainerWidth == 0 && l.ContainerHeight == 0 {
		// Automatically assigning a work area
		containerX0 = posX
		containerY0 = posY
		containerX1 = posX + boundsWidth
		containerY1 = posY + boundsHeight
	} else {
		containerX0 = posX
		containerY0 = posY
		containerX1 = posX + l.ContainerWidth
		containerY1 = posY + l.ContainerHeight
		if delta := boundsWidth - l.ContainerWidth; delta > 0 {
			switch l.GrowHorizontal {
			case GrowHorizontalRight:
				containerX1 += delta
			case GrowHorizontalLeft:
				containerX0 -= delta
			case GrowHorizontalNone:
			}
		}
		if delta := boundsHeight - l.ContainerHeight; delta > 0 {
			switch l.GrowVertical {
			case GrowVerticalDown:
				containerY1 += delta
			case GrowVerticalUp:
				containerY0 -= delta
				posY -= delta
			case GrowVerticalNone:
			}
		}
	}
	var (
		containerWidth  = containerX1 - containerX0
		containerHeight = containerY1 - containerY0
	)
	if l.ContainerColor.A != 0 {
		x0 := containerX0
		y0 := containerY0 - l.capHeight
		w := containerWidth
		h := containerHeight
		ebitenutil.DrawRect(screen, x0, y0, w, h, l.ContainerColor)
	}
	numLines := strings.Count(l.text, "\n") + 1
	switch l.AlignVertical {
	case AlignVerticalTop:
	case AlignVerticalCenter:
		posY += (containerHeight - l.estimateHeight(numLines)) / 2
	case AlignVerticalBottom:
		posY += containerHeight - l.estimateHeight(numLines)
	}

	if l.text == "" {
		return
	}
	var opts ebiten.DrawImageOptions
	opts.ColorM.ScaleWithColor(l.Color)

	if l.AlignHorizontal == AlignHorizontalLeft {
		opts.GeoM.Translate(posX, posY)
		text.DrawWithOptions(screen, l.text, l.fontFace, &opts)
		return
	}
	if l.lineSpacing != 1 {
		h := float64(l.fontFace.Metrics().Height.Round()) * l.lineSpacing
		l.fontFace = text.FaceWithLineHeight(l.fontFace, math.Round(h))
	}
	// We need to process the text line by line, aligning each
	// a separate line
	textRemaining := l.text
	offsetY := 0.0
	for {
		nextLine := strings.IndexByte(textRemaining, '\n')
		lineText := textRemaining
		if nextLine != -1 {
			lineText = textRemaining[:nextLine]
			textRemaining = textRemaining[nextLine+len("\n"):]
		}
		lineBounds := text.BoundString(l.fontFace, lineText)
		lineBoundsWidth := float64(lineBounds.Dx())
		offsetX := 0.0
		switch l.AlignHorizontal {
		case AlignHorizontalCenter:
			offsetX = (containerWidth - lineBoundsWidth) / 2
		case AlignHorizontalRight:
			offsetX = containerWidth - lineBoundsWidth
		}
		opts.GeoM.Reset()
		opts.GeoM.Translate(posX+offsetX, posY+offsetY)
		text.DrawWithOptions(screen, lineText, l.fontFace, &opts)
		if nextLine == -1 {
			break
		}
		offsetY += l.lineHeight
	}
}

func (l *Label) estimateHeight(numLines int) float64 {
	// We start with the height we need for the first line
	estimatedHeight := l.capHeight
	if numLines >= 2 {
		// Add height for all other rows
		estimatedHeight += (float64(numLines) - 1) * l.lineHeight
	}
	return estimatedHeight
}
