package game

import (
	"fmt"
	"os"
)

type GameStartPayload struct {
	PlayerSeeds       []int
	MiddleGroundSeeds []int
}

type PlayerMovementPayload struct {
	CurrPos Position
	//LastPos  Position
	RoomCode string
	PlayerId string
}

type GameWindowInputPayload struct {
	windowHeight int
	windowWidth  int
	innerWindowX int
	innerWindowY int
}

func GetTerrainIndexUsingSeed(terrains []Terrain, seed int) int {

	for i, tr := range terrains {
		if tr.Seed == seed {
			return i
		}
	}

	return 0
}

func logAtTop(str any) {
	MoveCursor(1, 1)
	fmt.Print("\033[2K")
	fmt.Print(str)
}

func itemStartingCol(WIDTH int, text string) int {
	return (WIDTH - len(text)) / 2
}

func TextInput() string {
	var input []rune

	for {
		char := ReadRune()

		// Handle Enter key (newline)
		if char == '\r' || char == '\n' {
			break
		}

		// Handle Backspace (delete last character)
		if char == '\b' || char == 127 { // '\b' is backspace, 127 is delete
			if len(input) > 0 {
				// Remove last character
				input = input[:len(input)-1]
				// Move cursor back, overwrite character with space, move cursor back again
				fmt.Print("\b \b")
			}
			continue
		}

		// Append character to input and display it
		input = append(input, char)
		fmt.Print(string(char)) // Display the character immediately
	}

	return string(input)
}

func pressEnterToProceed() {
	for {
		char := ReadRune()

		if char == '\r' || char == '\n' {
			break
		}
	}
}

func ClearScreen() {
	fmt.Print("\033[2J")
}

func MoveCursorToStartingPos() {
	fmt.Print("\033[H")
}

func MoveCursor(x, y int) {
	fmt.Printf("\033[%d;%dH", x, y)
}

func HideCursor() {
	fmt.Print("\033[?25l")
}

func RestoreEverything() {
	fmt.Print("\033[?25h")
}

func ReadRune() rune {
	var buf = make([]byte, 1)
	os.Stdin.Read(buf) // Read a single byte
	return rune(buf[0])
}
