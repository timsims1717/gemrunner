package debug

import (
	"fmt"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/typeface"
	"github.com/gopxl/pixel"
)

var (
	fpsText     *typeface.Text
	versionText *typeface.Text
	Release     int
	Version     int
	Build       int
)

func InitializeFPS() {
	fpsText = typeface.New("basic").
		WithAnchor(pixel.TopRight).
		WithScalar(2.)
	versionText = typeface.New("basic").
		WithAnchor(pixel.TopLeft).
		WithScalar(2.)
}

func DrawFPS(target pixel.Target, dim pixel.Vec) {
	fpsText.SetText(fmt.Sprintf("FPS: %d", timing.FPS))
	fpsText.Obj.Pos = winV.Add(pixel.V(dim.X*-0.5+2., dim.Y*-0.5+2))
	fpsText.Obj.Update()
	fpsText.Draw(target)
	versionText.SetText(fmt.Sprintf("%d.%d.%d", Release, Version, Build))
	versionText.Obj.Pos = winV.Add(pixel.V(dim.X*0.5-2., dim.Y*-0.5+2))
	versionText.Obj.Update()
	versionText.Draw(target)
}
