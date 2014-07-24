package arch

type Opcode struct {
	Value   uint16
	Xreg    uint8  // "X" register in opcode table. 0 if unused, max 16.
	Yreg    uint8  // 0 if unused.
	Literal uint16 // Last three hex digits of our opcode.
}

func MakeOpcode(opcode uint16) Opcode {
	op := Opcode{}
	op.Value = opcode
	op.Xreg = uint8((opcode >> 8) & 0xF) // Safe: F fits in uint8.
	op.Yreg = uint8((opcode >> 4) & 0xF) // Likewise.
	op.Literal = opcode & 0xFFF

	return op
}
