package arch

import "testing"

func TestSetup(t *testing.T) {
	c8 := MakeChip8()
	if c8.PC != 0x200 {
		t.Errorf("c8 Opcode was not initialized properly! Was: ",
			c8.Opcode)
	}

	// Check fontset.
	for i := 0; i < 80; i++ {
		if c8.Fontset[i] == 0x00 {
			t.Errorf("c8 Fontset was not loaded!")
		}
	}

	// TODO: test loading game
}

func TestClearScreen(t *testing.T) {
	c8 := MakeChip8()

	// Draw something to the screen, and see that it is not empty.
	c8.Opcode = 0xD324
	c8.DecodeExecute()

	clear := true
	for _, value := range c8.GFX {
		if value {
			clear = false
			break
		}
	}
	if clear {
		t.Errorf("DrawSprite failed to draw the screen!")
	}

	// Now, clear the screen, and check that it is empty.
	c8.Opcode = 0x00E0
	c8.DecodeExecute()

	clear = true
	for _, value := range c8.GFX {
		if value {
			clear = false
			break
		}
	}
	if !clear {
		t.Errorf("ClearScreen failed to clear the screen!")
	}
}

func TestDrawSprite(t *testing.T) {
	// Do this one later. Phew.
}

func TestSetIndexToSprite(t *testing.T) {
	// Do this one later too.
}

func TestCallReturn(t *testing.T) {
	c8 := MakeChip8()

	// Check initial stack.
	if c8.SP != 0 {
		t.Errorf("Stack not properly initalized!")
	}
	for index, val := range c8.Stack {
		if val != 0 {
			t.Errorf("Stack not empty! Val at index %v is %X\n",
				index, val)
		}
	}

	// Call a program at 789, and check stack.
	c8.Opcode = 0x2789
	c8.DecodeExecute()

	if c8.SP != 1 {
		t.Errorf("Stack not updated after call!")
	} else if c8.Stack[0] != 0x200 { // Starting PC
		t.Errorf("Old PC was not saved!")
	}
	for index, val := range c8.Stack[1:] {
		if val != 0 {
			t.Errorf("Stack not empty! Val at index %v is %X\n",
				index, val)
		}
	}

	// Return from that program, and see where we are.
	c8.Opcode = 0x00EE
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
	c8 := MakeChip8()

	// Test a simple add from 0 to register 2.
	if c8.Registers[2] != 0 {
		t.Errorf("Register did not start at 0!")
	}

	c8.Opcode = 0x7212
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[2] != 0x12 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[2])
	}

	// Test the max value in a different register.
	c8.Opcode = 0x73FF
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0xFF {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}

	// Now, add one, and test overflow.
	c8.Opcode = 0x7301
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0x00 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}
}

func TestAddWithCarry(t *testing.T) {
	c8 := MakeChip8()

	// Test adding the max value without overflow.
	c8.Opcode = 0x73FF // Add FF to reg 3 (0).
	c8.DecodeExecute()
	c8.Opcode = 0x8374 // Add reg 7 (0) to reg 3.
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0xFF {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	} else if c8.Registers[0xF] != 0 {
		t.Error("Carry register was not 0 on non-overflowing calculation!")
	}

	// Now, add one, and test overflow.
	c8.Opcode = 0x7401 // Add 1 to reg 4 (0).
	c8.DecodeExecute()
	c8.Opcode = 0x8344 // Add reg 4 (1) to reg 3 (FF).
	c8.DecodeExecute()

	// Test that value is now correct.
	if c8.Registers[3] != 0x00 {
		t.Errorf("Register value was not updated correctly! Val was %v\n",
			c8.Registers[3])
	}
	// Test that the overflow register is correctly set.
	if c8.Registers[0xF] != 1 {
		t.Error("Carry register was not 1 on overflowing calculation!")
	}

}
