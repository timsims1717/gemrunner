package ui

import (
	"fmt"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"math"
)

func CreateButtonElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Element {
	r := img.Batchers[element.Batch].GetSprite(element.SprKey).Frame()
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(r)
	spr := img.NewSprite(element.SprKey, element.Batch)
	cSpr := img.NewSprite(element.SprKey2, element.Batch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)

	bord := &Border{
		Rect:   pixel.R(0, 0, r.W(), r.H()),
		Style:  ThinBorderWhite,
		Hidden: true,
	}

	b := &Element{
		Key:         element.Key,
		Sprite:      spr,
		Sprite2:     cSpr,
		HelpText:    element.HelpText,
		Object:      obj,
		Entity:      e,
		ElementType: ButtonElement,
		Border:      bord,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	e.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, vp, func(hvc *data.HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock {
			click := hvc.Input.Get("click")
			if hvc.Hover && click.JustPressed() {
				dlg.Click = true
			}
			if hvc.Hover && click.Pressed() && dlg.Click {
				e.AddComponent(myecs.Drawable, cSpr)
				if b.OnHold != nil {
					b.OnHold()
				}
			} else {
				if hvc.Hover && click.JustReleased() && dlg.Click {
					dlg.Click = false
					if b.OnClick != nil {
						if b.Delay > 0. {
							dlg.Lock = true
							entity := myecs.Manager.NewEntity()
							entity.AddComponent(myecs.Update, data.NewTimerFunc(func() bool {
								hvc.Input.Get("click").Consume()
								hvc.Input.Get("rClick").Consume()
								b.OnClick()
								dlg.Lock = false
								myecs.Manager.DisposeEntity(entity)
								return false
							}, b.Delay))
						} else {
							hvc.Input.Get("click").Consume()
							hvc.Input.Get("rClick").Consume()
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

func CreateCheckboxElement(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Element {
	r := img.Batchers[element.Batch].GetSprite(element.SprKey).Frame()
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(r)
	spr := img.NewSprite(element.SprKey, element.Batch)
	cSpr := img.NewSprite(element.SprKey2, element.Batch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)

	bord := &Border{
		Rect:   pixel.R(0, 0, r.W(), r.H()),
		Style:  ThinBorderWhite,
		Hidden: true,
	}

	x := &Element{
		Key:         element.Key,
		Sprite:      spr,
		Sprite2:     cSpr,
		HelpText:    element.HelpText,
		Object:      obj,
		Border:      bord,
		Entity:      e,
		ElementType: CheckboxElement,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	e.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, vp, func(hvc *data.HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock && !dlg.Click {
			click := hvc.Input.Get("click")
			if hvc.Hover && click.JustPressed() {
				SetChecked(x, !x.Checked)
			}
		}
	}))
	return x
}

func SetChecked(x *Element, c bool) {
	x.Checked = c
	if x.Checked {
		x.Entity.AddComponent(myecs.Drawable, x.Sprite2)
	} else {
		x.Entity.AddComponent(myecs.Drawable, x.Sprite)
	}
}

func CreateContainer(element ElementConstructor, dlg *Dialog, vp *viewport.ViewPort) *Element {
	ctvp := viewport.New(nil)
	ctvp.ParentView = vp
	ctvp.SetRect(pixel.R(0, 0, element.Width, element.Height))
	ctvp.CamPos = pixel.V(0, 0)
	ctvp.PortPos = element.Position

	vpObj := object.New()
	vpObj.SetRect(pixel.R(0, 0, element.Width+1, element.Height+1))
	vpObj.SetPos(element.Position)
	vpObj.Layer = 99

	bord := &Border{
		Rect:  pixel.R(0, 0, element.Width, element.Height),
		Style: ThinBorder,
	}

	e := myecs.Manager.NewEntity().AddComponent(myecs.Object, vpObj)
	ct := &Element{
		Key:         element.Key,
		Border:      bord,
		Object:      vpObj,
		Entity:      e,
		ViewPort:    ctvp,
		ElementType: ContainerElement,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	for _, ele := range element.SubElements {
		if ele.Key == "" {
			fmt.Println("WARNING: element constructor has no key")
		}
		switch ele.ElementType {
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
			i := CreateInputElement(ele, dlg, ct, ct.ViewPort, false)
			ct.Elements = append(ct.Elements, i)
		case MultiLineInputElement:
			i := CreateInputElement(ele, dlg, ct, ct.ViewPort, true)
			ct.Elements = append(ct.Elements, i)
		case ScrollElement:
			s := CreateScrollElement(ele, dlg, ct, ct.ViewPort)
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

func CreateInputElement(element ElementConstructor, dlg *Dialog, parent *Element, vp *viewport.ViewPort, multiline bool) *Element {
	ivp := viewport.New(nil)
	ivp.ParentView = vp
	ivp.SetRect(pixel.R(0, 0, element.Width, element.Height))
	ivp.CamPos = pixel.V(ivp.Rect.W()*0.5-2, ivp.Rect.H()*-0.5+8)
	ivp.PortPos = element.Position

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, element.Width+1, element.Height+1))
	bvp.CamPos = pixel.ZV
	bvp.PortPos = element.Position

	vpObj := object.New()
	vpObj.SetRect(pixel.R(0, 0, element.Width+1, element.Height+1))
	vpObj.SetPos(element.Position)
	vpObj.Layer = 99

	bord := &Border{
		Rect:  pixel.R(0, 0, element.Width, element.Height),
		Style: ThinBorder,
	}

	tf := typeface.New("main").
		WithAlign(typeface.NewAlign(typeface.Left, typeface.Top)).
		WithScalar(0.0625)
	tf.SetPos(pixel.V(0., 4.))
	if multiline {
		tf.SetWidth(element.Width - world.TileSize*0.75)
	}
	tf.SetColor(pixel.ToRGBA(element.Color))
	//tf.Debug = true
	tf.Update()
	tf.SetText(element.Text)
	te := myecs.Manager.NewEntity()
	te.AddComponent(myecs.Object, tf.Obj)
	te.AddComponent(myecs.Drawable, tf)
	te.AddComponent(myecs.DrawTarget, ivp.Canvas)

	// the cursor
	cObj := object.New()
	cObj.Pos = tf.GetEndPos()
	cObj.SetRect(img.Batchers[element.Batch].GetSprite(element.SprKey).Frame())
	cObj.Offset.Y = -cObj.HalfHeight
	cObj.Hidden = true
	cSpr := img.NewSprite(element.SprKey, element.Batch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, cObj)
	e.AddComponent(myecs.Drawable, cSpr)

	i := &Element{
		Key:         element.Key,
		Value:       element.Text,
		Text:        tf,
		Object:      vpObj,
		CaretObj:    cObj,
		Sprite:      cSpr,
		CaretIndex:  len(element.Text),
		ViewPort:    ivp,
		Border:      bord,
		Entity:      e,
		ElementType: InputElement,
		MultiLine:   multiline,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}

	flashTimer := timing.New(0.53)
	e.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ivp, func(hvc *data.HoverClick) {
		if !dlg.Open {
			i.CaretObj.Hidden = true
			return
		}
		flashTimer.Update()
		wasActive := i.InFocus
		click := hvc.Input.Get("click")
		if dlg.Open && dlg.Active && !dlg.Lock {
			if click.JustPressed() {
				if hvc.ViewHover {
					if !wasActive {
						click.Consume()
					}
					dlg.SetFocus(i, true)
				} else {
					dlg.SetFocus(i, false)
				}
			}
		} else {
			dlg.SetFocus(i, false)
		}
		if !wasActive && i.InFocus {
			flashTimer.Reset()
			i.CaretObj.Hidden = false
		}
		if hvc.ViewHover && i.MultiLine {
			if hvc.Input.ScrollV > 0. {
				i.ViewPort.CamPos.Y += ScrollSpeed * timing.DT
			} else if hvc.Input.ScrollV < 0. {
				i.ViewPort.CamPos.Y -= ScrollSpeed * timing.DT
			}
			RestrictScroll(i)
			AlignBarToView(i)
		}
		if i.InFocus {
			ci := i.CaretIndex
			changed := TextNavigation(i, hvc)
			changed = Typing(i, hvc) || changed
			if changed {
				update := true
				for update {
					tf.SetText(i.Value)
					update = i.MultiLine && UpdateInputScrollBounds(i)
				}
			}
			if ci != i.CaretIndex || changed {
				i.CaretObj.Pos = tf.GetDotPos(i.CaretIndex)
				flashTimer.Reset()
				i.CaretObj.Hidden = false
				if i.MultiLine {
					MoveScrollToInclude(i, i.CaretObj.Pos.Y, i.CaretObj.Pos.Y-i.CaretObj.Rect.H())
					RestrictScroll(i)
					AlignBarToView(i)
				}
				tf.Update()
			}
			if flashTimer.Done() {
				i.CaretObj.Hidden = !i.CaretObj.Hidden
				flashTimer.Reset()
			}
		} else {
			i.CaretObj.Hidden = true
		}
	}))
	if i.MultiLine {
		CreateScrollBars(dlg, i, parent, element)
		UpdateInputScrollBounds(i)
		MoveToScrollTop(i)
	}

	return i
}

func TextNavigation(i *Element, hvc *data.HoverClick) bool {
	changed := false
	// cursor navigation
	left := hvc.Input.Get("left")
	right := hvc.Input.Get("right")
	up := hvc.Input.Get("up")
	down := hvc.Input.Get("down")
	home := hvc.Input.Get("home")
	end := hvc.Input.Get("end")
	if i.MultiLine { // multiline navigation
		w := i.ViewPort.Rect.W() - world.TileSize*0.75
		if i.Text.MaxWidth != w { // check width
			i.Text.SetWidth(w)
			changed = true
		}
		if home.JustPressed() { // go to previous new line
			i.CaretIndex, _ = i.Text.GetStartOfLine(i.CaretIndex)
		} else if end.JustPressed() { // go to next new line
			i.CaretIndex, _ = i.Text.GetEndOfLine(i.CaretIndex)
		} else if up.JustPressed() || up.Repeated() { // go up one line
			j := i.CaretIndex
			pos := i.Text.GetDotPos(i.CaretIndex)
			cy := pos.Y
			cx := pos.X
			x1 := -2000.
			prevLine := false
			for ; j >= 0; j-- {
				if prevLine {
					tx := i.Text.GetDotPos(j).X
					if tx <= cx {
						if x1-cx > cx-tx {
							break
						} else {
							j++
							break
						}
					} else {
						x1 = tx
					}
				} else if cy < i.Text.GetDotPos(j).Y {
					prevLine = true
				}
			}
			i.CaretIndex = j
		} else if down.JustPressed() || down.Repeated() { // go down one line
			j := i.CaretIndex
			pos := i.Text.GetDotPos(i.CaretIndex)
			cy := pos.Y
			cx := pos.X
			x1 := -2000.
			nextLine := false
			for ; j < i.Text.Len(); j++ {
				if !nextLine && cy > i.Text.GetDotPos(j).Y {
					nextLine = true
				}
				if nextLine {
					tx := i.Text.GetDotPos(j).X
					if tx >= cx {
						if cx-x1 > tx-cx {
							break
						} else {
							j--
							break
						}
					} else {
						x1 = tx
					}
				}
			}
			i.CaretIndex = j
		}
	} else { // non-multiline navigation
		if home.JustPressed() || up.JustPressed() {
			i.CaretIndex = 0 // go to beginning
		} else if end.JustPressed() || down.JustPressed() {
			i.CaretIndex = i.Text.Len() - 1 // go to end
		}
	}
	// universal navigation
	if left.JustPressed() || left.Repeated() {
		i.CaretIndex-- // go one left
	} else if right.JustPressed() || right.Repeated() {
		i.CaretIndex++ // go one right
	} else if hvc.Input.Get("click").JustPressed() { // using the mouse
		closest := 0
		dist := -1.
		for j := 0; j < i.Text.Len(); j++ {
			cPos := i.Text.GetDotPos(j).Add(pixel.V(0., i.CaretObj.Offset.Y))
			d := math.Abs(util.Magnitude(cPos.Sub(hvc.Pos)))
			if dist == -1. || d < dist {
				dist = d
				closest = j
			}
		}
		i.CaretIndex = closest
		hvc.Input.Get("click").Consume()
	}
	// make sure cursor doesn't go off the end or beginning
	if i.CaretIndex > i.Text.Len()-1 {
		i.CaretIndex = i.Text.Len() - 1
	}
	if i.CaretIndex < 0 {
		i.CaretIndex = 0
	}
	return changed
}

func Typing(i *Element, hvc *data.HoverClick) bool {
	changed := false
	// backspace, delete, new line
	back := hvc.Input.Get("backspace")
	if (back.JustPressed() || back.Repeated()) && i.CaretIndex > 0 {
		i.Value = fmt.Sprintf("%s%s", i.Value[:i.CaretIndex-1], i.Value[i.CaretIndex:])
		changed = true
		i.CaretIndex--
	}
	del := hvc.Input.Get("delete")
	if (del.JustPressed() || del.Repeated()) && i.CaretIndex < i.Text.Len()-1 {
		i.Value = fmt.Sprintf("%s%s", i.Value[:i.CaretIndex], i.Value[i.CaretIndex+1:])
		changed = true
	}
	if i.MultiLine {
		enter := hvc.Input.Get("enter")
		if enter.JustPressed() || enter.Repeated() {
			i.Value = fmt.Sprintf("%s\n%s", i.Value[:i.CaretIndex], i.Value[i.CaretIndex:])
			changed = true
			i.CaretIndex++
		}
	}
	// typing
	typed := hvc.Input.Typed
	if typed != "" {
		switch i.InputType {
		case AlphaNumeric:
			typed = util.OnlyAlphaNumeric(typed)
		case Numeric:
			typed = util.OnlyNumbers(typed)
		case Special:
			typed = util.JustChars(typed)
		}
		i.Value = fmt.Sprintf("%s%s%s", i.Value[:i.CaretIndex], typed, i.Value[i.CaretIndex:])
		changed = true
		i.CaretIndex += len(typed)
	}
	return changed
}

func ChangeText(input *Element, rt string) {
	input.Value = rt
	input.Text.SetText(input.Value)
}

func UpdateInputScrollBounds(input *Element) bool {
	yTop := input.Text.Obj.Pos.Y + input.Text.Text.Orig.Y + 1
	yBot := yTop - input.Text.GetHeight() - input.Text.Text.LineHeight*input.Text.Scalar
	input.YTop = yTop
	input.YBot = yBot
	if input.YTop-input.YBot > input.ViewPort.Rect.H() && input.ScrollUp.Object.Hidden {
		input.ScrollUp.Object.Hidden = false
		input.ScrollDown.Object.Hidden = false
		input.Bar.Object.Hidden = false
		input.ViewPort.SetRect(pixel.R(0, 0, input.Object.Rect.W()-world.TileSize, input.Object.Rect.H()))
		input.Text.SetWidth(input.ViewPort.Rect.W() - world.TileSize*0.75)
		input.ViewPort.PortPos.X = input.Object.Pos.X - world.HalfSize
		return true
	} else if input.YTop-input.YBot <= input.ViewPort.Rect.H() && !input.ScrollUp.Object.Hidden {
		input.ScrollUp.Object.Hidden = true
		input.ScrollDown.Object.Hidden = true
		input.Bar.Object.Hidden = true
		input.ViewPort.SetRect(pixel.R(0, 0, input.Object.Rect.W(), input.Object.Rect.H()))
		input.Text.SetWidth(input.ViewPort.Rect.W() - world.TileSize*0.75)
		input.ViewPort.PortPos.X = input.Object.Pos.X
		return true
	}
	input.ViewPort.CamPos.X = input.ViewPort.Rect.W()*0.5 - 2
	return false
}

func CreateScrollElement(element ElementConstructor, dlg *Dialog, parent *Element, vp *viewport.ViewPort) *Element {
	svp := viewport.New(nil)
	svp.ParentView = vp
	svp.SetRect(pixel.R(0, 0, element.Width-world.TileSize, element.Height))
	svp.CamPos = pixel.V(0, 0)
	svp.PortPos = element.Position
	svp.PortPos.X -= world.HalfSize

	bvp := viewport.New(nil)
	bvp.SetRect(pixel.R(0, 0, element.Width+1, element.Height+1))
	bvp.CamPos = pixel.V(0, 0)
	bvp.PortPos = element.Position

	vpObj := object.New()
	vpObj.SetRect(pixel.R(0, 0, element.Width+1, element.Height+1))
	vpObj.SetPos(element.Position)
	vpObj.Layer = 99

	bord := &Border{
		Rect:  pixel.R(0, 0, element.Width, element.Height),
		Style: ThinBorder,
	}

	e := myecs.Manager.NewEntity().AddComponent(myecs.Object, vpObj)
	s := &Element{
		Key:         element.Key,
		Border:      bord,
		Object:      vpObj,
		Entity:      e,
		ViewPort:    svp,
		ElementType: ScrollElement,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	e.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, svp, func(hvc *data.HoverClick) {
		if hvc.ViewHover {
			if hvc.Input.ScrollV > 0. {
				s.ViewPort.CamPos.Y += ScrollSpeed * timing.DT
			} else if hvc.Input.ScrollV < 0. {
				s.ViewPort.CamPos.Y -= ScrollSpeed * timing.DT
			}
			RestrictScroll(s)
			AlignBarToView(s)
		}
	}))
	CreateScrollBars(dlg, s, parent, element)

	for _, ele := range element.SubElements {
		if ele.Key == "" {
			fmt.Println("WARNING: element constructor has no key")
		}
		switch ele.ElementType {
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
			i := CreateInputElement(ele, dlg, s, s.ViewPort, false)
			s.Elements = append(s.Elements, i)
		case MultiLineInputElement:
			i := CreateInputElement(ele, dlg, s, s.ViewPort, true)
			s.Elements = append(s.Elements, i)
		case ScrollElement:
			s2 := CreateScrollElement(ele, dlg, s, s.ViewPort)
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
	MoveToScrollTop(s)
	return s
}

func CreateScrollBars(dlg *Dialog, s, parent *Element, element ElementConstructor) {
	btnX := s.ViewPort.Rect.W() * 0.5
	if s.MultiLine {
		btnX -= world.TileSize * 0.5
	}
	for i := 0; i < 3; i++ {
		var pos pixel.Vec
		var key, sprKey, cSprKey string
		switch i {
		case 0:
			pos = element.Position.Add(pixel.V(btnX, s.ViewPort.Rect.H()*0.5))
			key = fmt.Sprintf("%s_scroll_up", element.Key)
			sprKey = "scroll_up"
			cSprKey = "scroll_up_click"
		case 1:
			pos = element.Position.Add(pixel.V(btnX, s.ViewPort.Rect.H()*-0.5))
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
			SprKey2:     cSprKey,
			Batch:       element.Batch,
			Position:    pos,
			ElementType: ButtonElement,
		}
		var b *Element
		if parent != nil {
			b = CreateButtonElement(btn, dlg, parent.ViewPort)
			parent.Elements = append(parent.Elements, b)
		} else {
			b = CreateButtonElement(btn, dlg, dlg.ViewPort)
			dlg.Elements = append(dlg.Elements, b)
		}
		switch i {
		case 0:
			s.ScrollUp = b
			s.ButtonHeight = b.Object.Rect.H()
			b.Object.Pos.Y -= b.Object.Rect.H() * 0.5
			b.OnHold = func() {
				s.ViewPort.CamPos.Y += ScrollSpeed * timing.DT
				RestrictScroll(s)
				AlignBarToView(s)
			}
		case 1:
			s.ScrollDown = b
			b.Object.Pos.Y += b.Object.Rect.H() * 0.5
			b.OnHold = func() {
				s.ViewPort.CamPos.Y -= ScrollSpeed * timing.DT
				RestrictScroll(s)
				AlignBarToView(s)
			}
		case 2:
			s.Bar = b
			offset := 0.
			barClick := false
			b.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dlg.ViewPort, func(hvc *data.HoverClick) {
				if dlg.Open && dlg.Active && !dlg.Lock {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						b.Entity.AddComponent(myecs.Drawable, b.Sprite2)
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
}

func AlignViewToBar(s *Element) {
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

func AlignBarToView(s *Element) {
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
}

func MoveScrollToInclude(scroll *Element, yTop, yBot float64) {
	if scroll.ViewPort.CamPos.Y-scroll.ViewPort.Rect.H()*0.5 > yBot {
		scroll.ViewPort.CamPos.Y = yBot + scroll.ViewPort.Rect.H()*0.5
	} else if scroll.ViewPort.CamPos.Y+scroll.ViewPort.Rect.H()*0.5 < yTop {
		scroll.ViewPort.CamPos.Y = yTop - scroll.ViewPort.Rect.H()*0.5
	}
}

func RestrictScroll(s *Element) {
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

func UpdateScrollBounds(scroll *Element) {
	yTop := 0.
	yBot := 0.
	for i, ele := range scroll.Elements {
		obj := object.New()
		obj.Rect = ele.Object.Rect
		obj.Pos = ele.Object.Pos
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
}

func MoveToScrollTop(scroll *Element) {
	scroll.ViewPort.CamPos.Y = scroll.YTop - scroll.ViewPort.Rect.H()*0.5
	scroll.Bar.Object.Pos.Y = scroll.ViewPort.PortPos.Y + scroll.ViewPort.Rect.H()*0.5 - scroll.ButtonHeight - scroll.Bar.Object.Rect.H()*0.5
}

func CreateSpriteElement(element ElementConstructor) *Element {
	obj := object.New()
	obj.Pos = element.Position
	obj.Layer = 99
	obj.SetRect(img.Batchers[element.Batch].GetSprite(element.SprKey).Frame())
	spr := img.NewSprite(element.SprKey, element.Batch)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj).
		AddComponent(myecs.Drawable, spr)
	s := &Element{
		Key:         element.Key,
		Sprite:      spr,
		Object:      obj,
		Entity:      e,
		ElementType: SpriteElement,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	return s
}

func CreateTextElement(element ElementConstructor, vp *viewport.ViewPort) *Element {
	tf := typeface.NewOld("main", typeface.NewAlign(typeface.Left, typeface.Top), 1, 0.0625, 0, 0)
	tf.SetPos(element.Position)
	tf.SetColor(element.Color)
	tf.SetText(element.Text)
	tf.WithAnchor(element.Anchor)
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, tf.Obj)
	e.AddComponent(myecs.Drawable, tf)
	e.AddComponent(myecs.DrawTarget, vp.Canvas)
	t := &Element{
		Key:         element.Key,
		Text:        tf,
		Object:      tf.Obj,
		Entity:      e,
		ElementType: TextElement,
		Left:        element.Left,
		Right:       element.Right,
		Up:          element.Up,
		Down:        element.Down,
	}
	return t
}
