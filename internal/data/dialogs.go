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
	OnCloseSpc   func()

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
		case CheckboxElement:
			x := CreateCheckboxElement(element, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, x)
		case ContainerElement:
			ct2 := CreateContainer(element, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, ct2)
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

func CreateButtonElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Button {
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
				if b.OnHeld != nil {
					b.OnHeld(hvc)
				}
			} else {
				if hvc.Hover && click.JustReleased() && dlg.Click {
					dlg.Click = false
					if b.OnClick != nil {
						if b.Delay > 0. {
							dlg.Lock = true
							entity := myecs.Manager.NewEntity()
							entity.AddComponent(myecs.Update, NewTimerFunc(func() bool {
								MenuInput.Get("click").Consume()
								MenuInput.Get("rClick").Consume()
								b.OnClick()
								dlg.Lock = false
								myecs.Manager.DisposeEntity(entity)
								return false
							}, b.Delay))
						} else {
							MenuInput.Get("click").Consume()
							MenuInput.Get("rClick").Consume()
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

func CreateCheckboxElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Checkbox {
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(img.Batchers[constants.UIBatch].GetSprite(element.SprKey).Frame())
	spr := img.NewSprite(element.SprKey, constants.UIBatch)
	cSpr := img.NewSprite(element.ClickSprKey, constants.UIBatch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)
	x := &Checkbox{
		Key:      element.Key,
		Sprite:   spr,
		CheckSpr: cSpr,
		HelpText: element.HelpText,
		Object:   obj,
		Entity:   e,
	}
	e.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, vp, func(hvc *HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock && !dlg.Click {
			click := hvc.Input.Get("click")
			if hvc.Hover && click.JustPressed() {
				SetChecked(x, !x.Checked)
			}
		}
	}))
	return x
}

func SetChecked(x *Checkbox, c bool) {
	x.Checked = c
	if x.Checked {
		x.Entity.AddComponent(myecs.Drawable, x.CheckSpr)
	} else {
		x.Entity.AddComponent(myecs.Drawable, x.Sprite)
	}
}

func CreateContainer(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Container {
	ctvp := viewport.New(nil)
	ctvp.ParentView = vp
	ctvp.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	ctvp.CamPos = pixel.V(0, 0)
	ctvp.PortPos = element.Position

	vpObj := object.New()
	vpObj.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	vpObj.SetPos(element.Position)
	vpObj.Layer = 99

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = element.Position

	bObj := object.New()
	bObj.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	bObj.Layer = 99
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, bObj).
		AddComponent(myecs.Border, &Border{
			Width:  int(element.Width),
			Height: int(element.Height),
			Style:  ThinBorder,
		})

	e := myecs.Manager.NewEntity().AddComponent(myecs.Object, vpObj)
	ct := &Container{
		Key:          element.Key,
		BorderVP:     bvp,
		BorderObject: bObj,
		Object:       vpObj,
		BorderEntity: be,
		Entity:       e,
		ViewPort:     ctvp,
	}
	for _, ele := range element.SubElements {
		if ele.Key == "" {
			fmt.Println("WARNING: element constructor has no key")
		}
		switch ele.Element {
		case ButtonElement:
			b := CreateButtonElement(ele, dlg, ct.ViewPort)
			ct.Elements = append(ct.Elements, b)
		case CheckboxElement:
			x := CreateCheckboxElement(ele, dlg, ct.ViewPort)
			ct.Elements = append(ct.Elements, x)
		case ContainerElement:
			ct2 := CreateContainer(ele, dlg, ct.ViewPort)
			ct.Elements = append(ct.Elements, ct2)
		case InputElement:
			i := CreateInputElement(ele, dlg, ct.ViewPort)
			ct.Elements = append(ct.Elements, i)
		case ScrollElement:
			s := CreateScrollElement(ele, dlg, ct.ViewPort)
			ct.Elements = append(ct.Elements, s)
		case SpriteElement:
			s := CreateSpriteElement(ele)
			ct.Elements = append(ct.Elements, s)
		case TextElement:
			t := CreateTextElement(ele, ct.ViewPort)
			ct.Elements = append(ct.Elements, t)
		}
	}
	return ct
}

func CreateInputElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Input {
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
	bObj.SetRect(pixel.R(0, 0, (element.Width)*world.TileSize, (element.Height)*world.TileSize))
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
				if i.NumbersOnly {
					typed = util.OnlyNumbers(typed)
				} else {
					typed = util.OnlyAlphaNumeric(typed)
				}
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

func ChangeText(input *Input, rt string) {
	input.Value = rt
	input.Text.SetText(input.Value)
	input.CaretIndex = input.Text.Len() - 1
	input.CaretObj.Pos = input.Text.GetDotPos(input.CaretIndex)
	input.CaretObj.Pos.Y = 0
	input.CaretObj.Hidden = false
}

func CreateScrollElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Scroll {
	svp := viewport.New(nil)
	svp.ParentView = vp
	svp.SetRect(pixel.R(0, 0, (element.Width-1)*world.TileSize, element.Height*world.TileSize))
	svp.CamPos = pixel.V(0, 0)
	svp.PortPos = element.Position
	svp.PortPos.X -= world.HalfSize

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, (element.Width+1)*world.TileSize, (element.Height+1)*world.TileSize))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = element.Position

	vpObj := object.New()
	vpObj.SetRect(pixel.R(0, 0, element.Width*world.TileSize, element.Height*world.TileSize))
	vpObj.SetPos(element.Position)
	vpObj.Layer = 99

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

	e := myecs.Manager.NewEntity().AddComponent(myecs.Object, vpObj)
	s := &Scroll{
		Key:          element.Key,
		BorderVP:     bvp,
		BorderObject: bObj,
		Object:       vpObj,
		BorderEntity: be,
		Entity:       e,
		ViewPort:     svp,
	}
	e.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, svp, func(hvc *HoverClick) {
		if hvc.ViewHover {
			if hvc.Input.ScrollV > 0. {
				s.ViewPort.CamPos.Y += constants.ScrollSpeed * timing.DT
			} else if hvc.Input.ScrollV < 0. {
				s.ViewPort.CamPos.Y -= constants.ScrollSpeed * timing.DT
			}
			RestrictScroll(s)
			AlignBarToView(s)
		}
	}))
	btnX := svp.Rect.W() * 0.5
	for i := 0; i < 3; i++ {
		var pos pixel.Vec
		var key, sprKey, cSprKey string
		switch i {
		case 0:
			pos = element.Position.Add(pixel.V(btnX, svp.Rect.H()*0.5))
			key = fmt.Sprintf("%s_scroll_up", element.Key)
			sprKey = "scroll_up"
			cSprKey = "scroll_up_click"
		case 1:
			pos = element.Position.Add(pixel.V(btnX, svp.Rect.H()*-0.5))
			key = fmt.Sprintf("%s_scroll_down", element.Key)
			sprKey = "scroll_down"
			cSprKey = "scroll_down_click"
		case 2:
			pos = element.Position.Add(pixel.V(btnX, 0.))
			key = fmt.Sprintf("%s_scroll_bar", element.Key)
			sprKey = "scroll_bar"
			cSprKey = "scroll_bar_click"
		}
		btn := ElementConstructor{
			Key:         key,
			SprKey:      sprKey,
			ClickSprKey: cSprKey,
			Position:    pos,
			Element:     ButtonElement,
		}
		b := CreateButtonElement(btn, dlg, dlg.ViewPort)
		dlg.Elements = append(dlg.Elements, b)
		switch i {
		case 0:
			s.ButtonHeight = b.Object.Rect.H()
			b.Object.Pos.Y -= b.Object.Rect.H() * 0.5
			b.OnHeld = func(_ *HoverClick) {
				s.ViewPort.CamPos.Y += constants.ScrollSpeed * timing.DT
				RestrictScroll(s)
				AlignBarToView(s)
			}
		case 1:
			b.Object.Pos.Y += b.Object.Rect.H() * 0.5
			b.OnHeld = func(_ *HoverClick) {
				s.ViewPort.CamPos.Y -= constants.ScrollSpeed * timing.DT
				RestrictScroll(s)
				AlignBarToView(s)
			}
		case 2:
			s.Bar = b
			offset := 0.
			barClick := false
			b.Entity.AddComponent(myecs.Update, NewHoverClickFn(MenuInput, dlg.ViewPort, func(hvc *HoverClick) {
				if dlg.Open && dlg.Active && !dlg.Lock {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						b.Entity.AddComponent(myecs.Drawable, b.ClickSpr)
						offset = hvc.Pos.Y - b.Object.Pos.Y
						barClick = true
					}
					if click.Pressed() && barClick {
						b.Object.Pos.Y = hvc.Pos.Y - offset
						RestrictScroll(s)
						AlignViewToBar(s)
					} else {
						barClick = false
						b.Entity.AddComponent(myecs.Drawable, b.Sprite)
					}
				}
			}))
		}
	}
	for _, ele := range element.SubElements {
		if ele.Key == "" {
			fmt.Println("WARNING: element constructor has no key")
		}
		switch ele.Element {
		case ButtonElement:
			b := CreateButtonElement(ele, dlg, s.ViewPort)
			s.Elements = append(s.Elements, b)
		case CheckboxElement:
			x := CreateCheckboxElement(ele, dlg, s.ViewPort)
			s.Elements = append(s.Elements, x)
		case ContainerElement:
			ct := CreateContainer(ele, dlg, s.ViewPort)
			s.Elements = append(s.Elements, ct)
		case InputElement:
			i := CreateInputElement(ele, dlg, s.ViewPort)
			s.Elements = append(s.Elements, i)
		case ScrollElement:
			s2 := CreateScrollElement(ele, dlg, s.ViewPort)
			s.Elements = append(s.Elements, s2)
		case SpriteElement:
			s2 := CreateSpriteElement(ele)
			s.Elements = append(s.Elements, s2)
		case TextElement:
			t := CreateTextElement(ele, s.ViewPort)
			s.Elements = append(s.Elements, t)
		}
	}
	UpdateScrollBounds(s)
	return s
}

func AlignViewToBar(s *Scroll) {
	barHeight := s.ViewPort.Rect.H() - s.ButtonHeight*2 - s.Bar.Object.Rect.H()
	viewHeight := s.ViewPort.Rect.H()
	if math.Abs(s.YTop-s.YBot) < viewHeight {
		s.Bar.Object.Pos.Y = s.ViewPort.PortPos.Y + s.ViewPort.Rect.H()*0.5 - s.ButtonHeight - s.Bar.Object.Rect.H()*0.5
		return
	}
	scrollHeight := s.YTop - s.YBot - s.ViewPort.Rect.H()
	barTop := s.ViewPort.PortPos.Y + s.ViewPort.Rect.H()*0.5 - s.ButtonHeight - s.Bar.Object.Rect.H()*0.5
	barPos := s.Bar.Object.Pos.Y
	barDist := barTop - barPos
	barRatio := barDist / barHeight
	scrollDist := barRatio * scrollHeight
	s.ViewPort.CamPos.Y = s.YTop - s.ViewPort.Rect.H()*0.5 - scrollDist
}

func AlignBarToView(s *Scroll) {
	barHeight := s.ViewPort.Rect.H() - s.ButtonHeight*2 - s.Bar.Object.Rect.H()
	viewHeight := s.ViewPort.Rect.H()
	if math.Abs(s.YTop-s.YBot) < viewHeight {
		s.Bar.Object.Pos.Y = s.ViewPort.PortPos.Y + s.ViewPort.Rect.H()*0.5 - s.ButtonHeight - s.Bar.Object.Rect.H()*0.5
		return
	}
	scrollHeight := s.YTop - s.YBot - s.ViewPort.Rect.H()
	scrollTop := s.YTop - s.ViewPort.Rect.H()*0.5
	viewPos := s.ViewPort.CamPos.Y
	scrollDist := scrollTop - viewPos
	scrollRatio := scrollDist / scrollHeight
	barDist := scrollRatio * barHeight
	barTop := s.ViewPort.PortPos.Y + s.ViewPort.Rect.H()*0.5 - s.ButtonHeight - s.Bar.Object.Rect.H()*0.5
	s.Bar.Object.Pos.Y = barTop - barDist

	//barPos := s.Bar.Object.Pos.Y
	//barDist := barTop - barPos
	//barRatio := barDist / barHeight
	//scrollDist := barRatio * scrollHeight
	//s.ViewPort.CamPos.Y = s.YTop - s.ViewPort.Rect.H()*0.5 - scrollDist
}

func RestrictScroll(s *Scroll) {
	if s.Bar.Object.Pos.Y > s.ViewPort.PortPos.Y+s.ViewPort.Rect.H()*0.5-s.ButtonHeight-s.Bar.Object.Rect.H()*0.5 {
		s.Bar.Object.Pos.Y = s.ViewPort.PortPos.Y + s.ViewPort.Rect.H()*0.5 - s.ButtonHeight - s.Bar.Object.Rect.H()*0.5
	}
	if s.Bar.Object.Pos.Y < s.ViewPort.PortPos.Y-s.ViewPort.Rect.H()*0.5+s.ButtonHeight+s.Bar.Object.Rect.H()*0.5 {
		s.Bar.Object.Pos.Y = s.ViewPort.PortPos.Y - s.ViewPort.Rect.H()*0.5 + s.ButtonHeight + s.Bar.Object.Rect.H()*0.5
	}
	if s.ViewPort.CamPos.Y-s.ViewPort.Rect.H()*0.5 < s.YBot {
		s.ViewPort.CamPos.Y = s.YBot + s.ViewPort.Rect.H()*0.5
	}
	if s.ViewPort.CamPos.Y+s.ViewPort.Rect.H()*0.5 > s.YTop {
		s.ViewPort.CamPos.Y = s.YTop - s.ViewPort.Rect.H()*0.5
	}
}

func UpdateScrollBounds(scroll *Scroll) {
	yTop := 0.
	yBot := 0.
	for i, ele := range scroll.Elements {
		obj := object.New()
		if spr, okS := ele.(*SprElement); okS {
			obj.Rect = spr.Object.Rect
			obj.Pos = spr.Object.Pos
		} else if btn, okB := ele.(*Button); okB {
			obj.Rect = btn.Object.Rect
			obj.Pos = btn.Object.Pos
		} else if x, okX := ele.(*Checkbox); okX {
			obj.Rect = x.Object.Rect
			obj.Pos = x.Object.Pos
		} else if ct, okC := ele.(*Container); okC {
			obj.Rect = ct.ViewPort.Canvas.Bounds()
			obj.Pos = ct.ViewPort.PortPos
		} else if in, okI := ele.(*Input); okI {
			obj.Rect = in.Text.Obj.Rect
			obj.Pos = in.Text.Obj.Pos
		} else if txt, okT := ele.(*Text); okT {
			obj.Rect = txt.Text.Obj.Rect
			obj.Pos = txt.Text.Obj.Pos
		} else if scr, okScr := ele.(*Scroll); okScr {
			obj.Rect = scr.BorderObject.Rect
			obj.Pos = scr.ViewPort.PortPos
		}
		oTop := obj.Pos.Y + obj.Rect.H()*0.5 + 1
		oBot := obj.Pos.Y - obj.Rect.H()*0.5 - 1
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
	scroll.Bar.Object.Pos.Y = scroll.ViewPort.PortPos.Y + scroll.ViewPort.Rect.H()*0.5 - scroll.ButtonHeight - scroll.Bar.Object.Rect.H()*0.5
}

func CreateSpriteElement(element ElementConstructor) *SprElement {
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

func CreateTextElement(element ElementConstructor, vp *viewport.ViewPort) *Text {
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

func SetCloseSpcFn(key string, fn func()) {
	dialog, ok := Dialogs[key]
	if !ok {
		fmt.Printf("Warning: SetCloseSpcFn: %s not registered\n", key)
		return
	}
	dialog.OnCloseSpc = fn
}

func SetOnClick(dlgKey, btnKey string, fn func()) {
	dialog, ok := Dialogs[dlgKey]
	if !ok {
		fmt.Printf("Warning: SetOnClick: Dialog %s not registered\n", dlgKey)
		return
	}
	for _, ele := range dialog.Elements {
		if btn, okB := ele.(*Button); okB && btn.Key == btnKey {
			btn.OnClick = fn
			return
		}
	}
	fmt.Printf("Warning: SetOnClick: Button %s not registered in Dialog %s\n", btnKey, dlgKey)
}

func SetTempOnClick(dlgKey, btnKey string, fn func()) {
	dialog, ok := Dialogs[dlgKey]
	if !ok {
		fmt.Printf("Warning: SetTempOnClick: Dialog %s not registered\n", dlgKey)
		return
	}
	for _, ele := range dialog.Elements {
		if btn, okB := ele.(*Button); okB && btn.Key == btnKey {
			oldFn := btn.OnClick
			btn.OnClick = func() {
				fn()
				btn.OnClick = oldFn
				oldFn()
			}
			return
		}
	}
	fmt.Printf("Warning: SetTempOnClick: Button %s not registered in Dialog %s\n", btnKey, dlgKey)
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
		if dialog.OnCloseSpc != nil {
			dialog.OnCloseSpc()
			dialog.OnCloseSpc = nil
		}
	}
}
