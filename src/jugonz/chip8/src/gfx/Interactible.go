package gfx

type Interactible interface {
	SetKeys()
	KeyPressed(key uint8) bool // Return whether the key number has been pressed.
	ShouldClose() bool
	Quit()
}
