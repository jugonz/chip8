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

	// Manipulate data registers
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
	SetRegisterDelayTimer()
	SetRegisterKeyPress()

	// Manipulate special registers
	AddToIndex()
	SetIndexLiteral()
	SetDelayTimer()
	SetSoundTimer()
	BinaryMagic()

	// Context switching
	SaveRegisters()
	RestoreRegisters()
}
