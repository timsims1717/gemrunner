package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func CollisionSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC && !obj.Hidden {
			setCollisionFlags(ch)
			currPos := ch.FauxObj.Pos
			x, y := world.WorldToMap(currPos.X, currPos.Y)
			currTile := data.CurrLevel.Tiles.Get(x, y)
			leftTile := data.CurrLevel.Tiles.Get(x-1, y)
			rightTile := data.CurrLevel.Tiles.Get(x+1, y)
			upTile := data.CurrLevel.Tiles.Get(x, y+1)
			downTile := data.CurrLevel.Tiles.Get(x, y-1)
			if currTile == nil {
				outsideOfMap(ch, x, y)
			} else {
				// check each direction for collision
				nonFloorCollision(ch, currTile, leftTile, rightTile, upTile, currPos)
				standOnBelow := !ch.Flags.DropDown && standOnSystem(downTile)
				touchingFloor := currPos.Y-obj.HalfHeight <= currTile.Object.Pos.Y-world.HalfSize && !ch.Flags.HighJump && !ch.Flags.LongJump
				if (downTile == nil || downTile.Solid() ||
					(downTile.Ladder && !ch.Flags.OnLadder) ||
					standOnBelow && !ch.Flags.OnLadder) && touchingFloor {
					ch.Flags.Floor = true
					ch.FauxObj.Pos.Y = currTile.Object.Pos.Y - world.HalfSize + obj.HalfHeight
				}
				// check if the character can use a ladder
				ch.Flags.LadderHere = currTile.Ladder
				ch.Flags.LadderDown = downTile != nil && downTile.Ladder
				// check if the character is on a ladder
				//ch.Flags.OnLadder = (ch.Flags.OnLadder && !ch.Flags.Floor) || (ch.Flags.LadderHere && !ch.Flags.Floor)
				ch.Flags.OnLadder = ch.Flags.OnLadder && !ch.Flags.Floor && (ch.Flags.LadderHere || ch.Flags.LadderDown)
				// check to see if the character can run along the floor
				// They can run if they are touching solid ground (Floor)
				//   or are at the top of a ladder or the ladder is on top of Turf
				//   or are holding a ladder and there is floor directly below
				ch.Flags.CanRun = ((downTile == nil || downTile.Solid() || standOnBelow) && touchingFloor) ||
					(ch.Flags.LadderDown && touchingFloor && (!ch.Flags.LadderHere || downTile.Block == data.BlockTurf)) ||
					((downTile == nil || downTile.Solid() || downTile.Block == data.BlockTurf) && ch.Flags.OnLadder)
				// reset drop down
				ch.Flags.DropDown = false
			}
			//debug.AddTruthText("LeftWall:  ", ch.Flags.LeftWall)
			//debug.AddTruthText("RightWall: ", ch.Flags.RightWall)
			//debug.AddTruthText("Ceiling:   ", ch.Flags.Ceiling)
			//debug.AddTruthText("Floor:     ", ch.Flags.Floor)
			//debug.AddTruthText("On Ladder: ", ch.Flags.OnLadder)
			//debug.AddTruthText("Ladder Here: ", ch.Flags.LadderHere)
			//debug.AddTruthText("Ladder Down: ", ch.Flags.LadderDown)
		}
	}
}

func setCollisionFlags(ch *data.Dynamic) {
	ch.Flags.LeftWall = false
	ch.Flags.RightWall = false
	ch.Flags.Ceiling = false
	ch.Flags.Floor = false
	ch.Flags.CanRun = false
}

func outsideOfMap(ch *data.Dynamic, x, y int) {
	// uh oh
	ch.Flags.Ceiling = true
	ch.Flags.Floor = true
	ch.Flags.CanRun = true
	v := world.MapToWorld(world.Coords{X: x, Y: y})
	ch.FauxObj.Pos.Y = v.Y
}

func nonFloorCollision(ch *data.Dynamic, tile, leftTile, rightTile, upTile *data.Tile, chPos pixel.Vec) {
	// for left and right, we stop the character if the next tile is solid and either
	//   if they run into the tile
	//   or if they are on a ladder (so they stay in the center of the ladder)
	if (leftTile == nil || leftTile.Solid()) &&
		(chPos.X-ch.Object.HalfWidth <= tile.Object.Pos.X-world.HalfSize || ch.Flags.OnLadder) {
		ch.Flags.LeftWall = true
		ch.FauxObj.Pos.X = tile.Object.Pos.X - world.HalfSize + ch.Object.HalfWidth
	}
	if (rightTile == nil || rightTile.Solid()) &&
		(chPos.X+ch.Object.HalfWidth >= tile.Object.Pos.X+world.HalfSize || ch.Flags.OnLadder) {
		ch.Flags.RightWall = true
		ch.FauxObj.Pos.X = tile.Object.Pos.X + world.HalfSize - ch.Object.HalfWidth
	}
	// for up, we just make sure they don't enter a solid tile from below
	if (upTile == nil || upTile.Solid()) && chPos.Y+ch.Object.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
		ch.Flags.Ceiling = true
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y + world.HalfSize - ch.Object.HalfHeight
	}
}

func standOnSystem(downTile *data.Tile) bool {
	if downTile == nil {
		return false
	}
	for _, result := range myecs.Manager.Query(myecs.IsStandOn) {
		if obj, okO := result.Components[myecs.Object].(*object.Object); okO {
			x, y := world.WorldToMap(obj.Pos.X, obj.Pos.Y)
			standCoords := world.Coords{X: x, Y: y}
			if downTile.Coords == standCoords {
				return true
			}
		}
	}
	return false
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
		if !ch.Held.HasComponent(myecs.Parent) {
			ch.Held.AddComponent(myecs.Parent, ch.Object)
		}
		oldXOffset := ch.HeldObj.Offset.X
		oldYOffset := ch.HeldObj.Offset.Y
		//ch.HeldObj.Offset.X = 0
		//ch.HeldObj.Offset.Y = 0
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
		if tile == nil || tile.Solid() {
			ch.HeldObj.Offset.X = oldXOffset
			ch.HeldObj.Offset.Y = oldYOffset
			ch.Flags.Drop = true
		} else {
			if (leftTile == nil || leftTile.Solid()) &&
				heldPos.X-ch.HeldObj.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
				//ch.HeldObj.Offset.X += tile.Object.Pos.X - world.HalfSize + ch.HeldObj.HalfWidth - heldPos.X
				ch.HeldObj.Offset.X += ch.HeldObj.HalfWidth - world.HalfSize + tile.Object.Pos.X - heldPos.X
			}
			if (rightTile == nil || rightTile.Solid()) &&
				(heldPos.X+ch.HeldObj.HalfWidth >= tile.Object.Pos.X+world.HalfSize) {
				ch.HeldObj.Offset.X -= ch.HeldObj.HalfWidth - world.HalfSize - tile.Object.Pos.X + heldPos.X
			}
			if (upTile == nil || upTile.Solid()) &&
				heldPos.Y+ch.HeldObj.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
				ch.HeldObj.Offset.Y -= ch.HeldObj.HalfHeight - world.HalfSize - tile.Object.Pos.Y + heldPos.Y
			}
			if (downTile == nil || downTile.Solid()) &&
				(heldPos.Y-ch.HeldObj.HalfHeight <= tile.Object.Pos.X-world.HalfSize) {
				ch.HeldObj.Offset.Y += ch.HeldObj.HalfHeight - world.HalfSize + tile.Object.Pos.Y - heldPos.Y
			}
		}
	}
}
