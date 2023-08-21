package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	pxginput "github.com/timsims1717/pixel-go-input"
)

func HoverSystem(in *pxginput.Input) {
	for _, result := range myecs.Manager.Query(myecs.HasHover) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		hover, okH := result.Components[myecs.Hover].(*data.HoverFunky)
		if okO && okH {
			pos := in.World
			if view, okV := result.Entity.GetComponentData(myecs.ViewPort); okV {
				if vp, okV2 := view.(*viewport.ViewPort); okV2 {
					pos = vp.Projected(in.World)
				}
			}
			hovered := obj.PointInside(pos)
			if hovered && hover.Fn != nil {
				hover.Fn(in)
			}
		}
	}
}
