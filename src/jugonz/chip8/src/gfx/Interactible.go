package gfx

type Interactible interface {
	SetKeys()
	ShouldClose() bool
	Quit()
}
