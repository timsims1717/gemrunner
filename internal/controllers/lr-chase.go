package controllers

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/world"
	"github.com/beefsack/go-astar"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"image/color"
	"math"
)

type LRChase struct {
	Ch     *data.Dynamic
	Target *data.Dynamic
	Entity *ecs.Entity
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
			if nextTile.IsLadder() ||
				belowTile == nil || belowTile.IsSolid() || belowTile.IsLadder() {
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
		belowTile := data.CurrLevel.Tiles.Get(x, y-1)
		if lr.Ch.State == data.OnLadder &&
			(math.Abs(lr.Target.Object.Pos.Y-lr.Ch.Object.Pos.Y) > 1. ||
				(lr.Target.State != data.OnLadder &&
					(belowTile == nil || belowTile.IsSolid() || belowTile.IsLadder()))) { // Enemy is on the same level as player, but is on a ladder and needs to adjust
			if lr.Target.Object.Pos.Y > lr.Ch.Object.Pos.Y {
				actions.PrevDirection = data.Up
			} else if lr.Target.Object.Pos.Y < lr.Ch.Object.Pos.Y {
				actions.PrevDirection = data.Down
			}
		}
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
		}
		return actions
	}

	//bestT := -1
	bestD := -1
	next := world.Coords{X: -1, Y: -1}
	var bPath []astar.Pather
	for i := 0; i < constants.PuzzleWidth; i++ {
		path, d, found := astar.Path(data.CurrLevel.Tiles.Get(x, y), data.CurrLevel.Tiles.Get(i, py))
		if len(path) > 1 && found && (bestD == -1 || int(d) < bestD) {
			bestD = int(d)
			//bestT = i
			bPath = path
			next = path[len(path)-2].(*data.Tile).Coords
		}
	}
	if debug.ShowDebug {
		col := color.RGBA{
			R: 0,
			G: 255,
			B: 0,
			A: 255,
		}
		var p *data.Tile
		for i, t := range bPath {
			tile := t.(*data.Tile)
			if i > 0 {
				debug.AddLine(col, imdraw.RoundEndShape, p.Object.Pos, tile.Object.Pos, 2)
				col.R += 40
			} else {
				debug.AddLine(col, imdraw.RoundEndShape, tile.Object.Pos, tile.Object.Pos, 3)
			}
			p = tile
		}
	}
	if next.X > -1 && next.Y > -1 {
		mPos := world.MapToWorld(next)
		mPos = mPos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
		if next.X < x {
			actions.Direction = data.Left
		} else if next.X > x {
			actions.Direction = data.Right
		} else if next.Y > y && lr.Ch.Object.Pos.Y < mPos.Y {
			actions.Direction = data.Up
		} else if next.Y < y && lr.Ch.Object.Pos.Y > mPos.Y {
			actions.Direction = data.Down
		}
		if lr.Ch.Object.Pos.Y < mPos.Y {
			actions.PrevDirection = data.Up
		} else if lr.Ch.Object.Pos.Y > mPos.Y {
			actions.PrevDirection = data.Down
		} else if lr.Ch.Object.Pos.X > mPos.X {
			actions.PrevDirection = data.Left
		} else if lr.Ch.Object.Pos.X < mPos.X {
			actions.PrevDirection = data.Right
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
//		if nextTile == nil || nextTile.IsSolid() {
//			break
//		}
//		belowTile := data.CurrLevel.Tiles.Get(sx-1, sy-1)
//		if nextTile.Ladder ||
//			belowTile == nil || belowTile.IsSolid() || belowTile.Ladder {
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
//		if nextTile == nil || nextTile.IsSolid() {
//			break
//		}
//		belowTile := data.CurrLevel.Tiles.Get(sx+1, sy-1)
//		if nextTile.Ladder ||
//			belowTile == nil || belowTile.IsSolid() || belowTile.Ladder {
//			sx++
//		} else {
//			sx++
//			break
//		}
//	}
//	right = sx
//
//}
