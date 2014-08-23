package gfx

type Drawable interface {
	Draw()
	XorPixel(x, y uint16)
	GetPixel(x, y uint16) bool
	InBounds(x, y uint16) bool
	ClearScreen()
}
