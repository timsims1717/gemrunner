package data

import (
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

type Element int

const (
	ButtonElement = iota
	SpriteElement
	CustomElement
)

type ElementConstructor struct {
	SprKey      string
	ClickSprKey string
	HelpText    string
	Position    pixel.Vec
	Element     Element
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

type SprElement struct {
	Key    string
	Sprite *img.Sprite
	Object *object.Object
	Entity *ecs.Entity
}
