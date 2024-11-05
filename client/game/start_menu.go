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
	fmt.Print("\033[H\033[?25l") // move the cursor to the starting pos and hide
	for i := 0; i < HEIGHT; i++ {
		fmt.Printf("\033[%d;1H", i+1) // move the cursor to the start of next row

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

func (g *StartMenuScreen) HandleServerUpdate(packet utils.Packet) {}

func itemStartingCol(WIDTH int, text string) int {
	return (WIDTH - len(text)) / 2
}
