package gfx

type Drawable interface {
	Draw()
	SetPixel(x, y uint16)
	ClearPixel(x, y uint16)
	GetPixel(x, y uint16) bool
	ClearScreen()
}
