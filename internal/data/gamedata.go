package data

type GameData struct {
	ScreenWidth  int
	ScreenHeight int
	TileSize     int
}

func NewGameData() GameData {
	g := GameData{
		ScreenWidth:  30,
		ScreenHeight: 30,
		TileSize:     32,
	}
	return g
}
