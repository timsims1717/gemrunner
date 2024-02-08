package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/controllers"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func CharacterSystem() {
	for _, result := range myecs.Manager.Query(myecs.IsCharacter) {
		obj, okO := result.Components[myecs.Object].(*object.Object)
		ch, okC := result.Components[myecs.Character].(*data.Character)
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
			debug.AddTruthText("Action: ", ch.Actions.Action)
			if !ch.Flags.Floor && !ch.Flags.LadderDown &&
				!ch.Flags.LadderHere && !ch.Flags.OnLadder {
				fall(ch, currTile)
			} else {
				upOrDown(ch, currTile)
				if ch.Flags.OnLadder {
					ch.FauxObj.Pos.X = currTile.Object.Pos.X
					leapOff(ch, currTile)
				} else {
					run(ch)
				}
			}
			if !ch.Flags.OnLadder {
				ch.Flags.GoingUp = false
				ch.Flags.WentUp = false
			}
			if reanimator.FrameSwitch {
				obj.Pos = ch.FauxObj.Pos
				obj.Flip = ch.FauxObj.Flip
			}
		}
	}
}

func fall(ch *data.Character, tile *data.Tile) {
	ch.FauxObj.Pos.X = tile.Object.Pos.X
	ch.FauxObj.Pos.Y -= float64(reanimator.FRate) * constants.PlayerGravity * timing.DT
}

func upOrDown(ch *data.Character, tile *data.Tile) {
	if ch.Actions.Up && !ch.Flags.Ceiling {
		if ch.Flags.LadderHere {
			ch.FauxObj.Pos.Y += float64(reanimator.FRate) * constants.PlayerClimbSpeed * timing.DT
			ch.FauxObj.Flip = false
			ch.Flags.OnLadder = true
			ch.Flags.GoingUp = true
			ch.Flags.WentUp = true
		} else if ch.Flags.OnLadder {
			ch.FauxObj.Pos.Y = tile.Object.Pos.Y
			ch.FauxObj.Flip = false
			ch.Flags.OnLadder = false
		}
	} else if ch.Actions.Down {
		if (ch.Flags.LadderHere || ch.Flags.LadderDown) &&
			!ch.Flags.Floor {
			ch.FauxObj.Pos.Y -= float64(reanimator.FRate) * constants.PlayerDownSpeed * timing.DT
			ch.FauxObj.Flip = false
			ch.Flags.OnLadder = true
			ch.Flags.GoingUp = false
		}
	}
}

func leapOff(ch *data.Character, tile *data.Tile) {
	if ch.Actions.Left && !ch.Flags.LeftWall {
		ch.Flags.OnLadder = false
		ch.FauxObj.Flip = true
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X -= float64(reanimator.FRate) * constants.PlayerWalkSpeed * timing.DT
		} else {
			// todo: replace with leap left
			ch.FauxObj.Pos.X -= world.TileSize
		}
	} else if ch.Actions.Right && !ch.Flags.RightWall {
		ch.Flags.OnLadder = false
		ch.FauxObj.Flip = false
		ch.FauxObj.Pos.Y = tile.Object.Pos.Y
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X += float64(reanimator.FRate) * constants.PlayerWalkSpeed * timing.DT
		} else {
			// todo: replace with leap right
			ch.FauxObj.Pos.X += world.TileSize
		}
	}
}

func run(ch *data.Character) {
	if ch.Actions.Left && !ch.Flags.LeftWall {
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X -= float64(reanimator.FRate) * constants.PlayerWalkSpeed * timing.DT
			ch.FauxObj.Flip = true
			ch.Flags.OnLadder = false
		}
	} else if ch.Actions.Right && !ch.Flags.RightWall {
		if ch.Flags.CanRun {
			ch.FauxObj.Pos.X += float64(reanimator.FRate) * constants.PlayerWalkSpeed * timing.DT
			ch.FauxObj.Flip = false
			ch.Flags.OnLadder = false
		}
	}
}

func PlayerCharacter(pos pixel.Vec, pIndex int) *data.Character {
	obj := object.New()
	obj.SetRect(pixel.R(0, 0, 14, 16))
	obj.Pos = pos
	obj.Layer = 11
	fObj := object.New()
	fObj.Pos = pos
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Object, obj)
	e.AddComponent(myecs.Controller, controllers.NewPlayerInput("left", "right", "up", "down", "jump", "action"))
	p1 := &data.Character{
		Object:      obj,
		FauxObj:     fObj,
		Entity:      e,
		PlayerIndex: pIndex,
	}
	p1.Anim = PlayerAnimation(p1)
	e.AddComponent(myecs.Animated, p1.Anim)
	e.AddComponent(myecs.Drawable, p1.Anim)
	e.AddComponent(myecs.Character, p1)
	e.AddComponent(myecs.Collector, struct{}{})
	e.AddComponent(myecs.Temp, myecs.ClearFlag(false))
	return p1
}

func PlayerAnimation(player *data.Character) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, "player_idle", reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, "player_idle", reanimator.Loop)
	run := reanimator.NewBatchAnimationCustom("run", batch, "player_run", []int{0, 1, 2, 1}, reanimator.Loop)
	wall := reanimator.NewBatchAnimationFrame("wall", batch, "player_run", 1, reanimator.Loop)
	fall := reanimator.NewBatchSprite("fall", batch, "player_fall", reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, "player_jump", reanimator.Hold)
	climb := reanimator.NewBatchAnimation("climb", batch, "player_climb", reanimator.Loop)
	climb.SetTriggerAll(func() {
		climb.Freeze = !player.Flags.WentUp
		player.Flags.WentUp = false
	})
	slide := reanimator.NewBatchSprite("slide", batch, "player_slide", reanimator.Hold)
	return reanimator.New(reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(breath).
		AddAnimation(run).
		AddAnimation(wall).
		AddAnimation(fall).
		AddAnimation(jump).
		AddAnimation(climb).
		AddAnimation(slide).
		SetChooseFn(func() string {
			if player.Flags.OnLadder {
				if player.Actions.Up || player.Flags.GoingUp {
					return "climb"
				} else {
					return "slide"
				}
			} else if player.Flags.CanRun {
				if player.Actions.Left || player.Actions.Right {
					if (player.Actions.Left && player.Flags.LeftWall) ||
						(player.Actions.Right && player.Flags.RightWall) {
						return "wall"
					} else {
						return "run"
					}
				} else {
					return "idle"
				}
			} else {
				return "fall"
			}
		}), "breath")
}
