package game

import (
	"fmt"
)

type GameMenu interface {
	Init()
	DrawWindow(firstTime bool)
	GetName() string
	HandleInput(input byte)
}

type GameMenuManager struct {
	Menus            []GameMenu
	CurrentMenuIndex int
	startingX        int
	startingY        int
	gameWindowHeight int
	gameWindowWidth  int
}

type Inventory struct {
	GameMenuManager *GameMenuManager
	Width, Height   int
	items           []string
	Box             [][]rune
	Name            string
	ActiveOption    int
	PrevOption      int
}

// ******************************** INVENTORY METHODS ********************************

func (i *Inventory) HandleInput(input byte) {
	switch input {
	case 'i':
		i.HideWindow()
	case 'w':
		i.PrevOption = i.ActiveOption
		i.ActiveOption = (i.ActiveOption - 1 + len(i.items)) % len(i.items) // Move up
		i.DrawWindow(false)
	case 's':
		i.PrevOption = i.ActiveOption
		i.ActiveOption = (i.ActiveOption + 1) % len(i.items)
		i.DrawWindow(false)
	}
}

func (i *Inventory) HideWindow() {
	startingMenuX := i.GameMenuManager.startingX + i.GameMenuManager.gameWindowHeight - i.Height
	startingMenuY := i.GameMenuManager.startingY + i.GameMenuManager.gameWindowWidth - i.Width

	// displaying the box
	for j := 0; j < i.Height; j++ {
		for k := 0; k < i.Width; k++ {
			MoveCursor(startingMenuX+j, startingMenuY+k)
			fmt.Printf("%c", ' ')
		}
	}
}

func (i *Inventory) DrawWindow(firstTime bool) {

	if firstTime {
		i.DrawOuterBorder()
	}

	startingMenuX := i.GameMenuManager.startingX + i.GameMenuManager.gameWindowHeight - i.Height
	startingMenuY := i.GameMenuManager.startingY + i.GameMenuManager.gameWindowWidth - i.Width

	// displaying the box
	for j := 1; j < i.Height-1; j++ {
		if j == i.ActiveOption+2 {
			fmt.Print("\033[7m") // start reverse effect
		}
		for k := 1; k < i.Width-1; k++ {
			MoveCursor(startingMenuX+j, startingMenuY+k)
			fmt.Printf("%c", i.Box[j][k])
		}

		if j == i.ActiveOption+2 {
			fmt.Print("\033[0m")
		}
	}

}

func NewInventory(name string, gameMenuManager *GameMenuManager) *Inventory {
	width := 20
	height := 8

	box := make([][]rune, height)
	for i := 0; i < height; i++ {
		box[i] = make([]rune, width)
	}

	items := []string{"hello", "hi", "goodbyte", "ok"}

	inv := &Inventory{Width: width, Height: height, Name: "inv", Box: box, items: items, GameMenuManager: gameMenuManager, ActiveOption: 0}
	inv.Init()

	return inv
}

func (i *Inventory) Init() {

	//filling the array
	for j := 0; j < i.Height; j++ {
		for k := 0; k < i.Width; k++ {
			if j == 0 || j == i.Height-1 {
				i.Box[j][k] = '─'
			} else if k == 0 || k == i.Width-1 {
				i.Box[j][k] = '│'
			} else {
				i.Box[j][k] = ' '
			}
		}
	}

	i.Box[0][0] = '╭'
	i.Box[0][i.Width-1] = '╮'
	i.Box[i.Height-1][0] = '╰'
	i.Box[i.Height-1][i.Width-1] = '╯'

	// placing the title in array
	i.placeText(0, itemStartingCol(i.Width, i.Name), i.Name)

	// placing all the items in array
	for j, item := range i.items {
		i.placeText(j+2, itemStartingCol(i.Width, item), item)
	}
}

func (i *Inventory) placeText(row, col int, text string) {
	// if len(text)%2 != 0 {
	// 	col = col + 1
	// }

	for j, char := range text {
		i.Box[row][col+j] = char
	}
}

func (i *Inventory) GetName() string {
	return i.Name
}

func (i *Inventory) DrawOuterBorder() {

	startingMenuX := i.GameMenuManager.startingX + i.GameMenuManager.gameWindowHeight - i.Height
	startingMenuY := i.GameMenuManager.startingY + i.GameMenuManager.gameWindowWidth - i.Width

	for j := 0; j < i.Height; j++ {
		MoveCursor(startingMenuX+j, startingMenuY)
		fmt.Printf("%c", i.Box[j][0])
	}

	for j := i.Height - 1; j >= 0; j-- {
		MoveCursor(startingMenuX+j, startingMenuY+i.Width-1)
		fmt.Printf("%c", i.Box[j][i.Width-1])
	}

	for k := 0; k < i.Width; k++ {
		MoveCursor(startingMenuX, startingMenuY+k)
		fmt.Printf("%c", i.Box[0][k])
	}

	for k := i.Width - 1; k >= 0; k-- {
		MoveCursor(startingMenuX+i.Height-1, startingMenuY+k)
		fmt.Printf("%c", i.Box[i.Height-1][k])
	}

}

// ******************************** GAME MANAGER METHODS ********************************

func NewGameMenuManager(startingX int, startingY int, height int, width int) *GameMenuManager {
	return &GameMenuManager{Menus: make([]GameMenu, 0), CurrentMenuIndex: 0, startingX: startingX, startingY: startingY, gameWindowHeight: height, gameWindowWidth: width}
}

func (g *GameMenuManager) AddGameMenu(menu GameMenu) {
	g.Menus = append(g.Menus, menu)
}

func (g *GameMenuManager) ShowGameMenu(name string) {
	// could have used map for menu instead of menu arrays

	for i, menu := range g.Menus {
		if menu.GetName() == name {
			g.CurrentMenuIndex = i
			break
		}
	}

	g.Menus[g.CurrentMenuIndex].DrawWindow(true)
}

// func (g *GameMenuManager) Exit() {}

func (g *GameMenuManager) HandleInput(input byte) {
	g.Menus[g.CurrentMenuIndex].HandleInput(input)
}
