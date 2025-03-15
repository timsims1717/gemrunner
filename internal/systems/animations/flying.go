package animations

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
)

func FlyAnimation(ch *data.Dynamic) *reanimator.Tree {
	batch := img.Batchers[constants.TileBatch]
	idle := reanimator.NewBatchSprite("idle", batch, "fly_idle", reanimator.Hold)
	flying := reanimator.NewBatchAnimation("flying", batch, "fly_flying", reanimator.Loop)
	boom := reanimator.NewBatchAnimation("boom", batch, "fly_boom", reanimator.Tran)
	boom.SetEndTrigger(func() {
		ch.State = data.Dead
		ch.Flags.Hit = false
		ch.Flags.Attack = false
		ch.Flags.Crush = false
	})
	crush := reanimator.NewBatchAnimation("crush", batch, "fly_crush", reanimator.Tran)
	crush.SetEndTrigger(func() {
		ch.Flags.Hit = false
		ch.Flags.Attack = false
		ch.Flags.Crush = false
		ch.Flags.Blow = false
	})
	regen := reanimator.NewBatchAnimation("regen", batch, "fly_regen", reanimator.Tran)
	regen.SetEndTrigger(func() {
		ch.Flags.Regen = false
	})
	transIn := reanimator.NewBatchAnimation("trans_in", batch, "fly_trans_in", reanimator.Hold)
	transExit := reanimator.NewBatchAnimation("trans_exit", batch, "fly_trans_out", reanimator.Tran)
	transIn.SetEndTrigger(func() {
		ch.Flags.Transport = true
	})
	transExit.SetEndTrigger(func() {
		ch.Flags.ItemAction = data.NoItemAction
		ch.Object.Layer = ch.Layer
	})
	return reanimator.New(reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(flying).
		AddAnimation(boom).
		AddAnimation(crush).
		AddAnimation(regen).
		AddAnimation(transIn).
		AddAnimation(transExit).
		AddNull("none").
		SetChooseFn(func() string {
			switch ch.State {
			case data.Attack:
				return "boom"
			case data.Hit:
				if ch.Flags.Crush {
					return "crush"
				} else if ch.Flags.Hit {
					return "boom"
				} else {
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
			case data.Regen:
				return "regen"
			case data.Dead:
				return "none"
			default:
				return "flying"
			}
		}), "flying")
}
