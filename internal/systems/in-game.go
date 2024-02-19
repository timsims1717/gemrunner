package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
)

func InGameSystem() {
	if data.P1Input.Get("speedUp").JustPressed() {
		constants.FrameRate += constants.FrameRateInt
		if constants.FrameRate > constants.FrameRateMax {
			constants.FrameRate = constants.FrameRateMax
		}
	} else if data.P1Input.Get("speedDown").JustPressed() {
		constants.FrameRate -= constants.FrameRateInt
		if constants.FrameRate < constants.FrameRateMin {
			constants.FrameRate = constants.FrameRateMin
		}
	}
	if data.P1Input.Get("kill").JustPressed() {
		p1 := data.CurrLevel.Players[0]
		if p1 != nil && p1.State != data.Dead {
			p1.Flags.Hit = true
			p1.State = data.Hit
		}
	}
	if data.P2Input.Get("kill").JustPressed() {
		p2 := data.CurrLevel.Players[1]
		if p2 != nil && p2.State != data.Dead {
			p2.Flags.Hit = true
			p2.State = data.Hit
		}
	}
	if data.P3Input.Get("kill").JustPressed() {
		p3 := data.CurrLevel.Players[2]
		if p3 != nil && p3.State != data.Dead {
			p3.Flags.Hit = true
			p3.State = data.Hit
		}
	}
	if data.P4Input.Get("kill").JustPressed() {
		p4 := data.CurrLevel.Players[3]
		if p4 != nil && p4.State != data.Dead {
			p4.Flags.Hit = true
			p4.State = data.Hit
		}
	}
}
