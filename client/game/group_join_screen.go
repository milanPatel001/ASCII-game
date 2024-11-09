package game

import (
	"ascii/utils"
	"fmt"
	"log"
)

type GroupJoinScreen struct {
	gameConfig     *GameConfig
	codeInput      []byte
	inputCursorPos int
	render         bool
}

func NewGroupJoinScreen(g *GameConfig) *GroupJoinScreen {
	return &GroupJoinScreen{gameConfig: g}
}

func (g *GroupJoinScreen) Init() {

	g.codeInput = []byte{}
	//g.inputCursorPos = 21
	g.render = true
}

func (g *GroupJoinScreen) Enter() {
	ClearScreen()
	MoveCursorToStartingPos()
	fmt.Print("\033[?25h")
	s := "Enter the room code: "
	fmt.Print(s)

	g.inputCursorPos = len(s) + 1

}

func (g *GroupJoinScreen) Exit() {
	g.codeInput = []byte{}
	fmt.Print("\033[?25l")
} // when screen becomes inactive

func (g *GroupJoinScreen) HandleInput(input byte) {
	switch rune(input) {
	case '\033':
		g.gameConfig.ScreenManager.ChangeScreen("start_menu")
	case '\r', '\n':
		code := string(g.codeInput)

		payload := fmt.Sprintf("%v | %v", code, g.gameConfig.PlayerId)
		// Send codeinput to server to get the roomdata
		pkt, err := utils.CreatePacketAndSerialize("127.0.0.1", utils.JOIN_GROUP, []byte(payload))
		if err != nil {
			log.Fatal(err)
		}

		_, err = g.gameConfig.conn.Write(pkt)
		if err != nil {
			log.Fatal(err)
		}

	case '\b', 127:
		if len(g.codeInput) > 0 {

			g.inputCursorPos--
			MoveCursor(1, g.inputCursorPos)
			fmt.Print(" ")

			// Remove last character
			g.codeInput = g.codeInput[:len(g.codeInput)-1]

		}
	default:
		g.codeInput = append(g.codeInput, input)
		fmt.Printf("%c", input)
		g.inputCursorPos++
	}

}

func (g *GroupJoinScreen) Render() {
	MoveCursor(1, g.inputCursorPos)
}

func (g *GroupJoinScreen) Update() {}

func (g *GroupJoinScreen) NeedsUpdate() bool {
	return g.render
}

func (g *GroupJoinScreen) HandleServerUpdate(packet utils.Packet) {

	if packet.MessageType == utils.NOT_FOUND {
		fmt.Println("No room found with this code...")
		return
	}

	if packet.MessageType != utils.JOIN_GROUP {
		return
	}

	// set the room data
	if err := utils.GetComplexPayloadFromBytes(packet.Payload, &g.gameConfig.Room); err != nil {
		log.Fatal(err)
	}

	g.gameConfig.ScreenManager.ChangeScreen("wait_menu")
}

func (s *GroupJoinScreen) DrawWindow() {
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
	//fmt.Print("\033[0m")

	// first column
	for i := 1; i < s.gameConfig.GameWindowHeight; i++ {
		MoveCursor(startingX+i, startingY)
		fmt.Print("\033[92m│")
	}
	fmt.Print("\033[0m")
}
