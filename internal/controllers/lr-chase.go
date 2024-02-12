package controllers

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/beefsack/go-astar"
	"math/rand"
)

type LRChase struct {
	Ch     *data.Dynamic
	Target *data.Dynamic
	Timer  *timing.Timer
}

func NewLRChase(dyn *data.Dynamic) *LRChase {
	return &LRChase{
		Ch:    dyn,
		Timer: timing.New(constants.WaitToSwitch + rand.Float64()*3.),
	}
}

func (lr *LRChase) GetActions() data.Actions {
	lr.Timer.Update()
	actions := data.Actions{}
	if data.CurrLevel == nil || lr.Ch == nil {
		return actions
	}
	if lr.Timer.Done() || lr.Target == nil {
		lr.Target = PickClosestPlayerXFirst(lr.Ch)
		lr.Timer = timing.New(constants.WaitToSwitch + rand.Float64()*3.)
	}
	if lr.Target == nil {
		return actions
	}
	x, y := world.WorldToMap(lr.Ch.Object.Pos.X, lr.Ch.Object.Pos.Y)
	px, py := world.WorldToMap(lr.Target.Object.Pos.X, lr.Target.Object.Pos.Y)
	// picked a player
	// enemy is on the same y level as player
	if py == y {
		sx := x
		for {
			var nextTile, belowTile *data.Tile
			nextTile = data.CurrLevel.Tiles.Get(sx, y)
			belowTile = data.CurrLevel.Tiles.Get(sx, y-1)
			if nextTile == nil {
				// failed to find
				fmt.Println("WARNING: LRChase searched off the map")
				break
			}
			if nextTile.Ladder ||
				belowTile == nil || belowTile.Solid() || belowTile.Ladder {
				if sx < px {
					sx++
				} else if sx > px {
					sx--
				} else {
					break
				}
			} else {
				break
			}
		}
		// path to player found
		if sx == px {
			if px > x { // player is right of enemy
				actions.Right = true
			} else if px < x { // player is left of enemy
				actions.Left = true
			} else { // player is in same tile as enemy
				// found them
				if lr.Target.Object.Pos.X > lr.Ch.Object.Pos.X {
					actions.Right = true
				} else if lr.Target.Object.Pos.X < lr.Ch.Object.Pos.X {
					actions.Left = true
				}
			}
			return actions
		}
	}

	//bestT := -1
	bestD := -1
	next := world.Coords{X: -1, Y: -1}
	for i := 0; i < constants.PuzzleWidth; i++ {
		path, d, found := astar.Path(data.CurrLevel.Tiles.Get(x, y), data.CurrLevel.Tiles.Get(i, py))
		if len(path) > 0 && found && (bestD == -1 || int(d) < bestD) {
			bestD = int(d)
			//bestT = i
			next = path[0].(*data.Tile).Coords
		}
	}
	if next.X > -1 && next.Y > -1 {
		if next.X > x {
			actions.Right = true
		} else if next.X < x {
			actions.Left = true
		} else if next.Y > y {
			actions.Up = true
		} else if next.Y < y {
			actions.Down = true
		}
	}
	return actions
}

//func scanFloor(startX, startY int, cp *data.Dynamic) {
//	sx := startX
//	sy := startY
//	best := -1
//	curr := -1
//	left := -1
//	right := -1
//
//	// scan to the left
//	for sx > 0 {
//		nextTile := data.CurrLevel.Tiles.Get(sx-1, sy)
//		if nextTile == nil || nextTile.Solid() {
//			break
//		}
//		belowTile := data.CurrLevel.Tiles.Get(sx-1, sy-1)
//		if nextTile.Ladder ||
//			belowTile == nil || belowTile.Solid() || belowTile.Ladder {
//			sx--
//		} else {
//			sx--
//			break
//		}
//	}
//	left = sx
//	sx = startX
//	// scan to the right
//	for sx < constants.PuzzleWidth-1 {
//		nextTile := data.CurrLevel.Tiles.Get(sx+1, sy)
//		if nextTile == nil || nextTile.Solid() {
//			break
//		}
//		belowTile := data.CurrLevel.Tiles.Get(sx+1, sy-1)
//		if nextTile.Ladder ||
//			belowTile == nil || belowTile.Solid() || belowTile.Ladder {
//			sx++
//		} else {
//			sx++
//			break
//		}
//	}
//	right = sx
//
//}
