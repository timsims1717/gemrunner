package reanimator

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
	"math/rand"
)

func HumanoidAnimation(ch *data.Dynamic, sprPre string) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Tran)
	breath.SetEndTrigger(func() {
		ch.Flags.Breath = false
	})
	var run, climb, slide *reanimator.Anim
	var leapOnI, leapOffI, leapToI []int
	if sprPre == "demon" {
		run = reanimator.NewBatchAnimationCustom("run", batch, fmt.Sprintf("%s_run", sprPre), []int{0, 1, 2, 3, 4, 1, 2, 3}, reanimator.Loop)
		climb = reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)
		slide = reanimator.NewBatchAnimationCustom("slide", batch, fmt.Sprintf("%s_climb", sprPre), []int{0, 4, 3, 1}, reanimator.Loop)
		slide.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			if pre == "climb" {
				switch f {
				case 0, 1:
					a.Step = 0
				case 2:
					a.Step = 3
				case 3, 4:
					a.Step = 2
				case 5:
					a.Step = 1
				}
			}
			slide.Freeze = !ch.Flags.Climbed
			ch.Flags.Climbed = false
		})
		leapOnI = []int{0, 1, 2}
		leapOffI = []int{3, 2, 1, 0}
		leapToI = []int{3, 2, 1, 1, 2, 3}
	} else {
		run = reanimator.NewBatchAnimationCustom("run", batch, fmt.Sprintf("%s_run", sprPre), []int{0, 1, 2, 1}, reanimator.Loop)
		climb = reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)
		slide = reanimator.NewBatchSprite("slide", batch, fmt.Sprintf("%s_slide", sprPre), reanimator.Hold)
		leapOnI = []int{2, 2}
		leapOffI = []int{0, 1, 2}
		leapToI = []int{0, 1, 2, 2}
	}
	climb.SetTriggerAll(func() {
		climb.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	wall := reanimator.NewBatchAnimationFrame("wall", batch, fmt.Sprintf("%s_run", sprPre), 1, reanimator.Hold)
	fall := reanimator.NewBatchSprite("fall", batch, fmt.Sprintf("%s_fall", sprPre), reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, fmt.Sprintf("%s_jump", sprPre), reanimator.Hold)
	leapOn := reanimator.NewBatchAnimationCustom("leap_on", batch, fmt.Sprintf("%s_leap", sprPre), leapOnI, reanimator.Tran)
	leapOn.SetEndTrigger(func() {
		ch.Flags.LeapOn = false
		ch.ACounter = 0
		ch.ATimer = timing.New(ch.Vars.LeapDelay / float64(reanimator.FRate))
	})
	leapOff := reanimator.NewBatchAnimationCustom("leap_off", batch, fmt.Sprintf("%s_leap", sprPre), leapOffI, reanimator.Tran)
	leapOff.SetEndTrigger(func() {
		ch.Flags.LeapOff = false
		ch.ACounter = 0
		ch.ATimer = timing.New(ch.Vars.LeapDelay / float64(reanimator.FRate))
	})
	leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, fmt.Sprintf("%s_leap", sprPre), leapToI, reanimator.Tran)
	leapTo.SetEndTrigger(func() {
		ch.Flags.LeapTo = false
		ch.ACounter = 0
		ch.ATimer = timing.New(ch.Vars.LeapDelay / float64(reanimator.FRate))
	})
	pickUp := reanimator.NewBatchSprite("pick_up", batch, fmt.Sprintf("%s_pick_up", sprPre), reanimator.Tran)
	pickUp.SetTrigger(0, func() {
		ch.Flags.PickUp = false
	})
	holdIdle := reanimator.NewBatchAnimationFrame("hold_idle", batch, fmt.Sprintf("%s_pick_up", sprPre), 1, reanimator.Hold)
	holdRun := reanimator.NewBatchAnimationCustom("hold_run", batch, fmt.Sprintf("%s_hold_run", sprPre), []int{0, 1, 2, 1}, reanimator.Loop)
	holdIdleSide := reanimator.NewBatchSprite("hold_idle_side", batch, fmt.Sprintf("%s_pick_up", sprPre), reanimator.Hold)
	holdRunSide := reanimator.NewBatchAnimationCustom("hold_run_side", batch, fmt.Sprintf("%s_hold_run_side", sprPre), []int{0, 1, 2, 1}, reanimator.Loop)
	fallHold := reanimator.NewBatchSprite("fall_hold", batch, fmt.Sprintf("%s_fall_hold", sprPre), reanimator.Hold)
	jumpHold := reanimator.NewBatchSprite("jump_hold", batch, fmt.Sprintf("%s_jump_hold", sprPre), reanimator.Hold)
	fallHoldSide := reanimator.NewBatchSprite("fall_hold_side", batch, fmt.Sprintf("%s_fall_hold_side", sprPre), reanimator.Hold)
	jumpHoldSide := reanimator.NewBatchSprite("jump_hold_side", batch, fmt.Sprintf("%s_jump_hold_side", sprPre), reanimator.Hold)
	fullHit := []int{0, 1, 2, 3, 3, 3, 3, 3, 3, 3}
	fullAttack := []int{0, 1, 2, 3, 2, 3, 2, 3, 2, 3}
	hit := reanimator.NewBatchAnimationCustom("hit", batch, fmt.Sprintf("%s_hit", sprPre), fullHit, reanimator.Tran)
	hit.SetEndTrigger(func() {
		ch.Flags.Hit = false
	})
	attack := reanimator.NewBatchAnimationCustom("attack", batch, fmt.Sprintf("%s_attack", sprPre), fullAttack, reanimator.Tran)
	attack.SetEndTrigger(func() {
		ch.Flags.Attack = false
	})
	sw := reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(breath).
		AddAnimation(run).
		AddAnimation(wall).
		AddAnimation(fall).
		AddAnimation(jump).
		AddAnimation(climb).
		AddAnimation(slide).
		AddAnimation(leapOn).
		AddAnimation(leapOff).
		AddAnimation(leapTo).
		AddAnimation(pickUp).
		AddAnimation(holdIdle).
		AddAnimation(holdRun).
		AddAnimation(holdIdleSide).
		AddAnimation(holdRunSide).
		AddAnimation(fallHold).
		AddAnimation(jumpHold).
		AddAnimation(fallHoldSide).
		AddAnimation(jumpHoldSide).
		AddAnimation(hit).
		AddAnimation(attack).
		AddNull("none").
		SetChooseFn(func() string {
			switch ch.State {
			case data.Dead:
				return "none"
			case data.Hit:
				return "hit"
			case data.Attack:
				return "attack"
			case data.Grounded:
				if ch.Flags.PickUp {
					return "pick_up"
				}
				if ch.Actions.Left() || ch.Actions.Right() {
					if (ch.Actions.Left() && ch.Flags.LeftWall) ||
						(ch.Actions.Right() && ch.Flags.RightWall) {
						if ch.Flags.HoldUp {
							return "hold_idle"
						} else if ch.Flags.HoldSide {
							return "hold_idle_side"
						} else {
							return "wall"
						}
					} else {
						if ch.Flags.HoldUp {
							return "hold_run"
						} else if ch.Flags.HoldSide {
							return "hold_run_side"
						} else {
							return "run"
						}
					}
				} else {
					if ch.Flags.HoldUp {
						return "hold_idle"
					} else if ch.Flags.HoldSide {
						return "hold_idle_side"
					} else {
						if ch.Flags.Breath {
							return "breath"
						} else {
							if !ch.Flags.Breath && rand.Intn(constants.IdleFrequency*timing.FPS) == 0 {
								ch.Flags.Breath = true
							}
							return "idle"
						}
					}
				}
			case data.Ladder:
				if ch.Actions.Up() || ch.Flags.GoingUp {
					return "climb"
				} else {
					return "slide"
				}
			case data.Jumping:
				if ch.Flags.PickUp {
					return "pick_up"
				} else if ch.Flags.HoldUp {
					return "jump_hold"
				} else if ch.Flags.HoldSide {
					return "jump_hold_side"
				} else {
					return "jump"
				}
			case data.Falling:
				if ch.Flags.PickUp {
					return "pick_up"
				} else if ch.Flags.HoldUp {
					return "fall_hold"
				} else if ch.Flags.HoldSide {
					return "fall_hold_side"
				} else {
					return "fall"
				}
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
	tree := reanimator.New(sw, "breath")
	return tree
}
