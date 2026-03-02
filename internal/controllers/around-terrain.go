package controllers

import (
	"gemrunner/internal/data"
	"github.com/bytearena/ecs"
)

type AroundTerrain struct {
	Ch     *data.Dynamic
	Left   bool
	Entity *ecs.Entity
}

func NewAroundTerrain(dyn *data.Dynamic, e *ecs.Entity, left bool) *AroundTerrain {
	return &AroundTerrain{
		Ch:     dyn,
		Left:   left,
		Entity: e,
	}
}

func (at *AroundTerrain) ClearPrev() {}

func (at *AroundTerrain) GetActions() data.Actions {
	actions := data.NewAction()
	if at.Left {
		switch at.Ch.Flags.Orientation {
		case data.Down:
			actions.Direction = data.Left
		case data.Left:
			actions.Direction = data.Up
		case data.Up:
			actions.Direction = data.Right
		case data.Right:
			actions.Direction = data.Down
		}
	} else {
		switch at.Ch.Flags.Orientation {
		case data.Down:
			actions.Direction = data.Right
		case data.Left:
			actions.Direction = data.Down
		case data.Up:
			actions.Direction = data.Left
		case data.Right:
			actions.Direction = data.Up
		}
	}
	return actions
}

func (at *AroundTerrain) GetEntity() *ecs.Entity {
	return at.Entity
}
