package camera

import (
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	_interface "github.com/VxVxN/game/pkg/entity/interface"
	"github.com/VxVxN/game/pkg/item"
	"github.com/hajimehoshi/ebiten/v2"
)

type Camera struct {
	cfg             *config.Config
	positionPlayer  base.Position
	playerImage     *ebiten.Image
	entities        []_interface.Entity
	backgroundImage *ebiten.Image
	frontImages     []*ebiten.Image
	zoom            float64
	items           []*item.Item
}

func NewCamera(cfg *config.Config) *Camera {
	return &Camera{
		cfg:  cfg,
		zoom: 1,
	}
}

func (camera *Camera) UpdatePlayer(position base.Position) {
	camera.positionPlayer = position
}

func (camera *Camera) UpdateEntities(entities []_interface.Entity) {
	camera.entities = entities
}

func (camera *Camera) Draw(screen *ebiten.Image) {
	var geoM ebiten.GeoM
	geoM.Scale(camera.zoom, camera.zoom)
	tileSize := float64(camera.cfg.Common.TileSize)
	windowWidth := float64(camera.cfg.Common.WindowWidth)
	windowHeight := float64(camera.cfg.Common.WindowHeight)

	op := &ebiten.DrawImageOptions{GeoM: geoM}
	op.GeoM.Translate(tileSize*-camera.positionPlayer.X+windowWidth/2,
		tileSize*-camera.positionPlayer.Y+windowHeight/2)
	screen.DrawImage(camera.backgroundImage, op)

	for _, entity := range camera.entities {
		if entity.IsDead() {
			continue
		}
		op = &ebiten.DrawImageOptions{GeoM: geoM}
		op.GeoM.Translate(tileSize*-camera.positionPlayer.X+tileSize*entity.Position().X+windowWidth/2,
			tileSize*-camera.positionPlayer.Y+tileSize*entity.Position().Y+windowHeight/2)
		screen.DrawImage(entity.Image(), op)
	}

	for _, gameItem := range camera.items {
		op = &ebiten.DrawImageOptions{GeoM: geoM}
		op.GeoM.Translate(tileSize*-camera.positionPlayer.X+tileSize*gameItem.Position().X+windowWidth/2,
			tileSize*-camera.positionPlayer.Y+tileSize*gameItem.Position().Y+windowHeight/2)
		gameItem.Draw(screen, op)
	}

	op = &ebiten.DrawImageOptions{GeoM: geoM}
	op.GeoM.Translate(float64(camera.cfg.Common.WindowWidth/2), float64(camera.cfg.Common.WindowHeight/2))
	screen.DrawImage(camera.playerImage, op)

	for _, frontImage := range camera.frontImages {
		op = &ebiten.DrawImageOptions{GeoM: geoM}
		op.GeoM.Translate(tileSize*-camera.positionPlayer.X+windowWidth/2,
			tileSize*-camera.positionPlayer.Y+windowHeight/2)
		screen.DrawImage(frontImage, op)
	}
}

func (camera *Camera) AddBackgroundImage(image *ebiten.Image) {
	camera.backgroundImage = image
}

func (camera *Camera) AddPlayerImage(image *ebiten.Image) {
	camera.playerImage = image
}

func (camera *Camera) SetItems(items []*item.Item) {
	camera.items = items
}

func (camera *Camera) AddFrontImages(images []*ebiten.Image) {
	camera.frontImages = images
}

func (camera *Camera) SetZoom(zoom float64) {
	camera.zoom = zoom
}
