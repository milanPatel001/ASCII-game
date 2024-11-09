package game

import (
	"ascii/utils"
	"fmt"
)

const (
	WIDTH  = 32
	HEIGHT = 7
)

type StartMenuScreen struct {
	gameConfig                     *GameConfig
	items                          []string
	selectedOption, previousOption int
	box                            [HEIGHT][WIDTH]rune
	startingRow                    int
	boxStartingPos                 Position
	render                         bool
}

func NewMenuScreen(g *GameConfig) *StartMenuScreen {
	return &StartMenuScreen{
		gameConfig: g,
		items:      []string{"Create Group", "Join Group", "Exit"},
	}

}

func (s *StartMenuScreen) Init() {
	// Initializes box
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			if i == 0 || i == HEIGHT-1 || j == 0 || j == WIDTH-1 {
				s.box[i][j] = '#'
			} else {
				s.box[i][j] = ' '
			}
		}
	}

	s.boxStartingPos.X = s.gameConfig.StartingInnerWindowPos.X + ((s.gameConfig.GameWindowHeight - 4 - HEIGHT) / 2)
	s.boxStartingPos.Y = s.gameConfig.StartingInnerWindowPos.Y + ((s.gameConfig.GameWindowWidth - 4 - WIDTH) / 2)

	s.startingRow = (HEIGHT - len(s.items)) / 2

	// Place menu items or other elements inside the box
	for i, item := range s.items {
		s.placeText(s.startingRow+i, itemStartingCol(WIDTH, item), item)
	}

	s.render = true
}

func (s *StartMenuScreen) Enter() {
	ClearScreen()
	MoveCursorToStartingPos()
	s.DrawWindow()
}

func (s *StartMenuScreen) Exit() {}

func (s *StartMenuScreen) HandleInput(input byte) {
	switch rune(input) {
	case 'w':
		s.previousOption = s.selectedOption
		s.selectedOption = (s.selectedOption - 1 + 3) % 3 // Move up
	case 's':
		s.previousOption = s.selectedOption
		s.selectedOption = (s.selectedOption + 1) % 3
	case '\r': // Enter key
		if s.selectedOption == 0 {
			s.gameConfig.ScreenManager.ChangeScreen("group_creation")

		} else if s.selectedOption == 1 {
			s.gameConfig.ScreenManager.ChangeScreen("group_join")

		} else {
			s.gameConfig.TermChan <- true
		}

	}
}

func (s *StartMenuScreen) Render() {
	s.displayBox()
}

func (s *StartMenuScreen) Update() {}

// *** UTILITY FUNCTION ****

func (s *StartMenuScreen) displayBox() {
	fmt.Print("\033[?25l") // hide the cursor

	for i := 0; i < HEIGHT; i++ {
		MoveCursor(i+s.boxStartingPos.X, s.boxStartingPos.Y)

		for j := 0; j < WIDTH; j++ {
			if i == s.selectedOption+s.startingRow {
				fmt.Print("\033[7m") // start reverse effect
				fmt.Print(string(s.box[i][j]))
				fmt.Print("\033[0m")
			} else {
				fmt.Print(string(s.box[i][j]))
			}

		}
		fmt.Println()
	}
}

func (s *StartMenuScreen) placeText(row, col int, text string) {
	for i, char := range text {
		if col+i < WIDTH-1 {
			s.box[row][col+i] = char
		}
	}
}

func (s *StartMenuScreen) NeedsUpdate() bool {
	return s.render
}

func (s *StartMenuScreen) HandleServerUpdate(packet utils.Packet) {}

func (s *StartMenuScreen) DrawWindow() {
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
