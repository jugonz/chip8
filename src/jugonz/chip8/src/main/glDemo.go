package main

import (
	"fmt"
	gl "github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		w.SetShouldClose(true)
	}
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyCallback)
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if gl.Init() != 0 {
		panic("failed to init gl")
	}

	w, h := 640, 480
	ratio := float64(640) / float64(480)
	for !window.ShouldClose() {
		//Do OpenGL stuff

		gl.Viewport(0, 0, w, h)       // Make viewport
		gl.Clear(gl.COLOR_BUFFER_BIT) // clear it
		gl.MatrixMode(gl.PROJECTION)
		gl.LoadIdentity()
		gl.Ortho(-ratio, ratio, float64(-1), float64(1), float64(1), float64(-1))
		gl.MatrixMode(gl.MODELVIEW)
		gl.LoadIdentity()
		angle := glfw.GetTime() * float64(50)
		gl.Rotatef(float32(angle), float32(0), float32(0), float32(1))

		gl.Begin(gl.TRIANGLES)
		gl.Color3f(float32(1), float32(0), float32(0))
		gl.Vertex3f(float32(-0.6), float32(-0.4), float32(0))
		gl.Color3f(float32(0), float32(1), float32(0))
		gl.Vertex3f(float32(0.6), float32(-0.4), float32(0))
		gl.Color3f(float32(0), float32(0), float32(1))
		gl.Vertex3f(float32(0), float32(0.6), float32(0))
		gl.End()

		window.SwapBuffers()
		glfw.PollEvents()
	}

}
