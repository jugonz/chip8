package arch

import "math/rand"
import "time"

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
	Keyboard   [16]bool  // 1 if pressed
	Rando      math.Rand // PRNG
}

func MakeChip8() *Chip8 { // and initialize
	c8 = Chip8{}
	c8.Opcode = 0x200
	Rando = rand.New(time.Now().UnixNano())
	// Load fonset

	return &c8
}

func (c8 *Chip8) LoadGame() {
	// open file and load into memory, else panic
}

func (c8 *Chip8) EmulateCycle() {
	// Fetch
	c8.Opcode = c8.FetchOpcode()

	switch c8.Opcode >> 12 { // Decode (big-ass switch statement)
	case 0xA: // 0xANNN: set index register to address NNN
		c8.SetIndexLiteral()
	case 0xB:
		c8.JumpIndexLiterallOffset()
	case 0xC:
		c8.SetRegisterRandomMask()

	}

	// Execute

}

func (c8 *Chip8) FetchOpcode() uint16 {
	newOp := c8.Memory[c8.PC]
	newOp << 8
	newOp |= c8.Memory[c8.PC+1]
	return newOp
}

func (c8 *Chip8) DrawScreen() {

}

func (c8 *Chip8) SetKeys() {

}

// INSTRUCTION SET FUNCTIONS //
func (c8 *Chip8) SetIndexLiteral() {
	c8.IndexReg = c8.Opcode & 0x0FFF

	c8.PC += 2
}

func (c8 *Chip8) JumpIndexLiterallOffset() {
	newAddr := (c8.Opcode & 0x0FFF) + c8.Reigsters[0]

	c8.PC = newAddr
}

func (c8 *Chip8) SetRegisterRandomMask() {
	targetReg := (c8.Opcode >> 8) & 0xF
	mask := c8.Opcode & 0x00FF
	randNum := c8.Rando.Uint32() % 256 // needs to fit in a uint8
	c8.Reigsters[targetReg] = mask & randNum

	c8.PC += 2
}
