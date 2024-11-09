package game

import (
	"ascii/utils"
)

type Screen interface {
	Init()  // init
	Enter() // when screen becomes active (can render or not)
	Exit()  // when screen becomes inactive
	HandleInput(input byte)
	Render()
	Update() // will handle update and all that stuff (no rendering stuff)
	NeedsUpdate() bool
	HandleServerUpdate(packet utils.Packet)
	DrawWindow() // will be used in Enter method to draw window
}

type ScreenManager struct {
	screens       map[string]Screen
	ActiveScreen  Screen
	nextScreen    Screen
	isTransiting  bool
	width, height int
}

func NewScreenManager() *ScreenManager {
	return &ScreenManager{
		screens: make(map[string]Screen),
		width:   0, // TODO
		height:  0,
	}
}

func (sm *ScreenManager) AddScreen(name string, screen Screen) {
	sm.screens[name] = screen
	screen.Init()
}

func (sm *ScreenManager) ChangeScreen(name string) {
	if screen, exists := sm.screens[name]; exists {
		if sm.ActiveScreen != nil {
			sm.ActiveScreen.Exit()
		}
		sm.ActiveScreen = screen
		sm.ActiveScreen.Enter()
	}
}
