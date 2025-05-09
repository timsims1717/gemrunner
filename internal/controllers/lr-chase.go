package controllers

import (
	"fmt"
	"gemrunner/internal/data"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"math"
)

// from https://github.com/SimonHung/LodeRunner/blob/master/lodeRunner.guard.js

type LRChase struct {
	Ch     *data.Dynamic
	Target *data.Dynamic
	Entity *ecs.Entity
	path   []world.Coords
}

func NewLRChase(dyn *data.Dynamic, e *ecs.Entity) *LRChase {
	return &LRChase{
		Ch:     dyn,
		Entity: e,
	}
}

func (lr *LRChase) GetEntity() *ecs.Entity {
	return lr.Entity
}

func (lr *LRChase) ClearPrev() {}

func (lr *LRChase) GetActions() data.Actions {
	actions := data.NewAction()
	if data.CurrLevel == nil || lr.Ch == nil {
		return actions
	}
	lr.Target = PickClosestPlayerYFirst(lr.Ch)
	//if lr.Timer.Done() || lr.Target == nil {
	//	lr.Target = PickClosestPlayerXFirst(lr.Ch)
	//	lr.Timer = timing.New(constants.WaitToSwitch + rand.Float64()*3.)
	//}
	if lr.Target == nil {
		return actions
	}
	x, y := world.WorldToMap(lr.Ch.Object.Pos.X, lr.Ch.Object.Pos.Y)
	px, py := world.WorldToMap(lr.Target.Object.Pos.X, lr.Target.Object.Pos.Y)
	// picked a player
	// enemy is on the same y level as player
	if py == y && lr.Target.State != data.Falling {
		sx := x
		for {
			var nextTile, belowTile *data.Tile
			nextTile = data.CurrLevel.Get(sx, y)
			belowTile = data.CurrLevel.Get(sx, y-1)
			if nextTile == nil {
				// failed to find
				fmt.Println("WARNING: LRChase searched off the map")
				break
			}
			if nextTile.IsLadder() || nextTile.Block == data.BlockBar ||
				!belowTile.IsEmpty() {
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
		//belowTile := data.CurrLevel.Get(x, y-1)
		//if lr.Ch.State == data.OnLadder &&
		//	(math.Abs(lr.Target.Object.Pos.Y-lr.Ch.Object.Pos.Y) > 1. ||
		//		(lr.Target.State != data.OnLadder &&
		//			(belowTile.IsSolid() || belowTile.IsLadder()))) { // Enemy is on the same level as player, but is on a ladder and needs to adjust
		//	if lr.Target.Object.Pos.Y > lr.Ch.Object.Pos.Y {
		//		actions.PrevDirection = data.Up
		//	} else if lr.Target.Object.Pos.Y < lr.Ch.Object.Pos.Y {
		//		actions.PrevDirection = data.Down
		//	}
		//}
		if sx == px { // path to player found
			if px > x { // player is right of enemy
				actions.Direction = data.Right
			} else if px < x { // player is left of enemy
				actions.Direction = data.Left
			} else if math.Abs(lr.Target.Object.Pos.X-lr.Ch.Object.Pos.X) > 2. { // player is in same tile as enemy
				if lr.Target.Object.Pos.X > lr.Ch.Object.Pos.X {
					actions.Direction = data.Right
				} else if lr.Target.Object.Pos.X < lr.Ch.Object.Pos.X {
					actions.Direction = data.Left
				}
			}
			//if debug.ShowDebug {
			//
			//}
			return actions
		}
	}

	return lr.scanFloor(x, y, actions)
}

func (lr *LRChase) scanFloor(startX, sy int, actions data.Actions) data.Actions {
	sx := startX
	bestScore := 10000
	actions.Direction = data.NoDirection
	var left, right int

	// scan to the left
	for sx > 0 {
		nextTile := data.CurrLevel.Get(sx-1, sy)
		if nextTile == nil || nextTile.IsSolid() {
			break
		}
		belowTile := data.CurrLevel.Get(sx-1, sy-1)
		if nextTile.IsLadder() || nextTile.Block == data.BlockBar ||
			!belowTile.IsEmpty() {
			sx--
		} else {
			sx--
			break
		}
	}
	left = sx

	sx = startX
	// scan to the right
	for sx < data.CurrLevel.Metadata.Width-1 {
		nextTile := data.CurrLevel.Get(sx+1, sy)
		if nextTile == nil || nextTile.IsSolid() {
			break
		}
		belowTile := data.CurrLevel.Get(sx+1, sy-1)
		if nextTile.IsLadder() || nextTile.Block == data.BlockBar ||
			!belowTile.IsEmpty() {
			sx++
		} else {
			sx++
			break
		}
	}
	right = sx

	// scan the current x up and down as the best option
	sx = startX
	belowTile := data.CurrLevel.Get(sx, sy-1)
	if sy > 0 && !belowTile.IsSolid() { // can move down
		currScore := lr.scanDown(sx, startX, sy)
		if currScore < bestScore {
			bestScore = currScore
			actions.Direction = data.Down
		}
	}
	nextTile := data.CurrLevel.Get(sx, sy)
	if nextTile.IsLadder() { // can move up a ladder
		currScore := lr.scanUp(sx, startX, sy)
		if currScore < bestScore {
			bestScore = currScore
			actions.Direction = data.Up
		}
	}

	// scan each direction up and down as next best option
	currDir := data.Direction(data.Left)
	sx = left
	for {
		if sx == startX {
			if currDir == data.Left && right != startX {
				currDir = data.Right
				sx = right
			} else {
				break
			}
		}
		belowTile = data.CurrLevel.Get(sx, sy-1)
		if sy > 0 && !belowTile.IsSolid() { // can move down
			currScore := lr.scanDown(sx, startX, sy)
			if currScore < bestScore {
				bestScore = currScore
				actions.Direction = currDir
			}
		}
		nextTile = data.CurrLevel.Get(sx, sy)
		if nextTile.IsLadder() { // can move up a ladder
			currScore := lr.scanUp(sx, startX, sy)
			if currScore < bestScore {
				bestScore = currScore
				actions.Direction = currDir
			}
		}
		if currDir == data.Left {
			sx++
		} else {
			sx--
		}
	}
	return actions
}

func (lr *LRChase) scanDown(x, startX, startY int) int {
	y := startY
	_, py := world.WorldToMap(lr.Target.Object.Pos.X, lr.Target.Object.Pos.Y)
	for y > 0 && !data.CurrLevel.Get(x, y-1).IsSolid() { // can move down
		if x > 0 { // not at left edge, check left side
			nextBelow := data.CurrLevel.Get(x-1, y-1)
			nextTile := data.CurrLevel.Get(x-1, y)
			if y <= py && (nextBelow.IsSolid() ||
				nextBelow.IsLadder() ||
				nextTile.Block == data.BlockBar) { // can move left and below or level with runner
				break
			}
		}
		if x < data.CurrLevel.Metadata.Width-1 { // not at right edge, check right side
			nextBelow := data.CurrLevel.Get(x+1, y-1)
			nextTile := data.CurrLevel.Get(x+1, y)
			if y <= py && (nextBelow.IsSolid() ||
				nextBelow.IsLadder() ||
				nextTile.Block == data.BlockBar) { // can move right and below or level with runner
				break
			}
		}
		y--
	}

	// score
	if y == py {
		return util.Abs(startX - x)
	} else if y < py {
		return py - y + 200
	} else {
		return y - py + 100
	}
}

func (lr *LRChase) scanUp(x, startX, startY int) int {
	y := startY
	_, py := world.WorldToMap(lr.Target.Object.Pos.X, lr.Target.Object.Pos.Y)
	for y < data.CurrLevel.Metadata.Height-1 && data.CurrLevel.Get(x, y).IsLadder() { // while can go up
		y++
		if x > 0 { // not at left edge, check left side
			nextBelow := data.CurrLevel.Get(x-1, y-1)
			nextTile := data.CurrLevel.Get(x-1, y)
			if y >= py && (nextBelow.IsSolid() ||
				nextBelow.IsLadder() ||
				nextTile.Block == data.BlockBar) { // can move left and above or level with runner
				break
			}
		}
		if x < data.CurrLevel.Metadata.Width-1 { // not at right edge, check right side
			nextBelow := data.CurrLevel.Get(x+1, y-1)
			nextTile := data.CurrLevel.Get(x+1, y)
			if y >= py && (nextBelow.IsSolid() ||
				nextBelow.IsLadder() ||
				nextTile.Block == data.BlockBar) { // can move right and above or level with runner
				break
			}
		}
	}

	// score
	if y == py {
		return util.Abs(startX - x)
	} else if y < py {
		return py - y + 200
	} else {
		return y - py + 100
	}
}

// old version with a*
//func (lr *LRChase) GetActions() data.Actions {
//	actions := data.NewAction()
//	if data.CurrLevel == nil || lr.Ch == nil {
//		return actions
//	}
//	lr.Target = PickClosestPlayerYFirst(lr.Ch)
//	//if lr.Timer.Done() || lr.Target == nil {
//	//	lr.Target = PickClosestPlayerXFirst(lr.Ch)
//	//	lr.Timer = timing.New(constants.WaitToSwitch + rand.Float64()*3.)
//	//}
//	if lr.Target == nil {
//		return actions
//	}
//	x, y := world.WorldToMap(lr.Ch.Object.Pos.X, lr.Ch.Object.Pos.Y)
//	px, py := world.WorldToMap(lr.Target.Object.Pos.X, lr.Target.Object.Pos.Y)
//	// picked a player
//	// enemy is on the same y level as player
//	if py == y {
//		sx := x
//		for {
//			var nextTile, belowTile *data.Tile
//			nextTile = data.CurrLevel.Get(sx, y)
//			belowTile = data.CurrLevel.Get(sx, y-1)
//			if nextTile == nil {
//				// failed to find
//				fmt.Println("WARNING: LRChase searched off the map")
//				break
//			}
//			if nextTile.IsLadder() || nextTile.Block == data.BlockBar ||
//				!belowTile.IsEmpty() {
//				if sx < px {
//					sx++
//				} else if sx > px {
//					sx--
//				} else {
//					break
//				}
//			} else {
//				break
//			}
//		}
//		belowTile := data.CurrLevel.Get(x, y-1)
//		if lr.Ch.State == data.OnLadder &&
//			(math.Abs(lr.Target.Object.Pos.Y-lr.Ch.Object.Pos.Y) > 1. ||
//				(lr.Target.State != data.OnLadder &&
//					(belowTile.IsSolid() || belowTile.IsLadder()))) { // Enemy is on the same level as player, but is on a ladder and needs to adjust
//			if lr.Target.Object.Pos.Y > lr.Ch.Object.Pos.Y {
//				actions.PrevDirection = data.Up
//			} else if lr.Target.Object.Pos.Y < lr.Ch.Object.Pos.Y {
//				actions.PrevDirection = data.Down
//			}
//		}
//		if sx == px { // path to player found
//			if px > x { // player is right of enemy
//				actions.Direction = data.Right
//			} else if px < x { // player is left of enemy
//				actions.Direction = data.Left
//			} else if math.Abs(lr.Target.Object.Pos.X-lr.Ch.Object.Pos.X) > 2. { // player is in same tile as enemy
//				if lr.Target.Object.Pos.X > lr.Ch.Object.Pos.X {
//					actions.Direction = data.Right
//				} else if lr.Target.Object.Pos.X < lr.Ch.Object.Pos.X {
//					actions.Direction = data.Left
//				}
//			}
//		}
//		return actions
//	}
//
//	next := world.Coords{X: -1, Y: -1}
//	var bPath []astar.Pather
//	dy := py
//	data.PlayerAbove = py > y
//outerLoop:
//	for dy != y {
//		i := 0
//		ix := 0
//		for i < constants.PuzzleWidth {
//			dx := x + ix
//			if dx > -1 && dx < constants.PuzzleWidth {
//				path, _, found := astar.Path(data.CurrLevel.Get(x, y), data.CurrLevel.Get(dx, dy))
//				if len(path) > 1 && found {
//					bPath = path
//					next = path[len(path)-2].(*data.Tile).Coords
//					break outerLoop
//				}
//				i++
//			}
//			if ix >= 0 {
//				ix++
//			}
//			ix *= -1
//		}
//		if dy > y {
//			dy--
//		} else {
//			dy++
//		}
//	}
//	if next.X > -1 && next.Y > -1 {
//		if debug.ShowDebug {
//			col := color.RGBA{
//				R: 0,
//				G: 255,
//				B: 0,
//				A: 255,
//			}
//			var p *data.Tile
//			for i, t := range bPath {
//				tile := t.(*data.Tile)
//				if i > 0 {
//					debug.AddLine(col, imdraw.RoundEndShape, p.Object.Pos, tile.Object.Pos, 2)
//					col.R += 40
//				} else {
//					debug.AddLine(col, imdraw.RoundEndShape, tile.Object.Pos, tile.Object.Pos, 3)
//				}
//				p = tile
//			}
//		}
//		mPos := world.MapToWorld(next)
//		mPos = mPos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
//		if next.X < x {
//			actions.Direction = data.Left
//		} else if next.X > x {
//			actions.Direction = data.Right
//		} else if next.Y > y && lr.Ch.Object.Pos.Y < mPos.Y {
//			actions.Direction = data.Up
//		} else if next.Y < y && lr.Ch.Object.Pos.Y > mPos.Y {
//			actions.Direction = data.Down
//		}
//		if lr.Ch.Object.Pos.Y < mPos.Y {
//			actions.PrevDirection = data.Up
//		} else if lr.Ch.Object.Pos.Y > mPos.Y {
//			actions.PrevDirection = data.Down
//		} else if lr.Ch.Object.Pos.X > mPos.X {
//			actions.PrevDirection = data.Left
//		} else if lr.Ch.Object.Pos.X < mPos.X {
//			actions.PrevDirection = data.Right
//		}
//	}
//	return actions
//}
