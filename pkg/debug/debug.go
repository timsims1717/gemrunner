package debug

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
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
	if Text {
		DrawText(win)
	}
	DrawFPS(win)
}

func Clear() {
	imd.Clear()
	lines = []string{}
}
