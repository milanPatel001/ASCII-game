package game

import "time"

type Position struct {
	X, Y int
}

type Player struct {
	Pos      Position
	LastPos  Position
	Symbol   rune
	Velocity int
	//arsenal []Weapon
	Name           string
	Id             string
	LastUpdateTime time.Time
}

type TerrainTile struct {
	Name                string
	TileType            string // enum of normal, muddy, icy, flamy, lava, grassy etc
	Symbol              rune
	TileSpeedMultiplier int
	CanHide             bool
	CanWalk             bool
}

type Terrain struct {
	Tiles [][]TerrainTile
	// Height     int
	// Width      int
	SpawnCoord Position
	ExitCoord  []Position
}

type GameState struct {
	RoomCode string
	Players  []Player
	//npcs         []GameObject
	innerWindowWidth  int
	innerWindowHeight int
	Terrains          []Terrain // multiple rooms with different sizes
}

func NewGameState(windowHeight, windowWidth int, room *Room) *GameState {
	terrain := CreateTerrain(windowHeight, windowWidth)

	var players []Player
	symbols := []rune{'P', 'Q', 'R', 'X'}

	for i, playerId := range room.PlayersJoined {

		players = append(players, Player{
			Pos:      Position{2 * i, 0}, // relative to terrain coords
			LastPos:  Position{2 * i, 0},
			Symbol:   symbols[i],
			Velocity: 1,
			Id:       playerId,
		})

	}

	gameState := &GameState{
		RoomCode:          room.Code,
		Players:           players,
		innerWindowWidth:  windowWidth,
		innerWindowHeight: windowHeight,
		Terrains:          []Terrain{terrain},
	}

	return gameState
}

func CreateTerrain(height, width int) Terrain {
	t := make([][]TerrainTile, height)

	for i := 0; i < height; i++ {
		t[i] = make([]TerrainTile, width)
		for j := 0; j < width; j++ {
			t[i][j] = CreateTerrainTile()
		}
	}

	return Terrain{
		Tiles:      t,
		SpawnCoord: Position{0, 0},
	}
}

func CreateTerrainTile() TerrainTile {
	return TerrainTile{
		TileType:            "normal",
		Symbol:              '.',
		TileSpeedMultiplier: 1,
		CanWalk:             true,
	}
}
