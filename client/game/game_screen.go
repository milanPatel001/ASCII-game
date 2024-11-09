package game

import (
	"ascii/utils"
	"fmt"
	"math"
	"time"
)

type GameScreen struct {
	gameConfig             *GameConfig
	GameState              *GameState
	GameMenuManager        *GameMenuManager
	render                 bool
	CurrentPlayer          *Player
	currentTerrain         *Terrain
	lastMoveTime           time.Time
	moveDelay              time.Duration
	lastMovementPacketTime time.Time
	movementPacketDelay    time.Duration
	interpolationDuration  time.Duration
	isMenuOpen             bool
	menuChan               chan byte
}

func (g *GameScreen) Init() {
	g.render = true
	g.lastMoveTime = time.Now()
	g.lastMovementPacketTime = time.Now()
	g.moveDelay = 3 * g.gameConfig.TickRate
	g.movementPacketDelay = 3 * g.gameConfig.TickRate // max 8 packets/sec/player
	g.interpolationDuration = 2 * g.gameConfig.TickRate
	g.menuChan = make(chan byte, 1)
}

func (g *GameScreen) Enter() {

	if g.GameState == nil {
		g.GameState = NewGameState(g.gameConfig.GameWindowHeight-4, g.gameConfig.GameWindowWidth-4, g.gameConfig.Room)

		for i := range g.GameState.Players {
			if g.GameState.Players[i].Id == g.gameConfig.PlayerId {
				g.CurrentPlayer = &g.GameState.Players[i]
				break
			}
		}

		g.currentTerrain = &g.GameState.Terrains[0]
	}

	ClearScreen()
	MoveCursorToStartingPos()

	g.DrawWindow()

	g.DrawTerrainAndPlayers()

	g.CreateGameMenuManager()
}

func (g *GameScreen) Exit() {}

func (g *GameScreen) HandleInput(input byte) {

	if time.Since(g.lastMoveTime) < g.moveDelay {
		return
	}

	if g.isMenuOpen {
		g.menuChan <- input
		return
	}

	newX := g.CurrentPlayer.Pos.X
	newY := g.CurrentPlayer.Pos.Y
	moved := false

	switch input {
	case 'q':
		//g.isRunning = false
	case 'w':
		if newX > 0 {
			newX--
			moved = true
		}
	case 's':
		newX++
		moved = true

	case 'a':
		if newY > 0 {
			newY--
			moved = true
		}
	case 'd':
		newY++
		moved = true
	case 'i':
		if !g.isMenuOpen {
			g.isMenuOpen = true
			g.GameMenuManager.ShowGameMenu("inv")
		}
	}

	if !moved || !g.isValidMove(newX, newY) {
		return
	}

	g.lastMoveTime = time.Now()
	g.CurrentPlayer.LastPos = g.CurrentPlayer.Pos
	g.CurrentPlayer.Pos.X = newX
	g.CurrentPlayer.Pos.Y = newY

	payload := PlayerMovementPayload{
		CurrPos:  g.CurrentPlayer.Pos,
		RoomCode: g.GameState.RoomCode,
		PlayerId: g.CurrentPlayer.Id,
	}

	if time.Since(g.lastMovementPacketTime) < g.movementPacketDelay {
		return
	}

	b, _ := utils.ConvComplexPayloadToBytes(payload)
	pkt, _ := utils.CreatePacketAndSerialize("127.0.0.1", utils.PLAYER_MOVE, b)

	g.gameConfig.conn.Write(pkt)
	g.lastMovementPacketTime = time.Now()
}

func (g *GameScreen) Render() {

	if g.CurrentPlayer.LastPos.X == g.CurrentPlayer.Pos.X && g.CurrentPlayer.LastPos.Y == g.CurrentPlayer.Pos.Y {
		return
	}

	startingX := g.gameConfig.StartingInnerWindowPos.X
	startingY := g.gameConfig.StartingInnerWindowPos.Y

	// Restore terrain at old position
	MoveCursor(g.CurrentPlayer.LastPos.X+startingX, g.CurrentPlayer.LastPos.Y+startingY)
	fmt.Printf("%c", g.currentTerrain.Tiles[g.CurrentPlayer.LastPos.X][g.CurrentPlayer.LastPos.Y].Symbol)

	// Draw at new position
	MoveCursor(g.CurrentPlayer.Pos.X+startingX, g.CurrentPlayer.Pos.Y+startingY)
	fmt.Printf("%c", g.CurrentPlayer.Symbol)

	for i, pl := range g.GameState.Players {
		if pl.Id == g.CurrentPlayer.Id {
			continue
		}

		if (pl.LastPos.X == -1 && pl.LastPos.Y == -1) || (pl.LastPos.X == pl.Pos.X && pl.LastPos.Y == pl.Pos.Y) {
			continue
		}

		newX, newY := g.renderInterpolatedPlayer(&pl)
		g.GameState.Players[i].LastPos.X = newX
		g.GameState.Players[i].LastPos.Y = newY

		// MoveCursor(pl.LastPos.X, pl.LastPos.Y)
		// fmt.Printf("%c", g.currentTerrain.Tiles[pl.LastPos.X][pl.LastPos.Y].Symbol)

		// MoveCursor(pl.Pos.X, pl.Pos.Y)
		// fmt.Printf("%c", pl.Symbol)
	}

}

func (g *GameScreen) Update() {

}

func (g *GameScreen) NeedsUpdate() bool {
	return g.render
}

func (g *GameScreen) HandleServerUpdate(packet utils.Packet) {
	if packet.MessageType == utils.PLAYER_MOVE {
		var playerMovementPayload PlayerMovementPayload

		err := utils.GetComplexPayloadFromBytes(packet.Payload, &playerMovementPayload)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, pl := range g.GameState.Players {
			if pl.Id == playerMovementPayload.PlayerId {
				g.GameState.Players[i].LastPos = g.GameState.Players[i].Pos
				g.GameState.Players[i].Pos = playerMovementPayload.CurrPos
				g.GameState.Players[i].LastUpdateTime = time.Now()
				break
			}
		}
	}
}

func (g *GameScreen) renderInterpolatedPlayer(player *Player) (int, int) {
	// Calculate the interpolation progress based on the time since the last update

	startingX := g.gameConfig.StartingInnerWindowPos.X
	startingY := g.gameConfig.StartingInnerWindowPos.Y

	progress := float64(time.Since(player.LastUpdateTime)) / float64(g.interpolationDuration)
	if progress > 1.0 {
		progress = 1.0
	}

	// Interpolate the player's position
	newX := int(math.Floor(float64(player.LastPos.X) + float64(player.Pos.X-player.LastPos.X)*progress))
	newY := int(math.Floor(float64(player.LastPos.Y) + float64(player.Pos.Y-player.LastPos.Y)*progress))

	MoveCursor(player.LastPos.X+startingX, player.LastPos.Y+startingY)
	fmt.Printf("%c", g.currentTerrain.Tiles[player.LastPos.X][player.LastPos.Y].Symbol)

	MoveCursor(newX+startingX, newY+startingY)
	fmt.Printf("%c", player.Symbol)

	return newX, newY
}

func (g *GameScreen) isValidMove(newX, newY int) bool {
	if g.currentTerrain == nil {
		return false
	}

	// Check boundaries
	if newX < 0 || newX >= len(g.currentTerrain.Tiles) ||
		newY < 0 || newY >= len(g.currentTerrain.Tiles[0]) {
		return false
	}

	// Check collision with other players
	for _, player := range g.GameState.Players {
		if player.Id == g.CurrentPlayer.Id {
			continue
		}

		if player.Pos.X == newX && player.Pos.Y == newY {
			return false
		}
	}

	return true
}

func NewGameScreen(g *GameConfig) *GameScreen {
	return &GameScreen{gameConfig: g}
}

func (g *GameScreen) CreateGameMenuManager() {
	g.GameMenuManager = NewGameMenuManager(g.gameConfig.StartingGameWindowPos.X, g.gameConfig.StartingGameWindowPos.Y, g.gameConfig.GameWindowHeight, g.gameConfig.GameWindowWidth)
	g.GameMenuManager.AddGameMenu(NewInventory("inv", g.GameMenuManager))

	go func() {
		for {
			select {
			case input := <-g.menuChan:

				if g.isMenuOpen {
					g.GameMenuManager.HandleInput(input)

					if input == 'i' {
						g.isMenuOpen = false
						g.DrawTerrainAndPlayers()
					}
				}

			}
		}
	}()

}

func (s *GameScreen) DrawTerrainAndPlayers() {
	startingX := s.gameConfig.StartingInnerWindowPos.X
	startingY := s.gameConfig.StartingInnerWindowPos.Y

	for i := 0; i < len(s.currentTerrain.Tiles); i++ {
		for j := 0; j < len(s.currentTerrain.Tiles[i]); j++ {
			MoveCursor(startingX+i, startingY+j)
			fmt.Print(string(s.currentTerrain.Tiles[i][j].Symbol))
		}
	}

	for _, pl := range s.GameState.Players {
		MoveCursor(startingX+pl.Pos.X, startingY+pl.Pos.Y)
		fmt.Print(string(pl.Symbol))
	}
}

func (s *GameScreen) DrawWindow() {
	lineChar := "─"

	startingX := s.gameConfig.StartingGameWindowPos.X
	startingY := s.gameConfig.StartingGameWindowPos.Y

	MoveCursor(startingX, startingY)

	// First row
	for i := 0; i < s.gameConfig.GameWindowWidth; i++ {
		fmt.Printf("\033[92m%v", lineChar)
	}

	// left corner
	MoveCursor(startingX, startingY)
	fmt.Print("╭")

	// first column
	for i := 1; i < s.gameConfig.GameWindowHeight; i++ {
		MoveCursor(startingX+i, startingY)
		fmt.Print("\033[92m│")
	}

	fmt.Print("\033[1B")
	fmt.Print("\033[1D")

	// Lower Left Corner
	fmt.Print("╰")

	// last row
	for i := 0; i < s.gameConfig.GameWindowWidth; i++ {
		fmt.Printf("\033[92m%v", lineChar)
	}

	// Lower right corner
	MoveCursor(startingX+s.gameConfig.GameWindowHeight, startingY+s.gameConfig.GameWindowWidth)
	fmt.Print("╯")

	// last column
	for i := s.gameConfig.GameWindowHeight - 1; i >= 0; i-- {
		MoveCursor(startingX+i, startingY+s.gameConfig.GameWindowWidth)
		fmt.Print("\033[92m│")
	}

	// Upper right corner
	fmt.Print("\033[1D")
	fmt.Print("╮")

	fmt.Print("\033[0m")
}
