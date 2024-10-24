package systems

import "gemrunner/internal/data"

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
