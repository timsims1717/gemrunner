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
			ch.Flags.Action = false
			//ch.Flags.Drop = false
			actions := ct.GetActions()
			if !data.CurrLevel.Start {
				if result.Entity.HasComponent(myecs.Player) {
					if actions.Direction != data.None ||
						actions.Jump ||
						actions.PickUp ||
						actions.Action ||
						!ch.Flags.Floor {
						data.CurrLevel.Start = true
					}
				}
			}
			if data.CurrLevel.Start {
				if ch.Flags.Frame {
					ch.Flags.Frame = false
					ch.Actions = data.NewAction()
					if ch.Flags.JumpBuff > 0 {
						ch.Actions.Jump = true
					}
					if ch.Flags.PickUpBuff > 0 {
						ch.Actions.PickUp = true
					}
					if ch.Flags.ActionBuff > 0 {
						ch.Actions.Action = true
					}
					ch.Flags.JumpBuff--
					ch.Flags.PickUpBuff--
					ch.Flags.ActionBuff--
				}
				if actions.Direction != data.None {
					ch.Actions.Direction = actions.Direction
				}
				if actions.PrevDirection != data.None {
					ch.Actions.PrevDirection = actions.PrevDirection
				}
				ch.Actions.Jump = ch.Actions.Jump || actions.Jump
				ch.Actions.PickUp = ch.Actions.PickUp || actions.PickUp
				if actions.Jump {
					ch.Flags.JumpBuff = constants.ButtonBuffer
				}
				if actions.PickUp {
					ch.Flags.PickUpBuff = constants.ButtonBuffer
				}
				if actions.Action {
					ch.Flags.ActionBuff = constants.ButtonBuffer
				}
				debug.AddText(fmt.Sprintf("Direction: %5s", ch.Actions.Direction))
				debug.AddText(fmt.Sprintf("Previous:  %5s", ch.Actions.PrevDirection))
				debug.AddTruthText("Jump:      ", ch.Actions.Jump)
				debug.AddTruthText("PickUp:    ", ch.Actions.PickUp)
				debug.AddTruthText("Action:    ", ch.Actions.Action)

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
	if tile != nil &&
		ch.Actions.Jump &&
		!ch.Flags.Ceiling { // jump time
		jump(ch, tile)
		return
	}
	if ch.Actions.Left() && !ch.Flags.LeftWall { // run left
		ch.Object.Pos.X -= ch.Vars.WalkSpeed
		ch.Object.Flip = true
	} else if ch.Actions.Right() && !ch.Flags.RightWall { // run right
		ch.Object.Pos.X += ch.Vars.WalkSpeed
		ch.Object.Flip = false
	}
}

func jump(ch *data.Dynamic, tile *data.Tile) {
	ch.Flags.JumpBuff = 0
	upLeft := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y+1)
	upRight := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y+1)
	left := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
	right := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
	// High Jump if:
	//  the character is not going left or right
	//  or they are going left/right and there is a wall left/right
	//  or they are going left/right and there is a wall up left or up right
	// Otherwise, it's a long jump
	if (!ch.Actions.Left() && !ch.Actions.Right()) ||
		(ch.Actions.Left() && (left == nil || left.IsSolid())) ||
		(ch.Actions.Right() && (right == nil || right.IsSolid())) ||
		(ch.Actions.Left() && (upLeft == nil || upLeft.IsSolid())) ||
		(ch.Actions.Right() && (upRight == nil || upRight.IsSolid())) {
		ch.Flags.HighJump = true
	} else {
		ch.Flags.LongJump = true
	}
	// for both kinds of jumps
	ch.ACounter = 0
	if ch.Actions.Left() {
		ch.Flags.JumpL = true
		ch.Object.Flip = true
	} else if ch.Actions.Right() {
		ch.Flags.JumpR = true
		ch.Object.Flip = false
	} else {
		ch.Flags.JumpL = false
		ch.Flags.JumpR = false
	}
}

func onLadder(ch *data.Dynamic, tile, below *data.Tile) {
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
	} else {
		ch.Object.Pos.X += ch.Vars.LeapSpeed
	}
}

func flying(ch *data.Dynamic, tile *data.Tile) {
	if ch.Actions.Direction == data.Left { // fly left
		ch.Object.Pos.X -= ch.Vars.WalkSpeed
		ch.Object.Flip = true
	} else if ch.Actions.Direction == data.Right { // fly right
		ch.Object.Pos.X += ch.Vars.WalkSpeed
		ch.Object.Flip = false
	} else if ch.Actions.Direction == data.Up { // fly up
		ch.Object.Pos.Y += ch.Vars.WalkSpeed
	} else if ch.Actions.Direction == data.Down { // fly down
		ch.Object.Pos.Y -= ch.Vars.WalkSpeed
	}
}
