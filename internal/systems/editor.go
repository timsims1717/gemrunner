package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/pkg/img"
	"gemrunner/pkg/object"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

func EditorInit() {
	// initialize imDraw
	data.IMDraw = imdraw.New(nil)

	// initialize editor panel
	data.NewEditorPane()
	data.EditorPanel.ViewPort = viewport.New(nil)
	data.EditorPanel.ViewPort.SetILock(true)
	data.EditorPanel.ViewPort.SetRect(pixel.R(0, 0, world.TileSize*3., world.TileSize*8.))
	data.EditorPanel.ViewPort.CamPos = pixel.V(world.TileSize*-2.-3., world.TileSize*4.5+3.)

	// border
	borderObj := object.New()
	borderObj.Pos = pixel.V(world.TileSize*-2., world.TileSize*4.5)
	borderObj.SetRect(pixel.R(0, 0, world.TileSize*3, world.TileSize*8))
	borderObj.Layer = 10

	// the viewport for the block selectors
	vp := viewport.New(nil)
	vp.SetRect(pixel.R(0, 0, world.TileSize*6.+2., world.TileSize*3.+2))
	vp.CamPos = pixel.V(world.TileSize*3., -world.TileSize*1.5)
	data.EditorPanel.BlockSelect = vp

	// the block selectors
	b := 0
	for ; b < data.Empty; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		//fmt.Printf("Tile %d: (%d,%d)\n", b, int(obj.Pos.X), int(obj.Pos.Y))
		obj.Layer = 11
		obj.SetRect(pixel.R(0., 0., 16., 16.))
		sprB := img.NewSprite("black_square_big", constants.UIBatch)
		spr := img.NewSprite(data.Block(b).String(), constants.BGBatch)
		bId := data.Block(b)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, []*img.Sprite{sprB, spr}).
			AddComponent(myecs.Block, data.Block(b)).
			AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.BlockSelect, func(hvc *data.HoverClick) {
				click := hvc.Input.Get("click")
				if hvc.Hover && data.EditorPanel.SelectVis {
					data.EditorPanel.SelectObj.Pos = obj.Pos
					if click.JustPressed() || click.JustReleased() {
						data.EditorPanel.CurrBlock = bId
						data.EditorPanel.SelectVis = false
						data.EditorPanel.SelectQuick = false
						data.EditorPanel.SelectTimer = nil
						data.EditorPanel.Consume = ""
						switch data.EditorPanel.Mode {
						case data.Brush, data.Line, data.Square, data.Fill:
						default:
							data.EditorPanel.Mode = data.Brush
						}
						click.Consume()
					}
				}
			}))
	}
	objOutline := object.New()
	objOutline.Pos = data.BlockSelectPlacement(0)
	objOutline.Layer = 12
	sprO := img.NewSprite("white_outline", constants.UIBatch)
	myecs.Manager.NewEntity().
		AddComponent(myecs.Object, objOutline).
		AddComponent(myecs.Drawable, sprO)
	data.EditorPanel.SelectObj = objOutline
	for ; b < 18; b++ {
		obj := object.New()
		obj.Pos = data.BlockSelectPlacement(b)
		obj.Layer = 11
		spr := img.NewSprite("black_square_big", constants.UIBatch)
		myecs.Manager.NewEntity().
			AddComponent(myecs.Object, obj).
			AddComponent(myecs.Drawable, spr)
	}

	// block select
	blockObj := object.New()
	blockObj.Pos = borderObj.Pos
	blockObj.Pos.Y -= world.TileSize * 2.5
	blockObj.Layer = 10
	blockObj.Rect = pixel.R(-16., -16., 16., 16.)
	beBG := img.NewSprite("editor_tile_bg", constants.UIBatch)
	beFG := img.NewSprite(data.Block(data.Turf).String(), constants.BGBatch)
	be := myecs.Manager.NewEntity()
	be.AddComponent(myecs.Object, blockObj).
		AddComponent(myecs.Drawable, []*img.Sprite{beBG, beFG}).
		AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.ViewPort, func(hvc *data.HoverClick) {
			beFG.Key = data.EditorPanel.CurrBlock.String()
			data.EditorPanel.Hover = hvc.Hover
			click := hvc.Input.Get("click")
			if hvc.Hover && (data.EditorPanel.Consume == "select" || data.EditorPanel.Consume == "") {
				if data.EditorPanel.Consume == "" {
					data.EditorPanel.SelectQuick = false
					if click.JustPressed() {
						data.EditorPanel.Consume = "select"
						data.EditorPanel.SelectTimer = timing.New(0.2)
					}
				} else if data.EditorPanel.Consume == "select" {
					if click.JustPressed() {
						data.EditorPanel.Consume = ""
						data.EditorPanel.SelectTimer = nil
						data.EditorPanel.SelectQuick = false
					} else if click.JustReleased() {
						if data.EditorPanel.SelectTimer != nil && !data.EditorPanel.SelectTimer.Done() {
							data.EditorPanel.SelectQuick = true
							data.EditorPanel.Consume = "select"
						}
					} else if !click.Pressed() && !data.EditorPanel.SelectQuick {
						data.EditorPanel.Consume = ""
						data.EditorPanel.SelectTimer = nil
					}
				}
				data.EditorPanel.SelectVis = data.EditorPanel.Consume == "select"
			}
		}))

	// border and editor panel movement
	e := myecs.Manager.NewEntity()
	e.AddComponent(myecs.Border, &data.Border{
		Width:  2,
		Height: 7,
		Empty:  false,
	}).
		AddComponent(myecs.Object, borderObj).
		AddComponent(myecs.Update, data.NewHoverClickFn(data.EditorInput, data.EditorPanel.ViewPort, func(hvc *data.HoverClick) {
			data.EditorPanel.Hover = hvc.Hover
			click := hvc.Input.Get("click")
			if hvc.Hover {
				if data.EditorPanel.Consume == "move" || data.EditorPanel.Consume == "" {
					if click.JustPressed() {
						data.EditorPanel.Offset = data.EditorPanel.ViewPort.PostPortPos.Sub(hvc.Input.World)
						data.EditorPanel.Consume = "move"
					}
				}
			}
			if click.JustReleased() && data.EditorPanel.Offset != pixel.ZV {
				data.EditorPanel.ViewPort.PortPos = hvc.Input.World.Add(data.EditorPanel.Offset)
				data.EditorPanel.Offset = pixel.ZV
				data.EditorPanel.Consume = ""
			}
		}))
	data.EditorPanel.Entity = e
	data.EditorPanel.BlockView = &data.BlockView{
		Entity: be,
		Object: blockObj,
	}

	PushUndoArray(true)
}

func PuzzleInit() {
	for _, result := range myecs.Manager.Query(myecs.IsTile) {
		myecs.Manager.DisposeEntity(result)
	}
	if data.CurrPuzzle != nil {
		for _, row := range data.CurrPuzzle.Tiles {
			for _, tile := range row {
				obj := object.New()
				obj.Pos = world.MapToWorld(tile.Coords)
				obj.Pos.X += world.TileSize * 0.5
				obj.Pos.Y += world.TileSize * 0.5
				obj.Layer = 2
				myecs.Manager.NewEntity().
					AddComponent(myecs.Object, obj).
					AddComponent(myecs.Tile, tile)
				tile.Object = obj
			}
		}
		data.CurrPuzzle.Update = true
	}

	data.PuzzleView = viewport.New(nil)
	data.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*constants.PuzzleWidth, world.TileSize*constants.PuzzleHeight))
	data.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
	data.PuzzleView.PortPos = viewport.MainCamera.CamPos
	data.BorderView = viewport.New(nil)
	data.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(constants.PuzzleWidth+1), world.TileSize*(constants.PuzzleHeight+1)))
	data.BorderView.CamPos = pixel.V(world.TileSize*0.5*(constants.PuzzleWidth), world.TileSize*0.5*(constants.PuzzleHeight))
}

func UpdateEditorModeHotKey() {
	oldMode := data.EditorPanel.Mode
	if data.EditorInput.Get("ctrl").Pressed() {
		if data.EditorInput.Get("ctrlCopy").JustPressed() {
			// copy
			data.EditorPanel.Mode = data.Copy
		} else if data.EditorInput.Get("ctrlCut").JustPressed() {
			// cut
			data.EditorPanel.Mode = data.Cut
		} else if data.EditorInput.Get("ctrlPaste").JustPressed() {
			// paste
			data.EditorPanel.Mode = data.Paste
		} else if data.EditorInput.Get("shift").Pressed() &&
			data.EditorInput.Get("ctrlShiftRedo").JustPressed() {
			// redo
			data.EditorPanel.Mode = data.Redo
		} else if data.EditorInput.Get("ctrlUndo").JustPressed() {
			// undo
			data.EditorPanel.Mode = data.Undo
		}
	} else {
		for i := 0; i < data.EndModeList; i++ {
			hotkey := data.EditorInput.Get(data.EditorMode(i).String())
			if hotkey != nil && hotkey.JustPressed() {
				hotkey.Consume()
				data.EditorPanel.Mode = data.EditorMode(i)
			}
		}
	}
	if data.EditorPanel.Mode >= data.EndModeList {
		data.EditorPanel.Mode = data.Brush
	}
	if oldMode != data.EditorPanel.Mode {
		data.EditorPanel.LastMode = oldMode
	}
	if data.EditorPanel.LastMode >= data.EndModeList {
		data.EditorPanel.LastMode = data.Brush
	}
}

func PuzzleEditSystem() {
	if data.EditorPanel.SelectTimer != nil {
		data.EditorPanel.SelectTimer.Update()
	}
	if data.EditorPanel.Consume != "select" {
		data.EditorPanel.SelectObj.Pos = data.BlockSelectPlacement(int(data.EditorPanel.CurrBlock))
	}
	if data.EditorPanel.Consume != "" {
		data.EditorPanel.LastCoords = world.Coords{X: -1, Y: -1}
		switch data.EditorPanel.Consume {
		case "move":
			if data.EditorPanel.Offset.X != 0 || data.EditorPanel.Offset.Y != 0 {
				data.EditorPanel.ViewPort.PortPos = data.EditorInput.World.Add(data.EditorPanel.Offset)
			}
		}
	} else {
		if !data.EditorPanel.Hover {
			//data.EditorPanel.SelectVis = false
			projPos := data.PuzzleView.ProjectWorld(data.EditorInput.World)
			// switch editor mode
			x, y := world.WorldToMap(projPos.X, projPos.Y)
			coords := world.Coords{X: x, Y: y}
			if coords.X < -1000 || coords.X > 1000 || coords.Y < -1000 || coords.Y > 1000 {
				return
			}
			legal := CoordsLegal(coords)
			lastLegal := CoordsLegal(data.EditorPanel.LastCoords)
			rClick := data.EditorInput.Get("rightClick")
			click := data.EditorInput.Get("click")
			data.CurrPuzzle.Click = rClick.Pressed() || click.Pressed()
			if !rClick.Pressed() && !rClick.JustReleased() &&
				!click.Pressed() && !click.JustReleased() && legal {
				CreateHighlight(coords)
			} else {
				switch data.EditorPanel.Mode {
				case data.Brush:
					if rClick.JustPressed() || click.JustPressed() || !lastLegal {
						data.EditorPanel.LastCoords = coords
						lastLegal = CoordsLegal(data.EditorPanel.LastCoords)
					}
					line := world.Line(data.EditorPanel.LastCoords, coords)
					if rClick.Pressed() {
						for _, c := range line {
							DeleteBlock(c)
						}
					} else if click.Pressed() {
						for _, c := range line {
							SetBlock(c, data.EditorPanel.CurrBlock)
						}
					}
					if rClick.JustReleased() || click.JustReleased() {
						PushUndoArray(true)
					}
					if legal {
						CreateHighlight(coords)
						data.EditorPanel.LastCoords = coords
					}
				case data.Line:
					if lastLegal {
						if click.Pressed() && !data.EditorPanel.NoInput {
							if click.JustPressed() {
								data.EditorPanel.LastCoords = coords
							} else if rClick.JustPressed() {
								data.EditorPanel.NoInput = true
								break
							}
							line := world.Line(data.EditorPanel.LastCoords, coords)
							for _, c := range line {
								CreateHighlight(c)
							}
						} else if click.JustReleased() && !data.EditorPanel.NoInput {
							line := world.Line(data.EditorPanel.LastCoords, coords)
							for _, c := range line {
								SetBlock(c, data.EditorPanel.CurrBlock)
							}
							PushUndoArray(true)
							if legal {
								CreateHighlight(coords)
							}
						} else if rClick.Pressed() && !data.EditorPanel.NoInput {
							if rClick.JustPressed() {
								data.EditorPanel.LastCoords = coords
							} else if click.JustPressed() {
								data.EditorPanel.NoInput = true
								break
							}
							line := world.Line(data.EditorPanel.LastCoords, coords)
							for _, c := range line {
								CreateHighlight(c)
							}
						} else if rClick.JustReleased() && !data.EditorPanel.NoInput {
							line := world.Line(data.EditorPanel.LastCoords, coords)
							for _, c := range line {
								DeleteBlock(c)
							}
							PushUndoArray(true)
							if legal {
								CreateHighlight(coords)
							}
						} else if !click.Pressed() && !rClick.Pressed() {
							data.EditorPanel.NoInput = false
							if legal {
								CreateHighlight(coords)
								data.EditorPanel.LastCoords = coords
							}
						}
					} else {
						if legal {
							data.EditorPanel.NoInput = false
							CreateHighlight(coords)
							data.EditorPanel.LastCoords = coords
						}
					}
				case data.Square:
					if lastLegal {
						if click.Pressed() && !data.EditorPanel.NoInput {
							if click.JustPressed() {
								data.EditorPanel.LastCoords = coords
							} else if rClick.JustPressed() {
								data.EditorPanel.NoInput = true
								break
							}
							square := world.Square(data.EditorPanel.LastCoords, coords)
							for _, c := range square {
								CreateHighlight(c)
							}
						} else if click.JustReleased() && !data.EditorPanel.NoInput {
							square := world.Square(data.EditorPanel.LastCoords, coords)
							for _, c := range square {
								SetBlock(c, data.EditorPanel.CurrBlock)
							}
							PushUndoArray(true)
							if legal {
								CreateHighlight(coords)
							}
						} else if rClick.Pressed() && !data.EditorPanel.NoInput {
							if rClick.JustPressed() {
								data.EditorPanel.LastCoords = coords
							} else if click.JustPressed() {
								data.EditorPanel.NoInput = true
								break
							}
							square := world.Square(data.EditorPanel.LastCoords, coords)
							for _, c := range square {
								CreateHighlight(c)
							}
						} else if rClick.JustReleased() && !data.EditorPanel.NoInput {
							square := world.Square(data.EditorPanel.LastCoords, coords)
							for _, c := range square {
								DeleteBlock(c)
							}
							PushUndoArray(true)
							if legal {
								CreateHighlight(coords)
							}
						} else if !click.Pressed() && !rClick.Pressed() {
							data.EditorPanel.NoInput = false
							if legal {
								CreateHighlight(coords)
								data.EditorPanel.LastCoords = coords
							}
						}
					} else {
						if legal {
							data.EditorPanel.NoInput = false
							CreateHighlight(coords)
							data.EditorPanel.LastCoords = coords
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
						data.EditorPanel.LastCoords = coords
						lastLegal = CoordsLegal(data.EditorPanel.LastCoords)
					}
					line := world.Line(data.EditorPanel.LastCoords, coords)
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
						data.EditorPanel.LastCoords = coords
					}
				case data.Eyedrop:
					if legal {
						if rClick.JustReleased() || click.JustReleased() {
							tile := data.CurrPuzzle.Tiles[coords.Y][coords.X]
							if tile.Block != data.Empty {
								data.EditorPanel.CurrBlock = tile.Block
							}
							data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.EditorPanel.Mode
						}
						CreateHighlight(coords)
					}
				case data.Select:
					if lastLegal {
						if click.Pressed() && !data.EditorPanel.NoInput {
							if click.JustPressed() {
								data.EditorPanel.LastCoords = coords
							} else if rClick.JustPressed() {
								data.EditorPanel.NoInput = true
								break
							}
							CreateSquareSelect(data.EditorPanel.LastCoords, coords)
						} else if click.JustReleased() && !data.EditorPanel.NoInput {
							CreateSelection(data.EditorPanel.LastCoords, coords)
							data.EditorPanel.Mode = data.Move
						} else if !click.Pressed() && !rClick.Pressed() {
							data.EditorPanel.NoInput = false
							if legal {
								CreateHighlight(coords)
								data.EditorPanel.LastCoords = coords
							}
						}
					} else {
						if legal {
							data.EditorPanel.NoInput = false
							CreateHighlight(coords)
							data.EditorPanel.LastCoords = coords
						}
					}
				case data.Move:
					if legal {
						CreateHighlight(coords)
					}
					if lastLegal {
						if click.Pressed() && !data.EditorPanel.NoInput {
							if click.JustPressed() {
								inTest := coords
								inTest.X -= data.CurrSelect.Offset.X
								inTest.Y -= data.CurrSelect.Offset.Y
								if !CoordsLegalSelection(inTest) {
									data.EditorPanel.Mode = data.Select
									data.EditorPanel.NoInput = true
									break
								}
								data.EditorPanel.LastCoords = coords
							} else if rClick.JustPressed() {
								data.EditorPanel.NoInput = true
								data.CurrSelect.Offset = data.CurrSelect.Origin
								break
							} else {
								var move world.Coords
								move.X = coords.X - data.EditorPanel.LastCoords.X
								move.Y = coords.Y - data.EditorPanel.LastCoords.Y
								data.CurrSelect.Offset.X = data.CurrSelect.Origin.X + move.X
								data.CurrSelect.Offset.Y = data.CurrSelect.Origin.Y + move.Y
							}
							end := world.Coords{
								X: data.CurrSelect.Offset.X + data.CurrSelect.Width - 1,
								Y: data.CurrSelect.Offset.Y + data.CurrSelect.Height - 1,
							}
							CreateSquareSelect(data.CurrSelect.Offset, end)
						} else if click.JustReleased() && !data.EditorPanel.NoInput {
							var move world.Coords
							move.X = coords.X - data.EditorPanel.LastCoords.X
							move.Y = coords.Y - data.EditorPanel.LastCoords.Y
							data.CurrSelect.Offset.X = data.CurrSelect.Origin.X + move.X
							data.CurrSelect.Offset.Y = data.CurrSelect.Origin.Y + move.Y
							data.CurrSelect.Origin = data.CurrSelect.Offset
							if legal {
								data.EditorPanel.LastCoords = coords
							}
						} else if !click.Pressed() && !rClick.Pressed() {
							data.EditorPanel.NoInput = false
							if legal {
								data.EditorPanel.LastCoords = coords
							}
						}
					} else {
						if legal {
							data.EditorPanel.NoInput = false
							data.EditorPanel.LastCoords = coords
						}
					}
				}
			}
		}
		switch data.EditorPanel.Mode {
		case data.Copy:
			if CreateClip() {
				data.EditorPanel.Mode = data.Move
			} else {
				data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
			}
		case data.Cut:
			if CreateClip() {
				data.CurrSelect = nil
			}
			data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
		case data.Paste:
			// place the current clip
			PlaceSelection()
			if PlaceClip() {
				data.EditorPanel.Mode = data.Move
			} else {
				data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
			}
		case data.Delete:
			data.CurrSelect = nil
			data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
			PushUndoArray(true)
		case data.Undo:
			if len(data.EditorPanel.UndoStack) > 0 {
				PushRedoStack()
				PopUndoArray()
			}
			data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
		case data.Redo:
			if len(data.EditorPanel.RedoStack) > 0 {
				PushUndoArray(false)
				PopRedoStack()
			}
			data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
		case data.FlipVertical:
			if data.CurrSelect != nil {
				FlipVertical()
				data.EditorPanel.Mode = data.Move
			} else {
				data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
			}
		case data.FlipHorizontal:
			if data.CurrSelect != nil {
				FlipHorizontal()
				data.EditorPanel.Mode = data.Move
			} else {
				data.EditorPanel.Mode, data.EditorPanel.LastMode = data.EditorPanel.LastMode, data.Brush
			}
		}
	}
	if data.EditorPanel.Mode == data.Move {
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
						data.CurrPuzzle.Tiles[c.Y][c.X].Object.Hidden = true
					}
				}
			}
		} else {
			data.EditorPanel.Mode = data.Select
		}
	} else {
		// place the current selection
		PlaceSelection()
	}
	if data.EditorPanel.Mode >= data.EndModeList {
		data.EditorPanel.Mode = data.Brush
	}
	if data.EditorPanel.LastMode >= data.EndModeList {
		data.EditorPanel.LastMode = data.Brush
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
		orig := data.CurrPuzzle.Tiles[a.Y][a.X]
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
					tile := data.CurrPuzzle.Tiles[n.Y][n.X]
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
				SetBlock(c, data.EditorPanel.CurrBlock)
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
	data.IMDraw.Color = constants.WhiteColor
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
			old := data.CurrPuzzle.Tiles[dy+o.Y][dx+o.X]
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
	for _, row := range data.CurrSelect.Tiles {
		for _, tile := range row {
			c := tile.Coords
			c.X += data.CurrSelect.Offset.X
			c.Y += data.CurrSelect.Offset.Y
			if CoordsLegal(c) {
				tile.CopyInto(data.CurrPuzzle.Tiles[c.Y][c.X])
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
	if data.EditorPanel.LastChange != nil {
		if len(data.EditorPanel.UndoStack) > 50 {
			data.EditorPanel.UndoStack = data.EditorPanel.UndoStack[1:]
		}
		data.EditorPanel.UndoStack = append(data.EditorPanel.UndoStack, data.EditorPanel.LastChange)
		if resetRedo {
			data.EditorPanel.RedoStack = []*data.Puzzle{}
		}
	}
	data.EditorPanel.LastChange = data.CurrPuzzle.Copy()
}

func PopUndoArray() {
	if len(data.EditorPanel.UndoStack) > 0 {
		puz := data.EditorPanel.UndoStack[len(data.EditorPanel.UndoStack)-1]
		for y, row := range puz.Tiles {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzle.Tiles[y][x])
			}
		}
		data.CurrPuzzle.Update = true
		data.EditorPanel.UndoStack = data.EditorPanel.UndoStack[:len(data.EditorPanel.UndoStack)-1]
		data.EditorPanel.LastChange = puz
	}
}

func PushRedoStack() {
	if len(data.EditorPanel.RedoStack) > 50 {
		data.EditorPanel.RedoStack = data.EditorPanel.RedoStack[1:]
	}
	data.EditorPanel.RedoStack = append(data.EditorPanel.RedoStack, data.CurrPuzzle.Copy())
}

func PopRedoStack() {
	if len(data.EditorPanel.RedoStack) > 0 {
		puz := data.EditorPanel.RedoStack[len(data.EditorPanel.RedoStack)-1]
		for y, row := range puz.Tiles {
			for x, tile := range row {
				tile.CopyInto(data.CurrPuzzle.Tiles[y][x])
			}
		}
		data.CurrPuzzle.Update = true
		data.EditorPanel.RedoStack = data.EditorPanel.RedoStack[:len(data.EditorPanel.RedoStack)-1]
		data.EditorPanel.LastChange = puz
	}
}
