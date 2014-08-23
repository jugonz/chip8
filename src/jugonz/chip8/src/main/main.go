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

	chip8.Run() // Terminates when the quit key is pressed.

	chip8.Quit()
	runtime.UnlockOSThread()
}
