package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func LiftOrDropItem(ch *data.Dynamic, p int, throw bool) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	if ch.Held != nil {
		DropLift(ch, throw)
	} else {
		AttemptLift(ch, p, ch.Object.Flip)
	}
}

func AttemptLift(ch *data.Dynamic, p int, facingLeft bool) {
	if p < 0 || p >= constants.MaxPlayers {
		return
	}
	cx, cy := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
	chCoords := world.Coords{X: cx, Y: cy}
	upTile := data.CurrLevel.Tiles.Get(cx, cy+1)
	if upTile == nil || upTile.IsSolid() {
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
	var sameSpace, facing, down, up, cycle, holdingYou bool
	for _, result := range myecs.Manager.Query(myecs.IsPickUp) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		pu, okP := result.Components[myecs.PickUp].(*data.PickUp)
		if result.Entity.HasComponent(myecs.Dynamic) {
			if d, okD := result.Entity.GetComponentData(myecs.Dynamic); okD {
				dyn := d.(*data.Dynamic)
				holdingYou = dyn.Held != nil && dyn.HeldObj.ID == ch.Object.ID
			}
		}
		if okO && okP && !obj.Hidden &&
			obj.ID != ch.Object.ID &&
			!(pu.Held == p || pu.Inventory > -1) &&
			!holdingYou {
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
			// check if item is held by someone else, and drop it for them
			if pickUp.Held > -1 {
				DropLift(data.CurrLevel.Players[pickUp.Held], false)
			}
			if cycle {
				pickUp.Cycle[p]++
			} else {
				pickUp.Cycle[p] = 0
			}
			pickUp.Inventory = -1
			pickUp.Held = p
			if c, okC := heldEntity.entity.GetComponentData(myecs.Dynamic); okC {
				chHeld := c.(*data.Dynamic)
				DropLift(chHeld, false)
				chHeld.State = data.Carried
			}
			ch.Held = heldEntity.entity
			ch.Flags.PickUp = true
			ch.Object.Pos.X = ch.Object.LastPos.X
			ch.HeldObj = heldEntity.obj
			ch.HeldObj.IntA = ch.HeldObj.Layer
			ch.HeldObj.Layer = ch.Object.Layer - 1
			ch.Flags.HeldNFlip = pickUp.NeverFlip
			ch.Flags.HeldFlip = ch.Object.Flip && !ch.HeldObj.Flip
		}
	}
}

func DropLift(ch *data.Dynamic, throw bool) {
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
		if throw {
			ch.HeldObj.Pos.X = t.Object.Pos.X
			ch.HeldObj.Pos.Y = t.Object.Pos.Y
			// update the state
			if c, okC := ch.Held.GetComponentData(myecs.Dynamic); okC {
				chHeld := c.(*data.Dynamic)
				chHeld.Flags.Throw = true
				chHeld.State = data.Thrown
				chHeld.ACounter = 0
				chHeld.ATimer = nil
				chHeld.Flags.Climbed = false
				chHeld.Flags.GoingUp = false
				if ch.Actions.Left() {
					chHeld.Flags.JumpL = true
				} else if ch.Actions.Right() {
					chHeld.Flags.JumpR = true
				}
			}
		} else {
			// update the state
			if c, okC := ch.Held.GetComponentData(myecs.Dynamic); okC {
				chHeld := c.(*data.Dynamic)
				chHeld.State = data.Falling
			}
		}
		// update the pickup data
		if p, okP := ch.Held.GetComponentData(myecs.PickUp); okP {
			pickUp := p.(*data.PickUp)
			pickUp.Inventory = -1
			pickUp.Held = -1
		}
		// "drop it"
		ch.HeldObj = nil
		ch.Held.RemoveComponent(myecs.Parent)
		ch.Held = nil
		ch.Flags.PickUp = false
		//ch.Flags.Drop = false
	}
}

func updateHeldItem(ch *data.Dynamic, facingLeft bool) {
	if ch.Held != nil {
		if !ch.Flags.PickUp {
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
			ch.HeldObj.Offset.Y = ch.Object.HalfHeight + ch.HeldObj.HalfHeight - 1
			ch.HeldObj.Offset.X = 0
			heldPos := ch.Object.Pos.Add(ch.HeldObj.Offset)
			x, y := world.WorldToMap(heldPos.X, heldPos.Y)
			tile := data.CurrLevel.Tiles.Get(x, y)
			leftTile := data.CurrLevel.Tiles.Get(x-1, y)
			rightTile := data.CurrLevel.Tiles.Get(x+1, y)
			upTile := data.CurrLevel.Tiles.Get(x, y+1)
			downTile := data.CurrLevel.Tiles.Get(x, y-1)
			if tile == nil || tile.IsSolid() || tile.IsSolid() {
				dPos := pixel.V(oldXPos, oldYPos).Add(pixel.V(oldXOffset, oldYOffset))
				dx, dy := world.WorldToMap(dPos.X, dPos.Y)
				dTile := data.CurrLevel.Tiles.Get(dx, dy)
				if !(tile != nil && tile.Block == data.BlockFall &&
					dTile != nil && dTile.Coords.X == tile.Coords.X) {
					//ch.Flags.Drop = true
					ch.HeldObj.Offset.X = oldXOffset
					ch.HeldObj.Offset.Y = oldYOffset
					ch.HeldObj.Pos.X = oldXPos
					ch.HeldObj.Pos.Y = oldYPos
					if dTile == nil || dTile.IsSolid() {
						ch.HeldObj.Pos = ch.Object.Pos
						ch.HeldObj.Offset = pixel.ZV
					}
					DropLift(ch, false)
					return
				}
			}
			if (leftTile == nil || leftTile.IsSolid()) &&
				heldPos.X-ch.HeldObj.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
				//ch.HeldObj.Offset.X += tile.Object.Pos.X - world.HalfSize + ch.HeldObj.HalfWidth - heldPos.X
				ch.HeldObj.Offset.X += ch.HeldObj.HalfWidth - world.HalfSize + tile.Object.Pos.X - heldPos.X
			}
			if (rightTile == nil || rightTile.IsSolid()) &&
				(heldPos.X+ch.HeldObj.HalfWidth >= tile.Object.Pos.X+world.HalfSize) {
				ch.HeldObj.Offset.X -= ch.HeldObj.HalfWidth - world.HalfSize - tile.Object.Pos.X + heldPos.X
			}
			if (upTile == nil || upTile.IsSolid()) &&
				heldPos.Y+ch.HeldObj.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
				ch.HeldObj.Offset.Y -= ch.HeldObj.HalfHeight - world.HalfSize - tile.Object.Pos.Y + heldPos.Y
			}
			if (downTile == nil || downTile.IsSolid()) &&
				(heldPos.Y-ch.HeldObj.HalfHeight <= tile.Object.Pos.X-world.HalfSize) {
				ch.HeldObj.Offset.Y += ch.HeldObj.HalfHeight - world.HalfSize + tile.Object.Pos.Y - heldPos.Y
			}
		}
	}
}

func DoAction(ch *data.Dynamic) {
	if ch.Player > -1 && ch.Player < constants.MaxPlayers &&
		ch.Held != nil && ch.Held.HasComponent(myecs.Action) {
		if fnA, ok := ch.Held.GetComponentData(myecs.Action); ok {
			if colFn, okC := fnA.(*data.Interact); okC {
				ch.Flags.Using = true
				colFn.Fn(data.CurrLevel, int(ch.Player), ch, ch.Held)
			}
		}
	}
}
