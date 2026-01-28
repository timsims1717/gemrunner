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
		Direction: data.NoDirection,
		Entity:    e,
	}
}

func (rw *RandomWalk) ClearPrev() {}

func (rw *RandomWalk) GetActions() data.Actions {
	actions := data.NewAction()
	change := false
	if reanimator.FrameSwitch {
		if random.Level.Intn(100) == 0 {
			change = true
		}
		if random.Level.Intn(1000) == 0 {
			x := random.Level.Intn(data.CurrLevel.Metadata.Width)
			y := random.Level.Intn(data.CurrLevel.Metadata.Height)
			tile := data.CurrLevel.Get(x, y)
			if tile.IsEmpty() {
				rw.Ch.Object.SetPos(tile.Object.Pos)
			}
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
	} else if rw.Direction == data.NoDirection {
		change = true
	}
	if change {
		x, y := world.WorldToMap(rw.Ch.Object.Pos.X, rw.Ch.Object.Pos.Y)
		currTile := data.CurrLevel.Get(x, y)
		belowTile := data.CurrLevel.Get(x, y-1)
		l := !rw.Ch.Flags.LeftWall
		r := !rw.Ch.Flags.RightWall
		u := !rw.Ch.Flags.Ceiling && currTile.IsLadder()
		d := !belowTile.IsSolid()
		rn := random.Level.Intn(4)
		rno := rn
		for {
			switch rn {
			case 0:
				if l {
					rw.Direction = data.Left
				}
			case 1:
				if r {
					rw.Direction = data.Right
				}
			case 2:
				if u {
					rw.Direction = data.Up
				}
			case 3:
				if d {
					rw.Direction = data.Down
				}
			}
			rn++
			rn %= 4
			if rn == rno {
				break
			}
		}
	}
	actions.Direction = rw.Direction
	return actions
}

func (rw *RandomWalk) GetEntity() *ecs.Entity {
	return rw.Entity
}
