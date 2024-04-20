package reanimator

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/random"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/timing"
)

func HumanoidAnimation(ch *data.Dynamic, sprPre string) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	regen := reanimator.NewBatchAnimation("regen", batch, fmt.Sprintf("%s_regen", sprPre), reanimator.Tran)
	regen.SetEndTrigger(func() {
		ch.Flags.Regen = false
	})
	idle := reanimator.NewBatchSprite("idle", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Hold)
	breath := reanimator.NewBatchAnimation("breath", batch, fmt.Sprintf("%s_idle", sprPre), reanimator.Tran)
	breath.SetEndTrigger(func() {
		ch.Flags.Breath = false
	})
	var wall, run, climb, slide *reanimator.Anim
	var leapOnI, leapOffI, leapToI []int
	if sprPre == "demon" {
		wall = reanimator.NewBatchSprite("wall", batch, fmt.Sprintf("%s_chase", sprPre), reanimator.Hold)
		run = reanimator.NewBatchAnimation("run", batch, fmt.Sprintf("%s_run", sprPre), reanimator.Loop)
		climb = reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)
		slide = reanimator.NewBatchAnimationCustom("slide", batch, fmt.Sprintf("%s_climb", sprPre), []int{0, 6, 5, 4, 2, 1}, reanimator.Loop)
		slide.SetTriggerCAll(func(a *reanimator.Anim, pre string, f int) {
			if pre == "climb" {
				switch f {
				case 0, 7:
					a.Step = 0
				case 1:
					a.Step = 5
				case 2:
					a.Step = 4
				case 3, 4:
					a.Step = 3
				case 5:
					a.Step = 2
				case 6:
					a.Step = 1
				}
			}
			slide.Freeze = !ch.Flags.Climbed
			ch.Flags.Climbed = false
		})
		leapOnI = []int{0, 4, 3, 2, 5}
		leapOffI = []int{6, 6, 4, 4, 3, 3, 3, 2, 1}
		leapToI = []int{6, 6, 4, 3, 3, 3, 2, 2, 5}
	} else {
		wall = reanimator.NewBatchAnimationFrame("wall", batch, fmt.Sprintf("%s_run", sprPre), 2, reanimator.Hold)
		run = reanimator.NewBatchAnimation("run", batch, fmt.Sprintf("%s_run", sprPre), reanimator.Loop)
		climb = reanimator.NewBatchAnimation("climb", batch, fmt.Sprintf("%s_climb", sprPre), reanimator.Loop)
		slide = reanimator.NewBatchSprite("slide", batch, fmt.Sprintf("%s_slide", sprPre), reanimator.Hold)
		leapOnI = []int{1, 2}
		leapOffI = []int{2, 0, 1, 2}
		leapToI = []int{2, 0, 1, 2, 2}
	}
	climb.SetTriggerAll(func() {
		climb.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	fall := reanimator.NewBatchSprite("fall", batch, fmt.Sprintf("%s_fall", sprPre), reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, fmt.Sprintf("%s_jump", sprPre), reanimator.Hold)
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
	fullAttack := []int{0, 1, 2, 3, 4, 5, 4, 5, 4, 5}
	var attack, hit *reanimator.Anim
	if sprPre == "demon" {
		attack = reanimator.NewBatchAnimationCustom("attack", batch, fmt.Sprintf("%s_attack", sprPre), fullAttack, reanimator.Tran)
		hit = reanimator.NewBatchAnimation("hit", batch, fmt.Sprintf("%s_hit", sprPre), reanimator.Tran)
	} else {
		attack = reanimator.NewBatchAnimation("attack", batch, fmt.Sprintf("%s_attack", sprPre), reanimator.Tran)
		hit = reanimator.NewBatchAnimationCustom("hit", batch, fmt.Sprintf("%s_hit", sprPre), fullHit, reanimator.Tran)
	}
	attack.SetEndTrigger(func() {
		ch.Flags.Attack = false
	})
	hit.SetEndTrigger(func() {
		ch.Flags.Hit = false
		ch.Flags.Crush = false
	})
	crush := reanimator.NewBatchAnimation("crush", batch, fmt.Sprintf("%s_crush", sprPre), reanimator.Tran)
	crush.SetEndTrigger(func() {
		ch.Flags.Hit = false
		ch.Flags.Crush = false
	})
	sw := reanimator.NewSwitch().
		AddAnimation(regen).
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
		AddAnimation(throw).
		AddAnimation(hit).
		AddAnimation(crush).
		AddAnimation(attack).
		AddNull("none").
		SetChooseFn(func() string {
			switch ch.State {
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
				case data.ThrowBox:
					return "throw"
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
