package game

import (
	"ascii/utils"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"strconv"
)

type GroupCreationScreen struct {
	gameConfig     *GameConfig
	numInput       []byte
	pauseFlag      bool
	inputCursorPos int
	render         bool
}

func NewGroupCreationScreen(g *GameConfig) *GroupCreationScreen {
	return &GroupCreationScreen{gameConfig: g}
}

func (g *GroupCreationScreen) Init() {
	g.pauseFlag = false
	g.numInput = []byte{}
	//g.inputCursorPos = 59
	g.render = true
}

func (g *GroupCreationScreen) Enter() {
	ClearScreen()
	MoveCursorToStartingPos()
	fmt.Print("\033[?25h")
	s := "Enter the max numbers of players allowed (including you): "
	fmt.Print(s)

	g.inputCursorPos = len(s) + 1

}

func (g *GroupCreationScreen) Exit() {
	g.numInput = []byte{}
	fmt.Print("\033[?25l")
} // when screen becomes inactive

func (g *GroupCreationScreen) HandleInput(input byte) {
	switch rune(input) {
	case '\033':
		g.gameConfig.ScreenManager.ChangeScreen("start_menu")
	case '\r', '\n':
		num, err := strconv.Atoi(string(g.numInput))

		if err != nil || num > 20 || num < 1 {
			fmt.Println("\n( Enter a valid number from 1 to 20 !!! ) ")
			return
		}

		g.gameConfig.Room.TotalPlayers = num
		g.gameConfig.Room.Id = generateRoomId()
		g.gameConfig.Room.Code = generateCode()
		g.gameConfig.Room.PlayersJoined = append(g.gameConfig.Room.PlayersJoined, g.gameConfig.PlayerId)

		//send this room info to game server
		room := g.gameConfig.Room
		var buf bytes.Buffer

		enc := gob.NewEncoder(&buf)
		err = enc.Encode(room)
		if err != nil {
			log.Fatal("gob.Encode failed:", err)
		}

		roomPkt, err := utils.CreatePacketAndSerialize("127.0.0.1", utils.CREATE_GROUP, buf.Bytes())

		if err != nil {
			log.Fatal(err)
		}

		if _, err = g.gameConfig.conn.Write(roomPkt); err != nil {
			log.Fatal(err)
		}

		g.gameConfig.ScreenManager.ChangeScreen("wait_menu")

	case '\b', 127:
		if len(g.numInput) > 0 {

			g.inputCursorPos--
			MoveCursor(1, g.inputCursorPos)
			fmt.Print(" ")

			// Remove last character
			g.numInput = g.numInput[:len(g.numInput)-1]

		}
	default:
		g.numInput = append(g.numInput, input)
		fmt.Printf("%c", input)
		g.inputCursorPos++
	}

}

func (g *GroupCreationScreen) Render() {

	MoveCursor(1, g.inputCursorPos)

}

func (g *GroupCreationScreen) Update() {}

func (g *GroupCreationScreen) HandleServerUpdate(packet utils.Packet) {}

// TODO fix this
func generateRoomId() string {
	return "abc123"
}

func generateCode() string {
	return "abcdef"
}

func (g *GroupCreationScreen) NeedsUpdate() bool {
	return g.render
}
