package data

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

var (
	DialogStack     []*Dialog
	DialogsOpen     []*Dialog
	DialogStackOpen bool
	Dialogs         = map[string]*Dialog{}

	OpenPuzzleConstructor  *DialogConstructor
	EditorPanelConstructor *DialogConstructor
	EditorOptConstructor   *DialogConstructor
)

type Dialog struct {
	Key          string
	Pos          pixel.Vec
	ViewPort     *viewport.ViewPort
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	Buttons      []*Button

	Open  bool
	Click bool
}

type DialogConstructor struct {
	Key     string
	Width   float64
	Height  float64
	Pos     pixel.Vec
	Buttons []ButtonConstructor
}

func NewDialog(dc *DialogConstructor) {
	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, dc.Width*world.TileSize, dc.Height*world.TileSize))
	vp.CamPos = pixel.V(0, 0)
	vp.PortPos = viewport.MainCamera.PostCamPos.Add(dc.Pos)

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (dc.Width+1)*world.TileSize, (dc.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = viewport.MainCamera.PostCamPos.Add(dc.Pos)

	bObj := object.New()
	bObj.Layer = 99
	//bObj.Pos = dc.Pos
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, bObj).
		AddComponent(myecs.Border, &Border{
			Width:  int(dc.Width),
			Height: int(dc.Height),
		})

	dlg := &Dialog{
		Key:          dc.Key,
		Pos:          dc.Pos,
		ViewPort:     vp,
		BorderVP:     bvp,
		BorderObject: bObj,
		BorderEntity: be,
	}

	var buttons []*Button
	for _, btn := range dc.Buttons {
		obj := object.New().WithID("cancel_btn")
		obj.Pos = btn.Position
		obj.Layer = 99
		obj.SetRect(img.Batchers[constants.UIBatch].GetSprite(btn.SprKey).Frame())
		spr := img.NewSprite(btn.SprKey, constants.UIBatch)
		cSpr := img.NewSprite(btn.ClickSprKey, constants.UIBatch)
		btnE := myecs.Manager.NewEntity()
		btnE.AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr)
		b := &Button{
			Key:      btn.SprKey,
			Sprite:   spr,
			ClickSpr: cSpr,
			HelpText: btn.HelpText,
			Object:   obj,
			Entity:   btnE,
		}
		btnE.AddComponent(myecs.Update, NewHoverClickFn(MainInput, vp, func(hvc *HoverClick) {
			click := hvc.Input.Get("click")
			if dlg.Open {
				if hvc.Hover && click.JustPressed() {
					dlg.Click = true
				}
				if hvc.Hover && click.Pressed() && dlg.Click {
					btnE.AddComponent(myecs.Drawable, cSpr)
				} else {
					if hvc.Hover && click.JustReleased() && dlg.Click {
						dlg.Click = false
						if b.OnClick != nil {
							b.OnClick()
						}
					} else if !click.Pressed() && !click.JustReleased() && dlg.Click {
						dlg.Click = false
						btnE.AddComponent(myecs.Drawable, spr)
					} else {
						btnE.AddComponent(myecs.Drawable, spr)
					}
				}
			}
		}))
		buttons = append(buttons, b)
	}
	dlg.Buttons = buttons

	Dialogs[dc.Key] = dlg
}

func ClearDialogStack() {
	DialogStack = []*Dialog{}
}

func ClearDialogsOpen() {
	DialogsOpen = []*Dialog{}
}
