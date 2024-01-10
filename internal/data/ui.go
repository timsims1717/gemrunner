package data

import (
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"github.com/bytearena/ecs"
	"github.com/faiface/pixel"
)

type ButtonConstructor struct {
	SprKey      string
	ClickSprKey string
	HelpText    string
	Position    pixel.Vec
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
