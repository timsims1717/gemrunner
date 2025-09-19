package debug

import (
	"github.com/gopxl/pixel"
)

var (
	ShowDebug = false
	ShowText  = false
	Verbose   = false
	winV      *pixel.Vec
)

func Initialize(v *pixel.Vec) {
	winV = v
	InitializeLines()
	InitializeText()
	InitializeFPS()
}

func Draw(target pixel.Target, dim pixel.Vec) {
	DrawLines(target)
	DrawText(target, dim)
	DrawFPS(target, dim)
}

func Clear() {
	imd.Clear()
	lines = []string{}
}
