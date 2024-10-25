package data

import (
	"gemrunner/internal/myecs"
	"gemrunner/pkg/typeface"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type FloatingTextData struct {
	Pos       pixel.Vec  `json:"pos"`
	Raw       string     `json:"text"`
	Prox      bool       `json:"prox"`
	Bob       bool       `json:"bob"`
	Timer     int        `json:"timer"`
	HasShadow bool       `json:"hasShadow"`
	Color     pixel.RGBA `json:"color"`
	ShadowCol pixel.RGBA `json:"shadow"`
}

type FloatingText struct {
	Tile        *Tile
	Pos         pixel.Vec
	Entity      *ecs.Entity
	ShEntity    *ecs.Entity
	Text        *typeface.Text
	Shadow      *typeface.Text
	Raw         string
	Temp        bool
	Prox        bool
	ProxCounter int
	Bob         bool
	BobCounter  int
	TempCounter int
	Timer       int
	HasShadow   bool
	Color       pixel.RGBA
	ShadowCol   pixel.RGBA
}

func NewFloatingText() *FloatingText {
	txt := typeface.New("main").
		WithAlign(typeface.NewAlign(typeface.Center, typeface.Center)).
		WithAnchor(pixel.Center).
		WithScalar(0.0625)
	txt.Obj.ILock = false
	txt.Obj.Layer = 37
	txt.Update()
	shTxt := typeface.New("main").
		WithAlign(typeface.NewAlign(typeface.Center, typeface.Center)).
		WithAnchor(pixel.Center).
		WithScalar(0.0625)
	shTxt.Obj.ILock = false
	shTxt.Obj.Layer = 36
	shTxt.Hide()
	ft := &FloatingText{
		Text:      txt,
		Shadow:    shTxt,
		HasShadow: false,
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

func (ft *FloatingText) WithText(raw string) *FloatingText {
	ft.Raw = raw
	ft.Text.SetText(raw)
	ft.Text.Update()
	if ft.Shadow != nil {
		ft.Shadow.SetText(raw)
		ft.Shadow.Update()
	}
	return ft
}

func (ft *FloatingText) WithTimer(time int) *FloatingText {
	ft.Temp = time > 0
	ft.Timer = time
	return ft
}

func (ft *FloatingText) WithTile(tile *Tile) *FloatingText {
	ft.Tile = tile
	if ft.Tile == nil {
		ft.Entity.AddComponent(myecs.Temp, myecs.ClearFlag(false))
		ft.ShEntity.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	} else {
		ft.Entity.RemoveComponent(myecs.Temp)
		ft.ShEntity.RemoveComponent(myecs.Temp)
	}
	return ft
}

func (ft *FloatingText) WithBools(prox, bob bool) *FloatingText {
	ft.Prox = prox
	ft.Bob = bob
	return ft
}

func (ft *FloatingText) WithPos(pos pixel.Vec) *FloatingText {
	ft.Pos = pos
	ft.Text.SetPos(pos)
	if ft.Shadow != nil {
		ft.Shadow.SetPos(pos.Add(pixel.V(0, 1)))
	}
	return ft
}

func (ft *FloatingText) WithColor(color pixel.RGBA) *FloatingText {
	ft.Color = color
	ft.Text.SetColor(color)
	return ft
}

func (ft *FloatingText) WithShadow(color pixel.RGBA) *FloatingText {
	ft.ShadowCol = color
	ft.Shadow.SetColor(color)
	ft.HasShadow = true
	ft.Shadow.Obj.Hidden = false
	return ft
}

func (ft *FloatingText) RemoveShadow() {
	ft.HasShadow = false
	ft.Shadow.Obj.Hidden = true
}

func (ft *FloatingText) Show() {
	ft.Text.Obj.Hidden = false
	ft.Shadow.Obj.Hidden = false
}

func (ft *FloatingText) Hide() {
	ft.Text.Obj.Hidden = true
	ft.Shadow.Obj.Hidden = true
}

func (ft *FloatingText) Init(tile *Tile) {
	if ft.Text == nil {
		txt := typeface.New("main").
			WithAlign(typeface.NewAlign(typeface.Center, typeface.Center)).
			WithAnchor(pixel.Center).
			WithScalar(0.0625)
		txt.Obj.Layer = 37
		txt.Obj.ILock = false
		ft.Text = txt
		ft.Entity = myecs.Manager.NewEntity().
			AddComponent(myecs.Object, txt.Obj).
			AddComponent(myecs.Drawable, txt).
			AddComponent(myecs.Text, ft).
			AddComponent(myecs.Temp, myecs.ClearFlag(false))
	}
	if ft.Shadow == nil {
		shTxt := typeface.New("main").
			WithAlign(typeface.NewAlign(typeface.Center, typeface.Center)).
			WithAnchor(pixel.Center).
			WithScalar(0.0625)
		shTxt.Obj.Layer = 36
		shTxt.Obj.ILock = false
		ft.Shadow = shTxt
		ft.ShEntity = myecs.Manager.NewEntity().
			AddComponent(myecs.Object, shTxt.Obj).
			AddComponent(myecs.Drawable, shTxt).
			AddComponent(myecs.Temp, myecs.ClearFlag(false))
	}
	ft.WithTile(tile).
		WithText(ft.Raw).
		WithPos(ft.Pos).
		WithColor(ft.Color).
		WithTimer(ft.Timer)
	if ft.HasShadow {
		ft.WithShadow(ft.ShadowCol)
	}
	ft.Text.Update()
	ft.Shadow.Update()
	ft.Tile = tile
}

func CreateFloatingText(tile *Tile, ftd *FloatingTextData) {
	if ftd == nil {
		if tile.FloatingText != nil {
			myecs.Manager.DisposeEntity(tile.FloatingText.Entity)
			myecs.Manager.DisposeEntity(tile.FloatingText.ShEntity)
			tile.FloatingText = nil
		}
		return
	}
	if tile.FloatingText == nil {
		tile.FloatingText = NewFloatingText()
	}
	tile.FloatingText.WithTile(tile).
		WithText(ftd.Raw).
		WithPos(tile.Object.Pos).
		WithColor(ftd.Color).
		WithTimer(ftd.Timer).
		WithBools(ftd.Prox, ftd.Bob)
	if ftd.HasShadow {
		tile.FloatingText.WithShadow(ftd.ShadowCol)
	} else {
		tile.FloatingText.RemoveShadow()
	}
	tile.FloatingText.Text.Update()
	tile.FloatingText.Shadow.Update()
}

func (ftd *FloatingTextData) Copy() *FloatingTextData {
	if ftd == nil {
		return nil
	}
	return &FloatingTextData{
		Pos:       ftd.Pos,
		Raw:       ftd.Raw,
		Prox:      ftd.Prox,
		Bob:       ftd.Bob,
		Timer:     ftd.Timer,
		HasShadow: ftd.HasShadow,
		Color:     ftd.Color,
		ShadowCol: ftd.ShadowCol,
	}
}
