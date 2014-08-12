package main

import (
	"jugonz/chip8/src/arch"
)

var chip8 arch.Arch

func main() {
	chip8 = arch.MakeChip8(true) // DEBUG on.
	chip8.LoadGame("../../c8games/PONG")

	for !chip8.ShouldClose() {
		chip8.EmulateCycle()

		chip8.DrawScreen() // Only draws if needed.
		chip8.SetKeys()
	}

	chip8.Quit()
}
