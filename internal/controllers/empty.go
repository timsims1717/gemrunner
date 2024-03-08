package controllers

import (
	"gemrunner/internal/data"
	"github.com/bytearena/ecs"
)

type Empty struct {
	Ch     *data.Dynamic
	Entity *ecs.Entity
}

func NewEmpty(dyn *data.Dynamic, e *ecs.Entity) *Empty {
	return &Empty{
		Ch:     dyn,
		Entity: e,
	}
}

func (e *Empty) ClearPrev() {}

func (e *Empty) GetActions() data.Actions {
	return data.NewAction()
}

func (e *Empty) GetEntity() *ecs.Entity {
	return e.Entity
}
