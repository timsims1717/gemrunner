package _archive

//func CharacterSystem() {
//	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
//		obj, okO := result.Components[myecs.Object].(*object.Object)
//		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
//		ct, okT := result.Components[myecs.Controller].(data.Controller)
//		if okO && okC && okT && !obj.Hidden && data.CurrLevel.Start {
//			actions := ct.GetActions()
//			//if (!ch.Actions.Left && actions.Left) || (!ch.Actions.Right && actions.Right) {
//			//	ch.Flags.PLeftRight = true
//			//	ch.Flags.PUpDown = false
//			//}
//			//if (!ch.Actions.Up && actions.Up) || (!ch.Actions.Down && actions.Down) {
//			//	ch.Flags.PLeftRight = false
//			//	ch.Flags.PUpDown = true
//			//}
//			//if !actions.Left && !actions.Right && !actions.Up && !actions.Down {
//			//	ch.Flags.PLeftRight = false
//			//	ch.Flags.PUpDown = false
//			//}
//			if ch.Flags.Frame {
//				ch.Flags.Frame = false
//				ch.Actions = data.NewAction()
//			}
//			if actions.Direction != data.None {
//				ch.Actions.Direction = actions.Direction
//			}
//			if actions.PrevDirection != data.None {
//				ch.Actions.PrevDirection = actions.PrevDirection
//			}
//			ch.Actions.Jump = ch.Actions.Jump || actions.Jump
//			ch.Actions.PickUp = ch.Actions.PickUp || actions.PickUp
//			ch.Actions.Action = ch.Actions.Action || actions.Action
//			debug.AddText(fmt.Sprintf("Direction: %5s", ch.Actions.Direction))
//			debug.AddText(fmt.Sprintf("Previous:  %5s", ch.Actions.PrevDirection))
//			debug.AddTruthText("Jump:      ", ch.Actions.Jump)
//			debug.AddTruthText("PickUp:    ", ch.Actions.PickUp)
//			debug.AddTruthText("Action:    ", ch.Actions.Action)
//
//			if reanimator.FrameSwitch {
//				currPos := ch.Object.Pos
//				x, y := world.WorldToMap(currPos.X, currPos.Y)
//				currTile := data.CurrLevel.Tiles.Get(x, y)
//				ch.ACounter++
//				if !(ch.Flags.Hit || ch.Flags.Dead || ch.Flags.Attack) {
//					if ch.Flags.LeapOn || ch.Flags.LeapOff || ch.Flags.LeapTo {
//						leaping(ch, currTile)
//					} else if !ch.Flags.Floor &&
//						!ch.Flags.HighJump &&
//						!ch.Flags.LongJump &&
//						!ch.Flags.LadderHere &&
//						!ch.Flags.OnLadder {
//						falling(ch, currTile)
//					} else if ch.Flags.HighJump || ch.Flags.LongJump {
//						jumpingOld(ch, currTile)
//					} else {
//						upOrDown(ch, currTile)
//						if !(ch.Flags.Climbed &&
//							(ch.Actions.Direction == data.Up || ch.Actions.Direction == data.Down)) {
//							if ch.Flags.OnLadder {
//								gettingOffLadder(ch, currTile)
//							} else {
//								onTheGround(ch, currTile)
//							}
//						}
//					}
//					if ch.Actions.PickUp {
//						if ch.Flags.HoldUp || ch.Flags.HoldSide {
//							ch.Flags.Drop = true
//						} else if !ch.Flags.LeapOn && !ch.Flags.LeapOff &&
//							!ch.Flags.LeapTo && !ch.Flags.OnLadder &&
//							ch.Player > -1 && ch.Player < constants.MaxPlayers {
//							AttemptPickUp(ch, int(ch.Player), ch.Object.Flip)
//						}
//					} else if ch.Actions.Action {
//						if ch.Flags.HoldUp || ch.Flags.HoldSide {
//							//ch.Flags.Drop = true
//							ch.Flags.Action = true
//						}
//					}
//					if !ch.Flags.OnLadder {
//						ch.Flags.GoingUp = false
//						ch.Flags.Climbed = false
//					} else {
//						ch.Object.Pos.X = currTile.Object.Pos.X
//					}
//				}
//				//obj.Pos.X = ch.Object.Pos.X
//				//obj.Pos.Y = ch.Object.Pos.Y
//				//obj.Flip = ch.Object.Flip
//				updateHeldItem(ch, obj.Flip)
//				if ch.Flags.Dead || ch.Flags.Hit || ch.Flags.Attack ||
//					ch.Flags.OnLadder || ch.Flags.Drop || ch.Flags.Action {
//					if ch.Flags.Action && !(ch.Flags.Dead || ch.Flags.Hit || ch.Flags.Attack) {
//						DoAction(ch)
//						ch.Flags.Action = false
//					}
//					DropItem(ch)
//				}
//				ch.Flags.Frame = true
//			}
//		}
//	}
//}
//
//func upOrDown(ch *data.Dynamic, tile *data.Tile) {
//	if !ch.Flags.HoldUp && !ch.Flags.HoldSide {
//		if ch.Actions.Up() && !ch.Flags.Ceiling {
//			if ch.Flags.LadderHere {
//				if !ch.Flags.OnLadder {
//					ch.Control.ClearPrev()
//					ch.Actions.PrevDirection = data.None
//				}
//				ch.Object.Pos.Y += ch.Vars.ClimbSpeed
//				ch.Object.Flip = false
//				ch.Flags.OnLadder = true
//				ch.Flags.GoingUp = true
//				ch.Flags.Climbed = true
//			} else if ch.Flags.OnLadder {
//				ch.Object.Pos.Y = tile.Object.Pos.Y
//				ch.Object.Flip = false
//				ch.Flags.OnLadder = false
//				ch.Control.ClearPrev()
//				ch.Actions.PrevDirection = data.None
//			}
//		} else if ch.Actions.Down() {
//			if (ch.Flags.LadderHere && !ch.Flags.Floor) || ch.Flags.LadderDown {
//				if !ch.Flags.OnLadder {
//					ch.Control.ClearPrev()
//					ch.Actions.PrevDirection = data.None
//				}
//				ch.Object.Pos.Y -= ch.Vars.SlideSpeed
//				ch.Object.Flip = false
//				ch.Flags.OnLadder = true
//				ch.Flags.GoingUp = false
//				ch.Flags.Climbed = true
//			} else if ch.Flags.Floor {
//				ch.Control.ClearPrev()
//				ch.Actions.PrevDirection = data.None
//			}
//		}
//	} else {
//		if ch.Flags.HoldUp && ch.Actions.Down() && (ch.Actions.Left() || ch.Actions.Right()) {
//			ch.Flags.HoldUp = false
//			ch.Flags.HoldSide = true
//		} else if ch.Flags.HoldSide && ch.Actions.Up() {
//			ch.Flags.HoldUp = true
//			ch.Flags.HoldSide = false
//		}
//	}
//}
//
//func gettingOffLadder(ch *data.Dynamic, tile *data.Tile) {
//	if ch.Actions.Left() && !ch.Flags.LeftWall &&
//		(ch.Flags.CanRun || ch.ACounter > int(ch.Vars.LeapDelay)) {
//		ch.Flags.OnLadder = false
//		ch.Object.Flip = true
//		ch.Object.Pos.Y = tile.Object.Pos.Y
//		if ch.Flags.CanRun {
//			ch.Control.ClearPrev()
//			ch.Actions.PrevDirection = data.None
//			ch.Object.Pos.X -= ch.Vars.WalkSpeed
//		} else if ch.ACounter > int(ch.Vars.LeapDelay) {
//			lWall := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
//			if lWall != nil && lWall.Ladder {
//				ch.Flags.LeapTo = true
//			} else {
//				ch.Flags.LeapOff = true
//			}
//		}
//	} else if ch.Actions.Right() && !ch.Flags.RightWall &&
//		(ch.Flags.CanRun || ch.ACounter > int(ch.Vars.LeapDelay)) {
//		ch.Flags.OnLadder = false
//		ch.Object.Flip = false
//		ch.Object.Pos.Y = tile.Object.Pos.Y
//		if ch.Flags.CanRun {
//			ch.Control.ClearPrev()
//			ch.Actions.PrevDirection = data.None
//			ch.Object.Pos.X += ch.Vars.WalkSpeed
//		} else if ch.ACounter > int(ch.Vars.LeapDelay) {
//			rWall := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
//			if rWall != nil && rWall.Ladder {
//				ch.Flags.LeapTo = true
//			} else {
//				ch.Flags.LeapOff = true
//			}
//		}
//	}
//}
//
//func onTheGround(ch *data.Dynamic, tile *data.Tile) {
//	ch.LastTile = tile
//	if ch.Actions.Jump && ch.Flags.CanRun && !ch.Flags.Ceiling {
//		upLeft := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y+1)
//		upRight := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y+1)
//		left := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
//		right := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
//		// High Jump if:
//		//  the character is not going left or right
//		//  or they are going left/right and there is a wall left/right
//		//  or they are going left/right and there is a wall up left or up right
//		// Otherwise, it's a long jump
//		if (!ch.Actions.Left() && !ch.Actions.Right()) ||
//			(ch.Actions.Left() && (left == nil || left.Solid())) ||
//			(ch.Actions.Right() && (right == nil || right.Solid())) ||
//			(ch.Actions.Left() && (upLeft == nil || upLeft.Solid())) ||
//			(ch.Actions.Right() && (upRight == nil || upRight.Solid())) {
//			ch.Flags.HighJump = true
//			ch.Object.Pos.X = tile.Object.Pos.X
//			ch.Object.Pos.Y = tile.Object.Pos.Y + ch.Vars.HiJumpVSpeed
//			ch.ACounter = 0
//		} else {
//			ch.Flags.LongJump = true
//			if ch.Actions.Left() {
//				ch.Object.Pos.X = tile.Object.Pos.X - ch.Vars.LgJumpHSpeed
//			} else {
//				ch.Object.Pos.X = tile.Object.Pos.X + ch.Vars.LgJumpHSpeed
//			}
//			ch.ACounter = 0
//		}
//		// for both kinds of jumps
//		if ch.Actions.Left() {
//			ch.Flags.JumpL = true
//			ch.Object.Flip = true
//		} else if ch.Actions.Right() {
//			ch.Flags.JumpR = true
//			ch.Object.Flip = false
//		} else {
//			ch.Flags.JumpL = false
//			ch.Flags.JumpR = false
//		}
//	} else {
//		if !ch.Flags.PickUp {
//			if ch.Actions.Left() && !ch.Flags.LeftWall {
//				if ch.Flags.CanRun {
//					ch.Object.Pos.X -= ch.Vars.WalkSpeed
//					ch.Object.Flip = true
//					ch.Flags.OnLadder = false
//				} else if tile.Ladder {
//					ch.Flags.LeapOn = true
//				}
//			} else if ch.Actions.Right() && !ch.Flags.RightWall {
//				if ch.Flags.CanRun {
//					ch.Object.Pos.X += ch.Vars.WalkSpeed
//					ch.Object.Flip = false
//					ch.Flags.OnLadder = false
//				} else if tile.Ladder {
//					ch.Flags.LeapOn = true
//				}
//			}
//		}
//	}
//}
//
//func jumpingOld(ch *data.Dynamic, tile *data.Tile) {
//	if (ch.Flags.HighJump && ch.ACounter > int(ch.Vars.HiJumpTimer)) ||
//		(ch.Flags.LongJump && ch.ACounter > int(ch.Vars.LgJumpTimer)) ||
//		(ch.Flags.LeftWall || ch.Flags.RightWall) {
//		ch.Flags.JumpL = false
//		ch.Flags.JumpR = false
//		ch.Flags.HighJump = false
//		ch.Flags.LongJump = false
//		ch.Object.Pos.X = tile.Object.Pos.X
//		ch.Object.Pos.Y = tile.Object.Pos.Y
//		if ch.Flags.LadderHere {
//			ch.Flags.OnLadder = true
//			ch.ACounter = 0
//		}
//	} else {
//		if ch.Flags.HighJump {
//			ch.Object.Pos.Y += ch.Vars.HiJumpVSpeed
//			if tile.Coords != ch.LastTile.Coords {
//				if ch.Flags.JumpR {
//					ch.Object.Flip = false
//					if !ch.Flags.RightWall {
//						ch.Object.Pos.X += ch.Vars.HiJumpHSpeed
//					}
//				} else if ch.Flags.JumpL {
//					ch.Object.Flip = true
//					if !ch.Flags.LeftWall {
//						ch.Object.Pos.X -= ch.Vars.HiJumpHSpeed
//					}
//				}
//			} else {
//				// You can change the direction if you want before reaching the higher tile
//				if ch.Actions.Left() && !ch.Actions.Right() {
//					if ch.Flags.JumpR {
//						ch.Flags.JumpR = false
//					} else {
//						ch.Flags.JumpL = true
//						ch.Object.Flip = true
//					}
//				} else if ch.Actions.Right() && !ch.Actions.Left() {
//					if ch.Flags.JumpL {
//						ch.Flags.JumpL = false
//					} else {
//						ch.Flags.JumpR = true
//						ch.Object.Flip = false
//					}
//				}
//			}
//		} else if ch.Flags.LongJump {
//			if ch.Flags.JumpR {
//				ch.Object.Flip = false
//				if !ch.Flags.RightWall {
//					ch.Object.Pos.X += ch.Vars.LgJumpHSpeed
//				}
//			} else if ch.Flags.JumpL {
//				ch.Object.Flip = true
//				if !ch.Flags.LeftWall {
//					ch.Object.Pos.X -= ch.Vars.LgJumpHSpeed
//				}
//			}
//			if tile.Coords != ch.LastTile.Coords {
//				if util.Abs(tile.Coords.X-ch.LastTile.Coords.X) > 1 {
//					ch.Object.Pos.Y -= ch.Vars.LgJumpVSpeed
//				}
//			} else {
//				ch.Object.Pos.Y += ch.Vars.LgJumpVSpeed
//			}
//		}
//	}
//}
