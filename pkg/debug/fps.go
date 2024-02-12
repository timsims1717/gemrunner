package debug

import (
	"fmt"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/typeface"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	fpsText     *typeface.Text
	versionText *typeface.Text
	Release     int
	Version     int
	Build       int
)

func InitializeFPS() {
	fpsText = typeface.New("basic", typeface.NewAlign(typeface.Left, typeface.Bottom), 1.0, 2.0, 0., 0.)
	versionText = typeface.New("basic", typeface.NewAlign(typeface.Right, typeface.Bottom), 1.0, 2.0, 0., 0.)
}

func DrawFPS(win *pixelgl.Window) {
	fpsText.SetText(fmt.Sprintf("FPS: %d", timing.FPS))
	fpsText.Obj.Pos = winV.Add(pixel.V(win.Bounds().W()*-0.5+2., win.Bounds().H()*-0.5+2))
	fpsText.Obj.Update()
	fpsText.Draw(win)
	versionText.SetText(fmt.Sprintf("%d.%d.%d", Release, Version, Build))
	versionText.Obj.Pos = winV.Add(pixel.V(win.Bounds().W()*0.5-2., win.Bounds().H()*-0.5+2))
	versionText.Obj.Update()
	versionText.Draw(win)
}
