package controllers

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

type AStarChase struct {
	Ch     *data.Dynamic
	Target *world.Coords
	Entity *ecs.Entity
	lastT  world.Coords
	path   []world.Coords
}

func NewAStarChase(dyn *data.Dynamic, e *ecs.Entity, target *world.Coords) *AStarChase {
	return &AStarChase{
		Ch:     dyn,
		Target: target,
		Entity: e,
	}
}

func (as *AStarChase) ClearPrev() {}

func (as *AStarChase) GetActions() data.Actions {
	actions := data.NewAction()
	if as.Target == nil {
		// we can put player picker here
		return actions
	}
	if *as.Target == as.lastT {
		if len(as.path) > 0 {
			x, y := world.WorldToMap(as.Ch.Object.Pos.X, as.Ch.Object.Pos.Y)
			next := as.path[0]
			if next.X < x {
				actions.Direction = data.Left
			} else if next.X > x {
				actions.Direction = data.Right
			} else if next.Y < y {
				actions.Direction = data.Down
			} else if next.Y > y {
				actions.Direction = data.Up
			}
		}
		return actions
	}
	return actions
}

func (as *AStarChase) GetEntity() *ecs.Entity {
	return as.Entity
}
