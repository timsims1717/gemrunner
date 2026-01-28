package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/data/death"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/debug"
	"gemrunner/pkg/img"
	"gemrunner/pkg/reanimator"
	"gemrunner/pkg/sfx"
	"github.com/gopxl/pixel/imdraw"
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
	wereDoorsOpen := data.CurrLevel.DoorsOpen
	if !wereDoorsOpen {
		data.CurrLevel.DoorsOpen = len(myecs.Manager.Query(myecs.IsGem)) < 1
		if data.CurrLevel.DoorsOpen {
			sfx.SoundPlayer.PlaySound(constants.SFXDoorsOpen, 0.)
		}
	}
	if data.CurrLevel.FakePlayer != nil && debug.ShowDebug {
		debug.AddLine(constants.ColorRed, imdraw.RoundEndShape, data.CurrLevel.FakePlayer.Object.Pos, data.CurrLevel.FakePlayer.Object.Pos, 4.)
	}
	PlayerMetaInput()
	UpdatePlayerInv()
	ControlEnemiesSystem()
	RecordingSystem()
}

func PlayerMetaInput() {
	if data.P1Input.Get("speedUp").JustPressed() {
		if constants.Configuration.Gameplay.FrameRate >= 40 {
			constants.Configuration.Gameplay.FrameRate += constants.FrameRateInt
		}
		constants.Configuration.Gameplay.FrameRate += constants.FrameRateInt
		if constants.Configuration.Gameplay.FrameRate > constants.FrameRateMax {
			constants.Configuration.Gameplay.FrameRate = constants.FrameRateMax
		}
	} else if data.P1Input.Get("speedDown").JustPressed() {
		if constants.Configuration.Gameplay.FrameRate > 40 {
			constants.Configuration.Gameplay.FrameRate -= constants.FrameRateInt
		}
		constants.Configuration.Gameplay.FrameRate -= constants.FrameRateInt
		if constants.Configuration.Gameplay.FrameRate < constants.FrameRateMin {
			constants.Configuration.Gameplay.FrameRate = constants.FrameRateMin
		}
	}
	if data.P1Input.Get("kill").JustPressed() {
		p1 := data.CurrLevel.Players[0]
		if p1 != nil && p1.State != data.Dead {
			p1.Flags.Death = death.Dying
			p1.State = data.Hit
		}
	}
	if data.P2Input.Get("kill").JustPressed() {
		p2 := data.CurrLevel.Players[1]
		if p2 != nil && p2.State != data.Dead {
			p2.Flags.Death = death.Dying
			p2.State = data.Hit
		}
	}
	if data.P3Input.Get("kill").JustPressed() {
		p3 := data.CurrLevel.Players[2]
		if p3 != nil && p3.State != data.Dead {
			p3.Flags.Death = death.Dying
			p3.State = data.Hit
		}
	}
	if data.P4Input.Get("kill").JustPressed() {
		p4 := data.CurrLevel.Players[3]
		if p4 != nil && p4.State != data.Dead {
			p4.Flags.Death = death.Dying
			p4.State = data.Hit
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

func UpdatePlayerInv() {
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
			bombs := data.CurrLevel.Players[i].SmallBombs
			if bombs > 4 {
				if bombs > 99 {
					bombs = 99
				}
			}
			scoreDialog := ui.Dialogs[dgKey]
			for _, ele := range scoreDialog.Elements {
				switch ele.Key {
				case "player_score":
					score := data.CurrLevelSess.PlayerStats[i].Score + data.CurrLevelSess.PlayerStats[i].LScore
					if score > 9999999 {
						score = 9999999
					}
					ele.Text.SetText(fmt.Sprintf("%07d", score))
				case "player_deaths":
					deaths := data.CurrLevelSess.PlayerStats[i].Deaths
					if deaths > 999 {
						deaths = 999
					}
					ele.Text.SetText(fmt.Sprintf("x%03d", deaths))
				case "player_gems":
					gems := data.CurrLevelSess.PlayerStats[i].Gems + data.CurrLevelSess.PlayerStats[i].LGems
					if gems > 99999 {
						gems = 99999
					}
					ele.Text.SetText(fmt.Sprintf("x%05d", gems))
				case "player_bombs":
					if bombs > 4 {
						ele.Text.Show()
						ele.Text.SetText(fmt.Sprintf("x%02d", bombs))
					} else {
						ele.Text.Hide()
					}
				case "player_bomb_count_1":
					ele.Object.Hidden = bombs < 1 || bombs > 4
				case "player_bomb_count_2":
					ele.Object.Hidden = bombs < 2 || bombs > 4
				case "player_bomb_count_3":
					ele.Object.Hidden = bombs < 3 || bombs > 4
				case "player_bomb_count_4":
					ele.Object.Hidden = bombs < 4
				case "player_inv_item":
					if p.Inventory == nil {
						ele.Object.Hidden = true
						ele.Sprite.Key = ""
					} else {
						if draw, ok := p.Inventory.Entity.GetComponentData(myecs.Drawable); ok {
							if a, okA := draw.(*reanimator.Tree); okA {
								ele.Object.Hidden = false
								ele.Sprite.Key = a.GetSprite(a.Default).SKey
							} else if s, okS := draw.(*img.Sprite); okS {
								ele.Object.Hidden = false
								ele.Sprite.Key = s.Key
							} else if sa, okSA := draw.([]interface{}); okSA && len(sa) > 0 {
								if saSpr, ok2 := sa[0].(*img.Sprite); ok2 {
									ele.Object.Hidden = false
									ele.Sprite.Key = saSpr.Key
								} else if saAnim, ok3 := sa[0].(*reanimator.Tree); ok3 {
									ele.Object.Hidden = false
									ele.Sprite.Key = saAnim.GetSprite(saAnim.Default).SKey
								}
							}
						}
					}
				case "player_inv_item_2":
					if p.Inventory == nil {
						ele.Object.Hidden = true
						ele.Sprite.Key = ""
					} else {
						if draw, ok := p.Inventory.Entity.GetComponentData(myecs.Drawable); ok {
							if sa, okSA := draw.([]interface{}); okSA && len(sa) > 1 {
								if saSpr, ok2 := sa[1].(*img.Sprite); ok2 {
									ele.Object.Hidden = false
									ele.Object.Offset = saSpr.Offset
									ele.Sprite.Key = saSpr.Key
								}
							}
						}
					}
				case "player_inv_cnt":
					for _, ele2 := range ele.Elements {
						switch ele2.Key {
						case "player_inv_regen":
							if p.Inventory != nil {
								ele2.Object.Hidden = !p.Inventory.Metadata.Regenerate
							} else {
								ele2.Object.Hidden = true
							}
						case "player_inv_timer":
							if p.Inventory != nil {
								if p.Inventory.Metadata.Timer == -1 {
									ele2.Object.Hidden = true
								} else if p.Inventory.Metadata.Timer == 0 {
									// draw infinity symbol
									ele2.Object.Hidden = true
								} else {
									ele2.Text.SetText(fmt.Sprintf("%d", p.Inventory.Metadata.Timer-p.Inventory.Uses))
									ele2.Object.Hidden = false
								}
							} else {
								ele2.Object.Hidden = true
							}
						case "player_inv_infinity":
							if p.Inventory != nil {
								ele2.Object.Hidden = p.Inventory.Metadata.Timer != 0
							} else {
								ele2.Object.Hidden = true
							}
						}
					}
				}
			}
		}
	}
}

func RecordingSystem() {
	if data.CurrLevel.Recording {
		if (data.CurrLevel.ReplayFrame.P1Actions != nil &&
			data.CurrLevel.ReplayFrame.P1Actions.Any()) ||
			(data.CurrLevel.ReplayFrame.P2Actions != nil &&
				data.CurrLevel.ReplayFrame.P2Actions.Any()) ||
			(data.CurrLevel.ReplayFrame.P3Actions != nil &&
				data.CurrLevel.ReplayFrame.P3Actions.Any()) ||
			(data.CurrLevel.ReplayFrame.P4Actions != nil &&
				data.CurrLevel.ReplayFrame.P4Actions.Any()) {
			if data.CurrLevel.FrameNumber > 0 {
				data.CurrLevel.ReplayFrame.Frame = data.CurrLevel.FrameNumber - 1
				data.CurrLevel.LevelReplay.Frames = append(data.CurrLevel.LevelReplay.Frames, data.CurrLevel.ReplayFrame)
			}
			data.CurrLevel.ReplayFrame = data.ReplayFrame{}
		}
	}
}
