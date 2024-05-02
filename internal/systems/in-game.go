package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/reanimator"
)

func InGameSystem() {
	if reanimator.FrameSwitch && data.CurrLevel.Start {
		data.CurrLevel.FrameNumber++
		data.CurrLevel.FrameCounter++
		if data.CurrLevel.FrameCounter == constants.FrameCycle {
			data.CurrLevel.FrameCounter = 0
			data.CurrLevel.FrameCycle++
			data.CurrLevel.FrameChange = true
		} else {
			data.CurrLevel.FrameChange = false
		}
	} else {
		data.CurrLevel.FrameChange = false
	}
	data.CurrLevel.DoorsOpen = len(myecs.Manager.Query(myecs.IsGem)) < 1
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
	ControlEnemiesSystem()
}

func ControlEnemiesSystem() {
	if len(data.CurrLevel.Enemies) > 0 {
		p1 := data.CurrLevel.Players[0]
		p2 := data.CurrLevel.Players[1]
		p3 := data.CurrLevel.Players[2]
		p4 := data.CurrLevel.Players[3]
		if p1 != nil && data.DebugInput.Get("beBadGuyP1").JustPressed() {
			i := p1.Enemy + 1
			for {
				if i > len(data.CurrLevel.Enemies)-1 {
					// if we reach the end of the enemies list, reset to control the player
					ResetControl(data.CurrLevel.Enemies[p1.Enemy])
					ResetControl(p1)
					p1.Enemy = -1
					break
				} else if (p2 != nil && p2.Enemy == i) ||
					(p3 != nil && p3.Enemy == i) ||
					(p4 != nil && p4.Enemy == i) {
					// someone else is controlling this enemy
					i++
				} else {
					// no one is controlling this enemy, switch to it
					if p1.Enemy > -1 {
						// but first reset the enemy that's being controlled
						ResetControl(data.CurrLevel.Enemies[p1.Enemy])
					}
					p1.Enemy = i
					SetPlayerControl(data.CurrLevel.Enemies[p1.Enemy], 0)
					break
				}
			}
		}
		if p2 != nil && data.DebugInput.Get("beBadGuyP2").JustPressed() {
			i := p2.Enemy + 1
			for {
				if i >= len(data.CurrLevel.Enemies)-1 {
					// if we reach the end of the enemies list, reset to control the player
					ResetControl(data.CurrLevel.Enemies[p2.Enemy])
					ResetControl(p2)
					p2.Enemy = -1
					break
				} else if (p1 != nil && p1.Enemy == i) ||
					(p3 != nil && p3.Enemy == i) ||
					(p4 != nil && p4.Enemy == i) {
					// someone else is controlling this enemy
					i++
				} else {
					// no one is controlling this enemy, switch to it
					if p2.Enemy > -1 {
						// but first reset the enemy that's being controlled
						ResetControl(data.CurrLevel.Enemies[p2.Enemy])
					}
					p2.Enemy = i
					SetPlayerControl(data.CurrLevel.Enemies[p2.Enemy], 1)
					break
				}
			}
		}
		if p3 != nil && data.DebugInput.Get("beBadGuyP3").JustPressed() {
			i := p3.Enemy + 1
			for {
				if i >= len(data.CurrLevel.Enemies)-1 {
					// if we reach the end of the enemies list, reset to control the player
					ResetControl(data.CurrLevel.Enemies[p3.Enemy])
					ResetControl(p3)
					p3.Enemy = -1
					break
				} else if (p1 != nil && p1.Enemy == i) ||
					(p2 != nil && p2.Enemy == i) ||
					(p4 != nil && p4.Enemy == i) {
					// someone else is controlling this enemy
					i++
				} else {
					// no one is controlling this enemy, switch to it
					if p3.Enemy > -1 {
						// but first reset the enemy that's being controlled
						ResetControl(data.CurrLevel.Enemies[p3.Enemy])
					}
					p3.Enemy = i
					SetPlayerControl(data.CurrLevel.Enemies[p3.Enemy], 2)
					break
				}
			}
		}
		if p4 != nil && data.DebugInput.Get("beBadGuyP4").JustPressed() {
			i := p4.Enemy + 1
			for {
				if i >= len(data.CurrLevel.Enemies)-1 {
					// if we reach the end of the enemies list, reset to control the player
					ResetControl(data.CurrLevel.Enemies[p4.Enemy])
					ResetControl(p4)
					p4.Enemy = -1
					break
				} else if (p1 != nil && p1.Enemy == i) ||
					(p2 != nil && p2.Enemy == i) ||
					(p3 != nil && p3.Enemy == i) {
					// someone else is controlling this enemy
					i++
				} else {
					// no one is controlling this enemy, switch to it
					if p4.Enemy > -1 {
						// but first reset the enemy that's being controlled
						ResetControl(data.CurrLevel.Enemies[p4.Enemy])
					}
					p4.Enemy = i
					SetPlayerControl(data.CurrLevel.Enemies[p4.Enemy], 3)
					break
				}
			}
		}
	}
}
