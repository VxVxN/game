package main

import (
	_ "image/png"
	"log"

	"github.com/VxVxN/game/internal/config"
	"github.com/VxVxN/game/internal/game"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	cfg, err := config.ParseConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	game, err := game.NewGame(cfg)
	if err != nil {
		log.Fatalf("Failed to init game: %v", err)
	}

	ebiten.SetWindowTitle(cfg.Common.WindowTitle)
	ebiten.SetWindowSize(cfg.Common.WindowWidth, cfg.Common.WindowHeight)

	if err = ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}
