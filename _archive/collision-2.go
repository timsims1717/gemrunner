package _archive

//func CollisionSystemOld() {
//	for _, result := range myecs.Manager.Query(myecs.IsDynamic) {
//		obj, okO := result.Components[myecs.Object].(*object.Object)
//		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
//		if okO && okC && !obj.Hidden {
//			setCollisionFlagsOld(ch)
//			currPos := ch.Object.Pos
//			x, y := world.WorldToMap(currPos.X, currPos.Y)
//			currTile := data.CurrLevel.Tiles.Get(x, y)
//			leftTile := data.CurrLevel.Tiles.Get(x-1, y)
//			rightTile := data.CurrLevel.Tiles.Get(x+1, y)
//			upTile := data.CurrLevel.Tiles.Get(x, y+1)
//			downTile := data.CurrLevel.Tiles.Get(x, y-1)
//			if currTile == nil {
//				outsideOfMap(ch, x, y)
//			} else {
//				// check each direction for collision
//				nonFloorCollision(ch, currTile, leftTile, rightTile, upTile, currPos)
//				standOnBelow := !ch.Actions.Down() && standOnSystem(downTile)
//				touchingFloor := currPos.Y-obj.HalfHeight <= currTile.Object.Pos.Y-world.HalfSize && !ch.Flags.HighJump && !ch.Flags.LongJump
//				if (downTile == nil || downTile.Solid() ||
//					(downTile.Ladder && !ch.Flags.OnLadder) ||
//					standOnBelow && !ch.Flags.OnLadder) && touchingFloor {
//					ch.Flags.Floor = true
//					ch.Object.Pos.Y = currTile.Object.Pos.Y - world.HalfSize + obj.HalfHeight
//				}
//				// check if the character can use a ladder
//				ch.Flags.LadderHere = currTile.Ladder
//				ch.Flags.LadderDown = downTile != nil && downTile.Ladder
//				// check if the character is on a ladder
//				//ch.Flags.OnLadder = (ch.Flags.OnLadder && !ch.Flags.Floor) || (ch.Flags.LadderHere && !ch.Flags.Floor)
//				ch.Flags.OnLadder = ch.Flags.OnLadder && !ch.Flags.Floor && (ch.Flags.LadderHere || ch.Flags.LadderDown)
//				// check to see if the character can run along the floor
//				// They can run if they are touching solid ground (Floor)
//				//   or are at the top of a ladder or the ladder is on top of Turf
//				//   or are holding a ladder and there is floor directly below
//				ch.Flags.CanRun = ((downTile == nil || downTile.Solid() || standOnBelow) && touchingFloor) ||
//					(ch.Flags.LadderDown && touchingFloor && (!ch.Flags.LadderHere || downTile.Block == data.BlockTurf)) ||
//					((downTile == nil || downTile.Solid() || downTile.Block == data.BlockTurf) && ch.Flags.OnLadder)
//			}
//			//debug.AddTruthText("LeftWall:  ", ch.Flags.LeftWall)
//			//debug.AddTruthText("RightWall: ", ch.Flags.RightWall)
//			//debug.AddTruthText("Ceiling:   ", ch.Flags.Ceiling)
//			//debug.AddTruthText("Floor:     ", ch.Flags.Floor)
//			//debug.AddTruthText("On Ladder: ", ch.Flags.OnLadder)
//			//debug.AddTruthText("Ladder Here: ", ch.Flags.LadderHere)
//			//debug.AddTruthText("Ladder Down: ", ch.Flags.LadderDown)
//		}
//	}
//}
//
//func setCollisionFlagsOld(ch *data.Dynamic) {
//	ch.Flags.LeftWall = false
//	ch.Flags.RightWall = false
//	ch.Flags.Ceiling = false
//	ch.Flags.Floor = false
//	ch.Flags.CanRun = false
//}
//
//func nonFloorCollision(ch *data.Dynamic, tile, leftTile, rightTile, upTile *data.Tile, chPos pixel.Vec) {
//	// for left and right, we stop the character if the next tile is solid and either
//	//   if they run into the tile
//	//   or if they are on a ladder (so they stay in the center of the ladder)
//	if leftTile == nil || leftTile.Solid() {
//		if ch.Flags.OnLadder {
//			ch.Flags.LeftWall = true
//		} else if chPos.X-ch.Object.HalfWidth <= tile.Object.Pos.X-world.HalfSize {
//			ch.Flags.LeftWall = true
//			ch.Object.Pos.X = tile.Object.Pos.X - world.HalfSize + ch.Object.HalfWidth
//		}
//	}
//	if rightTile == nil || rightTile.Solid() {
//		if ch.Flags.OnLadder {
//			ch.Flags.RightWall = true
//		} else if chPos.X+ch.Object.HalfWidth >= tile.Object.Pos.X+world.HalfSize {
//			ch.Flags.RightWall = true
//			ch.Object.Pos.X = tile.Object.Pos.X + world.HalfSize - ch.Object.HalfWidth
//		}
//	}
//	// for up, we just make sure they don't enter a solid tile from below
//	if (upTile == nil || upTile.Solid()) && chPos.Y+ch.Object.HalfHeight >= tile.Object.Pos.Y+world.HalfSize {
//		ch.Flags.Ceiling = true
//		ch.Object.Pos.Y = tile.Object.Pos.Y + world.HalfSize - ch.Object.HalfHeight
//	}
//}
