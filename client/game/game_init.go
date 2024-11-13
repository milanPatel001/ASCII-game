package game

import (
	"ascii/utils"
	"time"
)

const (
	PLAYER_BASE = iota
	MIDDLE_GROUND
	NORMAL_TILE
	GRASS_TILE
	MUDDY_TILE
)

type Position struct {
	X, Y int
}

type ExitPosition struct {
	Pos      Position
	ExitSeed int
}

type Player struct {
	Pos      Position
	LastPos  Position
	Color    int
	Velocity int
	//arsenal []Weapon
	Name           string
	Id             string
	LastUpdateTime time.Time
	BaseSeed       int
	CurrSeed       int
}

type TerrainTile struct {
	Name                string
	TileType            int // enum of normal, muddy, icy, flamy, lava, grassy etc
	Symbol              rune
	TileSpeedMultiplier int
	CanHide             bool
	CanWalk             bool
}

// Single Screen
type Terrain struct {
	Seed  int
	Tiles [][]TerrainTile
	// Height     int
	// Width      int
	SpawnCoord     Position       // only used for player bases
	ExitCoord      []ExitPosition // player base will have only one exit coord
	FlagSpawnCoord Position
	FlagHoldCoord  Position
	TerrainType    int
}

type GameState struct {
	Players []Player
	//npcs         []GameObject
	innerWindowWidth  int
	innerWindowHeight int
	Terrains          []Terrain // multiple rooms with different sizes
}

func NewGameState(windowHeight, windowWidth int, room *Room, playerSeeds, middleGroundSeeds []int) *GameState {

	var players []Player
	colors := [12]int{31, 32, 33, 34, 35, 36, 91, 92, 93, 94, 95, 96}

	selectedTerrains := SelectTerrainsBasedOnSeeds(windowHeight, windowWidth, len(room.PlayersJoined), playerSeeds, middleGroundSeeds)

	for i, playerId := range room.PlayersJoined {

		players = append(players, Player{
			Pos:      selectedTerrains[i].SpawnCoord, // relative to terrain coords
			LastPos:  selectedTerrains[i].SpawnCoord,
			Color:    colors[i],
			Velocity: 1,
			Id:       playerId,
			BaseSeed: selectedTerrains[i].Seed,
			CurrSeed: selectedTerrains[i].Seed,
		})

	}

	gameState := &GameState{
		Players:           players,
		innerWindowWidth:  windowWidth,
		innerWindowHeight: windowHeight,
		Terrains:          selectedTerrains,
	}

	return gameState
}

func CreateTerrain(height, width int, terrainType int, tilesType int) Terrain {
	symbol := '.'

	if tilesType == GRASS_TILE {
		symbol = '-'
	} else if tilesType == MUDDY_TILE {
		symbol = ','
	}

	t := make([][]TerrainTile, height)

	for i := 0; i < height; i++ {
		t[i] = make([]TerrainTile, width)
		for j := 0; j < width; j++ {
			t[i][j] = CreateTerrainTile(tilesType, symbol, 1, true)
		}
	}

	return Terrain{
		Tiles:       t,
		TerrainType: terrainType,
	}
}

func CreateTerrainTile(tileType int, symbol rune, multiplier int, canWalk bool) TerrainTile {
	return TerrainTile{
		TileType:            tileType,
		Symbol:              symbol,
		TileSpeedMultiplier: multiplier,
		CanWalk:             canWalk,
	}
}

// **************** TERRAIN SEEDS **************** (for now only 3 created -> 2 for players and 1 middle ground)
/*

	# 2 to 4 players - 1 middle ground, 5 to 8 - 2 mg, 9 to 12 - 3mg, 13 to 16 - 4mg
	# Player base seeds : 0 to 15
	# Middle Grounds seeds : 16 to 19

	# Equal number of player base terrain seeds will be alloted to middle grounds. PlayNum/middleGroundsNum to all and  PlayNum%middleGroundsNum to some one by one
*/
func SelectTerrainsBasedOnSeeds(windowHeight, windowWidth int, playerNum int, playerSeeds, middleGroundSeeds []int) []Terrain {

	// first portion will be of player bases and rest will be of middle grounds
	var selectedTerrains []Terrain

	middleSeedIndex := 0 // goes till middleGroundSeeds' length - 1

	// player exit seeds for each middle ground
	var middleExitSeeds [][]int = make([][]int, len(middleGroundSeeds))

	// Creating player base terrains
	for _, pl := range playerSeeds {
		terrain := CreateTerrainUsingSeed(windowHeight, windowWidth, pl)
		terrain.Seed = pl

		tHeight := len(terrain.Tiles)
		tWidth := len(terrain.Tiles[0])

		mSeed := middleGroundSeeds[middleSeedIndex%len(middleGroundSeeds)]

		middleExitSeeds[middleSeedIndex%len(middleGroundSeeds)] = append(middleExitSeeds[middleSeedIndex%len(middleGroundSeeds)], pl)

		middleSeedIndex++

		pos := Position{tHeight / 2, tWidth - 1}
		terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{pos, mSeed})
		terrain.Tiles[pos.X][pos.Y].Symbol = '🡪'
		terrain.SpawnCoord = Position{tHeight / 2, 1}

		terrain.FlagHoldCoord = Position{tHeight / 2, 0}
		terrain.Tiles[terrain.FlagHoldCoord.X][terrain.FlagHoldCoord.Y].Symbol = 'O'

		terrain.FlagSpawnCoord = Position{(tHeight / 2) - 1, 0}
		terrain.Tiles[terrain.FlagSpawnCoord.X][terrain.FlagSpawnCoord.Y].Symbol = 'F'

		selectedTerrains = append(selectedTerrains, terrain)
	}

	// Creating middle ground terrains
	for i, m := range middleGroundSeeds {
		terrain := CreateTerrainUsingSeed(windowHeight, windowWidth, m)
		terrain.Seed = m

		tHeight := len(terrain.Tiles)
		tWidth := len(terrain.Tiles[0])

		terrain.ExitCoord = []ExitPosition{}

		playerExitSeeds := middleExitSeeds[i]

		// adding exit coords of playerbases to middlegrounds
		for j, p := range playerExitSeeds {
			var pos Position
			if j == 0 {
				pos = Position{0, 5}
				terrain.Tiles[pos.X][pos.Y].Symbol = '🡩'
			} else if j == 1 {
				pos = Position{0, len(terrain.Tiles[0]) - 6}
				terrain.Tiles[pos.X][pos.Y].Symbol = '🡩'
			} else if j == 2 {
				pos = Position{len(terrain.Tiles) - 1, 5}
				terrain.Tiles[pos.X][pos.Y].Symbol = '🡫'
			} else if j == 3 {
				pos = Position{len(terrain.Tiles) - 1, len(terrain.Tiles[0]) - 6}
				terrain.Tiles[pos.X][pos.Y].Symbol = '🡫'
			}

			terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{Pos: pos, ExitSeed: p})

		}

		// adding exit coords of middle grounds to middle grounds
		if len(middleGroundSeeds) > 1 {
			if i == 0 {
				terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{Pos: Position{tHeight / 2, tWidth - 1}, ExitSeed: middleGroundSeeds[i+1]})
			} else if i == len(middleGroundSeeds)-1 {
				terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{Pos: Position{tHeight / 2, 0}, ExitSeed: middleGroundSeeds[i-1]})
			} else {
				terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{Pos: Position{tHeight / 2, 0}, ExitSeed: middleGroundSeeds[i-1]})
				terrain.ExitCoord = append(terrain.ExitCoord, ExitPosition{Pos: Position{tHeight / 2, tWidth - 1}, ExitSeed: middleGroundSeeds[i+1]})
			}
		}

		selectedTerrains = append(selectedTerrains, terrain)
	}

	return selectedTerrains
}

func RandomSeedAssigner(playerNum int) ([]int, []int) {
	playerSeeds := []int{}
	middleGroundSeeds := []int{}

	middleGroundsLen := 1

	if playerNum > 12 {
		middleGroundsLen = 4
	} else if playerNum > 8 {
		middleGroundsLen = 3
	} else if playerNum > 4 {
		middleGroundsLen = 2
	}

	// Generating each section seed
	for i := 0; i < middleGroundsLen; i++ {
		randNum := utils.RandomNumberRange(16, 19)
		middleGroundSeeds = append(middleGroundSeeds, int(randNum))
	}

	for i := 0; i < playerNum; i++ {
		randNum := utils.RandomNumberRange(0, 15)
		playerSeeds = append(playerSeeds, int(randNum))
	}

	return playerSeeds, middleGroundSeeds

}

// TODO : will create proper ones later
func CreateTerrainUsingSeed(height, width, seed int) Terrain {

	if seed >= 16 {
		return CreateTerrain(height, width, MIDDLE_GROUND, GRASS_TILE)
	}

	tileType := NORMAL_TILE
	if seed%2 == 0 {
		tileType = MUDDY_TILE
	}

	return CreateTerrain(height, width, PLAYER_BASE, tileType)

}
