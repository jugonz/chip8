package gfx

import (
	"fmt"
	gl "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Screen struct {
	Width  int
	Height int
	Title  string
	Window glfw.Window
}

func MakeScreen(width int, height int, title string) Screen {
	s := Screen{}
	s.Width = width
	s.Height = height
	s.Title = title
	s.Init()

	return s
}

func (s *Screen) Init() {
	// Setup: Set GL error callback function...
	glfw.SetErrorCallback(s.GFXError)

	if !glfw.Init() { // Init GLFW3...
		panic("GLFW3 failed to initialize!\n")
	}

	// Now, create a window!
	win, err := glfw.CreateWindow(s.Width, s.Height, s.Title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("GLFW could not create window! Error: %v\n", err))
	}

	win.SetKeyCallback(keyCallback)
	win.MakeContextCurrent()
	s.Window = *win

	glfw.SwapInterval(1) // Use videosync. (People say it's good.)

	// Now, init OpenGL.
	if gl.Init() != 0 {
		panic("OpenGL could not initialize!\n")
	}

	gl.ClearColor(0, 0, 0, 0)
	gl.MatrixMode(gl.PROJECTION)

	gl.Ortho(0, 64, 32, 0, 0, 1)
	fmt.Println("screen done init")
}

func (s *Screen) Draw(data [2048]bool) {
	// I have no idea what I'm doing with OpenGL, so
	// this code is adapted from https://github.com/nictuku/chip-8/blob/master/system/video.go

	//gl.Viewport(0, 0, s.Width, s.Height)
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.MatrixMode(gl.POLYGON)
	//gl.Begin(gl.POLYGON)

	for yline := 0; yline < 32; yline++ {
		for xline := 0; xline < 64; xline++ {

			x, y := float32(xline), float32(yline)
			if !data[xline+(yline*64)] { // False = 0.
				fmt.Println("drawing meeeee...")
				gl.Color3f(0, 0, 0)
			} else { // True = 1.
				fmt.Println("drawing youuuuu...")
				gl.Color3f(1, 1, 1)
			}
			gl.Rectf(x, y, x+1, y+1)
		}
	}

	//gl.End()
	s.Window.SwapBuffers()
	glfw.PollEvents()
}

func (s *Screen) Quit() {
	glfw.Terminate()
}

func (s *Screen) GFXError(err glfw.ErrorCode, msg string) {
	panic(fmt.Errorf("GLFW Error: %v: %v\n", err, msg))
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}
