package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
)

func DynamicSystem() {
	if reanimator.FrameSwitch {
		for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
			obj, okO := result.Components[myecs.Object].(*object.Object)
			d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			isControlled := result.Entity.HasComponent(myecs.Controller)
			if okO && okC && !obj.Hidden && !isControlled && data.CurrLevel.Start {
				d.ACounter++
				if !result.Entity.HasComponent(myecs.Parent) {
					currPos := d.Object.Pos
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					currTile := data.CurrLevel.Get(x, y)
					if !d.Flags.Floor {
						if d.Flags.Thrown {
							thrown(d, currTile)
						} else if !d.Flags.Flying {
							falling(d, currTile)
						}
					} else {
						d.Flags.Thrown = false
						d.Flags.JumpL = false
						d.Flags.JumpR = false
					}
				} else {
					currPos := d.Object.Pos.Add(d.Object.Offset)
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					tile := data.CurrLevel.Get(x, y)
					d.LastTile = tile
				}
			}
		}
	}
}

func PushySystem() {
	if reanimator.FrameSwitch {
		for _, result := range myecs.Manager.Query(myecs.IsPushy) {
			_, okO := result.Components[myecs.Object].(*object.Object)
			d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			push, okS := result.Components[myecs.Pushy].(*data.Pushy)
			if okO && okC && okS && !d.Object.Hidden && !d.Flags.Ignore {
				if PushNext(d, push) {
					for _, result1 := range myecs.Manager.Query(myecs.IsPushy) {
						_, okO1 := result1.Components[myecs.Object].(*object.Object)
						d1, okC1 := result1.Components[myecs.Dynamic].(*data.Dynamic)
						push1, okS1 := result1.Components[myecs.Pushy].(*data.Pushy)
						if okO1 && okC1 && okS1 && !d1.Object.Hidden && !d1.Flags.Ignore &&
							d1.Object.ID != d.Object.ID && d1.Type == "air_ring" && d.Type == "air_ring" &&
							push.Direction != push1.Direction &&
							d.Object.Rect.Moved(d.Object.Pos).Intersects(d1.Object.Rect.Moved(d1.Object.Pos)) {
							StopPushing(d, push)
							StopPushing(d1, push1)
						}
					}
				}
			}
		}
	}
}

func PushNext(d *data.Dynamic, push *data.Pushy) bool {
	if push.PushNext != nil {
		if push.PushNext.Pushing == nil {
			push.PushNext = nil
		} else {
			if !PushNext(push.PushNext.Pushing, push.PushNext) {
				StopPushing(d, push)
				return false
			} else {
				_, y := world.WorldToMap(d.Object.Pos.X, d.Object.Pos.Y)
				_, y1 := world.WorldToMap(push.PushNext.Pushing.Object.Pos.X, push.PushNext.Pushing.Object.Pos.Y)
				if !(y1 == y && !push.PushNext.Pushing.Flags.Goop &&
					(push.PushNext.Pushing.State == data.Grounded || push.PushNext.Pushing.State == data.Flying || push.PushNext.Pushing.State == data.Falling) &&
					((push.Direction == data.Left && d.Object.Pos.X > push.PushNext.Pushing.Object.Pos.X) || (push.Direction == data.Right && d.Object.Pos.X < push.PushNext.Pushing.Object.Pos.X))) {
					// no longer eligible to push
					StopPushing(push.PushNext.Pushing, push.PushNext)
					push.PushNext = nil
				}
			}
		}
	}
	dx := push.Speed
	if push.Direction == data.Left {
		dx = -dx
	}
	x, y := world.WorldToMap(d.Object.Pos.X, d.Object.Pos.Y)
	tile := data.CurrLevel.Get(x, y)
	if d.State == data.Falling && push.OrigTile != nil && push.OrigTile != tile {
		return true
	}
	xt, yt := world.WorldToMap(d.Object.Pos.X+dx, d.Object.Pos.Y)
	mTile := data.CurrLevel.Get(xt, yt)
	if push.Direction == data.Right {
		xt += 1
	} else {
		xt -= 1
	}
	// test whether we would run into a wall here
	nTile := data.CurrLevel.Get(xt, yt)
	if d.State == data.Falling && (nTile.IsSolid() || nTile.Block == data.BlockLiquid || (nTile.Block == data.BlockFall && !nTile.Flags.Collapse)) {
		StopPushing(d, push)
		return false
	} else if ((push.Direction == data.Right && d.Object.Pos.X+d.Object.HalfWidth >= mTile.Object.Pos.X+world.HalfSize) ||
		(push.Direction == data.Left && d.Object.Pos.X-d.Object.HalfWidth <= mTile.Object.Pos.X-world.HalfSize)) &&
		(nTile.IsSolid() || nTile.Block == data.BlockLiquid || (nTile.Block == data.BlockFall && !nTile.Flags.Collapse)) {
		StopPushing(d, push)
		return false
	}
	// move
	if !push.NoMove {
		d.Object.Pos.X += dx
	}
	if push.PushNext == nil {
		x, y := world.WorldToMap(d.Object.Pos.X, d.Object.Pos.Y)
		for _, result1 := range myecs.Manager.Query(myecs.IsDynamic) {
			obj1, okO1 := result1.Components[myecs.Object].(*object.Object)
			d1, okC1 := result1.Components[myecs.Dynamic].(*data.Dynamic)
			if okO1 && okC1 && d.Object.ID != obj1.ID && !d1.Flags.Ignore &&
				d1.Pushing == nil && !d1.Flags.NoPush && !d1.Flags.Goop &&
				(d1.State == data.Grounded || d1.State == data.Flying || d1.State == data.Falling ||
					d1.State == data.DoingAction) &&
				d.Object.Rect.Moved(d.Object.Pos).Intersects(obj1.Rect.Moved(obj1.Pos)) &&
				((push.Direction == data.Left && d.Object.Pos.X > obj1.Pos.X) || (push.Direction == data.Right && d.Object.Pos.X < obj1.Pos.X)) {
				if d1.State == data.DoingAction {
					switch d1.Flags.ItemAction {
					case data.FireFlamethrower, data.ThrowingGoop:
						d1.Flags.ItemAction = data.NoItemAction
					case data.UseAirCannon, data.Drilling, data.DrillStart, data.Hiding, data.TransportIn, data.TransportExit:
						continue
					}
				}
				x1, y1 := world.WorldToMap(obj1.Pos.X, obj1.Pos.Y)
				if util.Abs(x-x1) < 2 && y1 == y {
					//if push.Direction == data.Left {
					//	obj1.Pos.X = d.Object.Pos.X - d.Object.HalfWidth - obj1.HalfWidth + 0.01
					//} else {
					//	obj1.Pos.X = d.Object.Pos.X + d.Object.HalfWidth + obj1.HalfWidth - 0.01
					//}
					ot := data.CurrLevel.Get(x1, y1)
					newPush := &data.Pushy{
						Direction: push.Direction,
						Speed:     push.Speed,
						Pushing:   d1,
						OrigTile:  ot,
					}
					if !PushNext(d1, newPush) {
						StopPushing(d, push)
						return false
					}
					d1.Pushing = newPush
					push.PushNext = newPush
					return true
				}
			}
		}
	}
	return true
}

func StopPushing(d *data.Dynamic, push *data.Pushy) {
	if push.PushNext != nil {
		StopPushing(push.PushNext.Pushing, push.PushNext)
		push.PushNext = nil
	}
	d.Pushing = nil
	d.Entity.RemoveComponent(myecs.Pushy)
	x, y := world.WorldToMap(d.Object.Pos.X, d.Object.Pos.Y)
	tile := data.CurrLevel.Get(x, y)
	if tile != nil && !tile.IsSolid() {
		d.Object.SetPos(tile.Object.Pos)
	}
}

func SmashSystem() {
	if reanimator.FrameSwitch {
		for _, result := range myecs.Manager.Query(myecs.IsSmash) {
			obj, okO := result.Components[myecs.Object].(*object.Object)
			d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
			s, okS := result.Components[myecs.Smash].(*float64)
			if okO && okC && okS {
				if d.Flags.Floor ||
					obj.Hidden {
					*s = obj.Pos.Y
				}
			}
		}
	}
}

func thrown(ch *data.Dynamic, tile *data.Tile) {
	if (ch.ACounter > constants.ThrownCounter) ||
		(ch.Flags.LeftWall && ch.Flags.JumpL) ||
		(ch.Flags.RightWall && ch.Flags.JumpR) ||
		ch.Flags.Floor {
		ch.Flags.Thrown = false
		ch.Flags.JumpL = false
		ch.Flags.JumpR = false
	} else {
		if ch.Flags.JumpR {
			if !ch.Flags.RightWall {
				ch.Object.Pos.X += constants.ThrownHSpeed
			}
		} else if ch.Flags.JumpL {
			if !ch.Flags.LeftWall {
				ch.Object.Pos.X -= constants.ThrownHSpeed
			}
		}
		if tile.Coords != ch.LastTile.Coords {
			if util.Abs(tile.Coords.X-ch.LastTile.Coords.X) > 1 {
				ch.Object.Pos.Y -= constants.ThrownVSpeed
			}
		} else {
			ch.Object.Pos.Y += constants.ThrownVSpeed
		}
	}
}
