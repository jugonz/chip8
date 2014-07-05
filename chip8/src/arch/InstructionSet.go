package arch

/**
 * Datatype to describe an instruction set of a simple computer.
 */
type InstructionSet interface {
	// Graphics controls
	ClearScreen()
	DrawSprite()
	SetIndexToSprite()

	// Control flow
	CallRCA1802()
	Return()
	Jump()
	JumpIndexLiteralOffset()
	Call()
	SkipInstrEqualLiteral()
	SkipInstrNotEqualLiteral()
	SkipInstrEqualReg()
	SkipInstrNotEqualReg()
	SkipInstrKeyPressed()
	SkipInstrKeyNotPressed()

	// Manipulating data registers
	SetRegToLiteral()
	SetRegToReg()
	Add()
	AddWithCarry()
	Or()
	And()
	Xor()
	SubXFromY()
	SubYFromX()
	ShiftRight()
	ShiftLeft()
	SetRegisterRandomMask()
	SaveBinaryCodedDecimal()
	GetKeyPress()
	GetDelayTimer()

	// Manipulating special registers
	AddRegisterToIndex()
	SetIndexLiteral()
	SetDelayTimer()
	SetSoundTimer()

	// Context switching
	SaveRegisters()
	RestoreRegisters()

	// Special
	UnknownInstruction()
}
