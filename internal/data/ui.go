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
	InputElement
	ScrollElement
	SpriteElement
	TextElement
	CustomElement
)

type ElementConstructor struct {
	Key         string
	Width       float64
	Height      float64
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

type Input struct {
	Key          string
	Value        string
	Active       bool
	Text         *typeface.Text
	TextEntity   *ecs.Entity
	CaretObj     *object.Object
	CaretSpr     *img.Sprite
	CaretIndex   int
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	ViewPort     *viewport.ViewPort
	Entity       *ecs.Entity
	Layer        int
}

type Scroll struct {
	Key          string
	UpSprite     *img.Sprite
	UpSprClick   *img.Sprite
	UpObject     *object.Object
	UpEntity     *ecs.Entity
	DwnSprite    *img.Sprite
	DwnSprClick  *img.Sprite
	DwnObject    *object.Object
	DwnEntity    *ecs.Entity
	BarSprite    []*img.Sprite
	BarObject    *object.Object
	BarEntity    *ecs.Entity
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	BorderEntity *ecs.Entity
	Entity       *ecs.Entity
	ViewPort     *viewport.ViewPort
	Layer        int
	Elements     []interface{}
	YTop         float64
	YBot         float64
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
