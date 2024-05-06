package load

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
	"strings"
)

func Dialogs(win *pixelgl.Window) {
	data.NewDialog(openPuzzleConstructor)
	data.NewDialog(changeNameConstructor)
	data.NewDialog(noPlayersInPuzzleConstructor)
	data.NewDialog(areYouSureDeleteConstructor)
	data.NewDialog(unableToSaveConstructor)
	data.NewDialog(unableToSaveConfirmConstructor)
	data.NewDialog(worldDialogConstructor)
	data.NewDialog(crackedTileOptionsConstructor)
	data.NewDialog(bombOptionsConstructor)
	editorPanels()
	data.NewDialog(editorOptBottomConstructor)
	data.NewDialog(editorOptRightConstructor)
	customizeDialogs(win)
	worldDialogShaders()
}

func customizeDialogs(win *pixelgl.Window) {
	for key := range data.Dialogs {
		dialog := data.Dialogs[key]
		b := 0
		for _, e := range dialog.Elements {
			if btn, okB := e.(*data.Button); okB {
				switch btn.Key {
				case "open_puzzle":
					btn.OnClick = OpenPuzzle
				case "new_btn":
					btn.OnClick = NewPuzzle
				case "open_btn":
					btn.OnClick = OpenOpenPuzzleDialog
				case "quit_btn":
					btn.OnClick = QuitEditor(win)
				case "save_btn":
					btn.OnClick = SavePuzzleSet
				case "world_btn":
					btn.OnClick = OpenChangeWorldDialog
				case "name_btn":
					btn.OnClick = OpenDialog(constants.DialogChangeName)
				case "test_btn":
					btn.OnClick = TestPuzzle
				case "check_puzzle_name":
					btn.OnClick = ChangeName
				case "check_cracked_tile":
					btn.OnClick = ConfirmCrackTileOptions
				case "confirm_bomb_options":
					btn.OnClick = ConfirmBombOptions
				case "bomb_regenerate_delay_minus":
					btn.OnClick = func() {
						IncOrDecBombInput(false)
					}
				case "bomb_regenerate_delay_plus":
					btn.OnClick = func() {
						IncOrDecBombInput(true)
					}
				case "check_no_players":
					btn.OnClick = CloseDialog(dialog.Key)
				case "check_change_world":
					btn.OnClick = ConfirmChangeWorld
				case "confirm_unable_to_save":
					btn.OnClick = CloseDialog(dialog.Key)
				case "add_btn":
					btn.OnClick = AddPuzzle
				case "prev_btn":
					btn.OnClick = PrevPuzzle
				case "next_btn":
					btn.OnClick = NextPuzzle
				case "delete_btn":
					btn.OnClick = DeletePuzzle
				case "confirm_delete":
					btn.OnClick = ConfirmDelete
				default:
					switch dialog.Key {
					case constants.DialogEditorPanelTop, constants.DialogEditorPanelLeft:
						btn.OnClick = EditorMode(data.ModeFromSprString(btn.Sprite.Key), btn, dialog)
					case constants.DialogUnableToSaveConfirm:
						if strings.Contains(btn.Key, "cancel") {
							btn.OnClick = func() {
								data.SetOnClick(constants.DialogUnableToSaveConfirm, "confirm_unable_to_save", CloseDialog(dialog.Key))
								data.CloseDialog(dialog.Key)
							}
						}
					default:
						if strings.Contains(btn.Key, "cancel") {
							btn.OnClick = CloseDialog(dialog.Key)
						} else if btn.OnClick == nil && btn.OnHeld == nil {
							btn.OnClick = Test(fmt.Sprintf("pressed button %s", btn.Key))
						}
					}
				}
			} else if spr, okS := e.(*data.SprElement); okS {
				switch spr.Key {
				case "block_select":
					beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
					if dialog.Key == "editor_panel_top" {
						beBG = nil
					}
					beFG := img.NewSprite(data.Block(data.BlockTurf).String(), constants.TileBatch)
					beEx := img.NewSprite("", constants.TileBatch)
					spr.Entity.AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG, beEx})
					spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
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
					if b < data.BlockEmpty {
						bId := data.Block(b)
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
						obj := spr.Object
						spr.Entity.AddComponent(myecs.Drawable, sprs)
						spr.Entity.AddComponent(myecs.Block, bId)
						spr.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if data.Editor != nil && dialog.Open && !data.DialogStackOpen {
								sprS.Key = bId.String()
								click := hvc.Input.Get("click")
								if hvc.Hover && data.Editor.SelectVis {
									wo := dialog.Elements[len(dialog.Elements)-1]
									if outline, ok := wo.(*data.SprElement); ok {
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
			} else if x, okX := e.(*data.Checkbox); okX {
				switch x.Key {
				case "custom_world_check":
					x.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
						if dialog.Open && dialog.Active && !dialog.Lock && !dialog.Click {
							click := hvc.Input.Get("click")
							if hvc.Hover && click.JustPressed() {
								data.SetChecked(x, !x.Checked)
								data.CustomWorldSelected = x.Checked
								for _, ele := range dialog.Elements {
									if txt, okT := ele.(*data.Text); okT {
										if o, okO := txt.Entity.GetComponentData(myecs.Object); okO {
											if obj, okO1 := o.(*object.Object); okO1 {
												switch txt.Key {
												case "current_world":
													obj.Hidden = x.Checked
												case "primary_text", "secondary_text", "doodad_text":
													obj.Hidden = !x.Checked
												}
											}
										}
									} else if x1, okX1 := ele.(*data.Checkbox); okX1 {
										if o, okO := x1.Entity.GetComponentData(myecs.Object); okO {
											if obj, okO1 := o.(*object.Object); okO1 {
												if strings.Contains(x1.Key, "check_primary") ||
													strings.Contains(x1.Key, "check_secondary") ||
													strings.Contains(x1.Key, "check_doodad") {
													obj.Hidden = !x.Checked
												}
											}
										}
									} else if str1, okS1 := ele.(*data.SprElement); okS1 {
										if o, okO := str1.Entity.GetComponentData(myecs.Object); okO {
											if obj, okO1 := o.(*object.Object); okO1 {
												if strings.Contains(str1.Key, "color_primary") ||
													strings.Contains(str1.Key, "color_secondary") ||
													strings.Contains(str1.Key, "color_doodad") {
													obj.Hidden = !x.Checked
												}
											}
										}
									}
								}
								if x.Checked {
									for _, ele := range dialog.Elements {
										if x2, ok := ele.(*data.Checkbox); ok {
											if !data.CustomSelectedBefore {
												updateColorCheckbox(x2)
											} else if x2.Checked {
												changeSelectedColor(x2.Key)
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
					if strings.Contains(x.Key, "check_primary") ||
						strings.Contains(x.Key, "check_secondary") ||
						strings.Contains(x.Key, "check_doodad") {
						x.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, dialog.ViewPort, func(hvc *data.HoverClick) {
							if dialog.Open && dialog.Active && !dialog.Lock && !dialog.Click {
								click := hvc.Input.Get("click")
								if hvc.Hover && click.JustPressed() && !x.Checked {
									data.SetChecked(x, true)
									changeSelectedColor(x.Key)
									if strings.Contains(x.Key, "check_primary") {
										worldDialogCustomShadersPrimary()
									} else if strings.Contains(x.Key, "check_secondary") {
										worldDialogCustomShadersSecondary()
									} else if strings.Contains(x.Key, "check_doodad") {
										worldDialogCustomShadersDoodad()
									}
									for _, ele := range dialog.Elements {
										if x1, okX1 := ele.(*data.Checkbox); okX1 {
											if ((strings.Contains(x1.Key, "check_primary") && strings.Contains(x.Key, "check_primary")) ||
												(strings.Contains(x1.Key, "check_secondary") && strings.Contains(x.Key, "check_secondary")) ||
												(strings.Contains(x1.Key, "check_doodad") && strings.Contains(x.Key, "check_doodad"))) &&
												x1.Key != x.Key {
												data.SetChecked(x1, false)
											}
										}
									}
								}
							}
						}))
					}
				}
			} else if scroll, okSc := e.(*data.Scroll); okSc {
				switch scroll.Key {
				case "world_list":
					for i := 0; i < constants.WorldCustom; i++ {
						index := i
						entry := data.ElementConstructor{
							Key:      fmt.Sprintf(worldListEntry.Key, i),
							Width:    worldListEntry.Width,
							Height:   worldListEntry.Height,
							HelpText: fmt.Sprintf(worldListEntry.HelpText, constants.WorldNames[i]),
							Position: pixel.V(0, float64(i)*-18+7),
							Element:  worldListEntry.Element,
						}
						tti := data.ElementConstructor{
							Key:      turfTileItem.Key,
							SprKey:   constants.WorldSprites[i],
							Position: turfTileItem.Position,
							Element:  turfTileItem.Element,
						}
						entry.SubElements = append(entry.SubElements, tti)
						lti := data.ElementConstructor{
							Key:      ladderTileItem.Key,
							SprKey:   constants.TileLadderMiddle,
							Position: ladderTileItem.Position,
							Element:  ladderTileItem.Element,
						}
						entry.SubElements = append(entry.SubElements, lti)
						dti := data.ElementConstructor{
							Key:      doodadTileItem.Key,
							SprKey:   constants.WorldDoodads[i],
							Position: doodadTileItem.Position,
							Element:  doodadTileItem.Element,
						}
						entry.SubElements = append(entry.SubElements, dti)
						wti := data.ElementConstructor{
							Key:      worldTxtItem.Key,
							Text:     constants.WorldNames[i],
							Position: worldTxtItem.Position,
							Element:  worldTxtItem.Element,
						}
						entry.SubElements = append(entry.SubElements, wti)
						ct := data.CreateContainer(entry, dialog, scroll.ViewPort)
						ct.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, scroll.ViewPort, func(hvc *data.HoverClick) {
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
										if it, okIT := de.(*data.Text); okIT {
											if it.Key == "current_world" {
												it.Text.SetText(fmt.Sprintf("World - %s", constants.WorldNames[index]))
											}
										}
									}
									for _, ie := range scroll.Elements {
										if ctI, okC := ie.(*data.Container); okC {
											for _, cie := range ctI.Elements {
												if it, okIT := cie.(*data.Text); okIT {
													it.Text.SetColor(pixel.ToRGBA(constants.ColorWhite))
												}
											}
										}
									}
									for _, ce := range ct.Elements {
										if it, okIT := ce.(*data.Text); okIT {
											it.Text.SetColor(pixel.ToRGBA(constants.ColorBlue))
										}
									}
									click.Consume()
								}
							}
						}))
						scroll.Elements = append(scroll.Elements, ct)
					}
					data.UpdateScrollBounds(scroll)
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
		}
	}
}

func editorPanels() {
	data.NewDialog(editorPanelTopConstructor)
	data.NewDialog(editorPanelLeftConstructor)
	blockSelectConstructor := &data.DialogConstructor{
		Key:      constants.DialogEditorBlockSelect,
		Width:    constants.BlockSelectWidth,
		Height:   constants.BlockSelectHeight,
		Pos:      pixel.V(0, 0),
		NoBorder: true,
	}
	w := int(blockSelectConstructor.Width)
	h := int(blockSelectConstructor.Height)
	size := w * h
	for i, b := range data.BlockList {
		if i > size {
			break
		}
		if b == data.BlockEmpty {
			continue
		}
		blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
			Key:      "block_select_tile",
			SprKey:   "black_square",
			Position: data.BlockSelectPlacement(i, w, h),
			Element:  data.SpriteElement,
		})
	}
	blockSelectConstructor.Elements = append(blockSelectConstructor.Elements, data.ElementConstructor{
		Key:      "white_outline",
		SprKey:   "white_outline",
		Position: data.BlockSelectPlacement(0, w, h),
		Element:  data.SpriteElement,
	})
	data.NewDialog(blockSelectConstructor)
	editorPanelLeft := data.Dialogs[constants.DialogEditorPanelLeft]
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
	editorPanelTop := data.Dialogs[constants.DialogEditorPanelTop]
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
	blockSelect := data.Dialogs[constants.DialogEditorBlockSelect]
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
