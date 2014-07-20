package arch

import (
	"fmt"
	"math"
)

/**
 * This file contains the implementations of Chip8 instructions.
 */

// Graphics controls

func (c8 *Chip8) ClearScreen() {
	for index := 0; index < len(c8.GFX); index++ {
		c8.GFX[index] = false // Clear the pixel on the screen.
	}

	c8.DrawFlag = true
	c8.PC += 2
}

func (c8 *Chip8) DrawSprite() {
	// All variables are promoted to uint16 for easier manipulation.
	xCoord := uint16(c8.Registers[c8.Opcode.Xreg])
	yCoord := uint16(c8.Registers[c8.Opcode.Yreg])
	height := c8.Opcode.Value & 0xF
	width := uint16(8)         // Width is hardcoded.
	shiftConst := uint16(0x80) // Shifting 128 right allows us to check indiv bits.
	yScale := uint16(64)       // Don't quite understand this at the moment.

	c8.Registers[0xF] = 0 // Assume we don't unset any pixels.

	var yLine, xLine uint16
	for yLine = 0; yLine < height; yLine++ {
		pixel := uint16(c8.Memory[c8.IndexReg+yLine])

		for xLine = 0; xLine < width; xLine++ {

			// If we need to draw this pixel...
			if pixel&(shiftConst>>xLine) != 0 {
				index := xCoord + xLine + (yScale * (yCoord + yLine))

				// XOR the pixel, saving whether we set it.
				if c8.GFX[index] {
					c8.Registers[0xF] = 1
					c8.GFX[index] = false
				} else {
					c8.GFX[index] = true
				}

			}
		}
	}

	c8.DrawFlag = true
	c8.PC += 2
}

func (c8 *Chip8) SetIndexToSprite() {
	char := c8.Registers[c8.Opcode.Xreg]
	offset := uint8(len(c8.Fontset) / 16) // Number of sprites per character.

	// Set index register to location of the
	// first fontset sprite of the matching character.
	c8.IndexReg = uint16(offset * char)

	c8.PC += 2
}

// Control flow

func (c8 *Chip8) CallRCA1802() {
	address := c8.Opcode.Value & 0xFFF
	panic(fmt.Sprintf("Unimplemented opcode CallRCA1802 called with address: %v\n",
		address))
}

func (c8 *Chip8) Return() {
	// No return values, just stack movement.
	c8.SP--
	c8.PC = c8.Stack[c8.SP]

	c8.PC += 2
}

func (c8 *Chip8) Jump() {
	c8.PC = c8.Opcode.Literal
}

func (c8 *Chip8) JumpIndexLiteralOffset() {
	newAddr := c8.Opcode.Literal + uint16(c8.Registers[0])

	c8.PC = newAddr
}

func (c8 *Chip8) Call() {
	// Store the PC in the stack pointer.
	c8.Stack[c8.SP] = c8.PC
	c8.SP++ // TODO: Overflow?

	c8.PC = c8.Opcode.Literal
}

func (c8 *Chip8) SkipInstrEqualLiteral() {
	literal := c8.Opcode.Value & 0xFF

	// If the register contents equal the literal...
	if uint16(c8.Registers[c8.Opcode.Xreg]) == literal {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrNotEqualLiteral() {
	literal := c8.Opcode.Value & 0xFF

	// If the register contents don't equal the literal...
	if uint16(c8.Registers[c8.Opcode.Xreg]) != literal {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrEqualReg() {
	// If the register contents are equal...
	if c8.Registers[c8.Opcode.Xreg] == c8.Registers[c8.Opcode.Yreg] {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrNotEqualReg() {
	// If the register contents are not equal...
	if c8.Registers[c8.Opcode.Xreg] != c8.Registers[c8.Opcode.Yreg] {
		c8.PC += 4 // skip an instructuion.
	} else {
		c8.PC += 2
	}
}

func (c8 *Chip8) SkipInstrKeyPressed() {
	if c8.Keyboard[c8.Opcode.Xreg] {
		c8.PC += 4
		return
	}

	c8.PC += 2 // If we didn't match, just move to the next instr.
}

func (c8 *Chip8) SkipInstrKeyNotPressed() {
	if !c8.Keyboard[c8.Opcode.Xreg] {
		c8.PC += 4
		return
	}

	c8.PC += 2
}

// Manipulating data registers

func (c8 *Chip8) SetRegToLiteral() {
	literal := c8.Opcode.Value & 0xFF

	// WARNING, MAY NOT FIT!
	c8.Registers[c8.Opcode.Xreg] = uint8(literal)

	c8.PC += 2
}

func (c8 *Chip8) SetRegToReg() {
	c8.Registers[c8.Opcode.Xreg] = c8.Registers[c8.Opcode.Yreg]

	c8.PC += 2
}

func (c8 *Chip8) Add() {
	literal := c8.Opcode.Value & 0xFF

	// WARNING, MIGHT NOT FIT
	c8.Registers[c8.Opcode.Xreg] += uint8(literal)

	c8.PC += 2
}

func (c8 *Chip8) AddWithCarry() {
	sum := int(c8.Registers[c8.Opcode.Xreg]) +
		int(c8.Registers[c8.Opcode.Yreg])
	c8.Registers[c8.Opcode.Xreg] = uint8(sum)

	if sum > math.MaxUint8 {
		c8.Registers[0xF] = 1 // If overflow, save 1 into last reg.
	} else {
		c8.Registers[0xF] = 0 // Else, save 0 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) Or() {
	c8.Registers[c8.Opcode.Xreg] =
		c8.Registers[c8.Opcode.Xreg] | c8.Registers[c8.Opcode.Yreg]
	c8.PC += 2
}

func (c8 *Chip8) And() {
	c8.Registers[c8.Opcode.Xreg] =
		c8.Registers[c8.Opcode.Xreg] & c8.Registers[c8.Opcode.Yreg]
	c8.PC += 2
}

func (c8 *Chip8) Xor() {
	c8.Registers[c8.Opcode.Xreg] =
		c8.Registers[c8.Opcode.Xreg] ^ c8.Registers[c8.Opcode.Yreg]
	c8.PC += 2
}

func (c8 *Chip8) SubXFromY() {
	diff := int(c8.Registers[c8.Opcode.Yreg]) -
		int(c8.Registers[c8.Opcode.Xreg])
	c8.Registers[c8.Opcode.Xreg] = uint8(diff)

	if diff < 0 {
		c8.Registers[0xF] = 0 // If underflow, save 0 into last reg.
	} else {
		c8.Registers[0xF] = 1 // Else, save 1 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) SubYFromX() {
	diff := int(c8.Registers[c8.Opcode.Xreg]) -
		int(c8.Registers[c8.Opcode.Yreg])
	c8.Registers[c8.Opcode.Xreg] = uint8(diff)

	if diff < 0 {
		c8.Registers[0xF] = 0 // If underflow, save 0 into last reg.
	} else {
		c8.Registers[0xF] = 1 // Else, save 1 into last reg.
	}

	c8.PC += 2
}

func (c8 *Chip8) ShiftRight() {
	// Set VF to least significant bit of Xreg before shifting.
	c8.Registers[0xF] = c8.Registers[c8.Opcode.Xreg] & 0x1

	c8.Registers[c8.Opcode.Xreg] = c8.Registers[c8.Opcode.Xreg] >> 1

	c8.PC += 2
}

func (c8 *Chip8) ShiftLeft() {
	// Set VF to most significant bit of Xreg before shifting.
	c8.Registers[0xF] = (c8.Registers[c8.Opcode.Xreg] >> 15) & 0x1

	c8.Registers[c8.Opcode.Xreg] = c8.Registers[c8.Opcode.Xreg] << 1

	c8.PC += 2
}

func (c8 *Chip8) SetRegisterRandomMask() {
	mask := uint8(c8.Opcode.Value & 0xFF)
	randNum := uint8(c8.Rando.Uint32() % 256) // Needs to fit in a uint8.

	c8.Registers[c8.Opcode.Xreg] = mask & randNum

	c8.PC += 2
}

func (c8 *Chip8) SaveBinaryCodedDecimal() {
	valueToConvert := c8.Registers[c8.Opcode.Xreg]

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
	for key := 0; key < len(c8.Keyboard); key++ {
		if c8.Keyboard[c8.Opcode.Xreg] {
			c8.Registers[c8.Opcode.Xreg] = uint8(key)
			c8.PC += 2
			return
		}
	}
	// Else, don't increment PC, wait another cycle for the key.
}

func (c8 *Chip8) GetDelayTimer() {
	c8.Registers[c8.Opcode.Xreg] = c8.DelayTimer // Save delay timer in reg.

	c8.PC += 2
}

// Manipulating special registers

func (c8 *Chip8) AddRegisterToIndex() {
	c8.IndexReg += uint16(c8.Registers[c8.Opcode.Xreg])

	c8.PC += 2
}

func (c8 *Chip8) SetIndexLiteral() {
	c8.IndexReg = c8.Opcode.Literal

	c8.PC += 2
}

func (c8 *Chip8) SetDelayTimer() {
	c8.DelayTimer = c8.Registers[c8.Opcode.Xreg]

	c8.PC += 2
}

func (c8 *Chip8) SetSoundTimer() {
	c8.SoundTimer = c8.Registers[c8.Opcode.Xreg]

	c8.PC += 2
}

// Context Switching

func (c8 *Chip8) SaveRegisters() {
	// Store all registers up to last register in memory,
	// starting in memory at the location in the index register.
	for loc, reg := c8.IndexReg, uint16(0); reg <= uint16(c8.Opcode.Xreg); loc, reg = loc+1, reg+1 {
		c8.Memory[loc] = c8.Registers[reg] // TODO: check overflow
	}

	c8.PC += 2
}

func (c8 *Chip8) RestoreRegisters() {
	// Load all registers up to last register from memory,
	// starting in memory at the location in the index register.
	for loc, reg := c8.IndexReg, uint16(0); reg <= uint16(c8.Opcode.Xreg); loc, reg = loc+1, reg+1 {
		c8.Registers[reg] = c8.Memory[loc] // TODO: check overflow
	}

	c8.PC += 2
}

// Special

func (c8 *Chip8) UnknownInstruction() {
	panic(fmt.Sprintf("Unknown instruction: %v\n", c8.Opcode))
}
