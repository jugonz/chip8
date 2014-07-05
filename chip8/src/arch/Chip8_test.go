package arch

import "testing"

func TestSetup(t *testing.T) {

	c8 := MakeChip8()
	if c8.Opcode != 0x200 {
		t.Errorf("c8 Opcode was not initialized properly! Was: ", c8.Opcode)
	}

}
