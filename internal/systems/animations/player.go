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

func PlayerAnimation(ch *data.Dynamic, sprPre string) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Tran)
	breath.SetEndTrigger(func() {
		ch.Flags.Breath = false
	})
	regenFrames := []int{0, 1, 2, 3, 4, 5, 6, 6, 7}
	regen := reanimator.NewBatchAnimationCustom("regen", batch, fmt.Sprintf("%s_regen", sprPre), regenFrames, reanimator.Tran)
	regen.SetTriggerAll(func() {
		if regen.Step > 2 && !ch.Flags.Floor {
			ch.Flags.Regen = false
		}
	})
	regen.SetEndTrigger(func() {
		ch.Flags.Regen = false
	})
	wall := reanimator.NewBatchAnimationFrame("wall", batch, fmt.Sprintf("%s_run", sprPre), 2, reanimator.Hold)
	run := reanimator.NewBatchAnimation("run", batch, fmt.Sprintf("%s_run", sprPre), reanimator.Loop)
	climb := reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)
	climb.SetTriggerAll(func() {
		climb.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	slide := reanimator.NewBatchSprite("slide", batch, fmt.Sprintf("%s_slide", sprPre), reanimator.Hold)
	bar := reanimator.NewBatchAnimation("bar", batch, fmt.Sprintf("%s_bar", sprPre), reanimator.Loop)
	bar.SetTriggerAll(func() {
		bar.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	digFrames := []int{0, 0, 1, 2, 3, 3, 4, 4, 4, 4}
	dig := reanimator.NewBatchAnimationCustom("dig", batch, fmt.Sprintf("%s_dig", sprPre), digFrames, reanimator.Tran)
	dig.SetEndTrigger(func() {
		ch.Flags.ItemAction = data.NoItemAction
	})
	fall := reanimator.NewBatchSprite("fall", batch, fmt.Sprintf("%s_fall", sprPre), reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, fmt.Sprintf("%s_jump", sprPre), reanimator.Hold)
	leapOnI := []int{1, 2}
	leapOffI := []int{2, 0, 1, 2}
	leapToI := []int{2, 0, 1, 2, 2}
	leapOn := reanimator.NewBatchAnimationCustom("leap_on", batch, fmt.Sprintf("%s_leap", sprPre), leapOnI, reanimator.Tran)
	leapOn.SetEndTrigger(func() {
		ch.Flags.LeapOn = false
		ch.ACounter = 0
	})
	leapOff := reanimator.NewBatchAnimationCustom("leap_off", batch, fmt.Sprintf("%s_leap", sprPre), leapOffI, reanimator.Tran)
	leapOff.SetEndTrigger(func() {
		ch.Flags.LeapOff = false
		ch.ACounter = 0
	})
	leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, fmt.Sprintf("%s_leap", sprPre), leapToI, reanimator.Tran)
	leapTo.SetEndTrigger(func() {
		ch.Flags.LeapTo = false
		ch.ACounter = 0
	})
	throw := reanimator.NewBatchAnimation("throw", batch, fmt.Sprintf("%s_throw", sprPre), reanimator.Tran)
	throw.SetEndTrigger(func() {
		ch.Flags.ItemAction = data.NoItemAction
	})
	fullHit := []int{0, 1, 2, 3, 4, 5, 5, 5, 5, 5}
	hit := reanimator.NewBatchAnimationCustom("hit", batch, fmt.Sprintf("%s_hit", sprPre), fullHit, reanimator.Tran)
	hit.SetEndTrigger(func() {
		ch.Flags.Hit = false
		ch.Flags.Crush = false
	})
	crush := reanimator.NewBatchAnimation("crush", batch, fmt.Sprintf("%s_crush", sprPre), reanimator.Tran)
	crush.SetEndTrigger(func() {
		ch.Flags.Hit = false
		ch.Flags.Crush = false
	})
	portalWait := reanimator.NewBatchAnimation("portal", batch, "portal_magic", reanimator.Loop)
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
		AddAnimation(hit).
		AddAnimation(crush).
		AddAnimation(portalWait).
		AddNull("none").
		SetChooseFn(func() string {
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
				} else {
					return "hit"
				}
			case data.Attack:
				return "attack"
			case data.DoingAction:
				switch ch.Flags.ItemAction {
				case data.MagicDig:
					return "dig"
				case data.MagicPlace:
					return "dig"
				case data.ThrowBox:
					return "throw"
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
