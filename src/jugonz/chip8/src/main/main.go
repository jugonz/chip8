package main

import (
	"jugonz/chip8/src/arch"
	"runtime"
)

var chip8 arch.Arch

func main() {
	runtime.LockOSThread()       // OpenGL requires code to be run on main thread.
	chip8 = arch.MakeChip8(true) // DEBUG on.

	chip8.LoadGame("../../c8games/PONG2")

	for !chip8.ShouldClose() {
		chip8.EmulateCycle()

		chip8.DrawScreen() // Only draws if needed.
		chip8.SetKeys()
		chip8.UpdateTimers()
	}

	chip8.Quit()
	runtime.UnlockOSThread()
}
