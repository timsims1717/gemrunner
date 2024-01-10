package debug

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	Debug = false
	Text  = false
	winV  *pixel.Vec
)

func Initialize(v *pixel.Vec) {
	winV = v
	InitializeLines()
	InitializeText()
	InitializeFPS()
}

func Draw(win *pixelgl.Window) {
	if Debug {
		DrawLines(win)
	}
	if Text {
		DrawText(win)
	}
	DrawFPS(win)
}

func Clear() {
	imd.Clear()
	lines = []string{}
}
