package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
)

func CharacterActionSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		ct, okT := result.Components[myecs.Controller].(data.Controller)
		if okO && okC && okT && !obj.Hidden && ct != nil {
			actions := ct.GetActions()
			if !data.CurrLevel.Start &&
				result.Entity.HasComponent(myecs.Player) &&
				ch.State != data.Regen &&
				(actions.Any() || !ch.Flags.Floor) {
				data.CurrLevel.Start = true
			}
			if data.CurrLevel.Start {
				if ch.Flags.Frame {
					ch.Flags.Frame = false
					ch.Actions = data.NewAction()
					if ch.Flags.PickUpBuff > 0 {
						ch.Actions.PickUp = true
						ch.Flags.PickUpBuff--
					}
					if ch.Flags.ActionBuff > 0 {
						ch.Actions.Action = true
						ch.Flags.ActionBuff--
					}
					if ch.Flags.DigRightBuff > 0 {
						ch.Actions.DigRight = true
						ch.Flags.DigRightBuff--
					}
					if ch.Flags.DigLeftBuff > 0 {
						ch.Actions.DigLeft = true
						ch.Flags.DigLeftBuff--
					}
				}
				if actions.Direction != data.NoDirection {
					ch.Actions.Direction = actions.Direction
				}
				if actions.PrevDirection != data.NoDirection {
					ch.Actions.PrevDirection = actions.PrevDirection
				}
				ch.Actions.PickUp = ch.Actions.PickUp || actions.PickUp
				ch.Actions.Action = ch.Actions.Action || actions.Action
				ch.Actions.DigLeft = ch.Actions.DigLeft || actions.DigLeft
				ch.Actions.DigRight = ch.Actions.DigRight || actions.DigRight
				if actions.PickUp {
					ch.Flags.PickUpBuff = constants.ButtonBuffer
				}
				if actions.Action && ch.State != data.DoingAction {
					ch.Flags.ActionBuff = constants.ButtonBuffer
				}
				if actions.DigLeft {
					ch.Flags.DigLeftBuff = constants.ButtonBuffer
				}
				if actions.DigRight {
					ch.Flags.DigRightBuff = constants.ButtonBuffer
				}
				if ch.Player > -1 {
					debug.AddText(fmt.Sprintf("Direction: %5s", ch.Actions.Direction))
					debug.AddText(fmt.Sprintf("Previous:  %5s", ch.Actions.PrevDirection))
					p, a, l, r := "-", "-", "-", "-"
					if ch.Actions.PickUp {
						p = "P"
					}
					if ch.Actions.Action {
						a = "A"
					}
					if ch.Actions.DigLeft {
						l = "<"
					}
					if ch.Actions.DigRight {
						r = ">"
					}
					debug.AddText(fmt.Sprintf("%s|%s|%s|%s", p, a, l, r))
				}

				if reanimator.FrameSwitch {
					currPos := ch.Object.Pos.Add(ch.Object.Offset)
					x, y := world.WorldToMap(currPos.X, currPos.Y)
					tile := data.CurrLevel.Tiles.Get(x, y)
					below := data.CurrLevel.Tiles.Get(x, y-1)
					ch.ACounter++
					switch ch.State {
					case data.Grounded:
						grounded(ch, tile, below)
					case data.OnLadder:
						onLadder(ch, tile, below)
					case data.OnBar:
						onBar(ch, tile, below)
					case data.Falling:
						falling(ch, tile)
					case data.Jumping:
						jumping(ch, tile)
					case data.Leaping:
						leaping(ch, tile)
					case data.Flying:
						flying(ch, tile)
					case data.DoingAction:
					}
					ch.Flags.Frame = true
				}
			}
		}
	}
}

func grounded(ch *data.Dynamic, tile, below *data.Tile) {
	ch.LastTile = tile
	//if tile != nil &&
	//	ch.Actions.Jump { // jump time
	//	if jump(ch, tile) {
	//		return
	//	}
	//}
	if ch.Actions.Left() && !ch.Flags.LeftWall { // run left
		ch.Object.Pos.X -= ch.Vars.WalkSpeed
		ch.Object.Flip = true
	} else if ch.Actions.Right() && !ch.Flags.RightWall { // run right
		ch.Object.Pos.X += ch.Vars.WalkSpeed
		ch.Object.Flip = false
	}
}

//func jump(ch *data.Dynamic, tile *data.Tile) bool {
//	left := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
//	right := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
//	//if ch.Flags.Ceiling &&
//	//	((ch.Actions.Left() && left.IsSolid()) ||
//	//		(ch.Actions.Right() && right.IsSolid())) {
//	//	return false
//	//}
//	ch.Flags.JumpBuff = 0
//	// High Jump if:
//	//  there is no ceiling here
//	//  the character is not going left or right
//	//  or they are going left/right and there is a wall left/right
//	//  or they are going left/right and there is a wall up left or up right
//	// Otherwise, it's a long jump
//	if !ch.Flags.Ceiling &&
//		((!ch.Actions.Left() && !ch.Actions.Right()) ||
//			(ch.Actions.Left() && left.IsSolid()) ||
//			(ch.Actions.Right() && right.IsSolid())) {
//		ch.Flags.HighJump = true
//	} else if (ch.Actions.Left() && !left.IsSolid()) ||
//		(ch.Actions.Right() && !right.IsSolid()) {
//		ch.Flags.LongJump = true
//	} else {
//		return false
//	}
//	// for both kinds of jumps
//	ch.ACounter = 0
//	if ch.Actions.Left() {
//		ch.Flags.JumpL = true
//		ch.Object.Flip = true
//	} else if ch.Actions.Right() {
//		ch.Flags.JumpR = true
//		ch.Object.Flip = false
//	} else {
//		ch.Flags.JumpL = false
//		ch.Flags.JumpR = false
//	}
//	return true
//}

func onLadder(ch *data.Dynamic, tile, below *data.Tile) {
	ch.LastTile = tile
	if ch.Actions.Up() && !ch.Flags.Ceiling {
		if tile.IsLadder() || (below != nil && below.IsLadder()) { // still on the ladder
			ch.Object.Pos.Y += ch.Vars.ClimbSpeed
			ch.Object.Pos.X = tile.Object.Pos.X
			ch.Object.Flip = false
			ch.Flags.GoingUp = true
			ch.Flags.Climbed = true
		}
	} else if ch.Actions.Down() {
		if (tile != nil && tile.IsLadder()) ||
			(below != nil && below.IsLadder()) { // still on the ladder
			ch.Object.Pos.Y -= ch.Vars.SlideSpeed
			ch.Object.Pos.X = tile.Object.Pos.X
			ch.Object.Flip = false
			ch.Flags.GoingUp = false
			ch.Flags.Climbed = true
		}
	}
}

func onBar(ch *data.Dynamic, tile, below *data.Tile) {
	ch.LastTile = tile
	if ch.Actions.Down() && (!ch.Flags.Floor || below.IsLadder()) { // drop down
		ch.Object.Pos.Y -= ch.Vars.Gravity
	} else if ch.Actions.Left() && !ch.Flags.LeftWall { // go left
		ch.Object.Pos.X -= ch.Vars.BarSpeed
		ch.Object.Flip = true
		ch.Flags.Climbed = true
	} else if ch.Actions.Right() && !ch.Flags.RightWall { // go right
		ch.Object.Pos.X += ch.Vars.BarSpeed
		ch.Object.Flip = false
		ch.Flags.Climbed = true
	}
}

func falling(ch *data.Dynamic, tile *data.Tile) {
	if tile == nil {
		return
	}
	ch.Object.Pos.X = tile.Object.Pos.X
	ch.Object.Pos.Y -= ch.Vars.Gravity
	if ch.Actions.Left() {
		//ch.Object.Pos.X -= ch.Vars.WalkSpeed
		ch.Object.Flip = true
	} else if ch.Actions.Right() {
		//ch.Object.Pos.X += ch.Vars.WalkSpeed
		ch.Object.Flip = false
	}
}

func jumping(ch *data.Dynamic, tile *data.Tile) {
	if (ch.Flags.HighJump && ch.ACounter > int(ch.Vars.HiJumpCntr)) ||
		(ch.Flags.LongJump && ch.ACounter > int(ch.Vars.LgJumpCntr)) ||
		(ch.Flags.LeftWall || ch.Flags.RightWall) {
		ch.Flags.JumpL = false
		ch.Flags.JumpR = false
		ch.Flags.HighJump = false
		ch.Flags.LongJump = false
	} else {
		if ch.Flags.HighJump {
			ch.Object.Pos.Y += ch.Vars.HiJumpVSpeed
			if tile.Coords != ch.LastTile.Coords {
				if ch.Flags.JumpR {
					ch.Object.Flip = false
					if !ch.Flags.RightWall {
						ch.Object.Pos.X += ch.Vars.HiJumpHSpeed
					}
				} else if ch.Flags.JumpL {
					ch.Object.Flip = true
					if !ch.Flags.LeftWall {
						ch.Object.Pos.X -= ch.Vars.HiJumpHSpeed
					}
				}
			} else {
				// You can change the direction if you want before reaching the higher tile
				if ch.Actions.Left() && !ch.Actions.Right() {
					if ch.Flags.JumpR {
						ch.Flags.JumpR = false
					} else {
						ch.Flags.JumpL = true
						ch.Object.Flip = true
					}
				} else if ch.Actions.Right() && !ch.Actions.Left() {
					if ch.Flags.JumpL {
						ch.Flags.JumpL = false
					} else {
						ch.Flags.JumpR = true
						ch.Object.Flip = false
					}
				}
			}
		} else if ch.Flags.LongJump {
			if ch.Flags.JumpR {
				ch.Object.Flip = false
				if !ch.Flags.RightWall {
					ch.Object.Pos.X += ch.Vars.LgJumpHSpeed
				}
			} else if ch.Flags.JumpL {
				ch.Object.Flip = true
				if !ch.Flags.LeftWall {
					ch.Object.Pos.X -= ch.Vars.LgJumpHSpeed
				}
			}
			if tile.Coords != ch.LastTile.Coords {
				if util.Abs(tile.Coords.X-ch.LastTile.Coords.X) > 1 {
					ch.Object.Pos.Y -= ch.Vars.LgJumpVSpeed
				}
			} else {
				ch.Object.Pos.Y += ch.Vars.LgJumpVSpeed
			}
		}
	}
}

func leaping(ch *data.Dynamic, tile *data.Tile) {
	ch.Object.Pos.Y = tile.Object.Pos.Y
	if ch.Object.Flip {
		// going left
		ch.Object.Pos.X -= ch.Vars.LeapSpeed
		if ch.Object.Pos.X < ch.NextTile.Object.Pos.X { // the leap is complete
			ch.Object.Pos.X = ch.NextTile.Object.Pos.X
			ch.Flags.LeapOn = false
			ch.Flags.LeapTo = false
			ch.Flags.LeapOff = false
		}
	} else {
		ch.Object.Pos.X += ch.Vars.LeapSpeed
		if ch.Object.Pos.X > ch.NextTile.Object.Pos.X { // the leap is complete
			ch.Object.Pos.X = ch.NextTile.Object.Pos.X
			ch.Flags.LeapOn = false
			ch.Flags.LeapTo = false
			ch.Flags.LeapOff = false
		}
	}
}

func flying(ch *data.Dynamic, tile *data.Tile) {
	if ch.Actions.Direction == data.Left { // fly left
		ch.Object.Pos.X -= ch.Vars.WalkSpeed
		//ch.Object.Pos.Y = tile.Object.Pos.Y
		ch.Object.Flip = true
	} else if ch.Actions.Direction == data.Right { // fly right
		ch.Object.Pos.X += ch.Vars.WalkSpeed
		//ch.Object.Pos.Y = tile.Object.Pos.Y
		ch.Object.Flip = false
	} else if ch.Actions.Direction == data.Up { // fly up
		ch.Object.Pos.Y += ch.Vars.WalkSpeed
		//ch.Object.Pos.X = tile.Object.Pos.X
	} else if ch.Actions.Direction == data.Down { // fly down
		ch.Object.Pos.Y -= ch.Vars.WalkSpeed
		//ch.Object.Pos.X = tile.Object.Pos.X
	}
}
