package gfx

import (
	"fmt"
	gl "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Screen struct {
	Width     int
	Height    int
	ResWidth  int
	ResHeight int
	Pixels    [64][32]bool // Chip8 resolution is static.
	Title     string
	Window    glfw.Window
}

func MakeScreen(width int, height int, title string) Screen {
	s := Screen{}
	s.Width = width
	s.Height = height
	s.Title = title

	// Chip8 resolution is hardcoded.
	s.ResWidth = len(s.Pixels)
	s.ResHeight = len(s.Pixels[0])

	s.Init()
	return s
}

func (s *Screen) Init() {
	// 1. Initialize GLFW and save window context.
	glfw.SetErrorCallback(s.GFXError)

	if !glfw.Init() { // Init GLFW3...
		panic("GLFW3 failed to initialize!\n")
	}

	win, err := glfw.CreateWindow(s.Width, s.Height, s.Title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("GLFW could not create window! Error: %v\n", err))
	}

	win.SetKeyCallback(keyCallback)
	win.MakeContextCurrent()
	glfw.SwapInterval(1) // Use videosync. (People say it's good.)

	s.Window = *win

	// 2. Initalize OpenGL.
	if gl.Init() != 0 {
		panic("OpenGL failed to initialize!\n")
	}

	// 3. Draw a black screen and set the coordinate system.
	gl.ClearColor(0, 0, 0, 0)
	gl.MatrixMode(gl.PROJECTION)
	gl.Ortho(0, float64(s.ResWidth), float64(s.ResHeight), 0, 0, 1)

	fmt.Println("Screen successfully initialized.")
}

/**
 * Methods to implement the Drawable interface.
 */
func (s *Screen) Draw() {
	// I have no idea what I'm doing with OpenGL, so
	// this code is adapted from
	// https://github.com/nictuku/chip-8/blob/master/system/video.go

	//gl.Viewport(0, 0, s.Width, s.Height)
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.MatrixMode(gl.POLYGON)

	for xLine := 0; xLine < s.ResWidth; xLine++ {
		for yLine := 0; yLine < s.ResHeight; yLine++ {

			if !s.Pixels[xLine][yLine] {
				gl.Color3d(0, 0, 0)
			} else {
				gl.Color3d(1, 1, 1) // Draw white.
			}
			x, y := float64(xLine), float64(yLine)
			gl.Rectd(x, y, x+1, y+1)
		}
	}

	s.Window.SwapBuffers() // Display what we just drew.
}

func (s *Screen) ClearScreen() {
	for xLine := 0; xLine < s.ResWidth; xLine++ {
		for yLine := 0; yLine < s.ResHeight; yLine++ {
			s.Pixels[xLine][yLine] = false
		}
	}
}

func (s *Screen) SetPixel(x, y uint16) {
	s.Pixels[x][y] = true
}

func (s *Screen) ClearPixel(x, y uint16) {
	s.Pixels[x][y] = false
}

func (s *Screen) GetPixel(x, y uint16) bool {
	return s.Pixels[x][y]
}

/**
 * Methods to implement the Interactible interface.
 */
func (s *Screen) SetKeys() {
	glfw.PollEvents()
}

func (s *Screen) ShouldClose() bool {
	return s.Window.ShouldClose()
}

func (s *Screen) Quit() {
	glfw.Terminate()
}

/**
 * Utility methods used for internal GLFW state.
 */
func (s *Screen) GFXError(err glfw.ErrorCode, msg string) {
	panic(fmt.Errorf("GLFW Error: %v: %v\n", err, msg))
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}
