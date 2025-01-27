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
	return reanimator.New(reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(flying).
		AddAnimation(boom).
		AddAnimation(crush).
		AddAnimation(regen).
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
			case data.Regen:
				return "regen"
			case data.Dead:
				return "none"
			default:
				return "flying"
			}
		}), "flying")
}
