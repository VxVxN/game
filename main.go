package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/VxVxN/game/internal/game"
)

func main() {
	game, err := game.NewGame()
	if err != nil {
		log.Fatalf("Failed to init game: %v", err)
	}
	ebiten.SetWindowTitle("Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Failed to run game: %v", err)
	}
}
