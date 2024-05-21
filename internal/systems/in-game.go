package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"time"
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
	UpdateScoreUI()
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

func FormatTimePlayed() string {
	if data.CurrLevelSess == nil {
		return ""
	}
	d := data.CurrLevelSess.TotalTime + data.CurrLevelSess.TimePlayed
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func UpdateScoreUI() {
	for i, p := range data.CurrLevel.Players {
		if p != nil {
			dgKey := constants.DialogPlayer1Inv
			switch i {
			case 1:
				dgKey = constants.DialogPlayer2Inv
			case 2:
				dgKey = constants.DialogPlayer3Inv
			case 3:
				dgKey = constants.DialogPlayer4Inv
			}
			scoreDialog := ui.Dialogs[dgKey]
			for _, ele := range scoreDialog.Elements {
				switch ele.Key {
				case "player_score":
					ele.Text.SetText(fmt.Sprintf("%07d", data.CurrLevelSess.PlayerStats[i].Score+data.CurrLevelSess.PlayerStats[i].LScore))
				case "player_deaths":
					ele.Text.SetText(fmt.Sprintf("x%03d", data.CurrLevelSess.PlayerStats[i].Deaths))
				case "player_gems":
					ele.Text.SetText(fmt.Sprintf("x%04d", data.CurrLevelSess.PlayerStats[i].Gems+data.CurrLevelSess.PlayerStats[i].LGems))
				case "player_inv_item":
					if p.Inventory == nil {
						ele.Object.Hidden = true
						ele.Sprite.Key = ""
					} else {
						if draw, ok := p.Inventory.GetComponentData(myecs.Drawable); ok {
							if a, okA := draw.(*reanimator.Tree); okA {
								ele.Object.Hidden = false
								ele.Sprite.Key = a.Default
							} else if s, okS := draw.(*img.Sprite); okS {
								ele.Object.Hidden = false
								ele.Sprite.Key = s.Key
							} else if sa, okSA := draw.([]interface{}); okSA && len(sa) > 0 {
								if saSpr, ok2 := sa[0].(*img.Sprite); ok2 {
									ele.Object.Hidden = false
									ele.Sprite.Key = saSpr.Key
								} else if saAnim, ok3 := sa[0].(*reanimator.Tree); ok3 {
									ele.Object.Hidden = false
									ele.Sprite.Key = saAnim.Default
								}
							}
						}
					}
				case "player_inv_item_2":
					if p.Inventory == nil {
						ele.Object.Hidden = true
						ele.Sprite.Key = ""
					} else {
						if draw, ok := p.Inventory.GetComponentData(myecs.Drawable); ok {
							if sa, okSA := draw.([]interface{}); okSA && len(sa) > 1 {
								if saSpr, ok2 := sa[1].(*img.Sprite); ok2 {
									ele.Object.Hidden = false
									ele.Object.Offset = saSpr.Offset
									ele.Sprite.Key = saSpr.Key
								}
							}
						}
					}
				}
			}
		}
	}
}
