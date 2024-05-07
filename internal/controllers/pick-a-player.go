package controllers

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
)

// PickClosestPlayerXFirst chooses a player by checking the X distance first
func PickClosestPlayerXFirst(ch *data.Dynamic) *data.Dynamic {
	x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
	cx := -1
	cy := -1
	ci := -1
	for i, p := range data.CurrLevel.Players {
		if p == nil ||
			p.State == data.Hit ||
			p.State == data.Dead ||
			p.State == data.Waiting {
			continue
		}
		px, py := world.WorldToMap(p.Object.Pos.X, p.Object.Pos.Y)
		dx := util.Abs(px - x)
		dy := util.Abs(py - y)
		if dx == cx {
			if dy < cy {
				cx = dx
				cy = dy
				ci = i
				continue
			}
		}
		if cx == -1 || dx < cx {
			cx = dx
			cy = dy
			ci = i
			continue
		}
	}
	if ci == -1 {
		return nil
	} else {
		return data.CurrLevel.Players[ci]
	}
}

// PickClosestPlayerYFirst chooses a player by checking the X distance first
func PickClosestPlayerYFirst(ch *data.Dynamic) *data.Dynamic {
	x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
	cx := -1
	cy := -1
	ci := -1
	for i, p := range data.CurrLevel.Players {
		if p == nil ||
			p.State == data.Hit ||
			p.State == data.Dead ||
			p.State == data.Waiting {
			continue
		}
		px, py := world.WorldToMap(p.Object.Pos.X, p.Object.Pos.Y)
		dx := util.Abs(px - x)
		dy := util.Abs(py - y)
		if dy == cy {
			if dx < cx {
				cx = dx
				cy = dy
				ci = i
				continue
			}
		}
		if cy == -1 || dy < cy {
			cx = dx
			cy = dy
			ci = i
			continue
		}
	}
	if ci == -1 {
		return nil
	} else {
		return data.CurrLevel.Players[ci]
	}
}
