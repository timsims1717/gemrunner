package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
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
	data.Editor.BlockSelect = data.Dialogs["block_select"].ViewPort

	// open editor dialogs
	if data.Editor.PosTop {
		data.OpenDialog("editor_panel_top")
		data.OpenDialog("editor_options_bot")
	} else {
		data.OpenDialog("editor_panel_left")
		data.OpenDialog("editor_options_right")
	}
	UpdateWorldShaders()
	PushUndoArray(true)
}

func UpdateWorldShaders() {
	// set puzzle shader uniforms
	data.PuzzleView.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.PrimaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.PrimaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.PrimaryColor.B))
	data.PuzzleView.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.SecondaryColor.R))
	data.PuzzleView.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.SecondaryColor.G))
	data.PuzzleView.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.SecondaryColor.B))
	// set editor panel shader uniforms
	editorPanelLeft := data.Dialogs["editor_panel_left"]
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.PrimaryColor.R))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.PrimaryColor.G))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.PrimaryColor.B))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.SecondaryColor.R))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.SecondaryColor.G))
	editorPanelLeft.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.SecondaryColor.B))
	editorPanelTop := data.Dialogs["editor_panel_top"]
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.PrimaryColor.R))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.PrimaryColor.G))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.PrimaryColor.B))
	editorPanelTop.ViewPort.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.SecondaryColor.R))
	editorPanelTop.ViewPort.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.SecondaryColor.G))
	editorPanelTop.ViewPort.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.SecondaryColor.B))
	// set editor select shader uniforms
	data.Editor.BlockSelect.Canvas.SetUniform("uRedPrimary", float32(data.CurrPuzzle.PrimaryColor.R))
	data.Editor.BlockSelect.Canvas.SetUniform("uGreenPrimary", float32(data.CurrPuzzle.PrimaryColor.G))
	data.Editor.BlockSelect.Canvas.SetUniform("uBluePrimary", float32(data.CurrPuzzle.PrimaryColor.B))
	data.Editor.BlockSelect.Canvas.SetUniform("uRedSecondary", float32(data.CurrPuzzle.SecondaryColor.R))
	data.Editor.BlockSelect.Canvas.SetUniform("uGreenSecondary", float32(data.CurrPuzzle.SecondaryColor.G))
	data.Editor.BlockSelect.Canvas.SetUniform("uBlueSecondary", float32(data.CurrPuzzle.SecondaryColor.B))
}

func ChangeWorldToNext() {
	if data.CurrPuzzle != nil {
		data.CurrPuzzle.WorldNumber++
		if data.CurrPuzzle.WorldNumber >= constants.WorldCustom {
			data.CurrPuzzle.WorldNumber %= constants.WorldCustom
		}
		ChangeWorldTo(data.CurrPuzzle.WorldNumber)
		data.CurrPuzzle.Update = true
	}
}

func ChangeWorldTo(world int) {
	data.CurrPuzzle.WorldSprite = constants.WorldSprites[world]
	data.CurrPuzzle.PrimaryColor = pixel.ToRGBA(constants.WorldPrimary[world])
	data.CurrPuzzle.SecondaryColor = pixel.ToRGBA(constants.WorldSecondary[world])
	UpdateWorldShaders()
}

func UpdateEditorModeHotKey() {
	oldMode := data.Editor.Mode
	if data.EditorInput.Get("ctrl").Pressed() || data.EditorInput.Get("rCtrl").Pressed() {
		if data.EditorInput.Get("ctrlCopy").JustPressed() {
			// copy
			data.Editor.Mode = data.Copy
		} else if data.EditorInput.Get("ctrlCut").JustPressed() {
			// cut
			data.Editor.Mode = data.Cut
		} else if data.EditorInput.Get("ctrlPaste").JustPressed() {
			// paste
			data.Editor.Mode = data.Paste
		} else if (data.EditorInput.Get("shift").Pressed() || data.EditorInput.Get("rShift").Pressed()) &&
			data.EditorInput.Get("ctrlShiftRedo").JustPressed() {
			// redo
			data.Editor.Mode = data.Redo
		} else if data.EditorInput.Get("ctrlUndo").JustPressed() {
			// undo
			data.Editor.Mode = data.Undo
		} else if data.EditorInput.Get("ctrlSave").JustPressed() {
			// save
			data.Editor.Mode = data.Save
		} else if data.EditorInput.Get("ctrlOpen").JustPressed() {
			// load
			data.Editor.Mode = data.Open
		}
	} else {
		for i := 0; i < data.EndModeList; i++ {
			hotkey := data.EditorInput.Get(data.EditorMode(i).String())
			if hotkey != nil && hotkey.JustPressed() {
				hotkey.Consume()
				data.Editor.Mode = data.EditorMode(i)
			}
		}
	}
	if data.Editor.Mode >= data.EndModeList {
		data.Editor.Mode = data.Brush
	}
	if oldMode != data.Editor.Mode {
		data.Editor.LastMode = oldMode
		data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		data.Editor.SelectVis = false
	}
	if data.Editor.LastMode >= data.EndModeList {
		data.Editor.LastMode = data.Brush
	}
}

func PuzzleEditSystem() {
	EditorPanelButtons()
	if data.Editor.SelectTimer != nil {
		data.Editor.SelectTimer.Update()
	}
	projPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
	x, y := world.WorldToMap(projPos.X, projPos.Y)
	coords := world.Coords{X: x, Y: y}
	// special check for a weird bug
	if coords.X < -1000 || coords.X > 1000 || coords.Y < -1000 || coords.Y > 1000 {
		return
	}
	legal := CoordsLegal(coords)
	lastLegal := CoordsLegal(data.Editor.LastCoords)
	rClick := data.EditorInput.Get("rightClick")
	click := data.EditorInput.Get("click")
	inInside := data.Editor.BlockSelect.PointInside(data.Editor.BlockSelect.ProjectWorld(data.EditorInput.World))
	if data.Editor.SelectVis && legal && !inInside &&
		(click.JustPressed() || rClick.JustPressed()) {
		data.Editor.SelectVis = false
	}
	//if data.Editor.Consume != "select" {
	//	data.Editor.SelectObj.Pos = data.BlockSelectPlacement(int(data.Editor.CurrBlock))
	//}
	if data.Editor.SelectVis {
		data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		//switch data.Editor.Consume {
		//case "move":
		//	if data.Editor.Offset.X != 0 || data.Editor.Offset.Y != 0 {
		//		data.Editor.ViewPort.PortPos = data.EditorInput.World.Add(data.Editor.Offset)
		//	}
		//}
	} else {
		data.CurrPuzzle.Click = rClick.Pressed() || click.Pressed()
		if !rClick.Pressed() && !rClick.JustReleased() &&
			!click.Pressed() && !click.JustReleased() && legal {
			CreateHighlight(coords)
			data.Editor.NoInput = false
			data.Editor.LastCoords = world.Coords{X: -1, Y: -1}
		} else if !data.Editor.NoInput {
			switch data.Editor.Mode {
			case data.Brush:
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
					PushUndoArray(true)
				}
				if legal {
					CreateHighlight(coords)
					data.Editor.LastCoords = coords
				}
			case data.Line:
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
						PushUndoArray(true)
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
						PushUndoArray(true)
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
			case data.Square:
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
						PushUndoArray(true)
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
						PushUndoArray(true)
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
			case data.Fill:
				if rClick.JustReleased() {
					Fill(coords, true)
				} else if click.JustReleased() {
					Fill(coords, false)
				}
				if legal {
					CreateHighlight(coords)
				}
			case data.Erase:
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
					PushUndoArray(true)
				}
				if legal {
					CreateHighlight(coords)
					data.Editor.LastCoords = coords
				}
			case data.Eyedrop:
				if legal {
					if rClick.JustReleased() || click.JustReleased() {
						tile := data.CurrPuzzle.Tiles.T[coords.Y][coords.X]
						if tile.Block != data.Empty {
							data.Editor.CurrBlock = tile.Block
						}
						if data.Editor.LastMode > data.Fill {
							data.Editor.Mode, data.Editor.LastMode = data.Brush, data.Editor.Mode
						} else {
							data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Editor.Mode
						}
					}
					CreateHighlight(coords)
				}
			case data.Select:
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
						data.Editor.Mode = data.Move
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
			case data.Move:
				if legal {
					CreateHighlight(coords)
					if click.Pressed() {
						if rClick.JustPressed() {
							data.Editor.Mode = data.Select
							data.Editor.NoInput = true
							data.CurrSelect.Offset = data.CurrSelect.Origin
							break
						} else if click.JustPressed() {
							inTest := coords
							inTest.X -= data.CurrSelect.Offset.X
							inTest.Y -= data.CurrSelect.Offset.Y
							if !CoordsLegalSelection(inTest) {
								data.Editor.Mode = data.Select
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
						data.Editor.Mode = data.Brush
						data.Editor.NoInput = true
						data.CurrSelect.Offset = data.CurrSelect.Origin
						break
					}
				}
			}
		}
		switch data.Editor.Mode {
		case data.Copy:
			if CreateClip() {
				data.Editor.Mode = data.Move
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
			}
		case data.Cut:
			if CreateClip() {
				data.CurrSelect = nil
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
		case data.Paste:
			// place the current clip
			PlaceSelection()
			if PlaceClip() {
				data.Editor.Mode = data.Move
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
			}
		case data.Delete:
			data.CurrSelect = nil
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
			PushUndoArray(true)
		case data.Undo:
			if len(data.CurrPuzzle.UndoStack) > 0 {
				PushRedoStack()
				PopUndoArray()
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
		case data.Redo:
			if len(data.CurrPuzzle.RedoStack) > 0 {
				PushUndoArray(false)
				PopRedoStack()
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
		case data.FlipVertical:
			if data.CurrSelect != nil {
				FlipVertical()
				data.Editor.Mode = data.Move
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
			}
		case data.FlipHorizontal:
			if data.CurrSelect != nil {
				FlipHorizontal()
				data.Editor.Mode = data.Move
			} else {
				data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
			}
		case data.Save:
			PlaceSelection()
			if err := SavePuzzle(); err != nil {
				fmt.Println("ERROR:", err)
			}
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
		case data.Open:
			//if err := LoadPuzzle(); err != nil {
			//	fmt.Println("Error:", err)
			//}
			data.OpenDialogInStack("open_puzzle")
			data.Editor.Mode, data.Editor.LastMode = data.Editor.LastMode, data.Brush
		}
	}
	if data.Editor.Mode == data.Move {
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
					obj.Pos.X += world.TileSize * 0.5
					obj.Pos.Y += world.TileSize * 0.5
					spr := GetTileSpritesSelection(tile)
					if len(spr) > 0 {
						myecs.Manager.NewEntity().
							AddComponent(myecs.Object, obj).
							AddComponent(myecs.Drawable, spr).
							AddComponent(myecs.Temp, myecs.ClearFlag(true))
					}
					// don't draw anything below
					if CoordsLegal(c) {
						data.CurrPuzzle.Tiles.T[c.Y][c.X].Object.Hidden = true
					}
				}
			}
		} else {
			data.Editor.Mode = data.Select
		}
	} else {
		// place the current selection
		PlaceSelection()
	}
	if data.Editor.Mode >= data.EndModeList {
		data.Editor.Mode = data.Brush
	}
	if data.Editor.LastMode >= data.EndModeList {
		data.Editor.LastMode = data.Brush
	}
}

func EditorPanelButtons() {
	var panel *data.Dialog
	if data.Editor.PosTop {
		panel = data.Dialogs["editor_panel_top"]
	} else {
		panel = data.Dialogs["editor_panel_left"]
	}
	if !panel.Click {
		for _, e := range panel.Elements {
			if btn, ok := e.(*data.Button); ok {
				if data.ModeFromSprString(btn.Sprite.Key) == data.Editor.Mode ||
					(data.ModeFromSprString(btn.Sprite.Key) == data.Select &&
						data.Editor.Mode == data.Move) {
					btn.Entity.AddComponent(myecs.Drawable, btn.ClickSpr)
				} else {
					btn.Entity.AddComponent(myecs.Drawable, btn.Sprite)
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
		obj.Pos.X += world.TileSize * 0.5
		obj.Pos.Y += world.TileSize * 0.5
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
		orig := data.CurrPuzzle.Tiles.T[a.Y][a.X]
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
					tile := data.CurrPuzzle.Tiles.T[n.Y][n.X]
					if tile.Block == orig.Block && tile.Ladder == orig.Ladder {
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
		PushUndoArray(true)
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
			old := data.CurrPuzzle.Tiles.T[dy+o.Y][dx+o.X]
			tile := old.Copy()
			tile.Coords = world.Coords{X: dx, Y: dy}
			data.CurrSelect.Tiles[dy] = append(data.CurrSelect.Tiles[dy], tile)
			old.Empty()
		}
	}
	data.CurrPuzzle.Update = true
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
	hasP1 := false
	for _, row := range data.CurrSelect.Tiles {
		for _, tile := range row {
			if tile.Block == data.Player1 {
				hasP1 = true
			}
		}
	}
	for _, row := range data.CurrPuzzle.Tiles.T {
		for _, tile := range row {
			if tile.Block == data.Player1 && hasP1 {
				tile.Block = data.Empty
			}
		}
	}
	for _, row := range data.CurrSelect.Tiles {
		for _, tile := range row {
			c := tile.Coords
			c.X += data.CurrSelect.Offset.X
			c.Y += data.CurrSelect.Offset.Y
			if CoordsLegal(c) {
				tile.CopyInto(data.CurrPuzzle.Tiles.T[c.Y][c.X])
				data.CurrPuzzle.Update = true
			}
		}
	}
	data.CurrSelect = nil
	PushUndoArray(true)
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
	if data.CurrPuzzle.LastChange != nil {
		if len(data.CurrPuzzle.UndoStack) > 50 {
			data.CurrPuzzle.UndoStack = data.CurrPuzzle.UndoStack[1:]
		}
		data.CurrPuzzle.UndoStack = append(data.CurrPuzzle.UndoStack, data.CurrPuzzle.LastChange)
		if resetRedo {
			data.CurrPuzzle.RedoStack = []*data.Tiles{}
		}
	}
	data.CurrPuzzle.LastChange = data.CurrPuzzle.CopyTiles()
}

func PopUndoArray() {
	if len(data.CurrPuzzle.UndoStack) > 0 {
		puz := data.CurrPuzzle.UndoStack[len(data.CurrPuzzle.UndoStack)-1]
		for y, row := range puz.T {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzle.Tiles.T[y][x])
			}
		}
		data.CurrPuzzle.Update = true
		data.CurrPuzzle.UndoStack = data.CurrPuzzle.UndoStack[:len(data.CurrPuzzle.UndoStack)-1]
		data.CurrPuzzle.LastChange = puz
	}
}

func PushRedoStack() {
	if len(data.CurrPuzzle.RedoStack) > 50 {
		data.CurrPuzzle.RedoStack = data.CurrPuzzle.RedoStack[1:]
	}
	data.CurrPuzzle.RedoStack = append(data.CurrPuzzle.RedoStack, data.CurrPuzzle.CopyTiles())
}

func PopRedoStack() {
	if len(data.CurrPuzzle.RedoStack) > 0 {
		puz := data.CurrPuzzle.RedoStack[len(data.CurrPuzzle.RedoStack)-1]
		for y, row := range puz.T {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzle.Tiles.T[y][x])
			}
		}
		data.CurrPuzzle.Update = true
		data.CurrPuzzle.RedoStack = data.CurrPuzzle.RedoStack[:len(data.CurrPuzzle.RedoStack)-1]
		data.CurrPuzzle.LastChange = puz
	}
}

func DisposeEditor() {

}
