package arch

/**
 * Datatype to describe the architecture of a simple emulator.
 */
type Arch interface {
	LoadGame()
	EmulateCycle()
	DrawScreen()
	SetKeys()
}
