package controllers

import (
	"gemrunner/internal/data"
	"github.com/bytearena/ecs"
)

type BackAndForth struct {
	Ch     *data.Dynamic
	Left   bool
	Entity *ecs.Entity
}

func NewBackAndForth(dyn *data.Dynamic, e *ecs.Entity, left bool) *BackAndForth {
	return &BackAndForth{
		Ch:     dyn,
		Left:   left,
		Entity: e,
	}
}

func (bf *BackAndForth) ClearPrev() {}

func (bf *BackAndForth) GetActions() data.Actions {
	actions := data.NewAction()
	if bf.Ch.Flags.RightWall {
		bf.Left = true
		actions.Direction = data.Left
	} else if bf.Ch.Flags.LeftWall {
		bf.Left = false
		actions.Direction = data.Right
	} else if bf.Left {
		actions.Direction = data.Left
	} else {
		actions.Direction = data.Right
	}
	return actions
}

func (bf *BackAndForth) GetEntity() *ecs.Entity {
	return bf.Entity
}
