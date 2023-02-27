package config

import (
	"encoding/json"
	"os"
)

type Mode string

const (
	ViewMode = "view"
)

type Config struct {
	Common struct {
		TileSize                   int    `json:"tileSize"`
		WindowTitle                string `json:"windowTitle"`
		WindowWidth                int    `json:"windowWidth"`
		WindowHeight               int    `json:"windowHeight"`
		RefreshRateFramesPerSecond int    `json:"refreshRateFramesPerSecond"`
		Mode                       Mode   `json:"mode"`
	} `json:"common"`
	Player struct {
		ImagePath  string `json:"imagePath"`
		FrameCount int    `json:"frameCount"`
	} `json:"player"`
	Map struct {
		Width       int    `json:"width"`
		Height      int    `json:"height"`
		TileSetPath string `json:"tileSetPath"`
	} `json:"map"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	if err = json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
