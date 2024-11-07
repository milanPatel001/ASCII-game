package game

import (
	"ascii/utils"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/term"
)

type Room struct {
	Id            string
	Code          string
	TotalPlayers  int
	PlayersJoined []string
}

type GameConfig struct {
	conn          net.Conn
	Room          *Room
	ScreenManager *ScreenManager
	PlayerId      string
	TermChan      chan bool
}

func InitializeGame(conn net.Conn, playerId string) {
	StartGame(&GameConfig{conn: conn, Room: &Room{}, ScreenManager: NewScreenManager(), PlayerId: playerId})
}

func StartGame(gameConfig *GameConfig) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	gameConfig.TermChan = handleTermination(oldState)

	gameConfig.ScreenManager.AddScreen("start_menu", NewMenuScreen(gameConfig))
	gameConfig.ScreenManager.AddScreen("group_creation", NewGroupCreationScreen(gameConfig))
	gameConfig.ScreenManager.AddScreen("wait_menu", NewWaitingScreen(gameConfig))
	gameConfig.ScreenManager.AddScreen("group_join", NewGroupJoinScreen(gameConfig))
	gameConfig.ScreenManager.AddScreen("game", NewGameScreen(gameConfig))

	gameConfig.ScreenManager.ChangeScreen("start_menu")

	// Game loop
	inputChan := make(chan byte)
	serverUpdateChan := make(chan utils.Packet)

	updateTicker := time.NewTicker(time.Second / 30)
	defer updateTicker.Stop()

	// Input handling goroutine
	go getInput(inputChan)
	go updateFromServer(serverUpdateChan, gameConfig.conn)

	// Main game loop
	for {
		select {
		case packet := <-serverUpdateChan:

			if gameConfig.ScreenManager.ActiveScreen != nil {
				gameConfig.ScreenManager.ActiveScreen.HandleServerUpdate(packet)
			}

		case input := <-inputChan:

			if input == 'q' {
				gameConfig.TermChan <- true
				break
			}

			if gameConfig.ScreenManager.ActiveScreen != nil {
				gameConfig.ScreenManager.ActiveScreen.HandleInput(input)
			}

		case <-updateTicker.C:

			if gameConfig.ScreenManager.ActiveScreen != nil {
				gameConfig.ScreenManager.ActiveScreen.Update()

				if gameConfig.ScreenManager.ActiveScreen.NeedsUpdate() {
					gameConfig.ScreenManager.ActiveScreen.Render()
				}

			}

		}
	}

}

func updateFromServer(serverUpdateChan chan<- utils.Packet, conn net.Conn) {
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		packet, err := utils.Deserialize(buffer[:n])

		if err != nil {
			log.Println("Packet deserializing error:", err)
			continue
		}

		serverUpdateChan <- packet

	}
}

func getInput(inputChan chan byte) {
	buffer := make([]byte, 1)
	for {
		os.Stdin.Read(buffer)
		inputChan <- buffer[0]
	}
}

func handleTermination(oldState *term.State) chan bool {
	c := make(chan bool, 1)

	go func() {
		<-c
		term.Restore(int(os.Stdin.Fd()), oldState)
		os.Exit(1)
	}()

	return c
}
