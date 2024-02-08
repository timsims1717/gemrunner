package data

import "github.com/bytearena/ecs"

type OnTouch struct {
	Fn func(*Level, *Character, *ecs.Entity)
}

func NewOnTouch(fn func(*Level, *Character, *ecs.Entity)) *OnTouch {
	return &OnTouch{Fn: fn}
}
