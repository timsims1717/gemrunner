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
	ScrollSpeed     float64
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
	Elements     []*Element
	NoBorder     bool
	OnOpen       func()
	OnClose      func()
	OnCloseSpc   func()

	Focused string

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

func (d *Dialog) ActionFocus() {

}

func (d *Dialog) LeftFocus() {

}

func (d *Dialog) RightFocus() {

}

func (d *Dialog) UpFocus() {

}

func (d *Dialog) DownFocus() {

}
