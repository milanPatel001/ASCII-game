package game

import (
	"ascii/utils"
	"fmt"
	"log"
)

type WaitingScreen struct {
	gameConfig      *GameConfig
	playerNumCol    int
	render          bool
	lastPlayerCount int
}

func NewWaitingScreen(g *GameConfig) *WaitingScreen {
	return &WaitingScreen{gameConfig: g, render: false}
}

func (w *WaitingScreen) Init() {
	w.lastPlayerCount = len(w.gameConfig.Room.PlayersJoined)
}

func (w *WaitingScreen) Enter() {
	ClearScreen()
	MoveCursorToStartingPos()
	fmt.Println("Current room id: ", w.gameConfig.Room.Id)
	fmt.Println("Invite code: ", w.gameConfig.Room.Code)
	fmt.Printf("Players Joined: %v/%v\n", len(w.gameConfig.Room.PlayersJoined), w.gameConfig.Room.TotalPlayers)
	fmt.Println("Waiting for the players...")
	fmt.Println()
	fmt.Println("Press Esc to Cancel this room....")

	w.playerNumCol = 17

}
func (w *WaitingScreen) Exit() {}

func (w *WaitingScreen) HandleInput(input byte) {
	switch rune(input) {
	case '\r':
		// only creator can start the game
		if len(w.gameConfig.Room.PlayersJoined) > 1 && w.gameConfig.Room.PlayersJoined[0] == w.gameConfig.PlayerId {
			// game screen
		}
	case '\033':

		packetType := utils.PLAYER_LEFT
		payload := fmt.Sprintf("%v | %v", w.gameConfig.Room.Code, w.gameConfig.PlayerId)

		if w.gameConfig.Room.PlayersJoined[0] == w.gameConfig.PlayerId {
			packetType = utils.DESTROY_ROOM
			payload = w.gameConfig.Room.Code
		}

		packet, err := utils.CreatePacketAndSerialize("127.0.0.1", uint8(packetType), []byte(payload))
		if err != nil {
			log.Fatal(err)
		}

		w.gameConfig.Room = &Room{}
		w.gameConfig.conn.Write(packet)
		w.gameConfig.ScreenManager.ChangeScreen("start_menu")
	}
}

func (w *WaitingScreen) Render() {
	if w.render {
		MoveCursor(3, w.playerNumCol)
		fmt.Printf("%v/%v\n", len(w.gameConfig.Room.PlayersJoined), w.gameConfig.Room.TotalPlayers)
		w.render = false
	}
}

func (w *WaitingScreen) Update() {

	if len(w.gameConfig.Room.PlayersJoined) != w.lastPlayerCount {
		w.lastPlayerCount = len(w.gameConfig.Room.PlayersJoined)
		w.render = true
	}
}

func (w *WaitingScreen) NeedsUpdate() bool {
	return w.render
}

func (g *WaitingScreen) HandleServerUpdate(packet utils.Packet) {
	if packet.MessageType == utils.BROADCAST {
		newPlayerJoinedId := string(packet.Payload)
		g.gameConfig.Room.PlayersJoined = append(g.gameConfig.Room.PlayersJoined, newPlayerJoinedId)
		return
	}

	if packet.MessageType == utils.DESTROY_ROOM {
		g.gameConfig.Room = &Room{}
		g.gameConfig.ScreenManager.ChangeScreen("start_menu")
		return
	}

	if packet.MessageType == utils.PLAYER_LEFT {
		playerLeftId := string(packet.Payload)

		i := 0
		for ind, id := range g.gameConfig.Room.PlayersJoined {
			if id == playerLeftId {
				i = ind
				break
			}
		}

		if i == len(g.gameConfig.Room.PlayersJoined)-1 {
			g.gameConfig.Room.PlayersJoined = g.gameConfig.Room.PlayersJoined[:i]
			return
		}

		g.gameConfig.Room.PlayersJoined = append(g.gameConfig.Room.PlayersJoined[:i], g.gameConfig.Room.PlayersJoined[i+1:]...)
	}

}
