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
				currPos := ch.Object.Pos
				x, y := world.WorldToMap(currPos.X, currPos.Y)
				tile := data.CurrLevel.Tiles.Get(x, y)
				below := data.CurrLevel.Tiles.Get(x, y-1)
				oldState := ch.State
				switch ch.State {
				case data.Grounded:
					if ch.Flags.HighJump || ch.Flags.LongJump {
						ch.State = data.Jumping
					} else if tile != nil && tile.Ladder { // a ladder is here
						if ch.Actions.Direction == data.Up { // climbed the ladder
							ch.State = data.Ladder
						} else if below != nil && below.Ladder && ch.Actions.Direction == data.Down { // down the ladder
							ch.State = data.Ladder
						} else if below != nil && below.Block != data.BlockTurf { // leaping onto ladder
							ch.State = data.Leaping
							ch.Flags.LeapOn = true
						}
					} else if below != nil && below.Ladder && ch.Actions.Direction == data.Down { // down the ladder
						ch.State = data.Ladder
					} else if !ch.Flags.Floor && below != nil && !below.Ladder {
						ch.State = data.Falling
					}
				case data.Ladder:
					if ch.Flags.Floor { // just got to the bottom or top
						ch.State = data.Grounded
						if ch.Actions.Direction == data.Left { // to the left
							ch.Object.Flip = true
						} else if ch.Actions.Direction == data.Right { // to the right
							ch.Object.Flip = false
						}
					} else if !tile.Ladder && (below == nil || !below.Ladder) {
						ch.State = data.Falling
					} else if tile.Ladder &&
						!(below == nil ||
							below.Solid() ||
							below.Block == data.BlockTurf) { // can only leap if the stuff below isn't solid
						if ch.Actions.Direction == data.Left &&
							!ch.Flags.LeftWall { // leaping to the left
							ch.State = data.Leaping
							ch.Object.Flip = true
							left := data.CurrLevel.Tiles.Get(x-1, y)
							if left != nil && left.Ladder { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						} else if ch.Actions.Direction == data.Right &&
							!ch.Flags.RightWall { // leaping to the right
							ch.State = data.Leaping
							ch.Object.Flip = false
							right := data.CurrLevel.Tiles.Get(x+1, y)
							if right != nil && right.Ladder { // to another ladder
								ch.Flags.LeapTo = true
							} else { // off the ladders
								ch.Flags.LeapOff = true
							}
						}
					} else if below == nil ||
						below.Solid() ||
						below.Block == data.BlockTurf ||
						(below.Ladder && !tile.Ladder) {
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
					} else if tile != nil && tile.Ladder {
						ch.State = data.Ladder
					}
				case data.Jumping:
					if !(ch.Flags.HighJump || ch.Flags.LongJump) {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.Ladder {
							ch.State = data.Ladder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Leaping:
					if !(ch.Flags.LeapOff || ch.Flags.LeapOn || ch.Flags.LeapTo) {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.Ladder {
							ch.State = data.Ladder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Flying:
					if !ch.Flags.Flying {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.Ladder {
							ch.State = data.Ladder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Hit:
					if !ch.Flags.Hit {
						ch.State = data.Dead
					}
				case data.Attack:
					if !ch.Flags.Attack {
						if ch.Flags.Floor {
							ch.State = data.Grounded
						} else if tile != nil && tile.Ladder {
							ch.State = data.Ladder
						} else {
							ch.State = data.Falling
						}
					}
				case data.Dead:
				}
				if oldState != ch.State { // a state change happened
					ch.ACounter = 0
					ch.ATimer = nil
					ch.Control.ClearPrev()
					ch.Actions.PrevDirection = data.None
					ch.Object.Pos.Y = tile.Object.Pos.Y
					if oldState != data.Grounded || ch.State != data.Leaping {
						ch.Object.Pos.X = tile.Object.Pos.X
					}
					ch.Flags.Climbed = false
					ch.Flags.GoingUp = false
				}
				updateHeldItem(ch, ch.Object.Flip)
				if ch.State == data.Dead ||
					ch.State == data.Hit ||
					ch.State == data.Attack ||
					ch.State == data.Ladder ||
					ch.State == data.Leaping ||
					ch.Flags.Drop ||
					ch.Flags.Action {
					if ch.Flags.Action {
						DoAction(ch)
					}
					DropItem(ch)
				}
			}
		}
	}
}
