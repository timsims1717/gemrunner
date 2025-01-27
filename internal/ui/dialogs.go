package ui

import (
	"fmt"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

var (
	ScrollSpeed        float64
	DialogStack        []*Dialog
	DialogsOpen        []*Dialog
	DialogStackOpen    bool
	Dialogs            = map[string]*Dialog{}
	DialogConstructors = map[string]*DialogConstructor{}
)

type Dialog struct {
	Key          string
	Pos          pixel.Vec
	ViewPort     *viewport.ViewPort
	Border       *Border
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	Elements     []*Element
	NoBorder     bool
	OnOpen       func()
	OnClose      func()
	OnCloseSpc   func()

	Focused *Element
	Default string

	Open   bool
	Active bool
	Loaded bool
	Click  bool
	Lock   bool
	Layer  int
}

type DialogConstructor struct {
	Key      string               `json:"key"`
	Width    float64              `json:"width"`
	Height   float64              `json:"height"`
	Pos      pixel.Vec            `json:"pos"`
	Elements []ElementConstructor `json:"elements,omitempty"`
	NoBorder bool                 `json:"noBorder,omitempty"`
	Default  string               `json:"default"`
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
		Default:  dc.Default,
	}

	if !dc.NoBorder {
		bvp := viewport.New(nil)
		bvp.SetRect(pixel.R(0, 0, (dc.Width+1)*world.TileSize, (dc.Height+1)*world.TileSize))
		bvp.CamPos = pixel.V(0, 0)
		bvp.PortPos = viewport.MainCamera.PostCamPos.Add(dc.Pos)

		bObj := object.New()
		bObj.Layer = 99
		//bObj.Pos = dc.Pos
		bord := &Border{
			Width:  int(dc.Width),
			Height: int(dc.Height),
		}
		be := myecs.Manager.NewEntity()
		be.AddComponent(myecs.Object, bObj).
			AddComponent(myecs.Border, bord)

		dlg.Border = bord
		dlg.BorderVP = bvp
		dlg.BorderObject = bObj
		dlg.BorderEntity = be
	}

	for _, element := range dc.Elements {
		if element.Key == "" {
			fmt.Printf("WARNING: element constructor in dialog %s has no key\n", dc.Key)
			continue
		} else if dlg.Default == "" {
			dlg.Default = element.Key
		}
		switch element.ElementType {
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
			i := CreateInputElement(element, dlg, nil, dlg.ViewPort, false)
			dlg.Elements = append(dlg.Elements, i)
		case MultiLineInputElement:
			i := CreateInputElement(element, dlg, nil, dlg.ViewPort, true)
			dlg.Elements = append(dlg.Elements, i)
		case ScrollElement:
			s := CreateScrollElement(element, dlg, nil, dlg.ViewPort)
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

func (d *Dialog) Get(key string) *Element {
	for _, e := range d.Elements {
		if e.Key == key {
			return e
		}
	}
	for _, e := range d.Elements {
		if se := e.Get(key); se != nil {
			return se
		}
	}
	return nil
}

func (d *Dialog) SetFocus(e *Element, isFocused bool) {
	eKey := ""
	if e != nil {
		e.InFocus = isFocused
		if !isFocused {
			d.Focused = nil
		} else {
			d.Focused = e
			eKey = e.Key
		}
	}
	for _, ele := range d.Elements {
		ele.setFocus(eKey)
	}
}

func (e *Element) setFocus(key string) {
	focused := e.Key == key
	switch e.ElementType {
	case ButtonElement, CheckboxElement:
		if e.Border != nil {
			e.Border.Hidden = !focused
		}
	}
	for _, ele := range e.Elements {
		ele.setFocus(key)
	}
	e.InFocus = focused
	if e.OnFocus != nil {
		e.OnFocus(focused)
	}
}

func (d *Dialog) DefaultFocus() bool {
	if d.Focused == nil {
		if d.Default == "" {
			return true
		}
		d.SetFocus(d.Get(d.Default), true)
		return true
	}
	return false
}

func (d *Dialog) ActionFocus() {
	if d.DefaultFocus() {
		return
	}
	if d.Focused.OnClick != nil {
		d.Focused.OnClick()
		return
	}
	switch d.Focused.ElementType {
	case CheckboxElement:
		SetChecked(d.Focused, !d.Focused.Checked)
	case ContainerElement:

	case InputElement:

	case MultiLineInputElement:

	case ScrollElement:

	}
}

func (d *Dialog) LeftFocus() {
	if d.DefaultFocus() {
		return
	}
	if d.Focused.Left == "" {
		return
	}
	if e := d.Get(d.Focused.Left); e != nil {
		d.SetFocus(e, true)
	}
}

func (d *Dialog) RightFocus() {
	if d.DefaultFocus() {
		return
	}
	if d.Focused.Right == "" {
		return
	}
	if e := d.Get(d.Focused.Right); e != nil {
		d.SetFocus(e, true)
	}
}

func (d *Dialog) UpFocus() {
	if d.DefaultFocus() {
		return
	}
	if d.Focused.Up == "" {
		return
	}
	if e := d.Get(d.Focused.Up); e != nil {
		d.SetFocus(e, true)
	}
}

func (d *Dialog) DownFocus() {
	if d.DefaultFocus() {
		return
	}
	if d.Focused.Down == "" {
		return
	}
	if e := d.Get(d.Focused.Down); e != nil {
		d.SetFocus(e, true)
	}
}

func Dispose(key string) {
	CloseDialog(key)
	for _, d := range Dialogs {
		if d.Key == key {
			DisposeDialog(d)
		}
	}
}

func DisposeDialog(d *Dialog) {
	DisposeSubElements(d.Elements)
	myecs.Manager.DisposeEntity(d.BorderEntity)
	delete(Dialogs, d.Key)
}

func DisposeSubElements(elements []*Element) {
	for _, e := range elements {
		DisposeSubElements(e.Elements)
		myecs.Manager.DisposeEntity(e.Entity)
	}
}
