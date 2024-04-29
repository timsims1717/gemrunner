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
	CheckboxElement
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

type UIElement interface {
	GetKey() string
	SetPos(pixel.Vec)
	GetEntity() *ecs.Entity
	GetObject() *object.Object
	GetHelp() string
	GetSprite() *img.Sprite
	GetSprite2() *img.Sprite
	GetElements() []UIElement
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
	OnHeld   func(*HoverClick)
}

type Checkbox struct {
	Key      string
	Sprite   *img.Sprite
	CheckSpr *img.Sprite
	HelpText string
	Object   *object.Object
	Entity   *ecs.Entity
	Checked  bool
}

type Container struct {
	Key          string
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	Object       *object.Object
	BorderEntity *ecs.Entity
	Entity       *ecs.Entity
	ViewPort     *viewport.ViewPort
	Layer        int
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
	Bar          *Button
	ButtonHeight float64
	BorderVP     *viewport.ViewPort
	BorderObject *object.Object
	Object       *object.Object
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
