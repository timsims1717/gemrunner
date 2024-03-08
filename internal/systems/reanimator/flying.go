package reanimator

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
	})
	return reanimator.New(reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(flying).
		AddAnimation(boom).
		AddAnimation(crush).
		AddNull("none").
		SetChooseFn(func() string {
			if ch.State == data.Hit || ch.State == data.Attack {
				if ch.Flags.Crush {
					return "crush"
				}
				return "boom"
			} else if ch.State == data.Dead {
				return "none"
			}
			return "flying"
		}), "flying")
}
