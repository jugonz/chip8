package arch

import (
	"math/rand"
	"time"
)

/**
 * Datatype to describe the architecture of a Chip8 system.
 */
type Chip8 struct {
	Opcode     uint16
	Memory     [4096]uint8
	Reigsters  [16]uint8
	IndexReg   uint16
	PC         uint16
	GFX        [64 * 32]bool // true if on
	DrawFlag   bool          // true if we just drew to the screen
	DelayTimer uint8
	SoundTimer uint8
	Stack      [16]uint16
	SP         uint16
	Keyboard   [16]bool  // true if pressed
	Rando      math.Rand // PRNG
	Fontset    [80]uint8
}

func MakeChip8() *Chip8 { // and initialize
	c8 = Chip8{}
	c8.Opcode = 0x200
	Rando = rand.New(time.Now().UnixNano())

	// Load fonset
	c8.Fontset = [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	return &c8
}

func (c8 *Chip8) LoadGame() {
	// open file and load into memory, else panic
}

func (c8 *Chip8) EmulateCycle() {
	c8.Opcode = c8.FetchOpcode() // Fetch instruction

	switch c8.Opcode >> 12 { // Decode (big-ass switch statement)
	case 0x0:
		switch c8.Opcode & 0xFF {
		case 0xE0:
			c8.ClearScreen()
		case 0xEE:
			c8.Return()
		default:
			c8.CallRCA1802()
		}
	case 0x1:
		c8.Jump()
	case 0x2:
		c8.Call()
	case 0x3 & 0x4:
		c8.SkipInstrEqualLiteral()
	case 0x4:
		c8.SkipInstrNotEqualLiteral()
	case 0x5:
		c8.SkipInstrEqualReg()
	case 0x6:
		c8.SetRegToLiteral()
	case 0x7:
		c8.Add()
	case 0x8:
		switch c8.Opcode & 0xF {
		case 0x0:
			c8.SetRegToReg()
		case 0x1:
			c8.Or()
		case 0x2:
			c8.And()
		case 0x3:
			c8.Xor()
		case 0x4:
			c8.AddWithCarry()
		case 0x5:
			c8.SubYFromX()
		case 0x6:
			c8.ShiftRight()
		case 0x7:
			c8.SubXFromY()
		case 0xE:
			c8.ShiftLeft()
		default:
			c8.UnknownInstruction()
		}
	case 0x9:
		c8.SkipInstrNotEqualReg()
	case 0xA: // 0xANNN: set index register to address NNN
		c8.SetIndexLiteral()
	case 0xB:
		c8.JumpIndexLiterallOffset()
	case 0xC:
		c8.SetRegisterRandomMask()
	case 0xD:
		c8.DrawSprite()
	case 0xE:
		c8.SkipInstrKey()
	case 0xF:
		switch c8.Opcode & 0xFF {
		case 0x07:
			c8.GetDelayTimer()
		case 0x0A:
			c8.GetKeyPress()
		case 0x15:
			c8.SetDelayTimer()
		case 0x18:
			c8.SetSoundTimer()
		case 0x1E:
			c8.AddRegisterToIndex()
		case 0x29:
			c8.SetIndexToSprite()
		case 0x33:
			c8.BinaryMagic()
		case 0x55:
			c8.SaveRegisters()
		case 0x65:
			c8.RestoreRegisters()
		default:
			c8.UnknownInstruction()
		}
	default:
		c8.UnknownInstruction()
	}

	// Execute

}

func (c8 *Chip8) FetchOpcode() uint16 {
	newOp := c8.Memory[c8.PC] << 8
	newOp |= c8.Memory[c8.PC+1]
	return newOp
}

func (c8 *Chip8) DrawScreen() {

}

func (c8 *Chip8) SetKeys() {

}
