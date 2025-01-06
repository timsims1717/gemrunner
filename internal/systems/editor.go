package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/util"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
)

func EditorInit() {
	// initialize imDraw
	data.IMDraw = imdraw.New(nil)

	// initialize editor panel
	data.NewEditor()
	//data.Editor.PosTop = true
	data.Editor.BlockSelect = ui.Dialogs[constants.DialogEditorBlockSelect].ViewPort

	// open editor dialogs
	if data.Editor.PosTop {
		ui.OpenDialog(constants.DialogEditorPanelTop)
		ui.OpenDialog(constants.DialogEditorOptionsBot)
	} else {
		ui.OpenDialog(constants.DialogEditorPanelLeft)
		ui.OpenDialog(constants.DialogEditorOptionsRight)
	}
}

func UpdateEditorShaders() {
	if data.Editor == nil {
		return
	}
	// set editor panel shader uniforms
	editorPanelLeft := ui.Dialogs[constants.DialogEditorPanelLeft]
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.R))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.G))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.B))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.R))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.G))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.B))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.R))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.G))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.B))
	editorPanelTop := ui.Dialogs[constants.DialogEditorPanelTop]
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.R))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.G))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.B))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.R))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.G))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.B))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.R))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.G))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.B))
	// set editor select shader uniforms
	data.Editor.BlockSelect.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.R))
	data.Editor.BlockSelect.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.G))
	data.Editor.BlockSelect.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor.B))
	data.Editor.BlockSelect.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.R))
	data.Editor.BlockSelect.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.G))
	data.Editor.BlockSelect.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor.B))
	data.Editor.BlockSelect.Canvas.SetUniform("uRedDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.R))
	data.Editor.BlockSelect.Canvas.SetUniform("uGreenDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.G))
	data.Editor.BlockSelect.Canvas.SetUniform("uBlueDoodad", float32(data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor.B))
}

func ChangeWorldToNext() {
	if data.CurrPuzzleSet != nil {
		data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber++
		if data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber >= constants.WorldCustom {
			data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber %= constants.WorldCustom
		}
		ChangeWorldTo(data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldNumber)
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.Changed = true
	}
}

func ChangeWorldTo(world int) {
	data.CurrPuzzleSet.CurrPuzzle.Metadata.WorldSprite = constants.WorldSprites[world]
	data.CurrPuzzleSet.CurrPuzzle.Metadata.PrimaryColor = pixel.ToRGBA(constants.WorldPrimary[world])
	data.CurrPuzzleSet.CurrPuzzle.Metadata.SecondaryColor = pixel.ToRGBA(constants.WorldSecondary[world])
	data.CurrPuzzleSet.CurrPuzzle.Metadata.DoodadColor = pixel.ToRGBA(constants.WorldDoodad[world])
	data.CurrPuzzleSet.CurrPuzzle.Metadata.MusicTrack = constants.WorldMusic[world]
	UpdateEditorShaders()
	UpdatePuzzleShaders()
}

func ChangeWorldShader(shaderMode int) {
	data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderMode = shaderMode
	data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderSpeed = constants.ShaderSpeeds[shaderMode]
	data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderX = constants.ShaderXs[shaderMode]
	data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderY = constants.ShaderYs[shaderMode]
	data.CurrPuzzleSet.CurrPuzzle.Metadata.ShaderCustom = constants.ShaderCustom[shaderMode]
	UpdateWorldShaders()
}

func UpdateEditorModeHotKey() {
	oldMode := data.Editor.Mode
	if data.MenuInput.Get("ctrl").Pressed() || data.MenuInput.Get("rCtrl").Pressed() {
		if data.MenuInput.Get("ctrlCopy").JustPressed() {
			// copy
			data.Editor.Mode = data.ModeCopy
		} else if data.MenuInput.Get("ctrlCut").JustPressed() {
			// cut
			data.Editor.Mode = data.ModeCut
		} else if data.MenuInput.Get("ctrlPaste").JustPressed() {
			// paste
			data.Editor.Mode = data.ModePaste
		} else if (data.MenuInput.Get("shift").Pressed() || data.MenuInput.Get("rShift").Pressed()) &&
			data.MenuInput.Get("ctrlShiftRedo").JustPressed() {
			// redo
			data.Editor.Mode = data.ModeRedo
		} else if data.MenuInput.Get("ctrlUndo").JustPressed() {
			// undo
			data.Editor.Mode = data.ModeUndo
		} else if data.MenuInput.Get("ctrlSave").JustPressed() {
			// save
			data.Editor.Mode = data.ModeSave
		} else if data.MenuInput.Get("ctrlOpen").JustPressed() {
			// open
			data.Editor.Mode = data.ModeOpen
		}
	} else {
		for i := 0; i < data.EndModeList; i++ {
			hotkey := data.MenuInput.Get(data.EditorMode(i).String())
			if hotkey != nil && hotkey.JustPressed() {
				hotkey.Consume()
				data.Editor.Mode = data.EditorMode(i)
			}
		}
	}
	if data.Editor.Mode >= data.EndModeList {
		data.Editor.Mode = data.ModeBrush
	}
	if oldMode != data.Editor.Mode {
		data.Editor.LastMode = oldMode
		data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		data.Editor.SelectVis = false
		data.CurrPuzzleSet.CurrPuzzle.Update = true
	}
	if data.Editor.LastMode >= data.EndModeList {
		data.Editor.LastMode = data.ModeBrush
	}
}

func PuzzleEditSystem() {
	EditorPanelButtons()
	if data.Editor.SelectTimer != nil {
		data.Editor.SelectTimer.Update()
	}
	projPos := data.PuzzleView.ProjectWorld(data.MenuInput.World)
	x, y := world.WorldToMap(projPos.X, projPos.Y)
	coords := world.Coords{X: x, Y: y}
	// special check for a weird bug
	if coords.X < -1000 || coords.X > 1000 || coords.Y < -1000 || coords.Y > 1000 {
		return
	}
	legal := CoordsLegal(coords)
	lastLegal := CoordsLegal(data.Editor.LastCoords)
	rClick := data.MenuInput.Get("rightClick")
	click := data.MenuInput.Get("click")
	inInside := data.Editor.BlockSelect.PointInside(data.Editor.BlockSelect.ProjectWorld(data.MenuInput.World))
	if data.Editor.SelectVis && legal && !inInside &&
		(click.JustPressed() || rClick.JustPressed()) {
		data.Editor.SelectVis = false
	}
	if data.Editor.SelectVis {
		data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
	} else {
		data.CurrPuzzleSet.CurrPuzzle.Click = rClick.Pressed() || click.Pressed()
		if !rClick.Pressed() && !rClick.JustReleased() &&
			!click.Pressed() && !click.JustReleased() && legal {
			CreateHighlight(coords)
			data.Editor.NoInput = false
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		} else if !data.Editor.NoInput {
		modeLabel:
			switch data.Editor.Mode {
			case data.ModeBrush:
				if rClick.JustPressed() || click.JustPressed() && !lastLegal {
					data.Editor.LastCoords = coords
					lastLegal = CoordsLegal(data.Editor.LastCoords)
				}
				line := world.Line(data.Editor.LastCoords, coords)
				if rClick.Pressed() {
					for _, c := range line {
						DeleteBlock(c)
					}
				} else if click.Pressed() {
					for _, c := range line {
						SetBlock(c, data.Editor.CurrBlock)
					}
				}
				if rClick.JustReleased() || click.JustReleased() {
					data.CurrPuzzleSet.CurrPuzzle.Update = true
					data.CurrPuzzleSet.CurrPuzzle.Changed = true
				}
				if legal {
					CreateHighlight(coords)
					data.Editor.LastCoords = coords
				}
			case data.ModeLine:
				if lastLegal {
					if click.Pressed() {
						if click.JustPressed() {
							data.Editor.LastCoords = coords
						} else if rClick.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						line := world.Line(data.Editor.LastCoords, coords)
						for _, c := range line {
							CreateHighlight(c)
						}
					} else if click.JustReleased() {
						line := world.Line(data.Editor.LastCoords, coords)
						for _, c := range line {
							SetBlock(c, data.Editor.CurrBlock)
						}
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
						if legal {
							CreateHighlight(coords)
						}
					} else if rClick.Pressed() {
						if rClick.JustPressed() {
							data.Editor.LastCoords = coords
						} else if click.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						line := world.Line(data.Editor.LastCoords, coords)
						for _, c := range line {
							CreateHighlight(c)
						}
					} else if rClick.JustReleased() {
						line := world.Line(data.Editor.LastCoords, coords)
						for _, c := range line {
							DeleteBlock(c)
						}
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
						if legal {
							CreateHighlight(coords)
						}
					}
				} else {
					if legal {
						data.Editor.NoInput = false
						CreateHighlight(coords)
						data.Editor.LastCoords = coords
					}
				}
			case data.ModeSquare:
				if lastLegal {
					if click.Pressed() {
						if click.JustPressed() {
							data.Editor.LastCoords = coords
						} else if rClick.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						square := world.Square(data.Editor.LastCoords, coords)
						for _, c := range square {
							CreateHighlight(c)
						}
					} else if click.JustReleased() {
						square := world.Square(data.Editor.LastCoords, coords)
						for _, c := range square {
							SetBlock(c, data.Editor.CurrBlock)
						}
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
						if legal {
							CreateHighlight(coords)
						}
					} else if rClick.Pressed() {
						if rClick.JustPressed() {
							data.Editor.LastCoords = coords
						} else if click.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						square := world.Square(data.Editor.LastCoords, coords)
						for _, c := range square {
							CreateHighlight(c)
						}
					} else if rClick.JustReleased() {
						square := world.Square(data.Editor.LastCoords, coords)
						for _, c := range square {
							DeleteBlock(c)
						}
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
						if legal {
							CreateHighlight(coords)
						}
					}
				} else {
					if legal {
						data.Editor.NoInput = false
						CreateHighlight(coords)
						data.Editor.LastCoords = coords
					}
				}
			case data.ModeFill:
				if rClick.JustReleased() {
					Fill(coords, true)
				} else if click.JustReleased() {
					Fill(coords, false)
				}
				if legal {
					CreateHighlight(coords)
				}
			case data.ModeErase:
				if rClick.JustPressed() || click.JustPressed() || !lastLegal {
					data.Editor.LastCoords = coords
					lastLegal = CoordsLegal(data.Editor.LastCoords)
				}
				line := world.Line(data.Editor.LastCoords, coords)
				if rClick.Pressed() || click.Pressed() {
					for _, c := range line {
						DeleteBlock(c)
					}
				}
				if rClick.JustReleased() || click.JustReleased() {
					data.CurrPuzzleSet.CurrPuzzle.Update = true
					data.CurrPuzzleSet.CurrPuzzle.Changed = true
				}
				if legal {
					CreateHighlight(coords)
					data.Editor.LastCoords = coords
				}
			case data.ModeEyedrop:
				if legal {
					if rClick.JustReleased() || click.JustReleased() {
						tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.T[coords.Y][coords.X]
						if tile.Block != data.BlockEmpty {
							data.Editor.CurrBlock = tile.Block
						}
						if data.Editor.LastMode > data.ModeFill {
							data.Editor.Mode, data.Editor.LastMode = data.ModeBrush, data.Editor.Mode
						} else {
							data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Editor.Mode
						}
					}
					CreateHighlight(coords)
				}
			case data.ModeWrench:
				if legal {
					if rClick.JustReleased() || click.JustReleased() {
						if data.CurrSelect != nil {
							inTest := coords
							inTest.X -= data.CurrSelect.Offset.X
							inTest.Y -= data.CurrSelect.Offset.Y
							if !CoordsLegalSelection(inTest) {
								PlaceSelection()
							} else {
								tile := data.CurrSelect.Tiles[inTest.Y][inTest.X]
								data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{}
								for _, row := range data.CurrSelect.Tiles {
									for _, t := range row {
										if t.Block == tile.Block {
											switch t.Block {
											case data.BlockFly:
												t.Metadata.Flipped = !t.Metadata.Flipped
												t.Metadata.Changed = true
											case data.BlockCracked, data.BlockLadderCracked,
												data.BlockBomb, data.BlockBombLit,
												data.BlockJetpack, data.BlockDisguise:
												data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = append(data.CurrPuzzleSet.CurrPuzzle.WrenchTiles, t)
											case data.BlockPhase:
												if rClick.JustReleased() {
													t.Metadata.Phase--
													if t.Metadata.Phase < 0 {
														t.Metadata.Phase = 7
													}
												} else if click.JustReleased() {
													t.Metadata.Phase++
													t.Metadata.Phase = t.Metadata.Phase % 8
												}
											}
										}
									}
								}
								switch tile.Block {
								case data.BlockFly, data.BlockPhase:
									data.CurrPuzzleSet.CurrPuzzle.Update = true
									data.CurrPuzzleSet.CurrPuzzle.Changed = true
								case data.BlockCracked, data.BlockLadderCracked:
									ui.OpenDialogInStack(constants.DialogCrackedTiles)
								case data.BlockBomb, data.BlockBombLit:
									ui.OpenDialogInStack(constants.DialogBomb)
								case data.BlockJetpack:
									ui.OpenDialogInStack(constants.DialogJetpack)
								case data.BlockDisguise:
									ui.OpenDialogInStack(constants.DialogDisguise)
								}
								break
							}
						}
						tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
						switch tile.Block {
						case data.BlockTurf:
							if tile.AltBlock == 0 {
								tile.AltBlock = 1
							} else {
								tile.AltBlock = 0
							}
						case data.BlockFly:
							tile.Metadata.Flipped = !tile.Metadata.Flipped
							tile.Object.Flip = tile.Metadata.Flipped
							tile.Metadata.Changed = true
							data.CurrPuzzleSet.CurrPuzzle.Update = true
							data.CurrPuzzleSet.CurrPuzzle.Changed = true
						case data.BlockCracked, data.BlockLadderCracked:
							data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{tile}
							ui.OpenDialogInStack(constants.DialogCrackedTiles)
						case data.BlockBomb, data.BlockBombLit:
							data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{tile}
							ui.OpenDialogInStack(constants.DialogBomb)
						case data.BlockJetpack:
							data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{tile}
							ui.OpenDialogInStack(constants.DialogJetpack)
						case data.BlockDisguise:
							data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{tile}
							ui.OpenDialogInStack(constants.DialogDisguise)
						case data.BlockPhase:
							if rClick.JustReleased() {
								tile.Metadata.Phase--
								if tile.Metadata.Phase < 0 {
									tile.Metadata.Phase = 7
								}
							} else if click.JustReleased() {
								tile.Metadata.Phase++
								tile.Metadata.Phase = tile.Metadata.Phase % 8
							}
							data.CurrPuzzleSet.CurrPuzzle.Update = true
							data.CurrPuzzleSet.CurrPuzzle.Changed = true
						}
					}
				}
			case data.ModeWire:
				if legal {
					if rClick.JustReleased() || click.JustReleased() {
						tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
						switch tile.Block {
						case data.BlockDemon, data.BlockDemonRegen,
							data.BlockFly, data.BlockFlyRegen:
							if rClick.JustReleased() {
								RemoveLinkedTiles(tile)
								tile.Metadata.Changed = true
							} else if click.JustReleased() {
								lt := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(data.Editor.LastCoords.X, data.Editor.LastCoords.Y)
								if lt != nil && lt.Block != tile.Block &&
									((lt.Block == data.BlockDemon && tile.Block == data.BlockDemonRegen) ||
										(tile.Block == data.BlockDemon && lt.Block == data.BlockDemonRegen) ||
										(lt.Block == data.BlockFly && tile.Block == data.BlockFlyRegen) ||
										(tile.Block == data.BlockFly && lt.Block == data.BlockFlyRegen)) {
									LinkTiles(tile, lt)
									//lt.Metadata.LinkedTiles = append(lt.Metadata.LinkedTiles, coords)
									//tile.Metadata.LinkedTiles = append(tile.Metadata.LinkedTiles, data.Editor.LastCoords)
									lt.Metadata.Changed = true
									tile.Metadata.Changed = true
								}
							}
						}
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
					} else if click.Pressed() {
						if click.JustPressed() && legal {
							data.Editor.LastCoords = coords
						} else if rClick.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						lt := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(data.Editor.LastCoords.X, data.Editor.LastCoords.Y)
						if lt != nil {
							switch lt.Block {
							case data.BlockDemon, data.BlockDemonRegen,
								data.BlockFly, data.BlockFlyRegen:
								data.IMDraw.Color = constants.ColorOrange
							default:
								break modeLabel
							}
							data.IMDraw.EndShape = imdraw.RoundEndShape
							if legal {
								tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
								if lt.Block != tile.Block &&
									((lt.Block == data.BlockDemon && tile.Block == data.BlockDemonRegen) ||
										(tile.Block == data.BlockDemon && lt.Block == data.BlockDemonRegen) ||
										(lt.Block == data.BlockFly && tile.Block == data.BlockFlyRegen) ||
										(tile.Block == data.BlockFly && lt.Block == data.BlockFlyRegen)) {
									data.IMDraw.Push(lt.Object.Pos, tile.Object.Pos)
									data.IMDraw.Line(2)
									break modeLabel
								}
							}
							data.IMDraw.Push(lt.Object.Pos, projPos)
							data.IMDraw.Line(2)
						}
					}
				}
			case data.ModeText:
				if legal {
					if click.JustReleased() { // open the text tool dialog
						tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)
						data.CurrPuzzleSet.CurrPuzzle.WrenchTiles = []*data.Tile{tile}
						ui.OpenDialogInStack(constants.DialogFloatingText)
						data.CurrPuzzleSet.CurrPuzzle.Update = true
						data.CurrPuzzleSet.CurrPuzzle.Changed = true
					} else if rClick.JustReleased() { // remove text
						//tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(coords.X, coords.Y)

						//data.CurrPuzzleSet.CurrPuzzle.Update = true
						//data.CurrPuzzleSet.CurrPuzzle.Changed = true
					}
				}
			case data.ModeSelect:
				if lastLegal {
					if click.Pressed() {
						if click.JustPressed() && legal {
							data.Editor.LastCoords = coords
						} else if rClick.JustPressed() {
							data.Editor.NoInput = true
							break
						}
						CreateSquareSelect(data.Editor.LastCoords, GetClosestLegal(coords))
					} else if click.JustReleased() {
						CreateSelection(data.Editor.LastCoords, GetClosestLegal(coords))
						data.Editor.Mode = data.ModeMove
					} else if !click.Pressed() && !rClick.Pressed() {
						data.Editor.NoInput = false
						if legal {
							CreateHighlight(coords)
							data.Editor.LastCoords = coords
						}
					}
				} else {
					if legal {
						data.Editor.NoInput = false
						CreateHighlight(coords)
						data.Editor.LastCoords = coords
					}
				}
			case data.ModeMove:
				if legal {
					CreateHighlight(coords)
					if click.Pressed() {
						if rClick.JustPressed() {
							data.Editor.Mode = data.ModeSelect
							data.Editor.NoInput = true
							data.CurrSelect.Offset = data.CurrSelect.Origin
							break
						} else if click.JustPressed() {
							inTest := coords
							inTest.X -= data.CurrSelect.Offset.X
							inTest.Y -= data.CurrSelect.Offset.Y
							if !CoordsLegalSelection(inTest) {
								data.Editor.Mode = data.ModeSelect
								break
							}
							data.Editor.LastCoords = coords
						} else {
							var move world.Coords
							move.X = coords.X - data.Editor.LastCoords.X
							move.Y = coords.Y - data.Editor.LastCoords.Y
							data.CurrSelect.Offset.X = data.CurrSelect.Origin.X + move.X
							data.CurrSelect.Offset.Y = data.CurrSelect.Origin.Y + move.Y
						}
						end := world.Coords{
							X: data.CurrSelect.Offset.X + data.CurrSelect.Width - 1,
							Y: data.CurrSelect.Offset.Y + data.CurrSelect.Height - 1,
						}
						CreateSquareSelect(data.CurrSelect.Offset, end)
					} else if click.JustReleased() {
						var move world.Coords
						move.X = coords.X - data.Editor.LastCoords.X
						move.Y = coords.Y - data.Editor.LastCoords.Y
						data.CurrSelect.Offset.X = data.CurrSelect.Origin.X + move.X
						data.CurrSelect.Offset.Y = data.CurrSelect.Origin.Y + move.Y
						data.CurrSelect.Origin = data.CurrSelect.Offset
					} else if rClick.JustPressed() {
						data.Editor.Mode = data.ModeBrush
						data.Editor.NoInput = true
						data.CurrSelect.Offset = data.CurrSelect.Origin
						break
					}
				}
			}
		}
		switch data.Editor.Mode {
		case data.ModeCopy:
			if CreateClip() {
				data.Editor.Mode = data.ModeMove
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
			}
		case data.ModeCut:
			if CreateClip() {
				data.CurrSelect = nil
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
		case data.ModePaste:
			// place the current clip
			PlaceSelection()
			if PlaceClip() {
				data.Editor.Mode = data.ModeMove
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
			}
		case data.ModeDelete:
			data.CurrSelect = nil
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
			data.CurrPuzzleSet.CurrPuzzle.Update = true
			data.CurrPuzzleSet.CurrPuzzle.Changed = true
		case data.ModeUndo:
			if len(data.CurrPuzzleSet.CurrPuzzle.UndoStack) > 0 {
				PushRedoStack()
				PopUndoArray()
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
		case data.ModeRedo:
			if len(data.CurrPuzzleSet.CurrPuzzle.RedoStack) > 0 {
				PushUndoArray(false)
				PopRedoStack()
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
		case data.ModeFlipVertical:
			if data.CurrSelect != nil {
				FlipVertical()
				data.Editor.Mode = data.ModeMove
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
			}
		case data.ModeFlipHorizontal:
			if data.CurrSelect != nil {
				FlipHorizontal()
				data.Editor.Mode = data.ModeMove
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
			}
		case data.ModeSave:
			PlaceSelection()
			if !SavePuzzleSet() {
				ui.OpenDialogInStack(constants.DialogUnableToSave)
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
		case data.ModeOpen:
			//if err := LoadPuzzle(); err != nil {
			//	fmt.Println("Error:", err)
			//}
			ui.OpenDialogInStack(constants.DialogOpenPuzzle)
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.ModeBrush
		}
	}
	if data.Editor.Mode == data.ModeMove || data.Editor.Mode == data.ModeWrench {
		if data.CurrSelect != nil {
			end := world.Coords{
				X: data.CurrSelect.Offset.X + data.CurrSelect.Width - 1,
				Y: data.CurrSelect.Offset.Y + data.CurrSelect.Height - 1,
			}
			CreateSquareSelect(data.CurrSelect.Offset, end)
			for _, row := range data.CurrSelect.Tiles {
				for _, tile := range row {
					c := tile.Coords
					c.X += data.CurrSelect.Offset.X
					c.Y += data.CurrSelect.Offset.Y
					obj := object.New()
					obj.Layer = 3
					obj.Pos = world.MapToWorld(c)
					obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
					obj.Flip = tile.Metadata.Flipped
					spr := GetTileSpritesSelection(tile)
					if len(spr) > 0 {
						myecs.Manager.NewEntity().
							AddComponent(myecs.Object, obj).
							AddComponent(myecs.Drawable, spr).
							AddComponent(myecs.Temp, myecs.ClearFlag(true))
					}
					// don't draw anything below
					if CoordsLegal(c) {
						data.CurrPuzzleSet.CurrPuzzle.Tiles.T[c.Y][c.X].Object.Hidden = true
					}
				}
			}
		} else if data.Editor.Mode == data.ModeMove {
			data.Editor.Mode = data.ModeSelect
		}
	} else {
		// place the current selection
		PlaceSelection()
	}
	if data.Editor.Mode == data.ModeWire && legal {
		// drawing line information
		tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.T[coords.Y][coords.X]
		switch tile.Block {
		case data.BlockDemon, data.BlockDemonRegen,
			data.BlockFly, data.BlockFlyRegen:
			data.IMDraw.Color = constants.ColorOrange
			data.IMDraw.EndShape = imdraw.RoundEndShape
			for _, rt := range tile.Metadata.LinkedTiles {
				data.IMDraw.Push(tile.Object.Pos, world.MapToWorld(rt).Add(pixel.V(world.HalfSize, world.HalfSize)))
				data.IMDraw.Line(2)
			}
		}
	}
	if data.Editor.Mode >= data.EndModeList {
		data.Editor.Mode = data.ModeBrush
	}
	if data.Editor.LastMode >= data.EndModeList {
		data.Editor.LastMode = data.ModeBrush
	}
}

func EditorPanelButtons() {
	var panel *ui.Dialog
	if data.Editor.PosTop {
		panel = ui.Dialogs[constants.DialogEditorPanelTop]
	} else {
		panel = ui.Dialogs[constants.DialogEditorPanelLeft]
	}
	if !panel.Click {
		for _, e := range panel.Elements {
			if e.ElementType == ui.ButtonElement {
				if data.ModeFromSprString(e.Sprite.Key) == data.Editor.Mode ||
					(data.ModeFromSprString(e.Sprite.Key) == data.ModeSelect &&
						data.Editor.Mode == data.ModeMove) {
					e.Entity.AddComponent(myecs.Drawable, e.Sprite2)
				} else {
					e.Entity.AddComponent(myecs.Drawable, e.Sprite)
				}
			}
		}
	}
}

func CreateHighlight(coords world.Coords) {
	if CoordsLegal(coords) {
		obj := object.New()
		obj.Layer = 4
		obj.Pos = world.MapToWorld(coords)
		obj.Pos = obj.Pos.Add(pixel.V(world.TileSize*0.5, world.TileSize*0.5))
		spr := img.NewSprite("white_checker", constants.UIBatch)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr).
			AddComponent(myecs.Temp, myecs.ClearFlag(true))
	}
}

// Fill uses a BFS algorithm to fill the space.
func Fill(a world.Coords, delete bool) {
	if CoordsLegal(a) {
		orig := data.CurrPuzzleSet.CurrPuzzle.Tiles.T[a.Y][a.X]
		if (delete && orig.Block != data.BlockEmpty) ||
			orig.Block != data.Editor.CurrBlock {
			var v []world.Coords
			var r []world.Coords
			var q []world.Coords
			v = append(v, a)
			r = append(r, a)
			q = append(q, a)
			for len(q) > 0 {
				top := q[0]
				q = q[1:]
				for _, n := range top.Neighbors() {
					if (n.X == top.X || n.Y == top.Y) && CoordsLegal(n) && !world.CoordsIn(n, v) {
						v = append(v, n)
						tile := data.CurrPuzzleSet.CurrPuzzle.Tiles.T[n.Y][n.X]
						if tile.Block == orig.Block {
							r = append(r, n)
							q = append(q, n)
						}
					}
				}
			}
			for _, c := range r {
				if delete {
					DeleteBlock(c)
				} else {
					SetBlock(c, data.Editor.CurrBlock)
				}
			}
			data.CurrPuzzleSet.CurrPuzzle.Update = true
			data.CurrPuzzleSet.CurrPuzzle.Changed = true
		}
	}
}

func CreateSquareSelect(a, b world.Coords) {
	a1, b1 := a, b
	a1.X = util.Min(a.X, b.X)
	a1.Y = util.Min(a.Y, b.Y)
	b1.X = util.Max(a.X, b.X)
	b1.Y = util.Max(a.Y, b.Y)
	b1.Y++
	b1.X++
	posA := world.MapToWorld(a1)
	posB := world.MapToWorld(b1)
	posA.X++
	posA.Y++
	posB.X--
	posB.Y--
	data.IMDraw.Color = constants.ColorWhite
	data.IMDraw.EndShape = imdraw.RoundEndShape
	data.IMDraw.Push(posA, pixel.V(posA.X, posB.Y), posB, pixel.V(posB.X, posA.Y))
	data.IMDraw.Polygon(2)
}

func CreateSelection(a, b world.Coords) {
	data.CurrSelect = new(data.Selection)
	var o world.Coords
	data.CurrSelect.Width = util.Abs(b.X-a.X) + 1
	data.CurrSelect.Height = util.Abs(b.Y-a.Y) + 1
	if a.X <= b.X {
		o.X = a.X
	} else {
		o.X = b.X
	}
	if a.Y <= b.Y {
		o.Y = a.Y
	} else {
		o.Y = b.Y
	}
	data.CurrSelect.Offset = o
	data.CurrSelect.Origin = o
	for dy := 0; dy < data.CurrSelect.Height; dy++ {
		data.CurrSelect.Tiles = append(data.CurrSelect.Tiles, []*data.Tile{})
		for dx := 0; dx < data.CurrSelect.Width; dx++ {
			old := data.CurrPuzzleSet.CurrPuzzle.Tiles.T[dy+o.Y][dx+o.X]
			tile := old.Copy()
			RemoveLinkedTiles(old)
			old.ToEmpty()
			tile.Coords = world.Coords{X: dx, Y: dy}
			data.CurrSelect.Tiles[dy] = append(data.CurrSelect.Tiles[dy], tile)
		}
	}
	data.CurrPuzzleSet.CurrPuzzle.Update = true
}

func CreateClip() bool {
	if data.CurrSelect == nil {
		return false
	}
	data.ClipSelect = new(data.Selection)
	data.ClipSelect.Width = data.CurrSelect.Width
	data.ClipSelect.Height = data.CurrSelect.Height
	data.ClipSelect.Offset = world.Coords{X: 0, Y: constants.PuzzleHeight - data.CurrSelect.Height}
	data.ClipSelect.Origin = data.CurrSelect.Offset
	for dy, row := range data.CurrSelect.Tiles {
		data.ClipSelect.Tiles = append(data.ClipSelect.Tiles, []*data.Tile{})
		for _, old := range row {
			tile := old.Copy()
			tile.Coords = old.Coords
			data.ClipSelect.Tiles[dy] = append(data.ClipSelect.Tiles[dy], tile)
		}
	}
	return true
}

func PlaceSelection() {
	if data.CurrSelect == nil {
		return
	}
	var hasP1, hasP2, hasP3, hasP4 bool
	for _, row := range data.CurrSelect.Tiles {
		for _, tile := range row {
			if tile.Block == data.BlockPlayer1 {
				hasP1 = true
			}
			if tile.Block == data.BlockPlayer2 {
				hasP2 = true
			}
			if tile.Block == data.BlockPlayer3 {
				hasP3 = true
			}
			if tile.Block == data.BlockPlayer4 {
				hasP4 = true
			}
		}
	}
	for _, row := range data.CurrPuzzleSet.CurrPuzzle.Tiles.T {
		for _, tile := range row {
			if tile.Block == data.BlockPlayer1 && hasP1 {
				tile.Block = data.BlockEmpty
			}
			if tile.Block == data.BlockPlayer2 && hasP2 {
				tile.Block = data.BlockEmpty
			}
			if tile.Block == data.BlockPlayer3 && hasP3 {
				tile.Block = data.BlockEmpty
			}
			if tile.Block == data.BlockPlayer4 && hasP4 {
				tile.Block = data.BlockEmpty
			}
		}
	}
	for _, row := range data.CurrSelect.Tiles {
		for _, tile := range row {
			c := tile.Coords
			c.X += data.CurrSelect.Offset.X
			c.Y += data.CurrSelect.Offset.Y
			if CoordsLegal(c) {
				tile.CopyInto(data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(c.X, c.Y))
				UpdateLinkedTiles(data.CurrPuzzleSet.CurrPuzzle.Tiles.Get(c.X, c.Y))
			}
		}
	}
	data.CurrSelect = nil
	data.CurrPuzzleSet.CurrPuzzle.Update = true
	data.CurrPuzzleSet.CurrPuzzle.Changed = true
}

func PlaceClip() bool {
	if data.ClipSelect == nil {
		return false
	}
	data.CurrSelect = new(data.Selection)
	data.CurrSelect.Width = data.ClipSelect.Width
	data.CurrSelect.Height = data.ClipSelect.Height
	data.CurrSelect.Offset = world.Coords{X: 0, Y: constants.PuzzleHeight - data.ClipSelect.Height}
	data.CurrSelect.Origin = data.CurrSelect.Offset
	for dy, row := range data.ClipSelect.Tiles {
		data.CurrSelect.Tiles = append(data.CurrSelect.Tiles, []*data.Tile{})
		for _, old := range row {
			tile := old.Copy()
			tile.Coords = old.Coords
			data.CurrSelect.Tiles[dy] = append(data.CurrSelect.Tiles[dy], tile)
		}
	}
	return true
}

func FlipVertical() bool {
	if data.CurrSelect == nil {
		return false
	}
	flip := new(data.Selection)
	flip.Width = data.CurrSelect.Width
	flip.Height = data.CurrSelect.Height
	flip.Offset = data.CurrSelect.Offset
	flip.Origin = data.CurrSelect.Offset
	for dy := len(data.CurrSelect.Tiles) - 1; dy >= 0; dy-- {
		l := len(flip.Tiles)
		flip.Tiles = append(flip.Tiles, data.CurrSelect.Tiles[dy])
		for i := 0; i < len(flip.Tiles[l]); i++ {
			flip.Tiles[l][i].Coords.Y = l
		}
	}
	data.CurrSelect = flip
	return true
}

func FlipHorizontal() bool {
	if data.CurrSelect == nil {
		return false
	}
	flip := new(data.Selection)
	flip.Width = data.CurrSelect.Width
	flip.Height = data.CurrSelect.Height
	flip.Offset = data.CurrSelect.Offset
	flip.Origin = data.CurrSelect.Offset
	for y, row := range data.CurrSelect.Tiles {
		flip.Tiles = append(flip.Tiles, []*data.Tile{})
		for dx := len(row) - 1; dx >= 0; dx-- {
			tile := row[dx]
			tile.Coords.X = len(flip.Tiles[y])
			flip.Tiles[y] = append(flip.Tiles[y], tile)
		}
	}
	data.CurrSelect = flip
	return true
}

func PushUndoArray(resetRedo bool) {
	if data.CurrPuzzleSet.CurrPuzzle.LastChange != nil {
		if len(data.CurrPuzzleSet.CurrPuzzle.UndoStack) > constants.UndoStackSize {
			data.CurrPuzzleSet.CurrPuzzle.UndoStack = data.CurrPuzzleSet.CurrPuzzle.UndoStack[1:]
		}
		data.CurrPuzzleSet.CurrPuzzle.UndoStack = append(data.CurrPuzzleSet.CurrPuzzle.UndoStack, data.CurrPuzzleSet.CurrPuzzle.LastChange)
		if resetRedo {
			data.CurrPuzzleSet.CurrPuzzle.RedoStack = []*data.Tiles{}
		}
	}
	data.CurrPuzzleSet.CurrPuzzle.LastChange = data.CurrPuzzleSet.CurrPuzzle.CopyTiles()
}

func PopUndoArray() {
	data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed = false
	if len(data.CurrPuzzleSet.CurrPuzzle.UndoStack) > 0 {
		puz := data.CurrPuzzleSet.CurrPuzzle.UndoStack[len(data.CurrPuzzleSet.CurrPuzzle.UndoStack)-1]
		for y, row := range puz.T {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzleSet.CurrPuzzle.Tiles.T[y][x])
			}
		}
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.UndoStack = data.CurrPuzzleSet.CurrPuzzle.UndoStack[:len(data.CurrPuzzleSet.CurrPuzzle.UndoStack)-1]
		data.CurrPuzzleSet.CurrPuzzle.LastChange = puz
	}
}

func PushRedoStack() {
	if len(data.CurrPuzzleSet.CurrPuzzle.RedoStack) > 50 {
		data.CurrPuzzleSet.CurrPuzzle.RedoStack = data.CurrPuzzleSet.CurrPuzzle.RedoStack[1:]
	}
	data.CurrPuzzleSet.CurrPuzzle.RedoStack = append(data.CurrPuzzleSet.CurrPuzzle.RedoStack, data.CurrPuzzleSet.CurrPuzzle.CopyTiles())
}

func PopRedoStack() {
	if len(data.CurrPuzzleSet.CurrPuzzle.RedoStack) > 0 {
		puz := data.CurrPuzzleSet.CurrPuzzle.RedoStack[len(data.CurrPuzzleSet.CurrPuzzle.RedoStack)-1]
		for y, row := range puz.T {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzleSet.CurrPuzzle.Tiles.T[y][x])
			}
		}
		data.CurrPuzzleSet.CurrPuzzle.Update = true
		data.CurrPuzzleSet.CurrPuzzle.RedoStack = data.CurrPuzzleSet.CurrPuzzle.RedoStack[:len(data.CurrPuzzleSet.CurrPuzzle.RedoStack)-1]
		data.CurrPuzzleSet.CurrPuzzle.LastChange = puz
	}
}

func UndoStackSystem() {
	if data.CurrPuzzleSet.CurrPuzzle.Changed {
		PushUndoArray(true)
		data.CurrPuzzleSet.CurrPuzzle.Changed = false
		data.CurrPuzzleSet.NeedToSave = true
		data.CurrPuzzleSet.CurrPuzzle.Metadata.Completed = false
	}
}

func DisposeEditor() {
	if data.Editor.PosTop {
		ui.CloseDialog(constants.DialogEditorPanelTop)
		ui.CloseDialog(constants.DialogEditorOptionsBot)
	} else {
		ui.CloseDialog(constants.DialogEditorPanelLeft)
		ui.CloseDialog(constants.DialogEditorOptionsRight)
	}
	data.Editor = nil
}
