package data

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
	"math"
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
	Elements     []interface{}
	NoBorder     bool
	OnOpen       func()
	OnClose      func()

	Open   bool
	Active bool
	Click  bool
	Lock   bool
	Layer  int
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

	for _, element := range dc.Elements {
		if element.Key == "" {
			fmt.Println("WARNING: element constructor has no key")
		}
		switch element.Element {
		case ButtonElement:
			b := CreateButtonElement(element, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, b)
		case InputElement:
			i := CreateInputElement(element, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, i)
		case ScrollElement:
			s := CreateScrollElement(element, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, s)
		case SpriteElement:
			s := CreateSpriteElement(element)
			dlg.Elements = append(dlg.Elements, s)
		case TextElement:
			t := CreateTextElement(element, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, t)
		}
	}

	Dialogs[dc.Key] = dlg
}

func CreateButtonElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) interface{} {
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(img.Batchers[constants.UIBatch].GetSprite(element.SprKey).Frame())
	spr := img.NewSprite(element.SprKey, constants.UIBatch)
	cSpr := img.NewSprite(element.ClickSprKey, constants.UIBatch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)
	b := &Button{
		Key:      element.Key,
		Sprite:   spr,
		ClickSpr: cSpr,
		HelpText: element.HelpText,
		Object:   obj,
		Entity:   e,
	}
	e.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, vp, func(hvc *HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock {
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
							entity := myecs.Manager.NewEntity()
							entity.AddComponent(myecs.Update, NewTimerFunc(func() bool {
								b.OnClick()
								dlg.Lock = false
								myecs.Manager.DisposeEntity(entity)
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
	return b
}

func CreateInputElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) interface{} {
	ivp := viewport.New(nil)
	ivp.ParentView = vp
	ivp.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	ivp.CamPos = pixel.V(ivp.Rect.W()*0.5-2, 0)
	ivp.PortPos = element.Position

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = element.Position

	bObj := object.New()
	bObj.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bObj.Layer = 99
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, bObj).
		AddComponent(myecs.Border, &Border{
			Width:  int(element.Width),
			Height: int(element.Height),
			Style:  ThinBorder,
		})

	tf := typeface.New("main", typeface.NewAlign(typeface.Left, typeface.Top), 1, 0.0625, 0, 0)
	tf.SetPos(pixel.ZV)
	tf.SetColor(pixel.ToRGBA(constants.ColorWhite))
	tf.SetText(element.Text)
	te := myecs.Manager.NewEntity()
	te.AddComponent(myecs.Object, tf.Obj)
	te.AddComponent(myecs.Drawable, tf)
	te.AddComponent(myecs.DrawTarget, ivp.Canvas)

	cObj := object.New()
	cObj.Pos = tf.GetEndPos()
	cObj.SetRect(img.Batchers[constants.UIBatch].GetSprite(constants.TextCaret).Frame())
	cSpr := img.NewSprite(constants.TextCaret, constants.UIBatch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, cObj)
	e.AddComponent(myecs.Drawable, cSpr)

	i := &Input{
		Key:          element.Key,
		Value:        element.Text,
		Text:         tf,
		TextEntity:   te,
		CaretObj:     cObj,
		CaretSpr:     cSpr,
		CaretIndex:   len(element.Text),
		BorderVP:     bvp,
		BorderObject: bObj,
		BorderEntity: be,
		ViewPort:     ivp,
		Entity:       e,
	}

	flashTimer := timing.New(0.53)
	e.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, ivp, func(hvc *HoverClick) {
		flashTimer.Update()
		wasActive := i.Active
		click := hvc.Input.Get("click")
		if dlg.Open && dlg.Active && !dlg.Lock {
			if click.JustPressed() {
				i.Active = hvc.ViewHover
				if hvc.ViewHover && !wasActive {
					click.Consume()
				}
			}
		} else {
			i.Active = false
		}
		if !wasActive && i.Active {
			flashTimer.Reset()
			cObj.Hidden = false
		}
		if i.Active {
			changed := false
			ci := i.CaretIndex
			left := hvc.Input.Get("left")
			right := hvc.Input.Get("right")
			if hvc.Input.Get("home").JustPressed() {
				i.CaretIndex = 0
			} else if hvc.Input.Get("end").JustPressed() {
				i.CaretIndex = tf.Len() - 1
			} else if left.JustPressed() || left.Repeated() {
				i.CaretIndex--
			} else if right.JustPressed() || right.Repeated() {
				i.CaretIndex++
			} else if click.JustPressed() {
				closest := 0
				dist := -1.
				for j := 0; j <= tf.Len(); j++ {
					d := math.Abs(util.Magnitude(tf.GetDotPos(j).Sub(hvc.Pos)))
					if dist == -1. || d < dist {
						dist = d
						closest = j
					}
				}
				i.CaretIndex = closest
			}
			if i.CaretIndex < 0 {
				i.CaretIndex = 0
			} else if i.CaretIndex > tf.Len()-1 {
				i.CaretIndex = tf.Len() - 1
			}
			back := hvc.Input.Get("backspace")
			if (back.JustPressed() || back.Repeated()) && i.CaretIndex > 0 {
				i.Value = fmt.Sprintf("%s%s", i.Value[:i.CaretIndex-1], i.Value[i.CaretIndex:])
				changed = true
				i.CaretIndex--
			}
			del := hvc.Input.Get("delete")
			if (del.JustPressed() || del.Repeated()) && i.CaretIndex < tf.Len()-1 {
				i.Value = fmt.Sprintf("%s%s", i.Value[:i.CaretIndex], i.Value[i.CaretIndex+1:])
				changed = true
			}
			typed := hvc.Input.Typed
			if typed != "" {
				i.Value = fmt.Sprintf("%s%s%s", i.Value[:i.CaretIndex], typed, i.Value[i.CaretIndex:])
				changed = true
				i.CaretIndex += len(typed)
			}
			if changed {
				tf.SetText(i.Value)
			}
			if ci != i.CaretIndex || changed {
				cObj.Pos = tf.GetDotPos(i.CaretIndex)
				cObj.Pos.Y = 0
				flashTimer.Reset()
				cObj.Hidden = false
			}
			if flashTimer.Done() {
				cObj.Hidden = !cObj.Hidden
				flashTimer.Reset()
			}
		} else {
			cObj.Hidden = true
		}
	}))

	return i
}

func CreateScrollElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) interface{} {

	svp := viewport.New(nil)
	svp.ParentView = vp
	svp.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	svp.CamPos = pixel.V(0, 0)
	svp.PortPos = element.Position

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = element.Position

	bObj := object.New()
	bObj.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bObj.Layer = 99
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, bObj).
		AddComponent(myecs.Border, &Border{
			Width:  int(element.Width),
			Height: int(element.Height),
			Style:  ThinBorder,
		})

	e := myecs.Manager.NewEntity().AddComponent(myecs.Object, bObj)
	s := &Scroll{
		Key:          element.Key,
		UpSprite:     nil,
		UpSprClick:   nil,
		UpObject:     nil,
		UpEntity:     nil,
		DwnSprite:    nil,
		DwnSprClick:  nil,
		DwnObject:    nil,
		DwnEntity:    nil,
		BarSprite:    nil,
		BarObject:    nil,
		BarEntity:    nil,
		BorderVP:     bvp,
		BorderObject: bObj,
		BorderEntity: be,
		Entity:       e,
		ViewPort:     svp,
		Elements:     nil,
	}
	e.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, svp, func(hvc *HoverClick) {

		if s.ViewPort.CamPos.Y-s.ViewPort.Rect.H()*0.5 < s.YBot {
			s.ViewPort.CamPos.Y = s.YBot + s.ViewPort.Rect.H()*0.5
		}
		if s.ViewPort.CamPos.Y+s.ViewPort.Rect.H()*0.5 > s.YTop {
			s.ViewPort.CamPos.Y = s.YTop - s.ViewPort.Rect.H()*0.5
		}
	}))
	UpdateScroll(s)
	return s
}

func UpdateScroll(scroll *Scroll) {
	yTop := 0.
	yBot := 0.
	for i, ele := range scroll.Elements {
		var obj *object.Object
		if spr, okS := ele.(*SprElement); okS {
			obj = spr.Object
		} else if btn, okB := ele.(*Button); okB {
			obj = btn.Object
		} else if txt, okT := ele.(*Text); okT {
			obj = txt.Text.Obj
		} else if scr, okScr := ele.(*Scroll); okScr {
			obj = scr.BorderObject
		}
		oTop := obj.Pos.Y + obj.Rect.H()*0.5
		oBot := obj.Pos.Y - obj.Rect.H()*0.5
		if i == 0 || yTop < oTop {
			yTop = oTop
		}
		if i == 0 || yBot > oBot {
			yBot = oBot
		}
	}
	scroll.YTop = yTop
	scroll.YBot = yBot
	scroll.ViewPort.CamPos.Y = scroll.YTop - scroll.ViewPort.Rect.H()*0.5
}

func CreateSpriteElement(element ElementConstructor) interface{} {
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(img.Batchers[constants.UIBatch].GetSprite(element.SprKey).Frame())
	spr := img.NewSprite(element.SprKey, constants.UIBatch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)
	s := &SprElement{
		Key:    element.Key,
		Sprite: spr,
		Object: obj,
		Entity: e,
	}
	return s
}

func CreateTextElement(element ElementConstructor, vp *viewport.ViewPort) interface{} {
	tf := typeface.New("main", typeface.NewAlign(typeface.Left, typeface.Top), 1, 0.0625, 0, 0)
	tf.SetPos(element.Position)
	tf.SetColor(pixel.ToRGBA(constants.ColorWhite))
	tf.SetText(element.Text)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, tf.Obj)
	e.AddComponent(myecs.Drawable, tf)
	e.AddComponent(myecs.DrawTarget, vp.Canvas)
	t := &Text{
		Key:    element.Key,
		Text:   tf,
		Entity: e,
	}
	return t
}

func ClearDialogStack() {
	DialogStack = []*Dialog{}
}

func ClearDialogsOpen() {
	DialogsOpen = []*Dialog{}
}

func OpenDialog(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	DialogsOpen = append(DialogsOpen, dialog)
	if dialog.OnOpen != nil {
		dialog.OnOpen()
	}
}

func OpenDialogInStack(key string) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: OpenDialog: %s not registered\n", key)
		return
	}
	dialog.Open = true
	DialogStack = append(DialogStack, dialog)
	if dialog.OnOpen != nil {
		dialog.OnOpen()
	}
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
		if dialog.OnClose != nil {
			dialog.OnClose()
		}
	}
}
