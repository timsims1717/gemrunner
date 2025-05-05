package animations

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"github.com/gopxl/pixel"
)

func DemonAnimation(ch *data.Dynamic) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, "demon_idle", reanimator.Hold)
	regen := reanimator.NewBatchAnimation("regen", batch, "demon_regen", reanimator.Tran)
	regen.SetEndTrigger(func() {
		ch.Flags.Regen = false
	})
	chase := reanimator.NewBatchSprite("chase", batch, "demon_chase", reanimator.Hold)
	run := reanimator.NewBatchAnimation("run", batch, "demon_run", reanimator.Loop)
	climb := reanimator.NewBatchAnimation("climb", batch, "demon_climb", reanimator.Loop)
	climb.SetTriggerAll(func() {
		climb.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	slide := reanimator.NewBatchAnimationCustom("slide", batch, "demon_climb", []int{0, 6, 5, 4, 2, 1}, reanimator.Loop)
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
	fall := reanimator.NewBatchSprite("fall", batch, "demon_fall", reanimator.Hold)
	jump := reanimator.NewBatchSprite("jump", batch, "demon_jump", reanimator.Hold)
	leapOnI := []int{0, 4, 3, 2, 5}
	leapOffI := []int{6, 6, 4, 4, 3, 3, 3, 2, 1}
	leapToI := []int{6, 6, 4, 3, 3, 3, 2, 2, 5}
	leapOn := reanimator.NewBatchAnimationCustom("leap_on", batch, "demon_leap", leapOnI, reanimator.Tran)
	leapOn.SetEndTrigger(func() {
		ch.Flags.LeapOn = false
		ch.ACounter = 0
	})
	leapOff := reanimator.NewBatchAnimationCustom("leap_off", batch, "demon_leap", leapOffI, reanimator.Tran)
	leapOff.SetEndTrigger(func() {
		ch.Flags.LeapOff = false
		ch.ACounter = 0
	})
	leapTo := reanimator.NewBatchAnimationCustom("leap_to", batch, "demon_leap", leapToI, reanimator.Tran)
	leapTo.SetEndTrigger(func() {
		ch.Flags.LeapTo = false
		ch.ACounter = 0
	})
	bar := reanimator.NewBatchAnimation("bar", batch, "demon_bar", reanimator.Loop)
	bar.SetTriggerAll(func() {
		bar.Freeze = !ch.Flags.Climbed
		ch.Flags.Climbed = false
	})
	fullAttack := []int{0, 1, 2, 3, 4, 5, 4, 5, 4, 5}
	attack := reanimator.NewBatchAnimationCustom("attack", batch, "demon_attack", fullAttack, reanimator.Tran)
	attack.SetEndTrigger(func() {
		ch.Flags.Attack = false
	})
	hit := reanimator.NewBatchAnimation("hit", batch, "demon_hit", reanimator.Tran)
	hit.SetEndTrigger(func() {
		ch.Flags.Death = death.None
	})
	drown := reanimator.NewBatchAnimation("drown", batch, "demon_hit", reanimator.Tran)
	drown.SetEndTrigger(func() {
		ch.Flags.Death = death.None
	})
	drown = drown.WithSpriteOffset(pixel.V(0, 6), 0)
	drown = drown.WithSpriteOffset(pixel.V(0, 4), 1)
	drown = drown.WithSpriteOffset(pixel.V(0, 2), 2)
	crush := reanimator.NewBatchAnimation("crush", batch, "demon_crush", reanimator.Tran)
	crush.SetEndTrigger(func() {
		ch.Flags.Death = death.None
	})
	transIn := reanimator.NewBatchAnimation("trans_in", batch, "demon_trans_in", reanimator.Hold)
	transExit := reanimator.NewBatchAnimation("trans_exit", batch, "demon_trans_out", reanimator.Tran)
	transIn.SetEndTrigger(func() {
		ch.Flags.Transport = true
	})
	transExit.SetEndTrigger(func() {
		ch.Flags.ItemAction = data.NoItemAction
		ch.Object.Layer = ch.Layer
	})
	sw := reanimator.NewSwitch().
		AddAnimation(regen).
		AddAnimation(idle).
		AddAnimation(run).
		AddAnimation(chase).
		AddAnimation(fall).
		AddAnimation(jump).
		AddAnimation(climb).
		AddAnimation(slide).
		AddAnimation(leapOn).
		AddAnimation(leapOff).
		AddAnimation(leapTo).
		AddAnimation(bar).
		AddAnimation(hit).
		AddAnimation(drown).
		AddAnimation(crush).
		AddAnimation(attack).
		AddAnimation(transIn).
		AddAnimation(transExit).
		AddNull("none").
		SetChooseFn(func() string {
			switch ch.State {
			case data.Regen:
				return "regen"
			case data.Dead:
				return "none"
			case data.Hit:
				switch ch.Flags.Death {
				case death.Crushed:
					return "crush"
				case death.Drowned:
					return "drown"
				case death.Exploded, death.Dying:
					return "hit"
				default:
					return "none"
				}
			case data.DoingAction:
				switch ch.Flags.ItemAction {
				//case data.MagicDig:
				//	return "dig"
				//case data.MagicPlace:
				//	return "dig"
				//case data.ThrowBox:
				//	return "throw"
				//case data.DonDisguise:
				//	return "don_disguise"
				//case data.DrillStart:
				//	return "drill_start"
				//case data.Drilling:
				//	return "drill"
				//case data.Hiding:
				//	return "hiding"
				//case data.FireFlamethrower:
				//	return "flamethrower"
				case data.TransportIn:
					return "trans_in"
				case data.TransportExit:
					return "trans_exit"
				default:
					return "idle"
				}
			case data.Attack:
				return "attack"
			case data.Grounded:
				if ch.Actions.Left() || ch.Actions.Right() {
					if (ch.Actions.Left() && (ch.Flags.LeftWall || ch.Flags.EnemyL)) ||
						(ch.Actions.Right() && (ch.Flags.RightWall || ch.Flags.EnemyR)) {
						return "chase"
					} else {
						return "run"
					}
				} else {
					return "idle"
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
			return "idle"
		})
	tree := reanimator.New(sw, "regen")
	return tree
}
