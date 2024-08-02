package data

import (
	"gemrunner/internal/myecs"
	"gemrunner/pkg/typeface"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type FloatingText struct {
	Tile     *Tile          `json:"-"`
	Entity   *ecs.Entity    `json:"-"`
	ShEntity *ecs.Entity    `json:"-"`
	Text     *typeface.Text `json:"-"`
	Shadow   *typeface.Text `json:"-"`
	Temp     bool           `json:"-"`
	Prox     bool           `json:"prox"`
	Bob      bool           `json:"bob"`
	Counter  int            `json:"-"`
}

func NewFloatingText(raw string, temp, prox, bob bool, pos pixel.Vec, color, shadow pixel.RGBA, tile *Tile) *FloatingText {
	txt := typeface.New("main", typeface.NewAlign(typeface.Center, typeface.Center), 1, 0.0625, 0, 0)
	txt.SetPos(pos)
	txt.Obj.Layer = 37
	txt.SetColor(color)
	txt.SetText(raw)
	shTxt := typeface.New("main", typeface.NewAlign(typeface.Center, typeface.Center), 1, 0.0625, 0, 0)
	shTxt.SetPos(pos.Add(pixel.V(0, 1)))
	shTxt.Obj.Layer = 36
	shTxt.SetColor(shadow)
	shTxt.SetText(raw)
	ft := &FloatingText{
		Text:   txt,
		Shadow: shTxt,
		Temp:   temp,
		Prox:   prox,
		Bob:    bob,
		Tile:   tile,
	}
	ft.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Object, txt.Obj).
		AddComponent(myecs.Drawable, txt).
		AddComponent(myecs.Text, ft).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	ft.ShEntity = myecs.Manager.NewEntity().
		AddComponent(myecs.Object, shTxt.Obj).
		AddComponent(myecs.Drawable, shTxt).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return ft
}

func (ft *FloatingText) SetText(raw string) {
	ft.Text.SetText(raw)
}

func (ft *FloatingText) Show() {
	ft.Text.Obj.Hidden = false
}

func (ft *FloatingText) Hide() {
	ft.Text.Obj.Hidden = true
}
