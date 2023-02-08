package player

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/VxVxN/game/internal/data"
)

type Player struct {
	Position base.Position
	image    *ebiten.Image
	gameData *data.GameData
}

func NewPlayer(position base.Position, imagePath string, gameData *data.GameData) (*Player, error) {
	image, _, err := ebitenutil.NewImageFromFile(imagePath)
	if err != nil {
		return nil, err
	}

	return &Player{
		Position: position,
		image:    image,
		gameData: gameData,
	}, nil
}

func (player *Player) Image() *ebiten.Image {
	return player.image
}

func (player *Player) Move(key ebiten.Key) {
	switch key {
	case ebiten.KeyUp:
		player.Position.Y--
	case ebiten.KeyDown:
		player.Position.Y++
	case ebiten.KeyLeft:
		player.Position.X--
	case ebiten.KeyRight:
		player.Position.X++
	default:
	}
}

func (player *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(player.gameData.TileSize*player.Position.X), float64(player.gameData.TileSize*player.Position.Y))
	screen.DrawImage(player.Image(), op)
}
