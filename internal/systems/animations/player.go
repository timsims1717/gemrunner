package animations

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
	"gemrunner/internal/random"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func PlayerAnimation(ch *data.Dynamic, sprPre string, triggers bool) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Tran)

	regenFrames := []int{0, 1, 2, 3, 4, 5, 5, 6, 6, 6, 7}
	regen := reanimator.NewBatchAnimationCustom("regen", batch, fmt.Sprintf("%s_regen", sprPre), regenFrames, reanimator.Tran)

	wall := reanimator.NewBatchAnimationFrame("wall", batch, fmt.Sprintf("%s_run", sprPre), 2, reanimator.Hold)
	run := reanimator.NewBatchAnimation("run", batch, fmt.Sprintf("%s_run", sprPre), reanimator.Loop)
	goop := reanimator.NewBatchAnimation("goop", batch, fmt.Sprintf("%s_goop", sprPre), reanimator.Loop)
	climb := reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)

	slide := reanimator.NewBatchSprite("slide", batch, fmt.Sprintf("%s_slide", sprPre), reanimator.Hold)
	bar := reanimator.NewBatchAnimation("bar", batch, fmt.Sprintf("%s_bar", sprPre), reanimator.Loop)

	digFrames := []int{0, 0, 1, 2, 3, 3, 4, 4, 4, 4}
	dig := reanimator.NewBatchAnimationCustom("dig", batch, fmt.Sprintf("%s_dig", sprPre), digFrames, reanimator.Tran)

	fall := reanimator.NewBatchSprite("fall", batch, fmt.Sprintf("%s_fall", sprPre), reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, fmt.Sprintf("%s_jump", sprPre), reanimator.Hold)

	landingFrames := []int{5, 5, 5, 5, 6, 6, 6, 6, 6, 6, 6, 7}
	landing := reanimator.NewBatchAnimationCustom("landing", batch, fmt.Sprintf("%s_regen", sprPre), landingFrames, reanimator.Tran)

	//leapOnI := []int{1, 2}
	//leapOffI := []int{2, 0, 1, 2}
	//leapToI := []int{2, 0, 1, 2, 2}
	//leapOn := reanimator.NewBatchAnimationCustom("leap_on", batch, fmt.Sprintf("%s_leap", sprPre), leapOnI, reanimator.Tran)
	//leapOff := reanimator.NewBatchAnimationCustom("leap_off", batch, fmt.Sprintf("%s_leap", sprPre), leapOffI, reanimator.Tran)
	//leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, fmt.Sprintf("%s_leap", sprPre), leapToI, reanimator.Tran)

	leapOn := reanimator.NewBatchAnimation("leap_on", batch, fmt.Sprintf("%s_leap", sprPre), reanimator.Tran)
	leapOff := reanimator.NewBatchAnimation("leap_off", batch, fmt.Sprintf("%s_leap_off", sprPre), reanimator.Tran)
	//leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, fmt.Sprintf("%s_leap", sprPre), leapToI, reanimator.Tran)

	throw := reanimator.NewBatchAnimation("throw", batch, fmt.Sprintf("%s_throw", sprPre), reanimator.Tran)

	jetpackLoop := []int{0, 0, 0, 1, 1, 1}
	jetpack := reanimator.NewBatchAnimationCustom("jetpack", batch, fmt.Sprintf("%s_jetpack", sprPre), jetpackLoop, reanimator.Loop)

	jetpackUp := reanimator.NewBatchAnimationCustom("jetpack_up", batch, fmt.Sprintf("%s_jetpack_up", sprPre), jetpackLoop, reanimator.Loop)
	jetpackDown := reanimator.NewBatchAnimationCustom("jetpack_down", batch, fmt.Sprintf("%s_jetpack_down", sprPre), jetpackLoop, reanimator.Loop)

	donDisguise := reanimator.NewBatchAnimation("don_disguise", batch, fmt.Sprintf("%s_don", sprPre), reanimator.Tran)

	drillStart := reanimator.NewBatchAnimation("drill_start", batch, fmt.Sprintf("%s_drill_start", sprPre), reanimator.Tran)
	drilling := reanimator.NewBatchAnimation("drill", batch, fmt.Sprintf("%s_drill", sprPre), reanimator.Loop)

	flameframes := []int{0, 1, 2, 3, 3, 2, 1, 0}
	flamethrower := reanimator.NewBatchAnimationCustom("flamethrower", batch, fmt.Sprintf("%s_flamethrower", sprPre), flameframes, reanimator.Loop)

	goopThrow := reanimator.NewBatchAnimation("goop_throw", batch, fmt.Sprintf("%s_goop_throw", sprPre), reanimator.Tran)

	airCannonFrames := []int{0, 0, 0, 1, 2, 3, 3, 3, 3, 4, 5}
	airCannon := reanimator.NewBatchAnimationCustom("air_cannon", batch, fmt.Sprintf("%s_air_cannon", sprPre), airCannonFrames, reanimator.Tran)

	hiding := reanimator.NewBatchAnimation("hiding", batch, fmt.Sprintf("%s_hiding", sprPre), reanimator.Tran)
	inHiding := reanimator.NewBatchSprite("in_hiding", batch, fmt.Sprintf("%s_in_hiding", sprPre), reanimator.Hold)

	fullHit := []int{0, 1, 2, 3, 4, 5, 5, 5, 5, 5}
	hit := reanimator.NewBatchAnimationCustom("hit", batch, fmt.Sprintf("%s_hit", sprPre), fullHit, reanimator.Tran)

	crush := reanimator.NewBatchAnimation("crush", batch, fmt.Sprintf("%s_crush", sprPre), reanimator.Tran)

	blow := reanimator.NewBatchAnimation("blow", batch, "exp_player", reanimator.Tran)
	drown := reanimator.NewBatchAnimation("drown", batch, fmt.Sprintf("%s_drown", sprPre), reanimator.Tran)
	drown = drown.WithOffset(pixel.V(0, 2))
	drown = drown.WithSpriteOffset(pixel.V(0, 3), 0)
	drown = drown.WithSpriteOffset(pixel.V(0, 1), 3)

	portalWait := reanimator.NewBatchAnimation("portal", batch, "portal_magic", reanimator.Loop)

	transIn := reanimator.NewBatchAnimation("trans_in", batch, fmt.Sprintf("%s_trans_in", sprPre), reanimator.Hold)
	transExit := reanimator.NewBatchAnimation("trans_exit", batch, fmt.Sprintf("%s_trans_out", sprPre), reanimator.Tran)
	// triggers
	if triggers {
		breath.SetEndTrigger(func() {
			ch.Flags.Breath = false
		})
		regen.SetTriggerAll(func() {
			if regen.Step > 2 && !ch.Flags.Floor {
				ch.Flags.Regen = false
			}
		})
		regen.SetEndTrigger(func() {
			ch.Flags.Regen = false
		})
		climb.SetTriggerAll(func() {
			climb.Freeze = !ch.Flags.Climbed
			ch.Flags.Climbed = false
		})
		bar.SetTriggerAll(func() {
			bar.Freeze = !ch.Flags.Climbed
			ch.Flags.Climbed = false
		})
		dig.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
		})
		landing.SetEndTrigger(func() {
			ch.Flags.Landing = false
		})
		leapOn.SetEndTrigger(func() {
			ch.Flags.LeapTo = false
			ch.Flags.LeapOn = false
			ch.ACounter = 0
		})
		leapOff.SetEndTrigger(func() {
			ch.Flags.LeapOff = false
			ch.ACounter = 0
		})
		//leapTo.SetEndTrigger(func() {
		//	ch.Flags.LeapTo = false
		//	ch.ACounter = 0
		//})
		throw.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
		})
		goopThrow.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
		})
		airCannon.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			switch f {
			case 3, 4, 5, 6:
				if ch.Object.Flip {
					if !ch.Flags.RightWall {
						ch.Object.Pos.X += 8. - float64(f)
					}
				} else {
					if !ch.Flags.LeftWall {
						ch.Object.Pos.X -= 8. - float64(f)
					}
				}
			case 0, 7:
				x, y := world.WorldToMap(ch.Object.Pos.X, ch.Object.Pos.Y)
				tile := data.CurrLevel.Get(x, y)
				if tile != nil {
					ch.Object.SetPos(tile.Object.Pos)
				}
			}
		})
		airCannon.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
		})
		donDisguise.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
			// set the player to disguised
			ch.Flags.Disguised = true
		})
		drillStart.SetEndTrigger(func() {
			ch.Flags.CheckAction = true
		})
		//drilling.SetTriggerAll(func() {
		//	ch.Object.Pos.Y -= 1.
		//})
		drilling.SetTrigger(1, func() {
			ch.Object.Pos.Y -= 2.
		})
		hiding.SetEndTrigger(func() {
			ch.State = data.InHiding
			ch.Flags.ItemAction = data.NoItemAction
		})
		jetpack.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			switch pre {
			case "jetpack_up", "jetpack_down":
				a.Step = f + 1
				a.Step %= len(jetpackLoop)
			case "jetpack":
			default:
				if ch.AnInt > -1 && ch.AnInt < len(jetpackLoop) {
					a.Step = ch.AnInt
				}
			}
		})
		jetpackUp.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			switch pre {
			case "jetpack", "jetpack_down":
				a.Step = f + 1
				a.Step %= len(jetpackLoop)
			case "jetpack_up":
			default:
				if ch.AnInt > -1 && ch.AnInt < len(jetpackLoop) {
					a.Step = ch.AnInt
				}
			}
		})
		jetpackDown.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			switch pre {
			case "jetpack", "jetpack_up":
				a.Step = f + 1
				a.Step %= len(jetpackLoop)
			case "jetpack_down":
			default:
				if ch.AnInt > -1 && ch.AnInt < len(jetpackLoop) {
					a.Step = ch.AnInt
				}
			}
		})
		hit.SetEndTrigger(func() {
			ch.Flags.NextStep = true
		})
		crush.SetEndTrigger(func() {
			ch.Flags.NextStep = true
		})
		blow.SetEndTrigger(func() {
			ch.Flags.NextStep = true
		})
		drown.SetEndTrigger(func() {
			ch.Flags.NextStep = true
		})
		transIn.SetEndTrigger(func() {
			ch.Flags.Transport = true
		})
		transExit.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
			ch.Object.Layer = ch.Layer
		})
	}
	// offsets
	idle.Offset.Y--
	breath.Offset.Y--
	run.Offset.Y--
	goop.Offset.Y--
	wall.Offset.Y--
	jump.Offset.Y--
	landing.Offset.Y--
	dig.Offset.Y--
	throw.Offset.Y--
	donDisguise.Offset.Y--
	drillStart.Offset.Y--
	drilling.Offset.Y--
	flamethrower.Offset.Y--
	goopThrow.Offset.Y--
	airCannon.Offset.Y--
	hit.Offset.Y--

	tree := reanimator.New().
		AddAnimation(regen).
		AddAnimation(idle).
		AddAnimation(breath).
		AddAnimation(run).
		AddAnimation(goop).
		AddAnimation(wall).
		AddAnimation(fall).
		AddAnimation(jump).
		AddAnimation(landing).
		AddAnimation(dig).
		AddAnimation(climb).
		AddAnimation(slide).
		AddAnimation(bar).
		AddAnimation(leapOn).
		AddAnimation(leapOff).
		//AddAnimation(leapTo).
		AddAnimation(throw).
		AddAnimation(jetpack).
		AddAnimation(jetpackUp).
		AddAnimation(jetpackDown).
		AddAnimation(donDisguise).
		AddAnimation(drillStart).
		AddAnimation(drilling).
		AddAnimation(flamethrower).
		AddAnimation(goopThrow).
		AddAnimation(airCannon).
		AddAnimation(hiding).
		AddAnimation(inHiding).
		AddAnimation(hit).
		AddAnimation(crush).
		AddAnimation(blow).
		AddAnimation(drown).
		AddAnimation(portalWait).
		AddAnimation(transIn).
		AddAnimation(transExit).
		AddNull("none").
		SetChooseFn(func() string {
			if sprPre == "disguise" && !ch.Flags.Disguised {
				return "none"
			}
			switch ch.State {
			case data.Waiting:
				return "portal"
			case data.Regen:
				return "regen"
			case data.Dead:
				return "none"
			case data.Hit:
				switch ch.Flags.Death {
				case death.Crushed:
					return "crush"
				case death.Exploded:
					return "blow"
				case death.Drowned:
					return "drown"
				case death.Dying:
					return "hit"
				default:
					return "none"
				}
			case data.Flying:
				if ch.Actions.Direction == data.Up {
					return "jetpack_up"
				} else if ch.Actions.Direction == data.Down {
					return "jetpack_down"
				} else {
					return "jetpack"
				}
			case data.DoingAction:
				switch ch.Flags.ItemAction {
				case data.MagicDig:
					return "dig"
				case data.MagicPlace:
					return "dig"
				case data.ThrowBox:
					return "throw"
				case data.DonDisguise:
					return "don_disguise"
				case data.DrillStart:
					return "drill_start"
				case data.Drilling:
					return "drill"
				case data.Hiding:
					return "hiding"
				case data.FireFlamethrower:
					return "flamethrower"
				case data.ThrowingGoop:
					return "goop_throw"
				case data.UseAirCannon:
					return "air_cannon"
				case data.TransportIn:
					return "trans_in"
				case data.TransportExit:
					return "trans_exit"
				default:
					return "idle"
				}
			case data.InHiding:
				return "in_hiding"
			case data.Grounded:
				if ch.Actions.Left() || ch.Actions.Right() {
					ch.Flags.Landing = false
					if (ch.Actions.Left() && (ch.Flags.LeftWall || ch.Flags.EnemyL)) ||
						(ch.Actions.Right() && (ch.Flags.RightWall || ch.Flags.EnemyR)) {
						return "wall"
					} else if ch.Flags.Goop {
						return "goop"
					} else {
						return "run"
					}
				} else {
					if ch.Flags.Landing {
						return "landing"
					} else if ch.Flags.Breath {
						return "breath"
					} else {
						if !ch.Flags.Breath && random.Effects.Intn(constants.IdleFrequency*timing.FPS) == 0 {
							ch.Flags.Breath = true
						}
						return "idle"
					}
				}
			case data.OnLadder:
				if ch.Actions.Up() || ch.Flags.GoingUp {
					return "climb"
				} else {
					return "slide"
				}
			case data.OnBar:
				return "bar"
			case data.Jumping:
				return "jump"
			case data.Falling:
				return "fall"
			case data.Leaping:
				if ch.Flags.LeapOff {
					return "leap_off"
				} else if ch.Flags.LeapOn || ch.Flags.LeapTo {
					return "leap_on"
					//} else if ch.Flags.LeapTo {
					//	return "leap_to"
				}
			}
			return "fall"
		}).SetDefault("regen").Finish()
	return tree
}
