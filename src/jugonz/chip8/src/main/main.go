package main

import (
	"flag"
	"fmt"
	"jugonz/chip8/src/arch"
	"runtime"
)

var path = flag.String("path", "", "path to a Chip8 ROM")
var debug = flag.Bool("debug", false, "debug mode")
var chip8 arch.Arch

func main() {
	flag.Parse()
	if *path == "" {
		fmt.Printf("No chip8 file path provided, quitting!\n")
		return
	}

	runtime.LockOSThread()         // OpenGL requires code to be run on main thread.
	chip8 = arch.MakeChip8(*debug) // DEBUG on.

	chip8.LoadGame(*path)

	chip8.Run() // Terminates when the quit key is pressed.

	chip8.Quit()
	runtime.UnlockOSThread()
}
