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
	})
	return reanimator.New(reanimator.NewSwitch().
		AddAnimation(idle).
		AddAnimation(flying).
		AddAnimation(boom).
		AddNull("none").
		SetChooseFn(func() string {
			if ch.State == data.Hit || ch.State == data.Attack {
				return "boom"
			} else if ch.State == data.Dead {
				return "none"
			}
			return "flying"
		}), "flying")
}
