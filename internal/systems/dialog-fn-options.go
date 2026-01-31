package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/content"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/sfx"
	"github.com/gopxl/pixel"
	"strconv"
)

func OpenOptions() {
	options := ui.Dialogs[constants.DialogOptions]
	fr := constants.Configuration.Gameplay.FrameRate
	var sliderValue int
	if fr <= 40 {
		sliderValue = fr/5 - 1
	} else {
		sliderValue = (fr + 30) / 10
	}
	for _, e := range options.Elements {
		ele := e
		switch ele.Key {
		case "gameplay_tab":
			ele.Get("label_shadow").Text.Show()
			ele.Border.Style = data.ThinBorderWhite
		case "graphics_tab", "audio_tab", "input_tab":
			ele.Get("label_shadow").Text.Hide()
			ele.Border.Style = data.ThinBorderBlue
		case "gameplay_tab_display", "graphics_tab_display", "audio_tab_display", "input_tab_display":
			ele.Object.Hidden = ele.Key != "gameplay_tab_display"
			for _, e1 := range ele.Elements {
				ele1 := e1
				switch ele1.Key {
				case "game_speed_value":
					ui.SetText(ele1, strconv.Itoa(sliderValue))
				case "game_speed_slider":
					ui.SetSliderValue(ele1, sliderValue)
				case "timer_check":
					ui.SetChecked(ele1, constants.Configuration.Gameplay.ShowTimer)
				case "screen_shake_check":
					ui.SetChecked(ele1, constants.Configuration.Gameplay.ScreenShake)
				case "record_check":
					ui.SetChecked(ele1, constants.Configuration.Gameplay.AlwaysRecord)
				case "scan_lines_check":
					ui.SetChecked(ele1, constants.Configuration.Graphics.Scanlines)
				case "bilinear_check":
					ui.SetChecked(ele1, constants.Configuration.Graphics.BilinearFilter)
				case "vsync_check":
					ui.SetChecked(ele1, constants.Configuration.Graphics.VSync)
				case "full_screen_check":
					ui.SetChecked(ele1, constants.Configuration.Graphics.Fullscreen)
				case "resolution_dropdown_scroll":
					ui.UpdateScrollBounds(ele1)
				case "set_color_check":
					ui.SetChecked(ele1, constants.Configuration.Graphics.SetColorMode)
				case "music_on":
					ui.SetChecked(ele1, constants.Configuration.Audio.MusicOn)
				case "music_value":
					ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.MusicVolume))
				case "music_slider":
					ui.SetSliderValue(ele1, constants.Configuration.Audio.MusicVolume)
				case "sfx_on":
					ui.SetChecked(ele1, constants.Configuration.Audio.SfxOn)
				case "sfx_value":
					ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.SfxVolume))
				case "sfx_slider":
					ui.SetSliderValue(ele1, constants.Configuration.Audio.SfxVolume)
				case "master_on":
					ui.SetChecked(ele1, constants.Configuration.Audio.MasterOn)
				case "master_value":
					ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.MasterVolume))
				case "master_slider":
					ui.SetSliderValue(ele1, constants.Configuration.Audio.MasterVolume)
				case "mute_check":
					ui.SetChecked(ele1, constants.Configuration.Audio.MuteUnfocus)
				}
			}
		}
	}
	data.OriginalConfiguration = constants.Configuration.Copy()
	ui.OpenDialogInStack(constants.DialogOptions)
}

func ResetOptions() {
	options := ui.Dialogs[constants.DialogOptions]
	fr := constants.Configuration.Gameplay.FrameRate
	var sliderValue int
	if fr <= 40 {
		sliderValue = (fr - 5) / 5
	} else {
		sliderValue = (fr + 30) / 10
	}
	for _, e := range options.Elements {
		ele := e
		switch ele.Key {
		case "gameplay_tab_display", "graphics_tab_display", "audio_tab_display", "input_tab_display":
			if !ele.Object.Hidden {
				switch ele.Key {
				case "gameplay_tab_display":
					constants.Configuration.Gameplay = constants.DefaultConfiguration.Gameplay.Copy()
				case "graphics_tab_display":
					constants.Configuration.Graphics = constants.DefaultConfiguration.Graphics.Copy()
				case "audio_tab_display":
					constants.Configuration.Audio = constants.DefaultConfiguration.Audio.Copy()
				case "input_tab_display":
					// by player?
					// constants.Configuration.Input = constants.DefaultConfiguration.Input.Copy()
				}
				for _, e1 := range ele.Elements {
					ele1 := e1
					switch ele1.Key {
					case "game_speed_value":
						ui.SetText(ele1, strconv.Itoa(sliderValue))
					case "game_speed_slider":
						ui.SetSliderValue(ele1, sliderValue)
					case "timer_check":
						ui.SetChecked(ele1, constants.Configuration.Gameplay.ShowTimer)
					case "screen_shake_check":
						ui.SetChecked(ele1, constants.Configuration.Gameplay.ScreenShake)
					case "record_check":
						ui.SetChecked(ele1, constants.Configuration.Gameplay.AlwaysRecord)
					case "scan_lines_check":
						ui.SetChecked(ele1, constants.Configuration.Graphics.Scanlines)
					case "bilinear_check":
						ui.SetChecked(ele1, constants.Configuration.Graphics.BilinearFilter)
					case "vsync_check":
						ui.SetChecked(ele1, constants.Configuration.Graphics.VSync)
					case "full_screen_check":
						ui.SetChecked(ele1, constants.Configuration.Graphics.Fullscreen)
					case "set_color_check":
						ui.SetChecked(ele1, constants.Configuration.Graphics.SetColorMode)
					case "music_value":
						ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.MusicVolume))
					case "music_slider":
						ui.SetSliderValue(ele1, constants.Configuration.Audio.MusicVolume)
					case "music_on":
						ui.SetChecked(ele1, constants.Configuration.Audio.MusicOn)
					case "sfx_value":
						ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.SfxVolume))
					case "sfx_slider":
						ui.SetSliderValue(ele1, constants.Configuration.Audio.SfxVolume)
					case "sfx_on":
						ui.SetChecked(ele1, constants.Configuration.Audio.SfxOn)
					case "master_value":
						ui.SetText(ele1, strconv.Itoa(constants.Configuration.Audio.MasterVolume))
					case "master_slider":
						ui.SetSliderValue(ele1, constants.Configuration.Audio.MasterVolume)
					case "master_on":
						ui.SetChecked(ele1, constants.Configuration.Audio.MasterOn)
					case "mute_check":
						ui.SetChecked(ele1, constants.Configuration.Audio.MuteUnfocus)
					}
				}
			}
		}
	}
	content.UpdateConfiguration()
}

func CancelOptions() {
	constants.Configuration = data.OriginalConfiguration.Copy()
	content.UpdateConfiguration()
	ui.CloseDialog(constants.DialogOptions)
}

func ConfirmOptions() {
	options := ui.Dialogs[constants.DialogOptions]
	for _, e := range options.Elements {
		ele := e
		switch ele.Key {
		case "gameplay_tab_display", "graphics_tab_display", "audio_tab_display", "input_tab_display":
			for _, e1 := range ele.Elements {
				ele1 := e1
				switch ele1.Key {
				case "game_speed_slider":
					sv := ele1.IntValue
					var fr int
					if sv <= 8 {
						fr = (sv + 1) * 5
					} else {
						fr = sv*10 - 30
					}
					constants.Configuration.Gameplay.FrameRate = fr
				case "timer_check":
					constants.Configuration.Gameplay.ShowTimer = ele1.Checked
				case "screen_shake_check":
					constants.Configuration.Gameplay.ScreenShake = ele1.Checked
				case "record_check":
					constants.Configuration.Gameplay.AlwaysRecord = ele1.Checked
				case "scan_lines_check":
					constants.Configuration.Graphics.Scanlines = ele1.Checked
				case "bilinear_check":
					constants.Configuration.Graphics.BilinearFilter = ele1.Checked
				case "vsync_check":
					constants.Configuration.Graphics.VSync = ele1.Checked
				case "full_screen_check":
					constants.Configuration.Graphics.Fullscreen = ele1.Checked
				case "set_color_check":
					constants.Configuration.Graphics.SetColorMode = ele1.Checked
				case "music_slider":
					constants.Configuration.Audio.MusicVolume = ele1.IntValue
				case "sfx_slider":
					constants.Configuration.Audio.SfxVolume = ele1.IntValue
				case "master_slider":
					constants.Configuration.Audio.MasterVolume = ele1.IntValue
				case "mute_check":
					constants.Configuration.Audio.MuteUnfocus = ele1.Checked
				}
			}
		}
	}
	go content.SaveConfig()
	content.UpdateConfiguration()
	ui.CloseDialog(constants.DialogOptions)
}

func customizeOptionsDialog() {
	optionsDlg := ui.Dialogs[constants.DialogOptions]
	for _, e := range optionsDlg.Elements {
		ele := e
		switch ele.Key {
		case "confirm":
			ele.OnClick = ConfirmOptions
		case "default":
			ele.OnClick = ResetOptions
		case "cancel":
			ele.OnClick = CancelOptions
		case "gameplay_tab", "graphics_tab", "audio_tab", "input_tab":
			ele.OnClick = func() {
				for _, dle := range optionsDlg.Elements {
					switch dle.Key {
					case "gameplay_tab", "graphics_tab", "audio_tab", "input_tab":
						if dle.ElementType == ui.ContainerElement {
							if ele.Key == dle.Key {
								dle.Border.Style = data.ThinBorderWhite
							} else {
								dle.Border.Style = data.ThinBorderBlue
							}
							for _, txt1 := range dle.Elements {
								if ele.Key == dle.Key && txt1.Key == "label_shadow" {
									txt1.Text.Show()
								} else if txt1.Key == "label_shadow" {
									txt1.Text.Hide()
								}
							}
						}
					case "gameplay_tab_display":
						dle.Object.Hidden = ele.Key != "gameplay_tab"
					case "graphics_tab_display":
						dle.Object.Hidden = ele.Key != "graphics_tab"
					case "audio_tab_display":
						dle.Object.Hidden = ele.Key != "audio_tab"
					case "input_tab_display":
						dle.Object.Hidden = ele.Key != "input_tab"
					}
				}
			}
			ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, optionsDlg.ViewPort, func(hvc *data.HoverClick) {
				if optionsDlg.Open && optionsDlg.Active {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						ele.OnClick()
						click.Consume()
					}
				}
			}))
		case "gameplay_tab_display", "graphics_tab_display", "audio_tab_display", "input_tab_display":
			for _, e1 := range ele.Elements {
				ele1 := e1
				switch ele1.Key {
				case "game_speed_value":
					ele1.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if optionsDlg.Open && optionsDlg.Active && !ele.Object.Hidden {
							slider := ele.Get("game_speed_slider")
							iVal := strconv.Itoa(slider.IntValue)
							if slider != nil && ele1.Value != iVal {
								ui.SetText(ele1, iVal)
								constants.Configuration.Gameplay.FrameRate = slider.IntValue * 5
								content.UpdateConfiguration()
							}
						}
					}))
				case "timer_check":
					ele1.OnClick = func() {
						constants.Configuration.Gameplay.ShowTimer = ele1.Checked
						if data.CurrLevel != nil {
							if ele1.Checked {
								ui.OpenDialog(constants.DialogPuzzleTimer)
							} else {
								ui.CloseDialog(constants.DialogPuzzleTimer)
							}
						}
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "scan_lines_check":
					ele1.OnClick = func() {
						constants.Configuration.Graphics.Scanlines = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "bilinear_check":
					ele1.OnClick = func() {
						constants.Configuration.Graphics.BilinearFilter = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "full_screen_check":
					ele1.OnClick = func() {
						constants.Configuration.Graphics.Fullscreen = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "vsync_check":
					ele1.OnClick = func() {
						constants.Configuration.Graphics.VSync = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "set_color_check":
					ele1.OnClick = func() {
						constants.Configuration.Graphics.SetColorMode = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "music_on":
					ele1.OnClick = func() {
						constants.Configuration.Audio.MusicOn = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "music_value":
					ele1.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if optionsDlg.Open && optionsDlg.Active && !ele.Object.Hidden {
							slider := ele.Get("music_slider")
							iVal := strconv.Itoa(slider.IntValue)
							if slider != nil && ele1.Value != iVal {
								ui.SetText(ele1, iVal)
								constants.Configuration.Audio.MusicVolume = slider.IntValue
								content.UpdateConfiguration()
							}
						}
					}))
				case "sfx_on":
					ele1.OnClick = func() {
						constants.Configuration.Audio.SfxOn = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "sfx_value":
					ele1.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if optionsDlg.Open && optionsDlg.Active && !ele.Object.Hidden {
							slider := ele.Get("sfx_slider")
							iVal := strconv.Itoa(slider.IntValue)
							if slider != nil && ele1.Value != iVal {
								ui.SetText(ele1, iVal)
								constants.Configuration.Audio.SfxVolume = slider.IntValue
								content.UpdateConfiguration()
							}
						}
					}))
				case "master_on":
					ele1.OnClick = func() {
						constants.Configuration.Audio.MasterOn = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "master_value":
					ele1.Entity.AddComponent(myecs.Update, data.NewFn(func() {
						if optionsDlg.Open && optionsDlg.Active && !ele.Object.Hidden {
							slider := ele.Get("master_slider")
							iVal := strconv.Itoa(slider.IntValue)
							if slider != nil && ele1.Value != iVal {
								ui.SetText(ele1, iVal)
								constants.Configuration.Audio.MasterVolume = slider.IntValue
								content.UpdateConfiguration()
							}
						}
					}))
				case "sfx_slider", "master_slider":
					ele1.OnClick = func() {
						sfx.SoundPlayer.PlaySound(constants.SFXGem, 0.)
					}
				case "mute_check":
					ele1.OnClick = func() {
						constants.Configuration.Audio.MuteUnfocus = ele1.Checked
					}
					SetOptionsCheckboxFunc(ele1, optionsDlg, ele)
				case "resolution_dropdown_scroll":
					ele1.Elements = []*ui.Element{}
					for i, res := range constants.Resolutions {
						j := i
						t := fmt.Sprintf("%dx%d", int(res.X), int(res.Y))
						x := -34.
						y := float64(i)*-18 + 7
						tt := ui.ElementConstructor{
							Key:         fmt.Sprintf("resolution_%d", i),
							Text:        t,
							Position:    pixel.V(x, y),
							ElementType: ui.TextElement,
							Anchor:      pixel.Right,
						}
						tte := ui.CreateTextElement(tt, ele1.ViewPort)
						tte.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
						tte.OnClick = func() {
							constants.Configuration.Graphics.Resolution = j
							content.UpdateConfiguration()
						}
						ele1.Elements = append(ele1.Elements, tte)
					}
					ui.UpdateDropdownElements(ele1, optionsDlg.Get("resolution_dropdown"), optionsDlg)
				}
			}
		}
	}
}

func SetOptionsCheckboxFunc(ele *ui.Element, options *ui.Dialog, cnt *ui.Element) {
	ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, options.ViewPort, func(hvc *data.HoverClick) {
		if options.Open && options.Active && !options.Lock && !options.Click && !cnt.Object.Hidden {
			click := hvc.Input.Get("click")
			if hvc.Hover && click.JustPressed() {
				ui.SetChecked(ele, !ele.Checked)
				ele.OnClick()
				content.UpdateConfiguration()
			}
		}
	}))
}
