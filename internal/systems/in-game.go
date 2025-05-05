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
	data.CurrLevel.DoorsOpen = len(myecs.Manager.Query(myecs.IsGem)) < 1
	if !wereDoorsOpen && data.CurrLevel.DoorsOpen {
		sfx.SoundPlayer.PlaySound(constants.SFXDoorsOpen, 0.)
	}
	if data.CurrLevel.FakePlayer != nil && debug.ShowDebug {
		debug.AddLine(constants.ColorRed, imdraw.RoundEndShape, data.CurrLevel.FakePlayer.Object.Pos, data.CurrLevel.FakePlayer.Object.Pos, 3.)
	}
	PlayerMetaInput()
	UpdatePlayerInv()
	ControlEnemiesSystem()
	RecordingSystem()
}

func PlayerMetaInput() {
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
						if draw, ok := p.Inventory.Entity.GetComponentData(myecs.Drawable); ok {
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
