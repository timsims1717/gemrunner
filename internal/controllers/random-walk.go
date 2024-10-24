package controllers

import (
	"gemrunner/internal/data"
	"gemrunner/internal/random"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
)

type RandomWalk struct {
	Ch        *data.Dynamic
	Direction data.Direction
	Entity    *ecs.Entity
}

func NewRandomWalk(dyn *data.Dynamic, e *ecs.Entity) *RandomWalk {
	return &RandomWalk{
		Ch:        dyn,
		Direction: data.None,
		Entity:    e,
	}
}

func (rw *RandomWalk) ClearPrev() {}

func (rw *RandomWalk) GetActions() data.Actions {
	actions := data.NewAction()
	change := false
	if reanimator.FrameSwitch {
		r := random.Level.Intn(100)
		if r < 1 {
			change = true
		}
	}
	if rw.Ch.Flags.RightWall && rw.Direction == data.Right {
		change = true
	} else if rw.Ch.Flags.LeftWall && rw.Direction == data.Left {
		change = true
	} else if rw.Ch.Flags.Ceiling && rw.Direction == data.Up {
		change = true
	} else if rw.Ch.Flags.Floor && rw.Direction == data.Down {
		change = true
	} else if rw.Direction == data.None {
		change = true
	}
	if change {
		x, y := world.WorldToMap(rw.Ch.Object.Pos.X, rw.Ch.Object.Pos.Y)
		currTile := data.CurrLevel.Tiles.Get(x, y)
		belowTile := data.CurrLevel.Tiles.Get(x, y-1)
		l := !rw.Ch.Flags.LeftWall
		r := !rw.Ch.Flags.RightWall
		u := !rw.Ch.Flags.Ceiling && currTile.IsLadder()
		d := !belowTile.IsNilOrSolid()
		rn := random.Level.Intn(4)
		if rn == 0 && l {
			rw.Direction = data.Left
		} else if rn == 1 && r {
			rw.Direction = data.Right
		} else if rn == 2 && u {
			rw.Direction = data.Up
		} else if rn == 3 && d {
			rw.Direction = data.Down
		}
	}
	actions.Direction = rw.Direction
	return actions
}

func (rw *RandomWalk) GetEntity() *ecs.Entity {
	return rw.Entity
}
