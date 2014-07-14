package main

import "arch"

var chip8 arch.Arch

func main() {

	chip8 = arch.MakeChip8()
	// Load game
	for {
		chip8.FetchOpcode()
		chip8.EmulateCycle()

		if chip8.DrawFlag {
			chip8.DrawScreen()
		}

		chip8.SetKeys()
	}
}
