package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/timing"
	"gemrunner/pkg/util"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
	"math"
	"strings"
)

//var keepTheseKeys1 = []string{
//	"puzzle_center", "puzzle_left", "puzzle_right",
//	"puzzle_top_center", "puzzle_top_left", "puzzle_top_right",
//	"puzzle_bot_center", "puzzle_bot_left", "puzzle_bot_right",
//	"puzzle_float_center", "puzzle_float_left", "puzzle_float_right",
//}

func SetAdventureViewMovement(element *ui.Element, dlg *ui.Dialog) {
	element.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, element.ViewPort, func(hvc *data.HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock {
			if hvc.ViewHover {
				//if hvc.Input.ScrollV > 0. {
				//	Zoom(dlg.Key, rearrange, true)
				//} else if hvc.Input.ScrollV < 0. {
				//	Zoom(dlg.Key, rearrange, false)
				//} else {
				halfW := element.ViewPort.Canvas.Bounds().W() * 0.5
				halfH := element.ViewPort.Canvas.Bounds().H() * 0.5
				inPos := hvc.Pos
				inPos.X -= element.ViewPort.CamPos.X
				inPos.Y -= element.ViewPort.CamPos.Y
				if inPos.X > halfW-world.TileSize || hvc.Input.Get("right").Pressed() { // move view right
					element.ViewPort.CamPos.X += constants.AdventureViewScrollSpeed * timing.DT
				} else if inPos.X < -halfW+world.TileSize || hvc.Input.Get("left").Pressed() { // move view left
					element.ViewPort.CamPos.X -= constants.AdventureViewScrollSpeed * timing.DT
				}
				if inPos.Y > halfH-world.TileSize || hvc.Input.Get("up").Pressed() { // move view up
					element.ViewPort.CamPos.Y += constants.AdventureViewScrollSpeed * timing.DT
				} else if inPos.Y < -halfH+world.TileSize || hvc.Input.Get("down").Pressed() { // move view down
					element.ViewPort.CamPos.Y -= constants.AdventureViewScrollSpeed * timing.DT
				}
				ClampAdventureView(element.ViewPort)
				//}
				dlg.SetFocus(element, true)
			}
		}
	}))
}

func SetPreviewMediumClick(element *ui.Element, dlg *ui.Dialog, pzlView *ui.Element) {
	clicked := false
	lastPos := element.Object.Pos
	offset := pixel.ZV
	drag := false
	element.Entity.AddComponent(myecs.Update, data.NewHoverClickFn(data.MenuInput, pzlView.ViewPort, func(hvc *data.HoverClick) {
		if dlg.Open && dlg.Active && !dlg.Lock {
			if hvc.ViewHover {
				click := hvc.Input.Get("click")
				if clicked {
					if click.JustReleased() {
						// highlight it
						element.Border.Style = data.ThickBorderWhite
						clicked = false
						drag = false
						grid := GetGridFromPzlView(element.Object.Pos.X, element.Object.Pos.Y)
						if swapEntry, ok := data.AdventureViewGridMap[grid]; ok { // if another puzzle is here, swap them
							// move grid map entry to new grid pos and vice versa
							prevGrid := data.AdventureViewGridArr[0]
							thisEntry := data.AdventureViewGridMap[prevGrid]
							data.AdventureViewGridMap[grid] = thisEntry
							data.AdventureViewGridMap[prevGrid] = swapEntry
							// update array entries
							data.AdventureViewGridArr[0] = grid
							data.AdventureViewGridArr[swapEntry.ViewIndex] = prevGrid
							// update positions
							swapEl := pzlView.Elements[swapEntry.ViewIndex]
							element.Object.Pos, swapEl.Object.Pos = swapEl.Object.Pos, lastPos
						} else {
							gridU := grid
							gridD := grid
							gridL := grid
							gridR := grid
							gridU.Y++
							gridD.Y--
							gridL.X--
							gridR.X++
							_, okU := data.AdventureViewGridMap[gridU]
							_, okD := data.AdventureViewGridMap[gridD]
							_, okL := data.AdventureViewGridMap[gridL]
							_, okR := data.AdventureViewGridMap[gridR]
							if okU || okD || okL || okR { // if an empty slot and next to a puzzle, put it here
								// move grid map entry to new grid pos
								mapEntry := data.AdventureViewGridMap[data.AdventureViewGridArr[0]]
								data.AdventureViewGridMap[grid] = mapEntry
								delete(data.AdventureViewGridMap, data.AdventureViewGridArr[0])
								// update array entry
								data.AdventureViewGridArr[0] = grid
								// update positions
								element.Object.Pos = AdvPuzzleViewPos(grid)
							} else { // just return it to its original position
								element.Object.Pos = lastPos
							}
						}
					} else if click.Pressed() {
						pos := hvc.Pos.Add(offset)
						if drag || util.Magnitude(pos.Sub(lastPos)) > 3. {
							element.Object.Pos = pos
							drag = true
						}
					}
				} else if click.JustPressed() && hvc.Hover {
					drag = false
					clicked = true
					lastPos = element.Object.Pos
					offset = element.Object.Pos.Sub(hvc.Pos)
					var index int
					// highlight it and un-highlight every other puzzle
					for i, ele := range pzlView.Elements {
						if ele.Key != element.Key {
							ele.Border.Style = data.ThinBorderWhite
						} else {
							index = i
							ele.Border.Style = data.ThickBorderWhite
						}
					}
					if index != 0 { // put this element at the top of the list
						pzlView.Elements[0], pzlView.Elements[index] = pzlView.Elements[index], pzlView.Elements[0]
						swapGrid := data.AdventureViewGridArr[0]
						selectedGrid := data.AdventureViewGridArr[index]
						// swap places in the grid array
						data.AdventureViewGridArr[0] = selectedGrid
						data.AdventureViewGridArr[index] = swapGrid
						// change the ViewIndex for each element
						swapMap := data.AdventureViewGridMap[swapGrid]
						swapMap.ViewIndex = index
						data.AdventureViewGridMap[swapGrid] = swapMap
						selectedMap := data.AdventureViewGridMap[selectedGrid]
						selectedMap.ViewIndex = 0
						data.AdventureViewGridMap[selectedGrid] = selectedMap
					}
				}
				if hvc.Hover && !clicked {
					name := element.Key[:strings.LastIndex(element.Key, "_")]
					title := dlg.Get("puzzle_set_view_title")
					if title.Text.Raw != name {
						shadow := dlg.Get("puzzle_set_view_title_shadow")
						title.Text.SetText(name)
						shadow.Text.SetText(name)
						title.Text.Obj.Pos.X = title.Text.GetWidth() * -0.5
						shadow.Text.Obj.Pos.X = shadow.Text.GetWidth() * -0.5
					}
				}
			} else if clicked && !hvc.Input.Get("click").Pressed() {
				element.Border.Style = data.ThickBorderWhite
				clicked = false
				drag = false
				// return it to its original position
				element.Object.Pos = lastPos
			}
		}
	}))
}

func GetGridFromPzlView(x, y float64) world.Coords {
	w := 29.
	h := 17.
	gx := int(math.Round(x / w))
	gy := int(math.Round(y / h))
	return world.NewCoords(gx, gy)
}

func ClampAdventureView(vp *viewport.ViewPort) {
	x := 29.
	y := 17.
	maxX := float64(data.CurrPuzzleSet.GridMax.X) * x
	maxY := float64(data.CurrPuzzleSet.GridMax.Y) * y
	minX := float64(data.CurrPuzzleSet.GridMin.X) * x
	minY := float64(data.CurrPuzzleSet.GridMin.Y) * y
	if vp.CamPos.X > maxX {
		vp.CamPos.X = maxX
	} else if vp.CamPos.X < minX {
		vp.CamPos.X = minX
	}
	if vp.CamPos.Y > maxY {
		vp.CamPos.Y = maxY
	} else if vp.CamPos.Y < minY {
		vp.CamPos.Y = minY
	}
	gridX := int(math.Round(vp.CamPos.X / x))
	gridY := int(math.Round(vp.CamPos.Y / y))
	grid := world.NewCoords(gridX, gridY)
	if pzl := data.CurrPuzzleSet.GetGrid(grid); pzl > -1 {
		data.AdventureViewGridPos = grid
	}
}

func AdventureViewZoomOne(key string) {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		var keepTheseKeys []string
		//data.AdventureViewZoomLevel = 1
		data.PuzzleSetViewIsMoving = false
		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
		pzlCoords := data.AdventureViewGridPos
		vi := 0
		for grid, i := range data.CurrPuzzleSet.PuzzGrid {
			pzl := data.CurrPuzzleSet.Puzzles[i]
			col := pzl.Metadata.PrimaryColor

			// positions
			pos1 := AdvPuzzleViewPos(grid)
			if pzlCoords == grid { // move camera here
				pzlView.ViewPort.CamPos = pos1
			}

			// medium, zoom level 1
			z1Key := fmt.Sprintf("%s_%d", pzl.Metadata.Name, i)
			pic := CreatePuzzlePreviewMedium(pzl)
			pSpr := pixel.NewSprite(pic, pic.Bounds())
			z1ec := ui.ElementConstructor{
				Key:         z1Key,
				Batch:       constants.UIBatch,
				Position:    pos1,
				ElementType: ui.SpriteElement,
			}
			z1Ele := ui.CreatePixelSpriteElement(z1ec, pSpr)
			z1Ele.Object.Mask = col
			z1Ele.Border = &data.Border{
				Rect:   pixel.R(0, 0, z1Ele.Object.Rect.W(), z1Ele.Object.Rect.H()),
				Style:  data.ThinBorderWhite,
				Hidden: false,
			}
			//z1Ele.Object.Hidden = data.AdventureViewZoomLevel != 1
			pzlView.Elements = append(pzlView.Elements, z1Ele)
			keepTheseKeys = append(keepTheseKeys, z1Key)

			data.AdventureViewGridMap[grid] = data.AdvViewPzl{
				SetIndex:  i,
				ViewIndex: vi,
			}
			data.AdventureViewGridArr[vi] = grid
			vi++
		}
		//ResetAdventureViewElements(key, keepTheseKeys)
	}
}

func AdvPuzzleViewPos(grid world.Coords) pixel.Vec {
	pos := pixel.V(float64(grid.X*29)+0.5, float64(grid.Y*17)+0.5)
	if grid.X < 0 {
		pos.X += 1.
	}
	if grid.Y < 0 {
		pos.Y += 1.
	}
	return pos
}

//func ClampAdventureView(vp *viewport.ViewPort) {
//	var maxX, maxY, minX, minY float64
//	var x, y float64
//	switch data.AdventureViewZoomLevel {
//	case 0:
//		x = 8.
//		y = 5.
//	case 1:
//		x = 29.
//		y = 17.
//	case 2:
//		ClampZoomTwo(vp)
//		return
//	}
//	maxX = float64(data.CurrPuzzleSet.GridMax.X) * x
//	maxY = float64(data.CurrPuzzleSet.GridMax.Y) * y
//	minX = float64(data.CurrPuzzleSet.GridMin.X) * x
//	minY = float64(data.CurrPuzzleSet.GridMin.Y) * y
//	if vp.CamPos.X > maxX {
//		vp.CamPos.X = maxX
//	} else if vp.CamPos.X < minX {
//		vp.CamPos.X = minX
//	}
//	if vp.CamPos.Y > maxY {
//		vp.CamPos.Y = maxY
//	} else if vp.CamPos.Y < minY {
//		vp.CamPos.Y = minY
//	}
//	gridX := int(math.Round(vp.CamPos.X / x))
//	gridY := int(math.Round(vp.CamPos.Y / y))
//	grid := world.NewCoords(gridX, gridY)
//	if pzl := data.CurrPuzzleSet.GetGrid(grid); pzl > -1 {
//		data.AdventureViewGridPos = grid
//	}
//}

//func ClampZoomTwo(vp *viewport.ViewPort) {
//
//}

//func AdventureViewZoomZero(key string, rearrange bool) {
//	if data.Editor != nil && data.CurrPuzzleSet != nil {
//		dialog := ui.Dialogs[key]
//		pzlView := dialog.Get("puzzle_set_view")
//		var keepTheseKeys []string
//		data.AdventureViewZoomLevel = 0
//		data.PuzzleSetViewIsMoving = false
//		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
//		pzlCoords := data.AdventureViewGridPos
//		for grid, i := range data.CurrPuzzleSet.PuzzGrid {
//			pzl := data.CurrPuzzleSet.Puzzles[i]
//			col := pzl.Metadata.PrimaryColor
//
//			// positions
//			pos0 := pixel.V(float64(grid.X*8), float64(grid.Y*5))
//			if grid.X < 0 {
//				pos0.X += 1.
//			}
//			if grid.Y < 0 {
//				pos0.Y += 1.
//			}
//			if pzlCoords == grid { // move camera here
//				pzlView.ViewPort.CamPos = pos0
//			}
//
//			// small, zoom level 0
//			z0Key := fmt.Sprintf("z0_%d_%d", grid.X, grid.Y)
//			z0ec := ui.ElementConstructor{
//				Key:         z0Key,
//				SprKey:      "white_level_preview",
//				Batch:       constants.UIBatch,
//				Position:    pos0,
//				ElementType: ui.SpriteElement,
//			}
//			z0Ele := ui.CreateSpriteElement(z0ec)
//			z0Ele.Object.Mask = col
//			pzlView.Elements = append(pzlView.Elements, z0Ele)
//			keepTheseKeys = append(keepTheseKeys, z0Key)
//			ResetAdventureViewElements(key, keepTheseKeys, 0)
//		}
//	}
//}

//func AdventureViewZoomTwo(key string) {
//	if data.Editor != nil && data.CurrPuzzleSet != nil {
//		dialog := ui.Dialogs[key]
//		pzlView := dialog.Get("puzzle_set_view")
//		data.AdventureViewZoomLevel = 2
//		data.PuzzleSetViewIsMoving = false
//		data.PuzzleSetViewIndex = data.CurrPuzzleSet.PuzzleIndex
//		pzlCoords := data.AdventureViewGridPos
//
//		CreatePuzzlePreview(pzlView.Get("puzzle_center"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X, pzlCoords.Y)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_left"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X-1, pzlCoords.Y)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_right"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X+1, pzlCoords.Y)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_top_center"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X, pzlCoords.Y+1)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_top_left"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X-1, pzlCoords.Y+1)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_top_right"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X+1, pzlCoords.Y+1)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_bot_center"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X, pzlCoords.Y-1)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_bot_left"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X-1, pzlCoords.Y-1)))
//		CreatePuzzlePreview(pzlView.Get("puzzle_bot_right"), data.CurrPuzzleSet.GetGrid(world.NewCoords(pzlCoords.X+1, pzlCoords.Y-1)))
//
//		pzlView.ViewPort.CamPos = pzlView.Get("puzzle_center").Object.Pos
//
//		ResetAdventureViewElements(key, nil, 2)
//	}
//}

func ResetAdventureViewElements(key string, keepTheseKeys []string) {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
		dialog := ui.Dialogs[key]
		pzlView := dialog.Get("puzzle_set_view")
		for i := len(pzlView.Elements) - 1; i >= 0; i-- {
			ele := pzlView.Elements[i]
			// remove extra elements
			//if util.ContainsStr(ele.Key, keepTheseKeys1) {
			//	if zoom != 2 { // dispose of puzzle view
			//		ui.DisposeSubElements(ele.Elements)
			//		ele.Object.Hidden = true
			//	}
			//} else
			if !util.ContainsStr(ele.Key, keepTheseKeys) {
				if i < len(pzlView.Elements)+1 {
					pzlView.Elements = append(pzlView.Elements[:i], pzlView.Elements[i+1:]...)
				} else if i == 0 {
					pzlView.Elements = []*ui.Element{}
				} else {
					pzlView.Elements = pzlView.Elements[:i]
				}
				ui.DisposeSubElements(ele.Elements)
				myecs.Manager.DisposeEntity(ele.Entity)
			}
		}
	}
}
