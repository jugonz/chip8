package arch

/**
 * Datatype to describe the architecture of a simple emulator.
 */
type Arch interface {
	LoadGame(filepath string)
	EmulateCycle()
	DrawScreen()
	SetKeys()
	UpdateTimers()
	ShouldClose() bool
	Quit()
}
