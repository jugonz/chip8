package gfx

import (
	"fmt"
	"github.com/go-gl/gl"
	"github.com/go-gl/glfw3"
)

type Screen struct {
	Width  int
	Height int
	Title  string
	Window glfw3.Window
}

func MakeScreen(width int, height int, title string) Screen {
	s := Screen{}
	s.Width = width
	s.Height = height
	s.Title = title
	s.Init()

	return s
}

func (s *Screen) GFXError(err glfw3.ErrorCode, msg string) {
	panic(fmt.Errorf("GLFW Error: %v: %v\n", err, msg))
}

func (s *Screen) Init() {
	// Setup: Set GL error callback function...
	glfw3.SetErrorCallback(s.GFXError)

	if !glfw3.Init() { // Init GLFW3...
		panic("GLFW3 failed to initialize!\n")
	}

	// Now, create a window!
	win, err := glfw3.CreateWindow(s.Width, s.Height, s.Title, nil, nil)
	if err != nil {
		panic(fmt.Errorf("GLFW could not create window! Error: %v\n", err))
	}
	win.MakeContextCurrent()
	s.Window = *win

	glfw3.SwapInterval(1) // Use videosync. (People say it's good.)

	// Now, init OpenGL.
	if gl.Init() != 0 {
		panic("OpenGL could not initialize!\n")
	}
}

func (s *Screen) Draw(data [2048]bool) {
	// I have no idea what I'm doing with OpenGL, so
	// this code is adapted from https://github.com/nictuku/chip-8/blob/master/system/video.go

	//gl.MatrixMode(gl.POLYGON)

	for yline := 0; yline < s.Height; yline++ {

		for xline := 0; xline < s.Width; xline++ {

			x, y := float32(xline), float32(yline)
			if !data[xline+yline*64] { // False = 0.
				gl.Color3f(0, 0, 0)
			} else { // True = 1.
				gl.Color3f(1, 1, 1)
			}
			gl.Rectf(x, y, x+1, y+1)
		}
	}

	s.Window.SwapBuffers()
}

func (s *Screen) Quit() {
	gl.End()
	glfw3.Terminate()
}
