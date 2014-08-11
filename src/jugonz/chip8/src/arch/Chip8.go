package arch

import (
	"fmt"
	"io"
	"jugonz/chip8/src/gfx"
	"math/rand"
	"os"
	"time"
)

/**
 * Datatype to describe the architecture of a Chip8 system.
 */
type Chip8 struct {
	// Core structural components.
	Opcode     Opcode
	Memory     [4096]uint8
	Registers  [16]uint8
	IndexReg   uint16
	PC         uint16
	DelayTimer uint8
	SoundTimer uint8
	Stack      [16]uint16
	SP         uint16
	Rando      *rand.Rand // PRNG

	// Interactive components.
	Keyboard   [16]bool // True if key pressed.
	Screen     gfx.Drawable
	Fontset    [80]uint8
	DrawFlag   bool // True if we just drew to the screen.
	Controller gfx.Interactible

	// Debug components.
	Debug bool
	Count int
}

func MakeChip8(debug bool) *Chip8 { // and initialize
	c8 := Chip8{}
	c8.Opcode = Opcode{}
	c8.PC = 0x200 // Starting PC address is static.
	c8.Rando = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define fonset.
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
	// Load fontset into memory.
	for char := 0; char < len(c8.Fontset); char++ {
		c8.Memory[char] = c8.Fontset[char]
	}

	screen := gfx.MakeScreen(640, 480, "Chip-8 Emulator")
	c8.Screen = &screen
	c8.Controller = &screen
	c8.Debug = debug
	return &c8
}

func (c8 *Chip8) LoadGame(filePath string) {
	// Open file and load into memory, else panic.
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf(
			"Error: File at %v could not be loaded! Error was: %v\n",
			filePath, err))
	}
	defer file.Close()

	// Now, we've opened the file, and can load
	// its contents into memory.
	// We must first get the file size.
	stat, err := file.Stat()
	if err != nil {
		panic(fmt.Sprintf(
			"Error: Info about file at %v could not be found! Error was: %v\n",
			filePath, err))
	}

	buffer := make([]byte, stat.Size()) // Make new buffer to store game.
	_, err = io.ReadFull(file, buffer)
	if err != nil {
		panic(fmt.Sprintf(
			"Error: File at %v could not be read completely! Error was: %v\n",
			filePath, err))
	}

	for index, value := range buffer {
		c8.Memory[index+0x200] = value
	}

}

func (c8 *Chip8) EmulateCycle() {
	c8.FetchOpcode() // Fetch instruction.
	if c8.Debug {
		fmt.Printf("On cycle %v, at mem loc %X\n", c8.Count, c8.PC)
		c8.Count++
	}
	c8.DecodeExecute()
}

func (c8 *Chip8) DecodeExecute() {
	switch c8.Opcode.Value >> 12 { // Decode (big-ass switch statement)
	case 0x0:
		switch c8.Opcode.Value & 0xFF {
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
		switch c8.Opcode.Value & 0xF {
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
	case 0xA:
		c8.SetIndexLiteral()
	case 0xB:
		c8.JumpIndexLiteralOffset()
	case 0xC:
		c8.SetRegisterRandomMask()
	case 0xD:
		c8.DrawSprite()
	case 0xE:
		switch c8.Opcode.Value & 0xFF {
		case 0x9E:
			c8.SkipInstrKeyPressed()
		case 0xA1:
			c8.SkipInstrKeyNotPressed()
		default:
			c8.UnknownInstruction()
		}
	case 0xF:
		switch c8.Opcode.Value & 0xFF {
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
			c8.SaveBinaryCodedDecimal()
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
}

func (c8 *Chip8) FetchOpcode() {
	newOp := uint16(c8.Memory[c8.PC]) << 8
	newOp |= uint16(c8.Memory[c8.PC+1])
	c8.Opcode = MakeOpcode(newOp)
}

func (c8 *Chip8) DrawScreen() {
	if c8.DrawFlag {
		c8.Screen.Draw()
		c8.DrawFlag = false
	}
}

func (c8 *Chip8) SetKeys() {
	c8.Controller.SetKeys()
}

func (c8 *Chip8) ShouldClose() bool {
	return c8.Controller.ShouldClose()
}

func (c8 *Chip8) Quit() {
	c8.Controller.Quit()
}
