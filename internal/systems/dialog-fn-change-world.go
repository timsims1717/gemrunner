package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"strings"
)

// change world dialog

func CustomizeWorldDialog() {
	dialog := ui.Dialogs[constants.DialogChangeWorld]
	for _, e := range dialog.Elements {
		ele := e
		switch ele.Key {
		case "confirm_change_world":
			ele.OnClick = ConfirmChangeWorld
		case "cancel_change_world":
			ele.OnClick = DisposeDialog(constants.DialogChangeWorld)
		case "custom_world_check":
			ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
				if dialog.Open && dialog.Active && !dialog.Lock && !dialog.Click {
					click := hvc.Input.Get("click")
					if hvc.Hover && click.JustPressed() {
						ui.SetChecked(ele, !ele.Checked)
						data.CustomWorldSelected = ele.Checked
						for _, ele2 := range dialog.Elements {
							if ele2.ElementType == ui.TextElement {
								if o, okO := ele2.Entity.GetComponentData(myecs.Object); okO {
									if obj, okO1 := o.(*object.Object); okO1 {
										switch ele2.Key {
										case "primary_text", "secondary_text", "doodad_text":
											obj.Hidden = !ele.Checked
										}
									}
								}
							} else if ele2.ElementType == ui.CheckboxElement {
								if o, okO := ele2.Entity.GetComponentData(myecs.Object); okO {
									if obj, okO1 := o.(*object.Object); okO1 {
										if strings.Contains(ele2.Key, "check_primary") ||
											strings.Contains(ele2.Key, "check_secondary") ||
											strings.Contains(ele2.Key, "check_doodad") {
											obj.Hidden = !ele.Checked
										}
									}
								}
							} else if ele2.ElementType == ui.SpriteElement {
								if o, okO := ele2.Entity.GetComponentData(myecs.Object); okO {
									if obj, okO1 := o.(*object.Object); okO1 {
										if strings.Contains(ele2.Key, "color_primary") ||
											strings.Contains(ele2.Key, "color_secondary") ||
											strings.Contains(ele2.Key, "color_doodad") {
											obj.Hidden = !ele.Checked
										}
									}
								}
							}
						}
						if ele.Checked {
							for _, ele2 := range dialog.Elements {
								if ele2.ElementType == ui.CheckboxElement {
									if !data.CustomSelectedBefore {
										updateColorCheckboxWorld(ele2)
									} else if ele2.Checked {
										changeSelectedColor(ele2.Key)
									}
								}
							}
							worldDialogCustomShaders()
						} else {
							worldDialogNormalShaders()
						}
						data.CustomSelectedBefore = true
					}
				}
			}))
		case "world_list":
			for i := 0; i < constants.WorldCustom; i++ {
				index := i
				y := float64(i)*-18 + 7
				entry := ui.ElementConstructor{
					Key:         fmt.Sprintf(load.WorldListEntry.Key, i),
					Width:       load.WorldListEntry.Width,
					Height:      load.WorldListEntry.Height,
					HelpText:    fmt.Sprintf(load.WorldListEntry.HelpText, constants.WorldNames[i]),
					Position:    pixel.V(0, y),
					ElementType: load.WorldListEntry.ElementType,
				}
				tti := ui.ElementConstructor{
					Key:         load.TurfTileItem.Key,
					SprKey:      constants.WorldSprites[i],
					Batch:       load.TurfTileItem.Batch,
					Position:    load.TurfTileItem.Position,
					ElementType: load.TurfTileItem.ElementType,
				}
				entry.SubElements = append(entry.SubElements, tti)
				lti := ui.ElementConstructor{
					Key:         load.LadderTileItem.Key,
					SprKey:      constants.TileLadderMiddle,
					Batch:       load.LadderTileItem.Batch,
					Position:    load.LadderTileItem.Position,
					ElementType: load.LadderTileItem.ElementType,
				}
				entry.SubElements = append(entry.SubElements, lti)
				dti := ui.ElementConstructor{
					Key:         load.DoodadTileItem.Key,
					SprKey:      constants.WorldDoodads[i],
					Batch:       load.DoodadTileItem.Batch,
					Position:    load.DoodadTileItem.Position,
					ElementType: load.DoodadTileItem.ElementType,
				}
				entry.SubElements = append(entry.SubElements, dti)
				wti := ui.ElementConstructor{
					Key:         load.WorldTxtItem.Key,
					Text:        constants.WorldNames[i],
					Position:    load.WorldTxtItem.Position,
					ElementType: load.WorldTxtItem.ElementType,
					Anchor:      pixel.Right,
				}
				entry.SubElements = append(entry.SubElements, wti)
				wtt := ui.ElementConstructor{
					Key:         load.WorldTxtItem.Key,
					Text:        constants.WorldNames[i],
					Position:    load.WorldTxtItem.Position.Add(pixel.V(0, y+1)),
					ElementType: load.WorldTxtItem.ElementType,
					Anchor:      pixel.Right,
				}
				wtte := ui.CreateTextElement(wtt, ele.ViewPort)
				wtte.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
				wtte.Text.Hide()
				ct := ui.CreateContainer(entry, dialog, ele.ViewPort)
				ct.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, ele.ViewPort, func(hvc *data.HoverClick) {
					if dialog.Open && dialog.Active {
						click := hvc.Input.Get("click")
						if hvc.Hover && click.JustPressed() {
							data.SelectedWorldIndex = index
							if !data.CustomWorldSelected {
								data.SelectedPrimaryColor = pixel.ToRGBA(constants.WorldPrimary[index])
								data.SelectedSecondaryColor = pixel.ToRGBA(constants.WorldSecondary[index])
								data.SelectedDoodadColor = pixel.ToRGBA(constants.WorldDoodad[index])
							}
							for _, de := range dialog.Elements {
								if de.ElementType == ui.ContainerElement {
									if de.Key == "world_container_selected" {
										for _, ce := range de.Elements {
											if ce.ElementType == ui.SpriteElement {
												switch ce.Key {
												case "turf_tile":
													ce.Sprite.Key = constants.WorldSprites[data.SelectedWorldIndex]
												case "doodad_tile":
													ce.Sprite.Key = constants.WorldDoodads[data.SelectedWorldIndex]
												}
											} else if ce.ElementType == ui.TextElement {
												ce.Text.SetText(constants.WorldNames[data.SelectedWorldIndex])
											}
										}
										pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
										sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
										dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
										de.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
										de.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
										de.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
										de.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
										de.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
										de.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
										de.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
										de.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
										de.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
									}
								}
							}
							for _, ie := range ele.Elements {
								//if ctI, okC := ie.(*data.Container); okC {
								//	for _, cie := range ctI.Elements {
								//		if it, okIT := cie.(*data.Text); okIT {
								//			it.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
								//		}
								//	}
								//}
								if ie.ElementType == ui.TextElement {
									ie.Text.Hide()
								}
							}
							//for _, ce := range ct.Elements {
							//	if it, okIT := ce.(*data.Text); okIT {
							//		it.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
							//	}
							//}
							wtte.Text.Show()
							click.Consume()
						}
					}
				}))
				ele.Elements = append(ele.Elements, ct)
				ele.Elements = append(ele.Elements, wtte)
			}
			ui.UpdateScrollBounds(ele)
			ui.MoveToScrollTop(ele)
		case "world_container_selected":
			tti := ui.ElementConstructor{
				Key:         load.TurfTileItem.Key,
				SprKey:      constants.WorldSprites[0],
				Batch:       load.TurfTileItem.Batch,
				Position:    load.TurfTileItem.Position.Add(pixel.V(-world.HalfSize, 0)),
				ElementType: load.TurfTileItem.ElementType,
			}
			s1 := ui.CreateSpriteElement(tti)
			ele.Elements = append(ele.Elements, s1)
			lti := ui.ElementConstructor{
				Key:         load.LadderTileItem.Key,
				SprKey:      constants.TileLadderMiddle,
				Batch:       load.LadderTileItem.Batch,
				Position:    load.LadderTileItem.Position.Add(pixel.V(-world.HalfSize, 0)),
				ElementType: load.LadderTileItem.ElementType,
			}
			s2 := ui.CreateSpriteElement(lti)
			ele.Elements = append(ele.Elements, s2)
			dti := ui.ElementConstructor{
				Key:         load.DoodadTileItem.Key,
				SprKey:      constants.WorldDoodads[0],
				Batch:       load.DoodadTileItem.Batch,
				Position:    load.DoodadTileItem.Position.Add(pixel.V(-world.HalfSize, 0)),
				ElementType: load.DoodadTileItem.ElementType,
			}
			s3 := ui.CreateSpriteElement(dti)
			ele.Elements = append(ele.Elements, s3)
			wti := ui.ElementConstructor{
				Key:         load.WorldTxtItem.Key,
				Text:        constants.WorldNames[0],
				Position:    load.WorldTxtItem.Position.Add(pixel.V(-world.HalfSize, 0)),
				ElementType: load.WorldTxtItem.ElementType,
			}
			t1 := ui.CreateTextElement(wti, ele.ViewPort)
			ele.Elements = append(ele.Elements, t1)
		default:
			if ele.ElementType == ui.ButtonElement {
				if ele.OnClick == nil && ele.OnHold == nil {
					ele.OnClick = Test(fmt.Sprintf("pressed button %s", ele.Key))
				}
			} else if ele.ElementType == ui.CheckboxElement {
				if strings.Contains(ele.Key, "check_primary") ||
					strings.Contains(ele.Key, "check_secondary") ||
					strings.Contains(ele.Key, "check_doodad") {
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if dialog.Open && dialog.Active && !dialog.Lock && !dialog.Click {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() && !ele.Checked {
								ui.SetChecked(ele, true)
								changeSelectedColor(ele.Key)
								if strings.Contains(ele.Key, "check_primary") {
									worldDialogCustomShadersPrimary()
								} else if strings.Contains(ele.Key, "check_secondary") {
									worldDialogCustomShadersSecondary()
								} else if strings.Contains(ele.Key, "check_doodad") {
									worldDialogCustomShadersDoodad()
								}
								for _, ele2 := range dialog.Elements {
									if ele2.ElementType == ui.CheckboxElement {
										if ((strings.Contains(ele2.Key, "check_primary") && strings.Contains(ele.Key, "check_primary")) ||
											(strings.Contains(ele2.Key, "check_secondary") && strings.Contains(ele.Key, "check_secondary")) ||
											(strings.Contains(ele2.Key, "check_doodad") && strings.Contains(ele.Key, "check_doodad"))) &&
											ele2.Key != ele.Key {
											ui.SetChecked(ele2, false)
										}
									}
								}
							}
						}
					}))
				}
			}
		}
	}
	worldDialogShaders()
}

func OpenChangeWorldDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.NewDialog(load.WorldDialogConstructor)
		CustomizeWorldDialog()
		changeWorld := ui.Dialogs[constants.DialogChangeWorld]
		UpdateDialogView(changeWorld)
		// check if this is a custom world
		data.CustomWorldSelected = data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber == constants.WorldCustom
		data.CustomSelectedBefore = data.CustomWorldSelected
		data.SelectedPrimaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor
		data.SelectedSecondaryColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor
		data.SelectedDoodadColor = data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor
		for _, ele := range changeWorld.Elements {
			if strings.Contains(ele.Key, "color_primary") ||
				strings.Contains(ele.Key, "color_secondary") ||
				strings.Contains(ele.Key, "color_doodad") ||
				strings.Contains(ele.Key, "check_primary") ||
				strings.Contains(ele.Key, "check_secondary") ||
				strings.Contains(ele.Key, "check_doodad") ||
				strings.Contains(ele.Key, "primary_text") ||
				strings.Contains(ele.Key, "secondary_text") ||
				strings.Contains(ele.Key, "doodad_text") {
				ele.Object.Hidden = !data.CustomWorldSelected
			}
			switch ele.ElementType {
			case ui.CheckboxElement:
				switch ele.Key {
				case "custom_world_check": // whether Custom World is checked
					ui.SetChecked(ele, data.CustomWorldSelected)
				default:
					updateColorCheckboxWorld(ele)
				}
			case ui.ScrollElement: // world list
				for ctI, ele2 := range ele.Elements {
					if ele2.ElementType == ui.TextElement {
						ele2.Text.Obj.Hidden = data.SelectedWorldIndex != ctI/2
					} else if ele2.ElementType == ui.ContainerElement && data.SelectedWorldIndex == ctI/2 {
						diff := ele2.Object.Rect.H()*0.5 + 1
						ui.MoveScrollToInclude(ele, ele2.Object.Pos.Y+diff, ele2.Object.Pos.Y-diff)
					}
				}
			case ui.ContainerElement: // selected world
				if ele.Key == "world_container_selected" {
					for _, ce := range ele.Elements {
						switch ce.Key {
						case "turf_tile":
							ce.Sprite.Key = constants.WorldSprites[data.SelectedWorldIndex]
						case "doodad_tile":
							ce.Sprite.Key = constants.WorldDoodads[data.SelectedWorldIndex]
						case "world_text":
							ce.Text.SetText(constants.WorldNames[data.SelectedWorldIndex])
						}
					}
					pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
					sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
					dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
					ele.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
					ele.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
					ele.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
					ele.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
					ele.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
					ele.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
					ele.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
				}
			}
		}
		if data.CustomWorldSelected {
			worldDialogCustomShaders()
		} else {
			worldDialogNormalShaders()
		}
		ui.OpenDialogInStack(constants.DialogChangeWorld)
	}
}

func ConfirmChangeWorld() {
	if data.Editor != nil && data.CurrPuzzleSet.CurrPuzzle != nil {
		//changeWorld := data.Dialogs[constants.DialogChangeWorld]
		if data.CustomWorldSelected {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = constants.WorldCustom
		} else {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber = data.SelectedWorldIndex
			data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack = constants.WorldMusic[data.SelectedWorldIndex]
		}
		data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor = data.SelectedPrimaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor = data.SelectedSecondaryColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor = data.SelectedDoodadColor
		data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[data.SelectedWorldIndex]
	}
	UpdateEditorShaders()
	UpdatePuzzleShaders()
	data.CurrPuzzleSet.CurrPuzzle.Update = true
	data.CurrPuzzleSet.CurrPuzzle.Changed = true
	ui.Dispose(constants.DialogChangeWorld)
}

// change world shaders

func worldDialogShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
					e2.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
			e1.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
		}
	}
}

func worldDialogNormalShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			i := 0
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					pc := pixel.ToRGBA(constants.WorldPrimary[i])
					sc := pixel.ToRGBA(constants.WorldSecondary[i])
					dc := pixel.ToRGBA(constants.WorldDoodad[i])
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
					i++
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			pc := pixel.ToRGBA(constants.WorldPrimary[data.SelectedWorldIndex])
			sc := pixel.ToRGBA(constants.WorldSecondary[data.SelectedWorldIndex])
			dc := pixel.ToRGBA(constants.WorldDoodad[data.SelectedWorldIndex])
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(pc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(pc.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(pc.B))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(sc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(sc.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(sc.B))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(dc.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(dc.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(dc.B))
		}
	}
}

func worldDialogCustomShaders() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
		}
	}
}

func worldDialogCustomShadersPrimary() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.SelectedPrimaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.SelectedPrimaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.SelectedPrimaryColor.B))
		}
	}
}

func worldDialogCustomShadersSecondary() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.SelectedSecondaryColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.SelectedSecondaryColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.SelectedSecondaryColor.B))
		}
	}
}

func worldDialogCustomShadersDoodad() {
	changeWorld := ui.Dialogs[constants.DialogChangeWorld]
	for _, e1 := range changeWorld.Elements {
		if e1.ElementType == ui.ScrollElement {
			for _, e2 := range e1.Elements {
				if e2.ElementType == ui.ContainerElement {
					e2.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
					e2.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
					e2.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
				}
			}
		} else if e1.ElementType == ui.ContainerElement {
			e1.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.SelectedDoodadColor.R))
			e1.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.SelectedDoodadColor.G))
			e1.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.SelectedDoodadColor.B))
		}
	}
}

func changeSelectedColor(key string) {
	switch key {
	case "red_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_primary":
		data.SelectedPrimaryColor = pixel.ToRGBA(constants.ColorBurnt)

	case "red_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_secondary":
		data.SelectedSecondaryColor = pixel.ToRGBA(constants.ColorBurnt)

	case "red_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_doodad":
		data.SelectedDoodadColor = pixel.ToRGBA(constants.ColorBurnt)

	case "white_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorWhite)
	case "red_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorPink)
	case "black_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorBlack)
	case "yellow_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_color":
		data.SelectedTextColor = pixel.ToRGBA(constants.ColorBurnt)

	case "red_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorRed)
	case "orange_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorOrange)
	case "green_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorGreen)
	case "cyan_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorCyan)
	case "blue_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorBlue)
	case "purple_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorPurple)
	case "pink_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorPink)
	case "yellow_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorYellow)
	case "gold_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorGold)
	case "brown_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorBrown)
	case "tan_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorTan)
	case "light_gray_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorLightGray)
	case "gray_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorGray)
	case "burnt_check_shadow":
		data.SelectedShadowColor = pixel.ToRGBA(constants.ColorBurnt)
	}
}

func updateColorCheckboxWorld(x *ui.Element) {
	switch x.Key {
	case "red_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_primary":
		ui.SetChecked(x, data.SelectedPrimaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_secondary":
		ui.SetChecked(x, data.SelectedSecondaryColor == pixel.ToRGBA(constants.ColorBurnt))
	case "red_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorRed))
	case "orange_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorOrange))
	case "green_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGreen))
	case "cyan_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorCyan))
	case "blue_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBlue))
	case "purple_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPurple))
	case "pink_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorPink))
	case "yellow_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorYellow))
	case "gold_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGold))
	case "brown_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBrown))
	case "tan_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorTan))
	case "light_gray_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorLightGray))
	case "gray_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorGray))
	case "burnt_check_doodad":
		ui.SetChecked(x, data.SelectedDoodadColor == pixel.ToRGBA(constants.ColorBurnt))
	}
}
