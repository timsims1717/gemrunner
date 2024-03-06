package systems

import (
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/world"
)

func CharacterStateSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		_, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		if okO && okC {
			if reanimator.FrameSwitch {
				currPos := ch.Object.Pos.Add(ch.Object.Offset)
				x, y := world.WorldToMap(currPos.X, currPos.Y)
				tile := data.CurrLevel.Tiles.Get(x, y)
				below := data.CurrLevel.Tiles.Get(x, y-1)
				oldState := ch.State
				switch ch.State {
				case data.Grounded:
					if ch.Flags.HighJump || ch.Flags.LongJump {
						ch.State = data.Jumping
					} else if tile != nil && tile.IsLadder() { // a ladder is here
						if ch.Actions.Direction == data.Up { // climbed the ladder
							ch.State = data.OnLadder
						} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
							ch.State = data.OnLadder
						} else if below != nil && below.Block != data.BlockTurf { // leaping onto ladder
							ch.State = data.Leaping
							ch.Flags.LeapOn = true
						}
					} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
						ch.State = data.OnLadder
					} else if !ch.Flags.Floor && below != nil && !below.IsLadder() {
						ch.State = data.Falling
					}
				case data.OnLadder:
					DropLift(ch, false)
					if ch.Flags.Floor { // just got to the bottom or top
						ch.State = data.Grounded
						if ch.Actions.Direction == data.Left { // to the left
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right { // to the right
							ch.Object.Flip = false
						}
					} else if !tile.IsLadder() && (below == nil || !below.IsLadder()) {
						ch.State = data.Falling
					} else if tile.IsLadder() &&
						!(below == nil ||
							below.SolidV() ||
							below.Block == data.BlockTurf) { // can only leap if the stuff below isn't solid
						if ch.Actions.Direction == data.Left &&
							!ch.Flags.LeftWall { // leaping to the left
							ch.State = data.Leaping
							ch.Object.Flip = true
							left := data.CurrLevel.Tiles.Get(x-1, y)
							if left != nil && left.IsLadder() { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						} else if ch.Actions.Direction == data.Right &&
							!ch.Flags.RightWall { // leaping to the right
							ch.State = data.Leaping
							ch.Object.Flip = false
							right := data.CurrLevel.Tiles.Get(x+1, y)
							if right != nil && right.IsLadder() { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						}
					} else if below == nil ||
						below.SolidV() ||
						below.Block == data.BlockTurf ||
						(below.IsLadder() && !tile.IsLadder()) {
						if ch.Actions.Direction == data.Left { // run to the left
							ch.State = data.Grounded
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right { // run to the right
							ch.State = data.Grounded
							ch.Object.Flip = false
						}
					}
				case data.Falling:
					if ch.Flags.Floor {
						ch.State = data.Grounded
					} else if tile != nil && tile.IsLadder() {
						ch.State = data.OnLadder
						ch.Object.Pos.X = tile.Object.Pos.X
					} else {
						ch.Object.Pos.X = tile.Object.Pos.X
					}
				case data.Jumping:
					if !(ch.Flags.HighJump || ch.Flags.LongJump) {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Leaping:
					if !(ch.Flags.LeapOff || ch.Flags.LeapOn || ch.Flags.LeapTo) {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Flying:
					if !ch.Flags.Flying {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Carried:
					if ch.Flags.HighJump || ch.Flags.LongJump {
						ch.State = data.Jumping
					} else if tile != nil && tile.IsLadder() { // a ladder is here
						if ch.Actions.Direction == data.Up { // climbed the ladder
							ch.State = data.OnLadder
						} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
							ch.State = data.OnLadder
						}
					} else if below != nil && below.IsLadder() && ch.Actions.Direction == data.Down { // down the ladder
						ch.State = data.OnLadder
					}
				case data.Hit:
					if !ch.Flags.Hit {
						ch.State = data.Dead
					}
				case data.Attack:
					if !ch.Flags.Attack {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.IsLadder() {
							ch.State = data.OnLadder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Dead:
				}
				if oldState != ch.State { // a state change happened
					if oldState == data.Carried {
						if pu, ok := ch.Entity.GetComponentData(myecs.PickUp); ok {
							pickUp := pu.(*data.PickUp)
							newState := ch.State
							DropLift(data.CurrLevel.Players[pickUp.Held], false)
							ch.State = newState
						}
					}
					ch.ACounter = 0
					ch.ATimer = nil
					ch.Control.ClearPrev()
					ch.Actions.PrevDirection = data.None
					ch.Object.Pos.Y = tile.Object.Pos.Y
					if !((oldState == data.Falling &&
						ch.State == data.Grounded) ||
						(oldState == data.Grounded &&
							ch.State == data.Leaping)) {
						ch.Object.Pos.X = tile.Object.Pos.X
					}
					ch.Flags.Climbed = false
					ch.Flags.GoingUp = false
				}
				updateHeldItem(ch, ch.Object.Flip)
				if ch.State != data.Dead &&
					ch.State != data.Hit &&
					ch.State != data.Attack {
					if ch.State != data.Leaping &&
						ch.Actions.Action {
						DoAction(ch)
						ch.Flags.ActionBuff = 0
					} else if ch.Actions.PickUp && !ch.Flags.Using {
						PickUpOrDropItem(ch, int(ch.Player))
						ch.Flags.PickUpBuff = 0
					} else if ch.State != data.Leaping &&
						ch.State != data.OnLadder &&
						ch.Actions.Lift && !ch.Flags.Using {
						LiftOrDropItem(ch, int(ch.Player), ch.Actions.Left() || ch.Actions.Right())
						ch.Flags.LiftBuff = 0
					}
				} else {
					if ch.State == data.Dead {
						DropItem(ch)
					}
					DropLift(ch, false)
				}
			}
		}
	}
}
