package data

import (
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/typeface"
	"gemrunner/pkg/viewport"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type Element int

const (
	ButtonElement = iota
	ContainerElement
	ScrollElement
	SpriteElement
	TextElement
	CustomElement
)

type ElementConstructor struct {
	Key         string
	SprKey      string
	ClickSprKey string
	Text        string
	HelpText    string
	Position    pixel.Vec
	Element     Element
	SubElements []ElementConstructor
}

type Button struct {
	Key      string
	Sprite   *img.Sprite
	ClickSpr *img.Sprite
	Delay    float64
	HelpText string
	Object   *object.Object
	Entity   *ecs.Entity
	OnClick  func()
}

type Container struct {
	Key          string
	Viewport     *viewport.ViewPort
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	Elements     []interface{}
}

type Scroll struct {
	Key         string
	UpSprite    *img.Sprite
	UpSprClick  *img.Sprite
	UpObject    *object.Object
	DwnSprite   *img.Sprite
	DwnSprClick *img.Sprite
	DwnObject   *object.Object
	BarSprite   []*img.Sprite
	BarObject   *object.Object
	BarEntity   *ecs.Entity
	Viewport    viewport.ViewPort
	Elements    []interface{}
}

type SprElement struct {
	Key    string
	Sprite *img.Sprite
	Object *object.Object
	Entity *ecs.Entity
}

type Text struct {
	Key    string
	Text   *typeface.Text
	Entity *ecs.Entity
}
