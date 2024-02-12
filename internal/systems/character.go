package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/systems/characters"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/bytearena/ecs"
	"github.com/gopxl/pixel"
)

func CharacterSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Dynamic].(*data.Dynamic)
		ct, okT := result.Components[myecs.Controller].(data.Controller)
		if okO && okC && okT && !obj.Hidden {
			ch.Actions = ct.GetActions()
			currPos := ch.FauxObj.Pos
			x, y := world.WorldToMap(currPos.X, currPos.Y)
			currTile := data.CurrLevel.Tiles.Get(x, y)
			debug.AddTruthText("Left:   ", ch.Actions.Left)
			debug.AddTruthText("Right:  ", ch.Actions.Right)
			debug.AddTruthText("Up:     ", ch.Actions.Up)
			debug.AddTruthText("Down:   ", ch.Actions.Down)
			debug.AddTruthText("Jump:   ", ch.Actions.Jump)
			debug.AddTruthText("PickUp: ", ch.Actions.PickUp)
			debug.AddTruthText("Action: ", ch.Actions.Action)
			ch.ATimer.Update()
			if !(ch.Flags.Hit || ch.Flags.Dead || ch.Flags.Attack) {
				if ch.Flags.LeapOn || ch.Flags.LeapOff || ch.Flags.LeapTo {
					leaping(ch, currTile)
				} else if !ch.Flags.Floor && !ch.Flags.HighJump && !ch.Flags.LongJump &&
					!ch.Flags.LadderHere && !ch.Flags.OnLadder {
					falling(ch, currTile)
				} else if ch.Flags.HighJump || ch.Flags.LongJump {
					jumping(ch, currTile)
				} else {
					upOrDown(ch, currTile)
					if ch.Flags.OnLadder {
						ch.FauxObj.Pos.X = currTile.Object.Pos.X
						gettingOffLadder(ch, currTile)
					} else {
						onTheGround(ch, currTile)
					}
				}
				if ch.Actions.PickUp {
					if ch.Flags.HoldUp || ch.Flags.HoldSide {
						ch.Flags.Drop = true
					} else if !ch.Flags.LeapOn && !ch.Flags.LeapOff &&
						!ch.Flags.LeapTo && !ch.Flags.OnLadder &&
						ch.Player > -1 && ch.Player < constants.MaxPlayers {
						AttemptPickUp(ch, int(ch.Player), ch.FauxObj.Flip)
					}
				} else if ch.Actions.Action {
					if ch.Flags.HoldUp || ch.Flags.HoldSide {
						//ch.Flags.Drop = true
						ch.Flags.Action = true
					}
				}
				if !ch.Flags.OnLadder {
					ch.Flags.GoingUp = false
					ch.Flags.Climbed = false
				}
				if ch.Actions.Down {
					ch.Flags.DropDown = true
				}
			}
			if reanimator.FrameSwitch {
				obj.Pos.X = ch.FauxObj.Pos.X
				obj.Pos.Y = ch.FauxObj.Pos.Y
				obj.Flip = ch.FauxObj.Flip
				updateHeldItem(ch, obj.Flip)
				if ch.Flags.Dead || ch.Flags.Hit || ch.Flags.Attack ||
					ch.Flags.OnLadder || ch.Flags.Drop || ch.Flags.Action {
					if ch.Flags.Action && !(ch.Flags.Dead || ch.Flags.Hit || ch.Flags.Attack) {
						DoAction(ch)
						ch.Flags.Action = false
					}
					DropItem(ch)
				}
			}
		}
	}
}

func leaping(ch *data.Dynamic, tile *data.Tile) {
	ch.FauxObj.Pos.Y = tile.Object.Pos.Y
	if ch.FauxObj.Flip {
		// going left
		ch.FauxObj.Pos.X -= float64(reanimator.FRate) * ch.Vars.LeapSpeed * timing.DT
	} else {
		ch.FauxObj.Pos.X += float64(reanimator.FRate) * ch.Vars.LeapSpeed * timing.DT
	}
}

func falling(ch *data.Dynamic, tile *data.Tile) {
	ch.FauxObj.Pos.X = tile.Object.Pos.X
	ch.FauxObj.Pos.Y -= float64(reanimator.FRate) * ch.Vars.Gravity * timing.DT
	ch.Flags.OnLadder = false
}

func upOrDown(ch *data.Dynamic, tile *data.Tile) {
	if !ch.Flags.HoldUp && !ch.Flags.HoldSide {
		if ch.Actions.Up && !ch.Flags.Ceiling {
			if ch.Flags.LadderHere {
				ch.FauxObj.Pos.Y += float64(reanimator.FRate) * ch.Vars.ClimbSpeed * timing.DT
				ch.FauxObj.Flip = false
				ch.Flags.OnLadder = true
				ch.Flags.GoingUp = true
				ch.Flags.Climbed = true
			} else if ch.Flags.OnLadder {
				ch.FauxObj.Pos.Y = tile.Object.Pos.Y
				ch.FauxObj.Flip = false
				ch.Flags.OnLadder = false
			}
		} else if ch.Actions.Down {
			if (ch.Flags.LadderHere && !ch.Flags.Floor) || ch.Flags.LadderDown {
				ch.FauxObj.Pos.Y -= float64(reanimator.FRate) * ch.Vars.SlideSpeed * timing.DT
				ch.FauxObj.Flip = false
				ch.Flags.OnLadder = true
				ch.Flags.GoingUp = false
				ch.Flags.Climbed = true
			}
		}
	} else {
		if ch.Flags.HoldUp && ch.Actions.Down && (ch.Actions.Left || ch.Actions.Right) {
			ch.Flags.HoldUp = false
			ch.Flags.HoldSide = true
		} else if ch.Flags.HoldSide && ch.Actions.Up {
			ch.Flags.HoldUp = true
			ch.Flags.HoldSide = false
		}
	}
}

func gettingOffLadder(ch *data.Dynamic, tile *data.Tile) {
	if ch.Actions.Left && !ch.Flags.LeftWall &&
		(ch.Flags.CanRun || ch.ATimer.Done()) {
		ch.Flags.OnLadder = false
		ch.FauxObj.Flip = true
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X -= float64(reanimator.FRate) * ch.Vars.WalkSpeed * timing.DT
		} else if ch.ATimer.Done() {
			lWall := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
			if lWall != nil && lWall.Ladder {
				ch.Flags.LeapTo = true
			} else {
				ch.Flags.LeapOff = true
			}
		}
	} else if ch.Actions.Right && !ch.Flags.RightWall &&
		(ch.Flags.CanRun || ch.ATimer.Done()) {
		ch.Flags.OnLadder = false
		ch.FauxObj.Flip = false
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X += float64(reanimator.FRate) * ch.Vars.WalkSpeed * timing.DT
		} else if ch.ATimer.Done() {
			rWall := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
			if rWall != nil && rWall.Ladder {
				ch.Flags.LeapTo = true
			} else {
				ch.Flags.LeapOff = true
			}
		}
	}
}

func onTheGround(ch *data.Dynamic, tile *data.Tile) {
	ch.LastTile = tile
	if ch.Actions.Jump && ch.Flags.CanRun && !ch.Flags.Ceiling {
		upLeft := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y+1)
		upRight := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y+1)
		left := data.CurrLevel.Tiles.Get(tile.Coords.X-1, tile.Coords.Y)
		right := data.CurrLevel.Tiles.Get(tile.Coords.X+1, tile.Coords.Y)
		// High Jump if:
		//  the character is not going left or right
		//  or they are going left/right and there is a wall left/right
		//  or they are going left/right and there is a wall up left or up right
		// Otherwise, it's a long jump
		if (!ch.Actions.Left && !ch.Actions.Right) ||
			(ch.Actions.Left && (left == nil || left.Solid())) ||
			(ch.Actions.Right && (right == nil || right.Solid())) ||
			(ch.Actions.Left && (upLeft == nil || upLeft.Solid())) ||
			(ch.Actions.Right && (upRight == nil || upRight.Solid())) {
			ch.Flags.HighJump = true
			ch.FauxObj.Pos.X = tile.Object.Pos.X
			ch.FauxObj.Pos.Y = tile.Object.Pos.Y + float64(reanimator.FRate)*ch.Vars.HiJumpVSpeed*timing.DT
			ch.ATimer = timing.New(ch.Vars.HiJumpTimer / float64(reanimator.FRate))
		} else {
			ch.Flags.LongJump = true
			if ch.Actions.Left {
				ch.FauxObj.Pos.X = tile.Object.Pos.X - float64(reanimator.FRate)*ch.Vars.LgJumpHSpeed*timing.DT
			} else {
				ch.FauxObj.Pos.X = tile.Object.Pos.X + float64(reanimator.FRate)*ch.Vars.LgJumpHSpeed*timing.DT
			}
			//ch.FauxObj.Pos.Y = tile.Object.Pos.Y + float64(reanimator.FRate)*constants.PlayerLongJumpSpeed*timing.DT
			ch.ATimer = timing.New(ch.Vars.LgJumpTimer / float64(reanimator.FRate))
		}
		// for both kinds of jumps
		if ch.Actions.Left && !ch.Actions.Right {
			ch.Flags.JumpL = true
			ch.FauxObj.Flip = true
		} else if ch.Actions.Right && !ch.Actions.Left {
			ch.Flags.JumpR = true
			ch.FauxObj.Flip = false
		} else {
			ch.Flags.JumpL = false
			ch.Flags.JumpR = false
		}
	} else {
		if !ch.Flags.PickUp {
			if ch.Actions.Left && !ch.Flags.LeftWall {
				if ch.Flags.CanRun {
					ch.FauxObj.Pos.X -= float64(reanimator.FRate) * ch.Vars.WalkSpeed * timing.DT
					ch.FauxObj.Flip = true
					ch.Flags.OnLadder = false
				} else if tile.Ladder {
					ch.Flags.LeapOn = true
				}
			} else if ch.Actions.Right && !ch.Flags.RightWall {
				if ch.Flags.CanRun {
					ch.FauxObj.Pos.X += float64(reanimator.FRate) * ch.Vars.WalkSpeed * timing.DT
					ch.FauxObj.Flip = false
					ch.Flags.OnLadder = false
				} else if tile.Ladder {
					ch.Flags.LeapOn = true
				}
			}
		}
	}
}

func jumping(ch *data.Dynamic, tile *data.Tile) {
	if ch.ATimer.Done() {
		ch.Flags.JumpL = false
		ch.Flags.JumpR = false
		ch.Flags.HighJump = false
		ch.Flags.LongJump = false
		ch.FauxObj.Pos.X = tile.Object.Pos.X
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y
		if ch.Flags.LadderHere {
			ch.Flags.OnLadder = true
			ch.ATimer = timing.New(ch.Vars.LeapDelay / float64(reanimator.FRate))
		}
	} else {
		if ch.Flags.HighJump {
			ch.FauxObj.Pos.Y += float64(reanimator.FRate) * ch.Vars.HiJumpVSpeed * timing.DT
			if tile.Coords != ch.LastTile.Coords {
				if ch.Flags.JumpR {
					ch.FauxObj.Flip = false
					if !ch.Flags.RightWall {
						ch.FauxObj.Pos.X += float64(reanimator.FRate) * ch.Vars.HiJumpHSpeed * timing.DT
					}
				} else if ch.Flags.JumpL {
					ch.FauxObj.Flip = true
					if !ch.Flags.LeftWall {
						ch.FauxObj.Pos.X -= float64(reanimator.FRate) * ch.Vars.HiJumpHSpeed * timing.DT
					}
				}
			} else {
				// You can change the direction if you want before reaching the higher tile
				if ch.Actions.Left && !ch.Actions.Right {
					if ch.Flags.JumpR {
						ch.Flags.JumpR = false
					} else {
						ch.Flags.JumpL = true
						ch.FauxObj.Flip = true
					}
				} else if ch.Actions.Right && !ch.Actions.Left {
					if ch.Flags.JumpL {
						ch.Flags.JumpL = false
					} else {
						ch.Flags.JumpR = true
						ch.FauxObj.Flip = false
					}
				}
			}
		} else if ch.Flags.LongJump {
			if ch.Flags.JumpR {
				ch.FauxObj.Flip = false
				if !ch.Flags.RightWall {
					ch.FauxObj.Pos.X += float64(reanimator.FRate) * ch.Vars.LgJumpHSpeed * timing.DT
				}
			} else if ch.Flags.JumpL {
				ch.FauxObj.Flip = true
				if !ch.Flags.LeftWall {
					ch.FauxObj.Pos.X -= float64(reanimator.FRate) * ch.Vars.LgJumpHSpeed * timing.DT
				}
			}
			if tile.Coords != ch.LastTile.Coords {
				if util.Abs(tile.Coords.X-ch.LastTile.Coords.X) > 1 {
					ch.FauxObj.Pos.Y -= float64(reanimator.FRate) * ch.Vars.LgJumpVSpeed * timing.DT
				}
			} else {
				ch.FauxObj.Pos.Y += float64(reanimator.FRate) * ch.Vars.LgJumpVSpeed * timing.DT
			}
		}
	}
}

func PlayerCharacter(pos pixel.Vec, pIndex int) *data.Dynamic {
	obj := object.New().WithID(fmt.Sprintf("player_%d", pIndex))
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Pos = pos
	obj.Layer = 27 - pIndex*2
	fObj := object.New()
	fObj.Pos = pos
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	switch pIndex {
	case 0:
		e.AddComponent(myecs.Controller, controllers.NewPlayerInput(data.P1Input, "p1_left", "p1_right", "p1_up", "p1_down", "p1_jump", "p1_pickUp", "p1_action"))
	case 1:
		e.AddComponent(myecs.Controller, controllers.NewPlayerInput(data.P2Input, "p2_left", "p2_right", "p2_up", "p2_down", "p2_jump", "p2_pickUp", "p2_action"))
	}
	player := data.NewDynamic()
	player.Object = obj
	player.FauxObj = fObj
	player.Entity = e
	player.Anim = characters.HumanoidAnimation(player, "player")
	player.Player = data.Player(pIndex)
	player.Vars = data.PlayerVars()
	e.AddComponent(myecs.Animated, player.Anim)
	e.AddComponent(myecs.Drawable, player.Anim)
	e.AddComponent(myecs.Dynamic, player)
	e.AddComponent(myecs.Player, player.Player)
	return player
}

func DemonCharacter(pos pixel.Vec, pIndex int) *data.Dynamic {
	obj := object.New().WithID("demon")
	obj.SetRect(pixel.R(0, 0, 12, 16))
	obj.Pos = pos
	obj.Layer = 29
	fObj := object.New()
	fObj.Pos = pos
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	//switch pIndex {
	//case 0:
	//	e.AddComponent(myecs.Controller, controllers.NewPlayerInput(data.P1Input, "p1_left", "p1_right", "p1_up", "p1_down", "p1_jump", "p1_pickUp", "p1_action"))
	//case 1:
	//	e.AddComponent(myecs.Controller, controllers.NewPlayerInput(data.P2Input, "p2_left", "p2_right", "p2_up", "p2_down", "p2_jump", "p2_pickUp", "p2_action"))
	//}
	demon := data.NewDynamic()
	demon.Object = obj
	demon.FauxObj = fObj
	demon.Entity = e
	demon.Anim = characters.HumanoidAnimation(demon, "demon")
	demon.Vars = data.DemonVars()
	e.AddComponent(myecs.Animated, demon.Anim)
	e.AddComponent(myecs.Drawable, demon.Anim)
	e.AddComponent(myecs.Dynamic, demon)
	e.AddComponent(myecs.OnTouch, data.NewInteract(KillPlayer))
	e.AddComponent(myecs.Controller, controllers.NewLRChase(demon))
	return demon
}

func KillPlayer(level *data.Level, p int, ch *data.Dynamic, entity *ecs.Entity) {
	if p < 0 || p >= constants.MaxPlayers || ch.Flags.Hit || ch.Flags.Dead {
		return
	}
	bg, ok := entity.GetComponentData(myecs.Dynamic)
	if ok {
		enemy := bg.(*data.Dynamic)
		if !enemy.Flags.Attack && (enemy.Flags.Floor || enemy.Flags.OnLadder || enemy.Flags.LeapTo || enemy.Flags.LeapOn || enemy.Flags.LeapOff) &&
			(ch.Flags.Floor || ch.Flags.OnLadder || ch.Flags.LeapTo || ch.Flags.LeapOn || ch.Flags.LeapOff) {
			ch.Flags.Hit = true
			enemy.Flags.Attack = true
		}
	}
}
