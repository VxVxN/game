package entity

import (
	"fmt"
	"github.com/VxVxN/game/internal/base"
	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/pkg/animation"
	"github.com/VxVxN/game/pkg/item"
	"github.com/VxVxN/game/pkg/label"
	"github.com/VxVxN/game/pkg/quest"
	"github.com/VxVxN/game/pkg/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Player struct {
	BaseEntity
	satiety    int
	experience int
	coins      int
	face       font.Face
	cfg        *config.Config
	items      []*item.Item
	handItem   *item.Item
	quests     []*quest.Quest
}

func NewPlayer(position base.Position, speed float64, x0, y0 int, cfg *config.Config) (*Player, error) {
	animation, err := animation.NewAnimation(cfg.Player.ImagePath, x0, y0, cfg.Player.FrameCount, cfg.Common.TileSize)
	if err != nil {
		return nil, err
	}

	font, err := sfnt.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF font: %v", err)
	}

	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: 16,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new status player face(font): %v", err)
	}

	return &Player{
		BaseEntity: BaseEntity{
			position:  position,
			xp:        10000,
			animation: animation,
			speed:     speed,
		},
		satiety: 10000,
		face:    face,
		cfg:     cfg,
	}, nil
}

func (player *Player) Move(key ebiten.Key) {
	if player.IsDead() {
		return
	}
	switch key {
	case ebiten.KeyUp:
		player.SetY(player.Position().Y - player.speed)
	case ebiten.KeyDown:
		player.SetY(player.Position().Y + player.speed)
	case ebiten.KeyLeft:
		player.SetX(player.Position().X - player.speed)
	case ebiten.KeyRight:
		player.SetX(player.Position().X + player.speed)
	default:
	}
	player.animation.Update(key)
}

func (player *Player) Stand() {
	player.animation.SetDefaultFrame()
}

func (player *Player) Satiety() int {
	return player.satiety / 100
}

func (player *Player) AddExperience(value int) {
	player.experience += value
}

func (player *Player) Experience() int {
	return player.experience
}

func (player *Player) Update(position base.Position) {
	if player.IsDead() {
		return
	}

	if player.satiety > 0 {
		player.satiety--
	} else {
		player.xp--
	}

	for _, quest := range player.quests {
		quest.UpdateProgress(player.items)
	}

	if player.handItem != nil && player.handItem.ItemType == item.AxeType {
		player.attack = 40
	}
}

func (player *Player) AddCoins(coins int) {
	if player.IsDead() {
		return
	}
	player.coins += coins
}

func (player *Player) Draw(screen *ebiten.Image) {
	handItem := "nothing"
	if player.handItem != nil {
		handItem = player.handItem.ItemType.String()
	}
	text := fmt.Sprintf("XP: %d%%, Satiety: %d%%, Experience: %d, Coins: %d, Hand: %s", player.XP(), player.Satiety(), player.Experience(), player.coins, handItem)

	topLabel := label.NewLabel(player.face, 0, 10, text)
	topLabel.ContainerWidth = float64(player.cfg.Common.WindowWidth)
	topLabel.AlignVertical = label.AlignVerticalCenter
	topLabel.AlignHorizontal = label.AlignHorizontalCenter
	topLabel.Draw(screen)
}

func (player *Player) TakeItem(item *item.Item) {
	if player.IsDead() {
		return
	}
	player.items = append(player.items, item)
}

func (player *Player) TakeItemInHand(item *item.Item) {
	if player.IsDead() {
		return
	}
	player.handItem = item
}

func (player *Player) DeleteItem(ItemType item.ItemType) {
	if player.IsDead() {
		return
	}

	for i, item := range player.items {
		if item.ItemType == ItemType {
			player.items = utils.DeleteElemSlice(player.items, i)
			return
		}
	}
}

func (player *Player) Items() []*item.Item {
	return player.items
}

func (player *Player) TakeQuest(quest *quest.Quest) {
	if player.IsDead() {
		return
	}
	player.quests = append(player.quests, quest)
}

func (player *Player) Quests() []*quest.Quest {
	return player.quests
}
