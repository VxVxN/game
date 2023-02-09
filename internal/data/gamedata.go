package data

type GameData struct {
	ScreenWidthPx  int
	ScreenHeightPx int
	MapWidth       int
	MapHeight      int
	TileSize       int
}

func NewGameData() GameData {
	g := GameData{
		MapWidth:  30,
		MapHeight: 30,
		TileSize:  32,
	}
	return g
}
