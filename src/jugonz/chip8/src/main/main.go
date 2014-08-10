package main

import (
	"jugonz/chip8/src/arch"
)

var chip8 arch.Arch

func main() {
	count := 0
	chip8 = arch.MakeChip8(true) // DEBUG on.
	chip8.LoadGame("../../c8games/PONG2")

	for !chip8.ShouldClose() {
		count++
		chip8.EmulateCycle()

		if chip8.ShouldDraw() {
			chip8.DrawScreen()
		}

		chip8.SetKeys()
	}
}
