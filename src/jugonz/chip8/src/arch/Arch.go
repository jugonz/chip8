package arch

/**
 * Datatype to describe the architecture of a simple emulator.
 */
type Arch interface {
	LoadGame(filepath string)
	EmulateCycle()
	DrawScreen()
	SetKeys()
	ShouldDraw() bool // Return whether or not the screen must be drawn.
	ShouldClose() bool
}
