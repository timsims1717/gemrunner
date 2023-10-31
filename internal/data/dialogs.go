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
	DialogStack []*dialog
	DialogOpen  bool
	Dialogs     = map[string]*dialog{}

	OpenPuzzleConstructor *DialogConstructor
)

type dialog struct {
	Key          string
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
	Buttons []ButtonConstructor
}

func NewDialog(dc *DialogConstructor) {
	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, dc.Width*world.TileSize, dc.Height*world.TileSize))
	vp.CamPos = pixel.V(0, 0)
	vp.PortPos = viewport.MainCamera.PostCamPos

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (dc.Width+1)*world.TileSize, (dc.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = viewport.MainCamera.PostCamPos

	bObj := object.New()
	bObj.Layer = 99
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, bObj).
		AddComponent(myecs.Border, &Border{
			Width:  int(dc.Width),
			Height: int(dc.Height),
		})

	dlg := &dialog{
		Key:          dc.Key,
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
			if dlg.Open && hvc.Input.Get("click").JustPressed() {
				if !dlg.Click && b.OnClick != nil && hvc.Hover {
					click := hvc.Input.Get("click")
					if click.JustPressed() {
						dlg.Click = true
						btnE.AddComponent(myecs.Drawable, cSpr)
						e := myecs.Manager.NewEntity()
						e.AddComponent(myecs.Update, NewTimerFunc(func() bool {
							dlg.Click = false
							btnE.AddComponent(myecs.Drawable, spr)
							b.OnClick()
							myecs.Manager.DisposeEntity(e)
							return false
						}, 0.2))
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
	DialogStack = []*dialog{}
}
