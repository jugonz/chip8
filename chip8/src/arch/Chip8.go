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
	c8.Opcode = c8.FetchOpcode() // Fetch instruction

	switch c8.Opcode >> 12 { // Decode (big-ass switch statement)
	case 0x0:
		switch c8.Opcode & 0xFF {
		case 0xE0:
			c8.ClearScreen()
		case 0xEE:
			c8.Return()
		default:
			c8.UnknownInstruction()
		}
	case 0x1:
		c8.Jump()
	case 0x2:
		c8.Call()
	case 0x3:
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
			c8.SetRegisterDelayTimer()
		case 0x0A:
			c8.SetRegisterKeyPress()
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

func (c8 *Chip8) UnknownInstruction() {
	panic(fmt.Sprintf("Unknown instruction: + ", c8.Opcode))
}

func (c8 *Chip8) SetIndexLiteral() {
	c8.IndexReg = c8.Opcode & 0x0FFF

	c8.PC += 2
}

func (c8 *Chip8) JumpIndexLiterallOffset() {
	// Store the PC in the stack pointer
	c8.Stack[c8.SP] = c8.PC
	c8.SP++ // Overflow?

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

func (c8 *Chip8) DrawSprite() {}

func (c8 *Chip8) SkipInstrKey() {
	key := (c8.Opcode >> 8) & 0xF

	switch c8.Opcode & 0xFF {
	case 0x9E: // Skip instr if key pressed
		if c8.Keyboard[key] {
			c8.PC += 4
			return
		}
	case 0xA1: // Skip instr if key not pressed
		if !c8.Keyboard[key] {
			c8.PC += 4
			return
		}
	default:
		c8.UnknownInstruction()
	}

	c8.PC += 2 // If we didn't match above, just move to next instr
}
