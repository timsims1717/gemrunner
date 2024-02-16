package debug

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	ShowDebug = false
	ShowText  = false
	winV      *pixel.Vec
)

func Initialize(v *pixel.Vec) {
	winV = v
	InitializeLines()
	InitializeText()
	InitializeFPS()
}

func Draw(win *pixelgl.Window) {
	DrawLines(win)
	DrawText(win)
	DrawFPS(win)
}

func Clear() {
	imd.Clear()
	lines = []string{}
}
