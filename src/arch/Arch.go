package arch

/**
 * Datatype to describe the architecture of a simple emulator.
 */
type Arch interface {
	LoadGame(filepath string)
	Run() // Returns when game or user quits.
	Quit()
}
