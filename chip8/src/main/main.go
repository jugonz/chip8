package main

import "arch"

var chip8 arch.Arch

func main() {

	chip8 = arch.MakeChip8()
	// init()
	for {
		chip8.EmulateCycle()

		// if drawFlag
		chip8.DrawScreen()

		chip8.SetKeys()
	}
}
