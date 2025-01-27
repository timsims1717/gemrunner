package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func PuzzleSetViewNextPuzzle(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if (!data.PuzzleSetViewAllowEnd && data.PuzzleSetViewIndex >= len(data.CurrPuzzleSet.Puzzles)-1) ||
				data.PuzzleSetViewIndex >= len(data.CurrPuzzleSet.Puzzles) { // we're at the end
				PuzzleSetViewEndBump(dlg)
			} else { // we're not at the end
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 8
				PuzzleSetViewNameAndNum(dlg, data.PuzzleSetViewIndex+1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						if ele.Key == "puzzle_float" {
							FillNext(ele)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								center := pzlView.Get("puzzle_center")
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range center.Elements {
									left.Elements = append(left.Elements, ie)
								}
								center.Elements = []*ui.Element{}
								for _, ie := range right.Elements {
									center.Elements = append(center.Elements, ie)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									right.Elements = append(right.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewIndex++
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewPrevPuzzle(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == 0 { // we're at the beginning
				PuzzleSetViewBeginBump(dlg)
			} else { // we're not at the beginning
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 8
				PuzzleSetViewNameAndNum(dlg, data.PuzzleSetViewIndex-1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						if ele.Key == "puzzle_float" {
							FillPrev(ele)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								center := pzlView.Get("puzzle_center")
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range center.Elements {
									right.Elements = append(right.Elements, ie)
								}
								center.Elements = []*ui.Element{}
								for _, ie := range left.Elements {
									center.Elements = append(center.Elements, ie)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									left.Elements = append(left.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewIndex--
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewGoToEndPuzzle(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if (!data.PuzzleSetViewAllowEnd && data.PuzzleSetViewIndex >= len(data.CurrPuzzleSet.Puzzles)-1) ||
				data.PuzzleSetViewIndex >= len(data.CurrPuzzleSet.Puzzles) { // we're at the end
				PuzzleSetViewEndBump(dlg)
			} else if data.PuzzleSetViewIndex == len(data.CurrPuzzleSet.Puzzles)-2 { // we're one away
				PuzzleSetViewNextPuzzle(dlg)()
			} else { // we're not at the end
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 16
				PuzzleSetViewNameAndNum(dlg, len(data.PuzzleSetViewPuzzles)-1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						if ele.Key == "puzzle_float" {
							FillPuzzle(ele, len(data.CurrPuzzleSet.Puzzles)-1, world.TileSize*16)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								center := pzlView.Get("puzzle_center")
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range center.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								center.Elements = []*ui.Element{}
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									center.Elements = append(center.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewIndex = len(data.CurrPuzzleSet.Puzzles) - 1
								CreatePuzzlePreview(left, data.PuzzleSetViewIndex-1)
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewGoToBeginPuzzle(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == 0 { // we're at the beginning
				PuzzleSetViewBeginBump(dlg)
			} else if data.PuzzleSetViewIndex == 1 { // we're one from the beginning
				PuzzleSetViewPrevPuzzle(dlg)()
			} else { // we're not at the beginning
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * -16
				PuzzleSetViewNameAndNum(dlg, 0)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						if ele.Key == "puzzle_float" {
							FillPuzzle(ele, 0, world.TileSize*-16)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								center := pzlView.Get("puzzle_center")
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range center.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								center.Elements = []*ui.Element{}
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									center.Elements = append(center.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewIndex = 0
								CreatePuzzlePreview(right, 1)
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewSwapNext(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == len(data.CurrPuzzleSet.Puzzles)-1 { // we're at the end
				PuzzleSetViewEndBump(dlg)
			} else { // we're not at the end
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 8
				PuzzleSetViewSetNum(dlg, data.PuzzleSetViewIndex+1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						switch ele.Key {
						case "puzzle_center":
							continue
						case "puzzle_left":
							move = world.TileSize * 8
						case "puzzle_right":
							move = world.TileSize * 16
						case "puzzle_float":
							move = world.TileSize * 8
							FillNext(ele)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range right.Elements {
									left.Elements = append(left.Elements, ie)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									right.Elements = append(right.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex+1] = data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex+1], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex]
								data.PuzzleSetViewIndex++
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewSwapPrev(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == 0 { // we're at the beginning
				PuzzleSetViewBeginBump(dlg)
			} else { // we're not at the beginning
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 8
				PuzzleSetViewSetNum(dlg, data.PuzzleSetViewIndex-1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						switch ele.Key {
						case "puzzle_center":
							continue
						case "puzzle_right":
							move = world.TileSize * 8
						case "puzzle_left":
							move = world.TileSize * 16
						case "puzzle_float":
							move = world.TileSize * 8
							FillPrev(ele)
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range left.Elements {
									right.Elements = append(right.Elements, ie)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									left.Elements = append(left.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex-1] = data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex-1], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex]
								data.PuzzleSetViewIndex--
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewSwapEnd(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == len(data.CurrPuzzleSet.Puzzles)-1 { // we're at the end
				PuzzleSetViewEndBump(dlg)
			} else if data.PuzzleSetViewIndex == len(data.CurrPuzzleSet.Puzzles)-2 { // we're one from the end
				PuzzleSetViewSwapNext(dlg)()
			} else { // we're not at the end
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 24
				PuzzleSetViewSetNum(dlg, len(data.PuzzleSetViewPuzzles)-1)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						switch ele.Key {
						case "puzzle_center":
							continue
						case "puzzle_float":
							floatX := world.TileSize * 16
							ele.ViewPort.PortPos.X = floatX
							ele.Object.Pos.X = floatX
							CreatePuzzlePreview(ele, data.PuzzleSetViewPuzzles[len(data.PuzzleSetViewPuzzles)-1])
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									left.Elements = append(left.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								tmp := data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex]
								data.PuzzleSetViewPuzzles = append(data.PuzzleSetViewPuzzles[:data.PuzzleSetViewIndex], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex+1:]...)
								data.PuzzleSetViewPuzzles = append(data.PuzzleSetViewPuzzles, tmp)
								data.PuzzleSetViewIndex = len(data.PuzzleSetViewPuzzles) - 1
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewSwapToBegin(dlg *ui.Dialog) func() {
	return func() {
		if data.Editor != nil && data.CurrPuzzleSet != nil && !data.PuzzleSetViewIsMoving {
			if data.PuzzleSetViewIndex == 0 { // we're at the beginning
				PuzzleSetViewBeginBump(dlg)
			} else if data.PuzzleSetViewIndex == 1 { // we're one from the beginning
				PuzzleSetViewSwapPrev(dlg)()
			} else { // we're not at the beginning
				data.PuzzleSetViewIsMoving = true
				pzlView := dlg.Get("puzzle_set_view")
				move := world.TileSize * 24
				PuzzleSetViewSetNum(dlg, 0)
				for i, ele := range pzlView.Elements {
					if ele.ElementType == ui.ContainerElement {
						switch ele.Key {
						case "puzzle_center":
							continue
						case "puzzle_float":
							floatX := world.TileSize * -16
							ele.ViewPort.PortPos.X = floatX
							ele.Object.Pos.X = floatX
							CreatePuzzlePreview(ele, data.PuzzleSetViewPuzzles[0])
						}
						interA := object.NewInterpolation(object.InterpolateCustom).
							SetValue(&ele.ViewPort.PortPos.X).
							SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						interC := object.NewInterpolation(object.InterpolateX).
							SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
						if i == len(pzlView.Elements)-1 {
							interC.SetOnComplete(func() {
								left := pzlView.Get("puzzle_left")
								right := pzlView.Get("puzzle_right")
								float := pzlView.Get("puzzle_float")
								for _, ie := range left.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								left.Elements = []*ui.Element{}
								for _, ie := range right.Elements {
									myecs.Manager.DisposeEntity(ie.Entity)
								}
								right.Elements = []*ui.Element{}
								for _, ie := range float.Elements {
									right.Elements = append(right.Elements, ie)
								}
								float.Elements = []*ui.Element{}
								tmp := []int{data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex]}
								tmp = append(tmp, append(data.PuzzleSetViewPuzzles[:data.PuzzleSetViewIndex], data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex+1:]...)...)
								data.PuzzleSetViewPuzzles = tmp
								data.PuzzleSetViewIndex = 0
								ResetPuzzleSetView(dlg)
								data.PuzzleSetViewIsMoving = false
							})
						}
						ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
					}
				}
			}
		}
	}
}

func PuzzleSetViewBeginBump(dlg *ui.Dialog) {
	data.PuzzleSetViewIsMoving = true
	pzlView := dlg.Get("puzzle_set_view")
	move := 3.
	for i, ele := range pzlView.Elements {
		if ele.ElementType == ui.ContainerElement {
			interA := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.ViewPort.PortPos.X).
				SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.ViewPort.PortPos.X+move, ele.ViewPort.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interC := object.NewInterpolation(object.InterpolateX).
				SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.Object.Pos.X+move, ele.Object.Pos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			if i == len(pzlView.Elements)-1 {
				interC.SetOnComplete(func() {
					ResetPuzzleSetView(dlg)
					data.PuzzleSetViewIsMoving = false
				})
			}
			ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
		}
	}
}

func PuzzleSetViewEndBump(dlg *ui.Dialog) {
	data.PuzzleSetViewIsMoving = true
	pzlView := dlg.Get("puzzle_set_view")
	move := -3.
	for i, ele := range pzlView.Elements {
		if ele.ElementType == ui.ContainerElement {
			interA := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.ViewPort.PortPos.X).
				SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.ViewPort.PortPos.X+move, ele.ViewPort.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interC := object.NewInterpolation(object.InterpolateX).
				SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.Object.Pos.X+move, ele.Object.Pos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			if i == len(pzlView.Elements)-1 {
				interC.SetOnComplete(func() {
					ResetPuzzleSetView(dlg)
					data.PuzzleSetViewIsMoving = false
				})
			}
			ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interC})
		}
	}
}

func FillNext(ele *ui.Element) {
	floatX := world.TileSize * 16
	ele.ViewPort.PortPos.X = floatX
	ele.Object.Pos.X = floatX
	if data.PuzzleSetViewIndex+2 >= len(data.PuzzleSetViewPuzzles) {
		CreatePuzzlePreview(ele, -1)
	} else {
		CreatePuzzlePreview(ele, data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex+2])
	}
}

func FillPrev(ele *ui.Element) {
	floatX := world.TileSize * -16
	ele.ViewPort.PortPos.X = floatX
	ele.Object.Pos.X = floatX
	if data.PuzzleSetViewIndex-2 < 0 {
		CreatePuzzlePreview(ele, -1)
	} else {
		CreatePuzzlePreview(ele, data.PuzzleSetViewPuzzles[data.PuzzleSetViewIndex-2])
	}
}

func FillPuzzle(ele *ui.Element, index int, floatX float64) {
	ele.ViewPort.PortPos.X = floatX
	ele.Object.Pos.X = floatX
	if index < 0 || index >= len(data.PuzzleSetViewPuzzles) {
		CreatePuzzlePreview(ele, -1)
	} else {
		CreatePuzzlePreview(ele, data.PuzzleSetViewPuzzles[index])
	}
}

func CreatePuzzlePreview(cnt *ui.Element, index int) {
	if index < 0 || index >= len(data.CurrPuzzleSet.Puzzles) {
		cnt.Object.Hidden = true
		return
	}
	if len(cnt.Elements) > 0 {
		for _, ie := range cnt.Elements {
			myecs.Manager.DisposeEntity(ie.Entity)
		}
		cnt.Elements = []*ui.Element{}
	}
	cnt.Object.Hidden = false
	//txt := ui.ElementConstructor{
	//	Key:         "new_float",
	//	Text:        data.CurrPuzzleSet.Puzzles[index].Metadata.Name,
	//	Position:    pixel.V(-50, 2),
	//	Color:       pixel.ToRGBA(constants.ColorWhite),
	//	ElementType: ui.TextElement,
	//}
	//cnt.Elements = []*ui.Element{ui.CreateTextElement(txt, cnt.ViewPort)}
	pzl := data.CurrPuzzleSet.Puzzles[index]
	for y, row := range pzl.Tiles.T {
		for x, tile := range row {
			posX := float64(x*4) - world.TileSize*3.5 + 2
			posY := float64(y*4) - world.TileSize*2 + 2
			pos := pixel.V(posX, posY)
			var key string
			switch tile.Block {
			case data.BlockTurf, data.BlockCracked, data.BlockFall, data.BlockSpike:
				key = constants.PreviewTurf
			case data.BlockBedrock, data.BlockPhase:
				key = constants.PreviewBedrock
			case data.BlockLadder, data.BlockLadderTurf,
				data.BlockLadderCracked, data.BlockLadderCrackedTurf:
				key = constants.PreviewLadder
			case data.BlockKey:
				key = constants.PreviewKey
			case data.BlockPlayer1, data.BlockPlayer2,
				data.BlockPlayer3, data.BlockPlayer4:
				key = constants.PreviewPlayer
			case data.BlockDemon, data.BlockFly:
				key = constants.PreviewEnemy
			case data.BlockDoorVisible, data.BlockDoorLocked:
				key = constants.PreviewDoor
			case data.BlockGem:
				key = constants.PreviewGem
			case data.BlockBar:
				key = constants.PreviewBar
			case data.BlockBomb, data.BlockBombLit:
				key = constants.PreviewBomb
			case data.BlockJetpack, data.BlockDisguise,
				data.BlockDrill, data.BlockFlamethrower,
				data.BlockJumpBoots, data.BlockBox:
				key = constants.PreviewTool
			}
			if key != "" {
				ec := ui.ElementConstructor{
					Key:         key,
					SprKey:      key,
					Batch:       constants.UIBatch,
					Position:    pos,
					ElementType: ui.SpriteElement,
				}
				cnt.Elements = append(cnt.Elements, ui.CreateSpriteElement(ec))
			}
		}
	}
}

func PuzzleSetViewNameAndNum(dlg *ui.Dialog, index int) {
	pzlView := dlg.Get("puzzle_set_view")
	// name and number
	num := pzlView.Get("puzzle_num")
	name := pzlView.Get("puzzle_name")
	num.Text.SetText(fmt.Sprintf("%04d", index+1))
	var theName string
	if index < len(data.PuzzleSetViewPuzzles) {
		theName = data.CurrPuzzleSet.Puzzles[data.PuzzleSetViewPuzzles[index]].Metadata.Name
	} else {
		theName = "End of Puzzle Set"
	}
	width := name.Text.Text.BoundsOf(theName).W() * name.Text.Scalar
	name.Text.Obj.Pos.X = width * -0.5
	name.Text.SetText(theName)
}

func PuzzleSetViewSetNum(dlg *ui.Dialog, index int) {
	pzlView := dlg.Get("puzzle_set_view")
	num := pzlView.Get("puzzle_num")
	num.Text.SetText(fmt.Sprintf("%04d", index+1))
}

func ResetPuzzleSetView(dlg *ui.Dialog) {
	pzlView := dlg.Get("puzzle_set_view")
	center := pzlView.Get("puzzle_center")
	left := pzlView.Get("puzzle_left")
	right := pzlView.Get("puzzle_right")
	float := pzlView.Get("puzzle_float")
	// remove all interpolation
	center.Entity.RemoveComponent(myecs.Interpolation)
	left.Entity.RemoveComponent(myecs.Interpolation)
	right.Entity.RemoveComponent(myecs.Interpolation)
	float.Entity.RemoveComponent(myecs.Interpolation)
	// set the previews to the right spot
	center.ViewPort.PortPos.X = 0
	center.Object.Pos.X = 0
	left.ViewPort.PortPos.X = data.RearrangeLeftX
	left.Object.Pos.X = data.RearrangeLeftX
	right.ViewPort.PortPos.X = data.RearrangeRightX
	right.Object.Pos.X = data.RearrangeRightX
	float.ViewPort.PortPos.X = data.RearrangeFloatX
	float.Object.Pos.X = data.RearrangeFloatX
	// hide previews if at the end
	left.Object.Hidden = data.PuzzleSetViewIndex <= 0
	center.Object.Hidden = data.PuzzleSetViewIndex >= len(data.PuzzleSetViewPuzzles)
	right.Object.Hidden = data.PuzzleSetViewIndex >= len(data.PuzzleSetViewPuzzles)-1
}
