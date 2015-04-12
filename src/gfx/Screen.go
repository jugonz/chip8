package gfx

import (
	"fmt"
	gl "github.com/go-gl/gl/v2.1/gl"
	glfw "github.com/go-gl/glfw/v3.0/glfw"
)

// Arrays cannot be const in Go, so the keyboard layout is a var.
var keyLayout = [16]glfw.Key{
	glfw.Key0, glfw.Key1, glfw.Key2, glfw.Key3,
	glfw.Key4, glfw.Key5, glfw.Key6, glfw.Key7,
	glfw.Key8, glfw.Key9, glfw.KeyA, glfw.KeyB,
	glfw.KeyC, glfw.KeyD, glfw.KeyE, glfw.KeyF,
}
var keyQuit = glfw.KeyEscape

type Screen struct {
	Width     int
	Height    int
	ResWidth  int
	ResHeight int
	Pixels    [][]bool
	Title     string
	Window    glfw.Window
	Keyboard  [16]bool // True if key pressed.
}

func MakeScreen(width int, height int, resWidth int, resHeight int,
	title string) Screen {
	s := Screen{}
	s.Width = width
	s.Height = height
	s.Title = title
	s.ResWidth = resWidth
	s.ResHeight = resHeight

	s.Pixels = make([][]bool, s.ResWidth)
	for col := range s.Pixels {
		s.Pixels[col] = make([]bool, s.ResHeight)
	}

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

	win.SetInputMode(glfw.StickyKeys, 1) // Turn on sticky keys to avoid callbacks!
	win.MakeContextCurrent()

	s.Window = *win

	// 2. Initalize OpenGL.
	err = gl.Init()
	if err != nil {
		panic("OpenGL failed to initialize!\n")
	}

	glfw.SwapInterval(1) // Use videosync. (People say it's good.)

	// 3. Draw a black screen and set the coordinate system.
	gl.ClearColor(0, 0, 0, 0)
	gl.MatrixMode(gl.PROJECTION)
	gl.Ortho(0, float64(s.ResWidth), float64(s.ResHeight), 0, 0, 1)
}

/**
 * Methods to implement the Drawable interface.
 */
func (s *Screen) Draw() {
	// I have no idea what I'm doing with OpenGL, so
	// this code is adapted from
	// https://github.com/nictuku/chip-8/blob/master/system/video.go

	//gl.Viewport(0, 0, s.Width, s.Height)
	//gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

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

func (s *Screen) XorPixel(x, y uint16) {
	s.Pixels[x][y] = s.Pixels[x][y] != true
}

func (s *Screen) GetPixel(x, y uint16) bool {
	return s.Pixels[x][y]
}

func (s *Screen) InBounds(x, y uint16) bool {
	return int(x) < s.ResWidth && int(y) < s.ResHeight
}

/**
 * Methods to implement the Interactible interface.
 */
func (s *Screen) SetKeys() {
	// Handle input ourselves!
	glfw.PollEvents()
	for keyNum, key := range keyLayout {
		s.ProcessKey(keyNum, key)
	}

	// Special case: if escape key is pressed, just quit.
	if quitState := s.Window.GetKey(keyQuit); quitState == glfw.Press {
		s.Window.SetShouldClose(true)
	}
}

func (s *Screen) ProcessKey(keyNum int, key glfw.Key) {
	action := s.Window.GetKey(key)

	switch action {
	case glfw.Press:
		//fmt.Printf("Key %X pressed!\n", keyNum)
		s.Keyboard[keyNum] = true
	case glfw.Release:
		s.Keyboard[keyNum] = false
	default:
		// Ignore key repeat.
	}
}

func (s *Screen) KeyPressed(key uint8) bool {
	return s.Keyboard[key]
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
