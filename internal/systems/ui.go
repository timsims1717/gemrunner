package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"github.com/gopxl/pixel"
)

func CursorInit() {
	if data.CursorEntity != nil {
		myecs.Manager.DisposeEntity(data.CursorEntity)
	}
	data.CursorObj = object.New()
	data.CursorObj.Layer = -10
	data.CursorObj.Offset = pixel.V(9, -9)
	data.CursorEntity = myecs.Manager.NewEntity()
	data.CursorEntity.AddComponent(myecs.Object, data.CursorObj)
	data.CursorEntity.AddComponent(myecs.Drawable, img.NewSprite("cursor", constants.UIBatch))
}

func CursorSystem(hideWhenIdle bool) {
	hide := false
	if hideWhenIdle {
		if data.HideCursorTimer == nil || data.MenuInput.MouseMoved {
			data.HideCursorTimer = timing.New(3.)
		}
		data.HideCursorTimer.Update()
		if data.HideCursorTimer.Done() {
			hide = true
		}
	}
	data.CursorObj.Hidden = hide
	data.CursorObj.Pos = data.MenuInput.Cursor
}
