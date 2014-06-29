package arch

/**
 * Datatype to describe the architecture of a Chip8 system.
 */
type Chip8 struct {
	Opcode     uint16
	Memory     [4096]uint8
	Reigsters  [16]uint8
	IndexReg   uint8
	PC         uint8
	GFX        [64 * 32]bool // 1 if on
	DelayTimer uint8
	SoundTimer uint8
	Stack      [16]uint8
	SP         uint8
	Keyboard   [16]bool // 1 if pressed
}

func (c8 *Chip8) LoadGame() {

}

func (c8 *Chip8) EmulateCycle() {
	// Fetch
	newOp := c8.Memory[c8.PC]
	newOp << 8
	newOp |= c8.Memory[c8.PC+1]
	c8.Opcode = newOp

	// Decode (big-ass switch statement)

	// Execute

}

func (c8 *Chip8) DrawScreen() {

}

func (c8 *Chip8) SetKeys() {

}
