package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func EditorDialogs(win *pixelgl.Window) {
	ui.NewDialog(load.OpenPuzzleConstructor)
	ui.NewDialog(load.ChangeNameConstructor)
	ui.NewDialog(load.PuzzleSettingsConstructor)
	ui.NewDialog(load.NoPlayersInPuzzleConstructor)
	ui.NewDialog(load.AreYouSureDeleteConstructor)
	ui.NewDialog(load.UnableToSaveConstructor)
	ui.NewDialog(load.UnableToSaveConfirmConstructor)
	ui.NewDialog(load.WorldDialogConstructor)
	ui.NewDialog(load.CrackedTileOptionsConstructor)
	ui.NewDialog(load.BombOptionsConstructor)
	ui.NewDialog(load.JetpackOptionsConstructor)
	editorPanels()
	ui.NewDialog(load.EditorOptBottomConstructor)
	ui.NewDialog(load.EditorOptRightConstructor)
	customizeEditorDialogs(win)
	worldDialogShaders()
}

func DisposeEditorDialogs() {
	for k, d := range ui.Dialogs {
		switch k {
		case constants.DialogOpenPuzzle,
			constants.DialogChangeName,
			constants.DialogPuzzleSettings,
			constants.DialogNoPlayersInPuzzle,
			constants.DialogAreYouSureDelete,
			constants.DialogUnableToSave,
			constants.DialogUnableToSaveConfirm,
			constants.DialogChangeWorld,
			constants.DialogCrackedTiles,
			constants.DialogBomb,
			constants.DialogJetpack,
			constants.DialogEditorPanelLeft,
			constants.DialogEditorPanelTop,
			constants.DialogEditorOptionsRight,
			constants.DialogEditorOptionsBot,
			constants.DialogEditorBlockSelect:
			ui.Dispose(d)
		}
	}
}

func customizeEditorDialogs(win *pixelgl.Window) {
	for key := range ui.Dialogs {
		dialog := ui.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			ele := e
			if ele.ElementType == ui.ButtonElement {
				switch ele.Key {
				case "open_puzzle":
					ele.OnClick = OnOpenPuzzle
				case "new_btn":
					ele.OnClick = NewPuzzle
				case "open_btn":
					ele.OnClick = OpenOpenPuzzleDialog
				case "exit_editor_btn":
					ele.OnClick = ExitEditor()
				case "save_btn":
					ele.OnClick = OnSavePuzzleSet
				case "world_btn":
					ele.OnClick = OpenChangeWorldDialog
				case "name_btn":
					ele.OnClick = OpenDialog(constants.DialogChangeName)
				case "test_btn":
					ele.OnClick = TestPuzzle
				case "check_puzzle_name":
					ele.OnClick = ChangeName
				case "puzzle_settings_btn":
					ele.OnClick = OpenPuzzleSettingsDialog
				case "confirm_puzzle_settings":
					ele.OnClick = ConfirmPuzzleSettings
				case "check_cracked_tile":
					ele.OnClick = ConfirmCrackTileOptions
				case "confirm_bomb_options":
					ele.OnClick = ConfirmBombOptions
				case "bomb_regenerate_delay_minus":
					ele.OnClick = func() {
						IncOrDecBombRegen(false)
					}
				case "bomb_regenerate_delay_plus":
					ele.OnClick = func() {
						IncOrDecBombRegen(true)
					}
				case "confirm_jetpack_options":
					ele.OnClick = ConfirmJetpackOptions
				case "jetpack_regenerate_delay_minus":
					ele.OnClick = func() {
						IncOrDecJetpackRegen(false)
					}
				case "jetpack_regenerate_delay_plus":
					ele.OnClick = func() {
						IncOrDecJetpackRegen(true)
					}
				case "jetpack_timer_minus":
					ele.OnClick = func() {
						IncOrDecJetpackTimer(false)
					}
				case "jetpack_timer_plus":
					ele.OnClick = func() {
						IncOrDecJetpackTimer(true)
					}
				case "check_no_players":
					ele.OnClick = CloseDialog(dialog.Key)
				case "check_change_world":
					ele.OnClick = ConfirmChangeWorld
				case "confirm_unable_to_save":
					ele.OnClick = CloseDialog(dialog.Key)
				case "add_btn":
					ele.OnClick = AddPuzzle
				case "prev_btn":
					ele.OnClick = PrevPuzzle
				case "next_btn":
					ele.OnClick = NextPuzzle
				case "delete_btn":
					ele.OnClick = OpenConfirmDelete
				case "confirm_delete":
					ele.OnClick = ConfirmDelete
				default:
					switch dialog.Key {
					case constants.DialogEditorPanelTop, constants.DialogEditorPanelLeft:
						ele.OnClick = EditorMode(data.ModeFromSprString(ele.Sprite.Key), ele, dialog)
					case constants.DialogUnableToSaveConfirm:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = func() {
								ui.SetOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", CloseDialog(dialog.Key))
								ui.CloseDialog(dialog.Key)
							}
						}
					default:
						if strings.Contains(ele.Key, "cancel") {
							ele.OnClick = CloseDialog(dialog.Key)
						} else if ele.OnClick == nil && ele.OnHold == nil {
							ele.OnClick = Test(fmt.Sprintf("pressed button %s", ele.Key))
						}
					}
				}
			} else if ele.ElementType == ui.SpriteElement {
				switch ele.Key {
				case "block_select":
					beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
					if dialog.Key == "editor_panel_top" {
						beBG = nil
					}
					beFG := img.NewSprite(data.Block(data.BlockTurf).String(), constants.TileBatch)
					beEx := img.NewSprite("", constants.TileBatch)
					ele.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
					ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
							beFG.Key = data.Editor.CurrBlock.String()
							switch data.Editor.CurrBlock {
							case data.BlockFall:
								beEx.Key = constants.TileFall
								beEx.Offset.Y = 0
							case data.BlockPhase:
								beEx.Key = constants.TilePhase
								beEx.Offset.Y = 0
							case data.BlockCracked:
								beEx.Key = constants.TileCracked
								beEx.Offset.Y = 0
							default:
								beEx.Key = ""
								beEx.Offset.Y = 0
							}
							data.Editor.Hover = hvc.Hover
							click := hvc.Input.Get("click")
							if hvc.Hover {
								if data.Editor.SelectVis {
									if click.JustPressed() {
										data.Editor.SelectVis = false
										data.Editor.SelectTimer = nil
										data.Editor.SelectQuick = false
										click.Consume()
									} else if click.JustReleased() {
										if data.Editor.SelectTimer != nil && !data.Editor.SelectTimer.Done() {
											data.Editor.SelectQuick = true
										}
									} else if !click.Pressed() && !data.Editor.SelectQuick {
										data.Editor.SelectVis = false
										data.Editor.SelectTimer = nil
									}
								} else {
									data.Editor.SelectQuick = false
									if click.JustPressed() {
										data.Editor.SelectVis = true
										data.Editor.SelectTimer = timing.New(0.2)
									}
								}
							}
						}
					}))
				case "block_select_tile":
					bId := data.BlockList[b]
					if bId != data.BlockEmpty {
						sprS := img.NewSprite(bId.String(), constants.TileBatch)
						sprs := []*img.Sprite{sprS}
						switch b {
						case data.BlockFall:
							sprs = append(sprs, img.NewSprite(constants.TileFall, constants.TileBatch))
						case data.BlockPhase:
							sprs = append(sprs, img.NewSprite(constants.TilePhase, constants.TileBatch))
						case data.BlockCracked:
							sprs = append(sprs, img.NewSprite(constants.TileCracked, constants.TileBatch))
						}
						obj := ele.Object
						ele.Entity.AddComponent(myecs.Drawable, sprs)
						ele.Entity.AddComponent(myecs.Block, bId)
						ele.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if data.Editor != nil && dialog.Open && !ui.DialogStackOpen {
								sprS.Key = bId.String()
								click := hvc.Input.Get("click")
								if hvc.Hover && data.Editor.SelectVis {
									outline := dialog.Elements[len(dialog.Elements)-1]
									if outline.ElementType == ui.SpriteElement {
										outline.Object.Pos = obj.Pos
									}
									if click.JustPressed() || click.JustReleased() {
										data.Editor.CurrBlock = bId
										data.Editor.SelectVis = false
										data.Editor.SelectQuick = false
										data.Editor.SelectTimer = nil
										switch data.Editor.Mode {
										case data.Brush, data.Line, data.Square, data.Fill:
										default:
											data.Editor.Mode = data.Brush
											data.CurrPuzzleSet.CurrPuzzle.Update = true
										}
										click.Consume()
									}
								}
							}
						}))
					}
					b++
				}
			} else if ele.ElementType == ui.CheckboxElement {
				switch ele.Key {
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
												updateColorCheckbox(ele2)
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
				default:
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
			} else if ele.ElementType == ui.ScrollElement {
				switch ele.Key {
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
						}
						entry.SubElements = append(entry.SubElements, wti)
						wtt := ui.ElementConstructor{
							Key:         load.WorldTxtItem.Key,
							Text:        constants.WorldNames[i],
							Position:    load.WorldTxtItem.Position.Add(pixel.V(0, y+1)),
							ElementType: load.WorldTxtItem.ElementType,
						}
						wtte := ui.CreateTextElement(wtt, ele.ViewPort)
						wtte.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
						wtte.Text.Hidden = true
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
											ie.Text.Hidden = true
										}
									}
									//for _, ce := range ct.Elements {
									//	if it, okIT := ce.(*data.Text); okIT {
									//		it.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
									//	}
									//}
									wtte.Text.Hidden = false
									click.Consume()
								}
							}
						}))
						ele.Elements = append(ele.Elements, ct)
						ele.Elements = append(ele.Elements, wtte)
					}
					ui.UpdateScrollBounds(ele)
				}
			} else if ele.ElementType == ui.ContainerElement {
				if ele.Key == "world_container_selected" {
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
				}
			}
		}
		switch dialog.Key {
		case constants.DialogOpenPuzzle:
			dialog.OnOpen = OnOpenPuzzleDialog
		case constants.DialogChangeName:
			dialog.OnOpen = OnChangeNameDialog
		case constants.DialogCrackedTiles:
			dialog.OnOpen = OnOpenCrackTileOptions
		case constants.DialogBomb:
			dialog.OnOpen = OnOpenBombOptions
		case constants.DialogJetpack:
			dialog.OnOpen = OnOpenJetpackOptions
		}
	}
}

func editorPanels() {
	ui.NewDialog(load.EditorPanelTopConstructor)
	ui.NewDialog(load.EditorPanelLeftConstructor)
	blockSelectConstructor := &ui.DialogConstructor{
		Key:      constants.DialogEditorBlockSelect,
		Width:    constants.BlockSelectWidth,
		Height:   constants.BlockSelectHeight,
		Pos:      pixel.V(0, 0),
		NoBorder: true,
	}
	w := int(blockSelectConstructor.Width)
	h := int(blockSelectConstructor.Height)
	size := w * h
	for i := range data.BlockList {
		if i > size {
			break
		}
		blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, ui.ElementConstructor{
			Key:         "block_select_tile",
			SprKey:      "black_square",
			Batch:       constants.UIBatch,
			Position:    data.BlockSelectPlacement(i, w, h),
			ElementType: ui.SpriteElement,
		})
	}
	blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, ui.ElementConstructor{
		Key:         "white_outline",
		SprKey:      "white_outline",
		Batch:       constants.UIBatch,
		Position:    data.BlockSelectPlacement(0, w, h),
		ElementType: ui.SpriteElement,
	})
	ui.NewDialog(blockSelectConstructor)
	editorPanelLeft := ui.Dialogs[constants.DialogEditorPanelLeft]
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	editorPanelLeft.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
	editorPanelTop := ui.Dialogs[constants.DialogEditorPanelTop]
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	editorPanelTop.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
	blockSelect := ui.Dialogs[constants.DialogEditorBlockSelect]
	blockSelect.ViewPort.Canvas.SetUniform("uRedPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBluePrimary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uRedDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(0))
	blockSelect.ViewPort.Canvas.SetFragmentShader(data.PuzzleShader)
}
