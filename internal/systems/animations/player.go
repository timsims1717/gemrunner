package animations

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/random"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
)

func PlayerAnimation(ch *data.Dynamic, sprPre string, triggers bool) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Tran)

	regenFrames := []int{0, 1, 2, 3, 4, 5, 6, 6, 7}
	regen := reanimator.NewBatchAnimationCustom("regen", batch, fmt.Sprintf("%s_regen", sprPre), regenFrames, reanimator.Tran)

	wall := reanimator.NewBatchAnimationFrame("wall", batch, fmt.Sprintf("%s_run", sprPre), 2, reanimator.Hold)
	run := reanimator.NewBatchAnimation("run", batch, fmt.Sprintf("%s_run", sprPre), reanimator.Loop)
	climb := reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)

	slide := reanimator.NewBatchSprite("slide", batch, fmt.Sprintf("%s_slide", sprPre), reanimator.Hold)
	bar := reanimator.NewBatchAnimation("bar", batch, fmt.Sprintf("%s_bar", sprPre), reanimator.Loop)

	digFrames := []int{0, 0, 1, 2, 3, 3, 4, 4, 4, 4}
	dig := reanimator.NewBatchAnimationCustom("dig", batch, fmt.Sprintf("%s_dig", sprPre), digFrames, reanimator.Tran)

	fall := reanimator.NewBatchSprite("fall", batch, fmt.Sprintf("%s_fall", sprPre), reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, fmt.Sprintf("%s_jump", sprPre), reanimator.Hold)
	leapOnI := []int{1, 2}
	leapOffI := []int{2, 0, 1, 2}
	leapToI := []int{2, 0, 1, 2, 2}
	leapOn := reanimator.NewBatchAnimationCustom("leap_on", batch, fmt.Sprintf("%s_leap", sprPre), leapOnI, reanimator.Tran)

	leapOff := reanimator.NewBatchAnimationCustom("leap_off", batch, fmt.Sprintf("%s_leap", sprPre), leapOffI, reanimator.Tran)

	leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, fmt.Sprintf("%s_leap", sprPre), leapToI, reanimator.Tran)

	throw := reanimator.NewBatchAnimation("throw", batch, fmt.Sprintf("%s_throw", sprPre), reanimator.Tran)

	jetpackLoop := []int{0, 0, 0, 1, 1, 1}
	jetpack := reanimator.NewBatchAnimationCustom("jetpack", batch, fmt.Sprintf("%s_jetpack", sprPre), jetpackLoop, reanimator.Loop)

	jetpackUp := reanimator.NewBatchAnimationCustom("jetpack_up", batch, fmt.Sprintf("%s_jetpack_up", sprPre), jetpackLoop, reanimator.Loop)
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
	jetpackDown := reanimator.NewBatchAnimationCustom("jetpack_down", batch, fmt.Sprintf("%s_jetpack_down", sprPre), jetpackLoop, reanimator.Loop)

	donDisguise := reanimator.NewBatchAnimation("don_disguise", batch, fmt.Sprintf("%s_don", sprPre), reanimator.Tran)

	fullHit := []int{0, 1, 2, 3, 4, 5, 5, 5, 5, 5}
	hit := reanimator.NewBatchAnimationCustom("hit", batch, fmt.Sprintf("%s_hit", sprPre), fullHit, reanimator.Tran)

	crush := reanimator.NewBatchAnimation("crush", batch, fmt.Sprintf("%s_crush", sprPre), reanimator.Tran)

	blow := reanimator.NewBatchAnimation("blow", batch, "exp_player", reanimator.Tran)

	portalWait := reanimator.NewBatchAnimation("portal", batch, "portal_magic", reanimator.Loop)
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
		leapOn.SetEndTrigger(func() {
			ch.Flags.LeapOn = false
			ch.ACounter = 0
		})
		leapOff.SetEndTrigger(func() {
			ch.Flags.LeapOff = false
			ch.ACounter = 0
		})
		leapTo.SetEndTrigger(func() {
			ch.Flags.LeapTo = false
			ch.ACounter = 0
		})
		throw.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
		})
		donDisguise.SetEndTrigger(func() {
			ch.Flags.ItemAction = data.NoItemAction
			// set the player to disguised
			ch.Flags.Disguised = true
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
			ch.Flags.Hit = false
			ch.Flags.Crush = false
			ch.Flags.Blow = false
		})
		crush.SetEndTrigger(func() {
			ch.Flags.Hit = false
			ch.Flags.Crush = false
			ch.Flags.Blow = false
		})
		blow.SetEndTrigger(func() {
			ch.Flags.Hit = false
			ch.Flags.Crush = false
			ch.Flags.Blow = false
		})
	}
	sw := reanimator.NewSwitch().
		AddAnimation(regen).
		AddAnimation(idle).
		AddAnimation(breath).
		AddAnimation(run).
		AddAnimation(wall).
		AddAnimation(fall).
		AddAnimation(jump).
		AddAnimation(dig).
		AddAnimation(climb).
		AddAnimation(slide).
		AddAnimation(bar).
		AddAnimation(leapOn).
		AddAnimation(leapOff).
		AddAnimation(leapTo).
		AddAnimation(throw).
		AddAnimation(jetpack).
		AddAnimation(jetpackUp).
		AddAnimation(jetpackDown).
		AddAnimation(donDisguise).
		AddAnimation(hit).
		AddAnimation(crush).
		AddAnimation(blow).
		AddAnimation(portalWait).
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
				if ch.Flags.Crush {
					return "crush"
				} else if ch.Flags.Blow {
					return "blow"
				} else if ch.Flags.Hit {
					return "hit"
				} else {
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
				default:
					return "idle"
				}
			case data.Grounded:
				if ch.Actions.Left() || ch.Actions.Right() {
					if (ch.Actions.Left() && (ch.Flags.LeftWall || ch.Flags.EnemyL)) ||
						(ch.Actions.Right() && (ch.Flags.RightWall || ch.Flags.EnemyR)) {
						return "wall"
					} else {
						return "run"
					}
				} else {
					if ch.Flags.Breath {
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
				} else if ch.Flags.LeapOn {
					return "leap_on"
				} else if ch.Flags.LeapTo {
					return "leap_to"
				}
			}
			return "fall"
		})
	tree := reanimator.New(sw, "regen")
	return tree
}
