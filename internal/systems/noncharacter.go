package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func DynamicSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		d, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		isControlled := result.Entity.HasComponent(myecs.Controller)
		if okO && okC && !obj.Hidden && !isControlled && data.CurrLevel.Start {
			if !result.Entity.HasComponent(myecs.Parent) {
				if reanimator.FrameSwitch {
					currPos := d.Object.Pos
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					currTile := data.CurrLevel.Tiles.Get(x, y)
					if !d.Flags.Floor {
						falling(d, currTile)
					}
				}
			}
		}
	}
}

type pickUpCheck struct {
	priority int
	cycle    int
	entity   *ecs.Entity
	obj      *object.Object
}

func AttemptPickUp(ch *data.Dynamic, p int, facingLeft bool) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	cx, cy := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
	chCoords := world.Coords{X: cx, Y: cy}
	leftTile := data.CurrLevel.Tiles.Get(cx-1, cy)
	rightTile := data.CurrLevel.Tiles.Get(cx+1, cy)
	upTile := data.CurrLevel.Tiles.Get(cx, cy+1)
	holdUp := !(upTile == nil || upTile.SolidV())
	holdSide := (facingLeft && !(leftTile == nil || leftTile.SolidH())) || (!facingLeft && !(rightTile == nil || rightTile.SolidH()))
	if !holdUp && !holdSide {
		return
	}
	// PickUp location priority:
	//  Same Space,
	//  left/right (in facing direction),
	//  down, then up.
	// No picking up things behind you, above you, or diagonal
	// PickUp player priority is used to allow players to cycle through things to pick up
	// If the player priority for all items is 0, the item priority decides
	//  Item priority generally starts with items like bombs and keys, and then to boxes.
	// Whenever an item is picked up with other items, its priority increases
	// Whenever an item is picked up without other items, its priority drops to 0
	var heldEntity pickUpCheck
	var sameSpace, facing, down, up, cycle bool
	for _, result := range myecs.Manager.Query(myecs.IsPickUp) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		pu, okP := result.Components[myecs.PickUp].(*data.PickUp)
		if okO && okP && !obj.Hidden {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			pickUpCoords := world.Coords{X: x, Y: y}
			// Check:
			//  No higher priority spaces
			//  correct space
			//  Either:
			//    No held item found yet
			//    or nothing in this space priority
			//    or cycle priority is lower than current cycle priority
			//    or item priority is lower than current item priority
			if chCoords == pickUpCoords &&
				(heldEntity.entity == nil || !sameSpace ||
					pu.Cycle[p] < heldEntity.cycle ||
					pu.Priority < heldEntity.priority) {
				if sameSpace {
					cycle = true
				} else {
					cycle = false
				}
				sameSpace = true
				heldEntity = pickUpCheck{
					priority: pu.Priority,
					cycle:    pu.Cycle[p],
					entity:   result.Entity,
					obj:      obj,
				}
			} else if !sameSpace &&
				pickUpCoords.Y == chCoords.Y &&
				((facingLeft && chCoords.X-1 == pickUpCoords.X) ||
					(!facingLeft && chCoords.X+1 == pickUpCoords.X)) &&
				(heldEntity.entity == nil || !facing ||
					pu.Cycle[p] < heldEntity.cycle ||
					pu.Priority < heldEntity.priority) {
				if facing {
					cycle = true
				} else {
					cycle = false
				}
				facing = true
				heldEntity = pickUpCheck{
					priority: pu.Priority,
					cycle:    pu.Cycle[p],
					entity:   result.Entity,
					obj:      obj,
				}
			} else if !sameSpace &&
				!facing &&
				chCoords.Y-1 == pickUpCoords.Y &&
				chCoords.X == pickUpCoords.X &&
				(heldEntity.entity == nil || !down ||
					pu.Cycle[p] < heldEntity.cycle ||
					pu.Priority < heldEntity.priority) {
				if down {
					cycle = true
				} else {
					cycle = false
				}
				down = true
				heldEntity = pickUpCheck{
					priority: pu.Priority,
					cycle:    pu.Cycle[p],
					entity:   result.Entity,
					obj:      obj,
				}
			} else if !sameSpace &&
				!facing &&
				!down &&
				chCoords.Y+1 == pickUpCoords.Y &&
				chCoords.X == pickUpCoords.X &&
				(heldEntity.entity == nil || !up ||
					pu.Cycle[p] < heldEntity.cycle ||
					pu.Priority < heldEntity.priority) {
				if up {
					cycle = true
				} else {
					cycle = false
				}
				up = true
				heldEntity = pickUpCheck{
					priority: pu.Priority,
					cycle:    pu.Cycle[p],
					entity:   result.Entity,
					obj:      obj,
				}
			}
		}
	}
	if heldEntity.entity != nil {
		if pu, ok := heldEntity.entity.GetComponentData(myecs.PickUp); ok {
			pickUp := pu.(*data.PickUp)
			if cycle {
				pickUp.Cycle[p]++
			} else {
				pickUp.Cycle[p] = 0
			}
			//if !heldEntity.entity.HasComponent(myecs.Dynamic) {
			//	d := data.NewDynamic()
			//	d.Object = heldEntity.obj
			//	d.FauxObj = object.New()
			//	d.FauxObj.Pos = d.Object.Pos
			//	d.Entity = heldEntity.entity
			//	d.Entity.AddComponent(myecs.Dynamic, d)
			//}
			ch.Held = heldEntity.entity
			ch.Flags.PickUp = true
			if holdUp {
				ch.Flags.HoldUp = true
			} else {
				ch.Flags.HoldSide = true
			}
			ch.HeldObj = heldEntity.obj
			ch.HeldObj.IntA = ch.HeldObj.Layer
			ch.HeldObj.Layer = ch.Object.Layer - 1
			ch.Flags.HeldNFlip = pickUp.NeverFlip
			ch.Flags.HeldFlip = ch.Object.Flip && !ch.HeldObj.Flip
		}
	}
}

func DropItem(ch *data.Dynamic) {
	if ch.Held != nil {
		// change the layer back
		ch.HeldObj.Layer = ch.HeldObj.IntA
		// normalize the pos (so it drops where you held it)
		ch.HeldObj.Pos = ch.HeldObj.Pos.Add(ch.HeldObj.Offset)
		ch.HeldObj.Offset = pixel.ZV
		// change the pos to center on a tile
		x, y := world.WorldToMap(ch.HeldObj.Pos.X, ch.HeldObj.Pos.Y)
		t := data.CurrLevel.Tiles.Get(x, y)
		if t != nil {
			ch.HeldObj.Pos = t.Object.Pos
		}
		// "drop it"
		ch.HeldObj = nil
		ch.Held.RemoveComponent(myecs.Parent)
		ch.Held = nil
		ch.Flags.PickUp = false
		ch.Flags.HoldUp = false
		ch.Flags.HoldSide = false
		ch.Flags.Drop = false
	}
}

func DoAction(ch *data.Dynamic) {
	if ch.Player > -1 && ch.Player < constants.MaxPlayers &&
		ch.Held != nil && ch.Held.HasComponent(myecs.Action) {
		if fnA, ok := ch.Held.GetComponentData(myecs.Action); ok {
			if colFn, okC := fnA.(*data.Interact); okC {
				colFn.Fn(data.CurrLevel, int(ch.Player), ch, ch.Held)
			}
		}
	}
}

func updateHeldItem(ch *data.Dynamic, facingLeft bool) {
	if !ch.Flags.PickUp && (ch.Flags.HoldUp || ch.Flags.HoldSide) {
		if !ch.Flags.HeldNFlip {
			if ch.Flags.HeldFlip {
				ch.HeldObj.Flip = !ch.Object.Flip
			} else {
				ch.HeldObj.Flip = ch.Object.Flip
			}
		}
		oldXOffset := ch.HeldObj.Offset.X
		oldYOffset := ch.HeldObj.Offset.Y
		oldXPos := ch.HeldObj.Pos.X
		oldYPos := ch.HeldObj.Pos.Y
		if !ch.Held.HasComponent(myecs.Parent) {
			ch.Held.AddComponent(myecs.Parent, ch.Object)
			oldXPos = ch.Object.Pos.X
			oldYPos = ch.Object.Pos.Y
		}
		if ch.Flags.HoldUp {
			ch.HeldObj.Offset.Y = ch.Object.HalfHeight + ch.HeldObj.HalfHeight - 1
			ch.HeldObj.Offset.X = 0
		} else if ch.Flags.HoldSide {
			offsetX := ch.Object.HalfWidth + ch.HeldObj.HalfHeight - 1
			if facingLeft {
				ch.HeldObj.Offset.X = -offsetX
			} else {
				ch.HeldObj.Offset.X = offsetX
			}
			ch.HeldObj.Offset.Y = 0
		}
		heldPos := ch.Object.Pos.Add(ch.HeldObj.Offset)
		x, y := world.WorldToMap(heldPos.X, heldPos.Y)
		tile := data.CurrLevel.Tiles.Get(x, y)
		leftTile := data.CurrLevel.Tiles.Get(x-1, y)
		rightTile := data.CurrLevel.Tiles.Get(x+1, y)
		upTile := data.CurrLevel.Tiles.Get(x, y+1)
		downTile := data.CurrLevel.Tiles.Get(x, y-1)
		if tile == nil || tile.SolidV() || tile.SolidH() {
			dPos := pixel.V(oldXPos, oldYPos).Add(pixel.V(oldXOffset, oldYOffset))
			dx, dy := world.WorldToMap(dPos.X, dPos.Y)
			dTile := data.CurrLevel.Tiles.Get(dx, dy)
			if !(tile != nil && tile.Block == data.BlockFall &&
				dTile != nil && dTile.Coords.X == tile.Coords.X) {
				ch.Flags.Drop = true
				ch.HeldObj.Offset.X = oldXOffset
				ch.HeldObj.Offset.Y = oldYOffset
				ch.HeldObj.Pos.X = oldXPos
				ch.HeldObj.Pos.Y = oldYPos
				if dTile == nil || dTile.SolidV() {
					ch.HeldObj.Pos = ch.Object.Pos
					ch.HeldObj.Offset = pixel.ZV
				}
				return
			}
		}
		if (leftTile == nil || leftTile.SolidH()) &&
			heldPos.X-ch.HeldObj.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
			//ch.HeldObj.Offset.X += tile.Object.Pos.X - world.HalfSize + ch.HeldObj.HalfWidth - heldPos.X
			ch.HeldObj.Offset.X += ch.HeldObj.HalfWidth - world.HalfSize + tile.Object.Pos.X - heldPos.X
		}
		if (rightTile == nil || rightTile.SolidH()) &&
			(heldPos.X+ch.HeldObj.HalfWidth >= tile.Object.Pos.X+world.HalfSize) {
			ch.HeldObj.Offset.X -= ch.HeldObj.HalfWidth - world.HalfSize - tile.Object.Pos.X + heldPos.X
		}
		if (upTile == nil || upTile.SolidV()) &&
			heldPos.Y+ch.HeldObj.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
			ch.HeldObj.Offset.Y -= ch.HeldObj.HalfHeight - world.HalfSize - tile.Object.Pos.Y + heldPos.Y
		}
		if (downTile == nil || downTile.SolidV()) &&
			(heldPos.Y-ch.HeldObj.HalfHeight <= tile.Object.Pos.X-world.HalfSize) {
			ch.HeldObj.Offset.Y += ch.HeldObj.HalfHeight - world.HalfSize + tile.Object.Pos.Y - heldPos.Y
		}
	}
}
