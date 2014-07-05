package arch

import (
	"math"
)

/**
 * This file contains the implementations of Chip8 instructions.
 */

// Graphics controls

func (c8 *Chip8) ClearScreen() {
	c8.Reigsters[0xF] = 0 // Assume we don't unset any pixels.

	for index := 0; index < len(c8.GFX); index++ {
		if c8.GFX[index] {
			c8.Reigsters[0xF] = 1 // Was previously drawn, note.
		}

		c8.GFX[index] = false // Clear the pixel on the screen.
	}

	c8.DrawFlag = true
	c8.PC += 2
}

func (c8 *Chip8) DrawSprite() {
	xCoord := c8.Reigsters[(c8.Opcode>>8)&0xF]
	yCoord := c8.Reigsters[(c8.Opcode>>4)&0xF]
	height := c8.Opcode & 0xF
	width := 8         // Width is hardcoded.
	shiftConst := 0x80 // Shifting 128 right allows us to check indiv bits.
	yScale := 64       // Don't quite understand this at the moment.

	c8.Reigsters[0xF] = 0 // Assume we don't unset any pixels.

	for yLine := 0; yLine < height; yLine++ {
		pixel := c8.Memory[c8.IndexReg+yLine]

		for xLine := 0; xLine < width; xLine++ {

			// If we need to draw this pixel...
			if pixel&(shiftConst>>xLine) != 0 {
				index := xCoord + xLine + (yScale * (yCoord + yLine))

				if c8.GFX[index] {
					c8.Reigsters[0xF] = 1
				}
				c8.GFX[index] ^= 1
			}
		}
	}

	c8.DrawFlag = true
	c8.PC += 2
}

func (c8 *Chip8) SetIndexToSprite() {
	char := c8.Reigsters[(c8.Opcode>>8)&0xF]

	// TODO: set index register to sprite of char

	c8.PC += 2
}

// Control flow

func (c8 *Chip8) CallRCA1802() {}

func (c8 *Chip8) Return() {
	// No return values, just stack movement.
	c8.SP--
	c8.PC = c8.Stack[c8.SP]

	c8.PC += 2
}

func (c8 *Chip8) Jump() {
	newAddr := c8.Opcode & 0xFFF

	c8.PC = newAddr
}

func (c8 *Chip8) JumpIndexLiteralOffset() {
	newAddr := (c8.Opcode & 0xFFF) + c8.Reigsters[0]

	c8.PC = newAddr
}

func (c8 *Chip8) Call() {
	newAddr := c8.Opcode & 0xFFF

	// Store the PC in the stack pointer.
	c8.Stack[c8.SP] = c8.PC
	c8.SP++ // TODO: Overflow?

	c8.PC = newAddr
}

func (c8 *Chip8) SkipInstrEqualLiteral() {
	sourceReg := (c8.Opcode >> 8) & 0xF
	literal := c8.Opcode & 0xFF

	// If the register contents equal the literal...
	if c8.Reigsters[sourceReg] == literal {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrNotEqualLiteral() {
	sourceReg := (c8.Opcode >> 8) & 0xF
	literal := c8.Opcode & 0xFF

	// If the register contents don't equal the literal...
	if c8.Reigsters[sourceReg] != literal {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrEqualReg() {
	source1 := (c8.Opcode >> 8) & 0xF
	source2 := (c8.Opcode >> 4) & 0xF

	// If the register contents are equal...
	if c8.Reigsters[source1] == c8.Reigsters[source2] {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrNotEqualReg() {
	source1 := (c8.Opcode >> 8) & 0xF
	source2 := (c8.Opcode >> 4) & 0xF

	// If the register contents are not equal...
	if c8.Reigsters[source1] != c8.Reigsters[source2] {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrKey() {
	key := (c8.Opcode >> 8) & 0xF

	switch c8.Opcode & 0xFF {
	case 0x9E: // Skip instr if key pressed.
		if c8.Keyboard[key] {
			c8.PC += 4
			return
		}
	case 0xA1: // Skip instr if key not pressed.
		if !c8.Keyboard[key] {
			c8.PC += 4
			return
		}
	default:
		c8.UnknownInstruction()
	}

	c8.PC += 2 // If we didn't match above, just move to next instr.
}

// Manipulating data registers

func (c8 *Chip8) SetRegToLiteral() {
	sourceReg := (c8.Opcode >> 8) & 0xF
	literal := c8.Opcode & 0xFF

	c8.Reigsters[sourceReg] = literal

	c8.PC += 2
}

func (c8 *Chip8) SetRegToReg() {
	sourceReg := (c8.Opcode >> 4) & 0xF
	destReg := (c8.Opcode >> 8) & 0xF

	c8.Reigsters[destReg] = c8.Reigsters[sourceReg]

	c8.PC += 2
}

func (c8 *Chip8) Add() {
	sourceReg := (c8.Opcode >> 8) & 0xF
	literal := c8.Opcode & 0xFF

	c8.Reigsters[sourceReg] += literal

	c8.PC += 2
}

func (c8 *Chip8) AddWithCarry() {
	regToAdd := (c8.Opcode >> 4) & 0xF
	destReg := (c8.Opcode >> 8) & 0xF

	sum := int(c8.Reigsters[destReg]) + int(c8.Reigsters[regToAdd])
	c8.Reigsters[destReg] = uint8(sum)

	if sum > math.MaxUint8 {
		c8.Reigsters[0xF] = 1 // If overflow, save 1 into last reg.
	} else {
		c8.Reigsters[0xF] = 0 // Else, save 0 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) Or() {
	sourceReg := (c8.Opcode >> 4) & 0xF
	destReg := (c8.Opcode >> 8) & 0xF

	c8.Reigsters[destReg] = c8.Reigsters[destReg] | c8.Reigsters[sourceReg]
	c8.PC += 2
}

func (c8 *Chip8) And() {
	sourceReg := (c8.Opcode >> 4) & 0xF
	destReg := (c8.Opcode >> 8) & 0xF

	c8.Reigsters[destReg] = c8.Reigsters[destReg] & c8.Reigsters[sourceReg]
	c8.PC += 2
}

func (c8 *Chip8) Xor() {
	sourceReg := (c8.Opcode >> 4) & 0xF
	destReg := (c8.Opcode >> 8) & 0xF

	c8.Reigsters[destReg] = c8.Reigsters[destReg] ^ c8.Reigsters[sourceReg]
	c8.PC += 2
}

func (c8 *Chip8) SubXFromY() {
	regY := (c8.Opcode >> 4) & 0xF
	regX := (c8.Opcode >> 8) & 0xF

	diff := int(c8.Reigsters[regY]) - int(c8.Reigsters[regX])
	c8.Reigsters[regX] = uint8(diff)

	if diff < 0 {
		c8.Registers[0xF] = 0 // If underflow, save 0 into last reg.
	} else {
		c8.Reigsters[0xF] = 1 // Else, save 1 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) SubYFromX() {
	regY := (c8.Opcode >> 4) & 0xF
	regX := (c8.Opcode >> 8) & 0xF

	diff := int(c8.Reigsters[regX]) - int(c8.Reigsters[regY])
	c8.Reigsters[regX] = uint8(diff)

	if diff < 0 {
		c8.Registers[0xF] = 0 // If underflow, save 0 into last reg.
	} else {
		c8.Reigsters[0xF] = 1 // Else, save 1 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) ShitftRight() {
	sourceReg := (c8.Opcode << 8) & 0xF

	// Set VF to least significant bit of sourceReg before shifting.
	c8.Reigsters[0xF] = c8.Reigsters[sourceReg] & 0x1

	c8.Registers[sourceReg] = c8.Registers[sourceReg] >> 1

	c8.PC += 2
}

func (c8 *Chip8) ShitftLeft() {
	sourceReg := (c8.Opcode << 8) & 0xF

	// Set VF to most significant bit of sourceReg before shifting.
	c8.Reigsters[0xF] = (c8.Reigsters[sourceReg] >> 15) & 0x1

	c8.Registers[sourceReg] = c8.Registers[sourceReg] << 1

	c8.PC += 2
}

func (c8 *Chip8) SetRegisterRandomMask() {
	targetReg := (c8.Opcode >> 8) & 0xF
	mask := c8.Opcode & 0xFF
	randNum := c8.Rando.Uint32() % 256 // Needs to fit in a uint8.

	c8.Reigsters[targetReg] = mask & randNum

	c8.PC += 2
}

func (c8 *Chip8) SaveBinaryCodedDecimal() {
	sourceReg := (c8.Opcode >> 8) & 0xF
	valueToConvert := c8.Reigsters[sourceReg]

	// Store the decimal representation of value in memory so that
	// the hundreths digit of the value is in Mem[Index],
	// the tenths digit is in Mem[Index+1], and
	// the ones digit is in Mem[Index+2].
	c8.Memory[c8.IndexReg] = valueToConvert / 100
	c8.Memory[c8.IndexReg+1] = (valueToConvert / 10) % 10
	c8.Memory[c8.IndexReg+2] = (valueToConvert % 100) % 10

	c8.PC += 2
}

func (c8 *Chip8) GetKeyPress() {
	targetReg := (c8.Opcode >> 8) & 0xF
	key := 0 // How to wait for key press?

	c8.Registers[targetReg] = key // TODO: save key press

	c8.PC += 2
}

func (c8 *Chip8) GetDelayTimer() {
	targetReg := (c8.Opcode >> 8) & 0xF

	c8.Registers[targetReg] = c8.DelayTimer // Save delay timer in reg.

	c8.PC += 2
}

// Manipulating special registers

func (c8 *Chip8) AddRegisterToIndex() {
	sourceReg := (c8.Opcode >> 8) & 0xF

	c8.IndexReg += c8.Reigsters[sourceReg]

	c8.PC += 2
}

func (c8 *Chip8) SetIndexLiteral() {
	c8.IndexReg = c8.Opcode & 0xFFF

	c8.PC += 2
}

func (c8 *Chip8) SetDelayTimer() {
	sourceReg := (c8.Opcode >> 8) & 0xF

	c8.DelayTimer = c8.Reigsters[sourceReg]

	c8.PC += 2
}

func (c8 *Chip8) SetSoundTimer() {
	sourceReg := (c8.Opcode >> 8) & 0xF

	c8.SoundTimer = c8.Reigsters[sourceReg]

	c8.PC += 2
}

// Context Switching

func (c8 *Chip8) SaveRegisters() {
	lastRegister := (c8.Opcode >> 8) & 0xF

	// Store all registers up to last register in memory,
	// starting in memory at the location in the index register.
	for loc, reg := c8.IndexReg, 0; reg <= lastRegister; loc, reg = loc+1, reg+1 {
		c8.Memory[loc] = c8.Reigsters[reg] // TODO: check overflow
	}

	c8.PC += 2
}

func (c8 *Chip8) RestoreRegisters() {
	lastRegister := (c8.Opcode >> 8) & 0xF

	// Load all registers up to last register from memory,
	// starting in memory at the location in the index register.
	for loc, reg := c8.IndexReg, 0; reg <= lastRegister; loc, reg = loc+1, reg+1 {
		c8.Reigsters[reg] = c8.Memory[loc] // TODO: check overflow
	}

	c8.PC += 2
}

// Special

func (c8 *Chip8) UnknownInstruction() {
	panic(fmt.Sprintf("Unknown instruction: + ", c8.Opcode))
}
