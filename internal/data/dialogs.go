package data

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

var (
	DialogStack     []*Dialog
	DialogsOpen     []*Dialog
	DialogStackOpen bool
	Dialogs         = map[string]*Dialog{}
)

type Dialog struct {
	Key          string
	Pos          pixel.Vec
	ViewPort     *viewport.ViewPort
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	Buttons      []*Button
	Sprites      []*SprElement
	NoBorder     bool

	Open  bool
	Click bool
	Lock  bool
}

type DialogConstructor struct {
	Key      string
	Width    float64
	Height   float64
	Pos      pixel.Vec
	Elements []ElementConstructor
	NoBorder bool
}

func NewDialog(dc *DialogConstructor) {
	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, dc.Width*world.TileSize, dc.Height*world.TileSize))
	vp.CamPos = pixel.V(0, 0)
	vp.PortPos = viewport.MainCamera.PostCamPos.Add(dc.Pos)

	dlg := &Dialog{
		Key:      dc.Key,
		Pos:      dc.Pos,
		ViewPort: vp,
		NoBorder: dc.NoBorder,
	}

	if !dc.NoBorder {
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

		dlg.BorderVP = bvp
		dlg.BorderObject = bObj
		dlg.BorderEntity = be
	}

	var buttons []*Button
	var sprites []*SprElement
	for _, element := range dc.Elements {
		obj := object.New()
		obj.Pos = element.Position
		obj.Layer = 99
		obj.SetRect(img.Batchers[constants.UIBatch].GetSprite(element.SprKey).Frame())
		spr := img.NewSprite(element.SprKey, constants.UIBatch)
		cSpr := img.NewSprite(element.ClickSprKey, constants.UIBatch)
		e := myecs.Manager.NewEntity()
		e.AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr)
		switch element.Element {
		case ButtonElement:
			b := &Button{
				Key:      element.SprKey,
				Sprite:   spr,
				ClickSpr: cSpr,
				HelpText: element.HelpText,
				Object:   obj,
				Entity:   e,
			}
			e.AddComponent(myecs.Update, NewHoverClickFn(MainInput, vp, func(hvc *HoverClick) {
				if dlg.Open && !dlg.Lock {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						dlg.Click = true
					}
					if hvc.Hover && click.Pressed() && dlg.Click {
						e.AddComponent(myecs.Drawable, cSpr)
					} else {
						if hvc.Hover && click.JustReleased() && dlg.Click {
							dlg.Click = false
							if b.OnClick != nil {
								if b.Delay > 0. {
									dlg.Lock = true
									e := myecs.Manager.NewEntity()
									e.AddComponent(myecs.Update, NewTimerFunc(func() bool {
										b.OnClick()
										dlg.Lock = false
										myecs.Manager.DisposeEntity(e)
										return false
									}, b.Delay))
								} else {
									b.OnClick()
								}
							}
						} else if !click.Pressed() && !click.JustReleased() && dlg.Click {
							dlg.Click = false
							e.AddComponent(myecs.Drawable, spr)
						} else {
							e.AddComponent(myecs.Drawable, spr)
						}
					}
				}
			}))
			buttons = append(buttons, b)
		case SpriteElement:
			s := &SprElement{
				Key:    element.SprKey,
				Sprite: spr,
				Object: obj,
				Entity: e,
			}
			sprites = append(sprites, s)
		}
	}
	dlg.Buttons = buttons
	dlg.Sprites = sprites

	Dialogs[dc.Key] = dlg
}

func ClearDialogStack() {
	DialogStack = []*Dialog{}
}

func ClearDialogsOpen() {
	DialogsOpen = []*Dialog{}
}

func CloseDialog(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: CloseDialog: %s not registered\n", key)
		return
	}
	dialog.Open = false
	index := -1
	stack := false
	for i, d := range DialogsOpen {
		if d.Key == key {
			index = i
			break
		}
	}
	for i, d := range DialogStack {
		if d.Key == key {
			index = i
			stack = true
			break
		}
	}
	if index == -1 {
		fmt.Printf("Warning: CloseDialog: %s not open\n", key)
		return
	} else {
		if stack {
			if len(DialogStack) == 1 {
				ClearDialogStack()
			} else {
				DialogStack = append(DialogStack[:index], DialogStack[index+1:]...)
			}
		} else {
			if len(DialogsOpen) == 1 {
				ClearDialogsOpen()
			} else {
				DialogsOpen = append(DialogsOpen[:index], DialogsOpen[index+1:]...)
			}
		}
	}
}
