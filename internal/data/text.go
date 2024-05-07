package data

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/typeface"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type FloatingText struct {
	Entity  *ecs.Entity
	Text    *typeface.Text
	Temp    bool
	Prox    bool
	Bob     bool
	Counter int
	Timer   int
}

func NewFloatingText(raw string, temp, prox, bob bool, timer int, pos pixel.Vec) *FloatingText {
	txt := typeface.New("main", typeface.NewAlign(typeface.Center, typeface.Center), 1, 0.0625, 0, 0)
	txt.SetPos(pos)
	txt.Obj.Layer = 36
	txt.SetColor(pixel.ToRGBA(constants.ColorWhite))
	txt.SetText(raw)
	ft := &FloatingText{
		Text:  txt,
		Temp:  temp,
		Prox:  prox,
		Bob:   bob,
		Timer: timer,
	}
	ft.Entity = myecs.Manager.NewEntity().
		AddComponent(myecs.Object, txt.Obj).
		AddComponent(myecs.Drawable, txt).
		AddComponent(myecs.Text, ft).
		AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return ft
}

func (ft *FloatingText) SetText(raw string) {
	ft.Text.SetText(raw)
}

func (ft *FloatingText) Show() {
	ft.Counter = 0
	ft.Text.Obj.Hidden = false
}
