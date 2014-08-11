package arch

import (
	"jugonz/chip8/src/gfx"
	"testing"
)

func TestSetup(t *testing.T) {
	c8 := MakeChip8(false)
	if c8.PC != 0x200 {
		t.Errorf("c8 Opcode was not initialized properly! Was: %v\n",
			c8.Opcode)
	}

	// Check fontset.
	for i := 0; i < 80; i++ {
		if c8.Fontset[i] == 0x00 {
			t.Errorf("c8 Fontset was not loaded!\n")
		}
	}

	c8.LoadGame("../../c8games/PONG2")
	// Check that byte 1 is 22 and last byte is EE.
	if c8.Memory[0x200] != 0x22 {
		t.Errorf(
			"c8 PONG2 was not properly loaded! Expected byte 1 to be 0x22, was: %v\n",
			c8.Memory[0x200])
	} else if c8.Memory[0x307] != 0xEE {
		t.Errorf(
			"c8 PONG2 was not properly loaded! Expected last byte to be 0xEE, was: %v\n",
			c8.Memory[0x307])
	}
}

func TestSkipInstr(t *testing.T) {
	c8 := MakeChip8(false)

	// First, add the literal (A3) to a register.
	c8.Opcode = MakeOpcode(0x71A3)
	c8.DecodeExecute()

	// Now, check the PC.
	if c8.PC != 0x202 {
		t.Error("PC not properly updated!\n")
	}

	// Now, check that an instruction is skipped when
	// comparing the literal.
	c8.Opcode = MakeOpcode(0x31A3)
	c8.DecodeExecute()

	if c8.PC != 0x206 {
		t.Error("PC did not skip an instruction after a literal compare!\n")
	}

	// Now, check that an instruction is NOT skipped
	// when comparing the same literal.
	c8.Opcode = MakeOpcode(0x41A3)
	c8.DecodeExecute()

	if c8.PC != 0x208 {
		t.Error("PC incorrectly skipped an instruction after a literal compare!\n")
	}

	// Now, check that comparing two identical registers
	// leads to an instruction skip.
	c8.Opcode = MakeOpcode(0x72A3) // Add literal to another register.
	c8.DecodeExecute()

	c8.Opcode = MakeOpcode(0x5120)
	c8.DecodeExecute()

	if c8.PC != 0x20E {
		t.Error("PC did not skip an instruction on register compare!\n")
	}
}

func TestClearScreen(t *testing.T) {
	c8 := MakeChip8(false)
	screen := c8.Screen.(*gfx.Screen)

	// Draw something to the screen, and see that it is not empty.
	c8.Opcode = MakeOpcode(0xD324)
	c8.DecodeExecute()

	clear := true
	for index, _ := range screen.Pixels {
		for _, value := range screen.Pixels[index] {
			if value {
				clear = false
				break
			}
		}
	}
	if clear {
		t.Errorf("DrawSprite failed to draw the screen!\n")
	}

	// Now, clear the screen, and check that it is empty.
	c8.Opcode = MakeOpcode(0x00E0)
	c8.DecodeExecute()

	clear = true
	for index, _ := range screen.Pixels {
		for _, value := range screen.Pixels[index] {
			if value {
				clear = false
				break
			}
		}
	}
	if !clear {
		t.Errorf("ClearScreen failed to clear the screen!\n")
	}
}

func TestDrawSprite(t *testing.T) {
	// Do this one later. Phew.
}

func TestSetIndexToSprite(t *testing.T) {
	// Do this one later too.
}

func TestCallReturn(t *testing.T) {
	c8 := MakeChip8(false)

	// Check initial stack.
	if c8.SP != 0 {
		t.Errorf("Stack not properly initalized!\n")
	}
	for index, val := range c8.Stack {
		if val != 0 {
			t.Errorf("Stack not empty! Val at index %v is %X\n",
				index, val)
		}
	}

	// Call a program at 789, and check stack.
	c8.Opcode = MakeOpcode(0x2789)
	c8.DecodeExecute()

	if c8.SP != 1 {
		t.Errorf("Stack not updated after call!\n")
	} else if c8.Stack[0] != 0x200 { // Starting PC
		t.Errorf("Old PC was not saved!\n")
	}
	for index, val := range c8.Stack[1:] {
		if val != 0 {
			t.Errorf("Stack not empty! Val at index %v is %X\n",
				index, val)
		}
	}

	// Return from that program, and see where we are.
	c8.Opcode = MakeOpcode(0x00EE)
	c8.DecodeExecute()

	if c8.SP != 0 {
		t.Errorf("Stack not updated after return!")
	}
	// Return doesn't erase stack on decrement.
	for index, val := range c8.Stack[1:] {
		if val != 0 {
			t.Errorf("Stack not empty! Val at index %v is %X\n",
				index, val)
		}
	}
}

func TestAdd(t *testing.T) {
	c8 := MakeChip8(false)

	// Test a simple add from 0 to register 2.
	if c8.Registers[2] != 0 {
		t.Errorf("Register did not start at 0!\n")
	}

	c8.Opcode = MakeOpcode(0x7212)
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[2] != 0x12 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	}

	// Test the max value in a different register.
	c8.Opcode = MakeOpcode(0x73FF)
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0xFF {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}

	// Now, add one, and test overflow.
	c8.Opcode = MakeOpcode(0x7301)
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0x00 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}
}

func TestAddWithCarry(t *testing.T) {
	c8 := MakeChip8(false)

	// Test adding the max value without overflow.
	c8.Opcode = MakeOpcode(0x73FF) // Add FF to reg 3 (0).
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x8374) // Add reg 7 (0) to reg 3.
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0xFF {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	} else if c8.Registers[0xF] != 0 {
		t.Error("Carry register was not 0 on non-overflowing calculation!\n")
	}

	// Now, add one, and test overflow.
	c8.Opcode = MakeOpcode(0x7401) // Add 1 to reg 4 (0).
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x8344) // Add reg 4 (1) to reg 3 (FF).
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0x00 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}
	// Test that the overflow register is correctly set.
	if c8.Registers[0xF] != 1 {
		t.Error("Carry register was not 1 on overflowing calculation!\n")
	}
}

func TestSub(t *testing.T) {
	c8 := MakeChip8(false)

	// Set initial register values.
	c8.Opcode = MakeOpcode(0x71A2) // Add A2 to reg 1 (0).
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x7203) // Add 03 to reg 2 (0).
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[1] != 0xA2 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[2] != 0x3 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	}

	// Now, subtract nothing and check the value.
	c8.Opcode = MakeOpcode(0x8135)
	c8.DecodeExecute()

	if c8.Registers[1] != 0xA2 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[2] != 0x3 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	}

	// Subtract 2 (3) from 1 (A2).
	c8.Opcode = MakeOpcode(0x8125)
	c8.DecodeExecute()

	if c8.Registers[1] != 0x9F {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[2] != 0x3 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	} else if c8.Registers[0xF] != 1 {
		t.Error("Register underflow was falsely reported!\n")
	}

	// Subtract 1 (9F) from 2 (3), check for underflow.
	c8.Opcode = MakeOpcode(0x8215)
	c8.DecodeExecute()

	if c8.Registers[2] != 0x64 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	} else if c8.Registers[1] != 0x9F {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[0xF] != 0 { // Check for underflow.
		t.Error("Carry register was not 0 on underflow!\n")
	}
}

func TestShift(t *testing.T) {
	c8 := MakeChip8(false)

	// Load register 1 with 1.
	c8.Opcode = MakeOpcode(0x7101) // Register 1 has 1.
	c8.DecodeExecute()

	if c8.Registers[1] != 0x1 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	}

	// Now, shift it left.
	c8.Opcode = MakeOpcode(0x819E) // 9 can be anything.
	c8.DecodeExecute()

	if c8.Registers[1] != 0x2 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[0xF] != 0 { // Check that the MSB was 0.
		t.Errorf("MSB of shifted number was not 0! Val was %v\n",
			c8.Registers[0xF])
	}

	// Now, shift the register right twice.
	c8.Opcode = MakeOpcode(0x8176)
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x8166)
	c8.DecodeExecute()

	if c8.Registers[1] != 0x0 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[0xF] != 1 { // Check that the MSB was 0.
		t.Errorf("MSB of shifted number was not 1! Val was %v\n",
			c8.Registers[0xF])
	}
}

func TestSaveRestoreRegs(t *testing.T) {
	c8 := MakeChip8(false)

	// First, load the reigsters with some data.
	c8.Opcode = MakeOpcode(0x71A1) // Reg 1 has A1.
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x7206) // Reg 2 has 06.
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x76D4) // Reg 3 has D4.
	c8.DecodeExecute()

	// Check that the registers are correctly filled with data.
	if c8.Registers[1] != 0xA1 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[2] != 0x06 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	} else if c8.Registers[6] != 0xD4 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[6])
	}

	// Now, set the index register to our memory save location
	// (here, just arbitrarily pick 0x345).
	c8.Opcode = MakeOpcode(0xA345)
	c8.DecodeExecute()

	if c8.IndexReg != 0x345 {
		t.Errorf("Index register value was not updated correctly! Val was %v\n",
			c8.IndexReg)
	}

	// Now, load our registers (up to register 6) into memory.
	c8.Opcode = MakeOpcode(0xF655)
	c8.DecodeExecute()

	// Check that the registers are correctly filled with data.
	if c8.Memory[0x346] != 0xA1 {
		t.Errorf("Memory was not loaded with register contents properly! Val was %v\n",
			c8.Memory[0x346])
	} else if c8.Memory[0x347] != 0x06 {
		t.Errorf("Memory was not loaded with register contents properly! Val was %v\n",
			c8.Memory[0x347])
	} else if c8.Memory[0x34B] != 0xD4 {
		t.Errorf("Memory was not loaded with register contents properly! Val was %v\n",
			c8.Memory[0x34B])
	}

	// Now, change registers 1 and 5.
	c8.Opcode = MakeOpcode(0x7101)
	c8.DecodeExecute()
	c8.Opcode = MakeOpcode(0x75DD)
	c8.DecodeExecute()

	// Check that they've been updated.
	if c8.Registers[1] != 0xA2 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[5] != 0xDD {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[5])
	} else if c8.IndexReg != 0x345 {
		t.Errorf("Index register value spuriously updated! Val was %v\n",
			c8.IndexReg)
	}

	// Now, reload our registers with memory contents and check them.
	c8.Opcode = MakeOpcode(0xF665)
	c8.DecodeExecute()

	if c8.Registers[1] != 0xA1 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[1])
	} else if c8.Registers[2] != 0x06 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	} else if c8.Registers[5] != 0x0 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[5])
	} else if c8.Registers[6] != 0xD4 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[6])
	} else if c8.IndexReg != 0x345 {
		t.Errorf("Index register value was spuriously updated! Val was %v\n",
			c8.IndexReg)
	}
}
