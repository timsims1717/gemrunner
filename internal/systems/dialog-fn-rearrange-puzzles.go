package systems

import (
	"fmt"
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/load"
	"gemrunner/internal/myecs"
	"gemrunner/internal/ui"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/object"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

// rearrange puzzles

func OpenRearrangePuzzlesDialog() {
	if data.Editor != nil && data.CurrPuzzleSet != nil {
		ui.NewDialog(load.RearrangePuzzleSetConstructor)
		rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
		for _, ele := range rearrangePzl.Elements {
			switch ele.Key {
			case "confirm_rearrange_puzzle":
				ele.OnClick = ConfirmRearrangedPuzzles
			case "rearrange_next":
				ele.OnHold = RearrangeNextPuzzle
				ele.OnClick = RearrangeNextPuzzle
			case "rearrange_prev":
				ele.OnHold = RearrangePrevPuzzle
				ele.OnClick = RearrangePrevPuzzle
			case "rearrange_swap_next":
				ele.OnHold = RearrangeSwapNextPuzzle
				ele.OnClick = RearrangeSwapNextPuzzle
			case "rearrange_swap_prev":
				ele.OnHold = RearrangeSwapPrevPuzzle
				ele.OnClick = RearrangeSwapPrevPuzzle
			case "rearrange_end":
				ele.OnHold = RearrangeSwapEndPuzzle
				ele.OnClick = RearrangeSwapEndPuzzle
			case "rearrange_begin":
				ele.OnHold = RearrangeSwapBeginPuzzle
				ele.OnClick = RearrangeSwapBeginPuzzle
			case "cancel_rearrange_puzzle":
				ele.OnClick = DisposeDialog(constants.DialogRearrangePuzzleSet)
			}
		}
		UpdateDialogView(rearrangePzl)
		data.RearrangePuzzleMove = false
		data.RearrangePuzzleIndex = data.CurrPuzzleSet.PuzzleIndex
		data.RearrangePuzzles = make([]int, len(data.CurrPuzzleSet.Puzzles))
		for i := range data.RearrangePuzzles {
			data.RearrangePuzzles[i] = i
		}
		pzlView := rearrangePzl.Get("rearrange_puzzle_view")
		CreatePuzzlePreview(pzlView.Get("puzzle_center"), data.RearrangePuzzleIndex)
		CreatePuzzlePreview(pzlView.Get("puzzle_left"), data.RearrangePuzzleIndex-1)
		CreatePuzzlePreview(pzlView.Get("puzzle_right"), data.RearrangePuzzleIndex+1)
		ResetRearrangePuzzleView()
		RearrangeSetNameAndNum(data.RearrangePuzzleIndex)
		ui.OpenDialogInStack(constants.DialogRearrangePuzzleSet)
	}
}

func ConfirmRearrangedPuzzles() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		// rearrange puzzle set
		var newPuzzles []*data.Puzzle
		for _, i := range data.RearrangePuzzles {
			newPuzzles = append(newPuzzles, data.CurrPuzzleSet.Puzzles[i])
		}
		data.CurrPuzzleSet.Puzzles = newPuzzles
		// go to currently selected puzzle
		data.CurrPuzzleSet.SetTo(data.RearrangePuzzleIndex)
		PuzzleInit()

		ui.Dispose(constants.DialogRearrangePuzzleSet)
	}
}

func RearrangeNextPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-1 { // we're at the end
			RearrangeEndBump()
		} else { // we're not at the end
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 8
			RearrangeSetNameAndNum(data.RearrangePuzzleIndex + 1)
			for i, ele := range pzlView.Elements {
				if ele.ElementType == ui.ContainerElement {
					if ele.Key == "puzzle_float" {
						FillNext(ele)
					}
					interA := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.ViewPort.PortPos.X).
						SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
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
							data.RearrangePuzzleIndex++
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangePrevPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == 0 { // we're at the beginning
			RearrangeBeginBump()
		} else { // we're not at the beginning
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 8
			RearrangeSetNameAndNum(data.RearrangePuzzleIndex - 1)
			for i, ele := range pzlView.Elements {
				if ele.ElementType == ui.ContainerElement {
					if ele.Key == "puzzle_float" {
						FillPrev(ele)
					}
					interA := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.ViewPort.PortPos.X).
						SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
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
							data.RearrangePuzzleIndex--
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangeSwapNextPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-1 { // we're at the end
			RearrangeEndBump()
		} else { // we're not at the end
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 8
			RearrangeSetNum(data.RearrangePuzzleIndex + 1)
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
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
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
							data.RearrangePuzzles[data.RearrangePuzzleIndex], data.RearrangePuzzles[data.RearrangePuzzleIndex+1] = data.RearrangePuzzles[data.RearrangePuzzleIndex+1], data.RearrangePuzzles[data.RearrangePuzzleIndex]
							data.RearrangePuzzleIndex++
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangeSwapPrevPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == 0 { // we're at the beginning
			RearrangeBeginBump()
		} else { // we're not at the beginning
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 8
			RearrangeSetNum(data.RearrangePuzzleIndex - 1)
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
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
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
							data.RearrangePuzzles[data.RearrangePuzzleIndex], data.RearrangePuzzles[data.RearrangePuzzleIndex-1] = data.RearrangePuzzles[data.RearrangePuzzleIndex-1], data.RearrangePuzzles[data.RearrangePuzzleIndex]
							data.RearrangePuzzleIndex--
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangeSwapEndPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-1 { // we're at the end
			RearrangeEndBump()
		} else if data.RearrangePuzzleIndex == len(data.CurrPuzzleSet.Puzzles)-2 { // we're one from the end
			RearrangeSwapNextPuzzle()
		} else { // we're not at the end
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 24
			RearrangeSetNum(len(data.RearrangePuzzles) - 1)
			for i, ele := range pzlView.Elements {
				if ele.ElementType == ui.ContainerElement {
					switch ele.Key {
					case "puzzle_center":
						continue
					case "puzzle_float":
						floatX := world.TileSize * 16
						ele.ViewPort.PortPos.X = floatX
						ele.BorderVP.PortPos.X = floatX
						ele.Object.Pos.X = floatX
						CreatePuzzlePreview(ele, data.RearrangePuzzles[len(data.RearrangePuzzles)-1])
					}
					interA := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.ViewPort.PortPos.X).
						SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X-move, constants.RearrangeMoveDur, ease.OutCubic)
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
							tmp := data.RearrangePuzzles[data.RearrangePuzzleIndex]
							data.RearrangePuzzles = append(data.RearrangePuzzles[:data.RearrangePuzzleIndex], data.RearrangePuzzles[data.RearrangePuzzleIndex+1:]...)
							data.RearrangePuzzles = append(data.RearrangePuzzles, tmp)
							data.RearrangePuzzleIndex = len(data.RearrangePuzzles) - 1
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangeSwapBeginPuzzle() {
	if data.Editor != nil && data.CurrPuzzleSet != nil && !data.RearrangePuzzleMove {
		if data.RearrangePuzzleIndex == 0 { // we're at the beginning
			RearrangeBeginBump()
		} else if data.RearrangePuzzleIndex == 1 { // we're one from the beginning
			RearrangeSwapPrevPuzzle()
		} else { // we're not at the beginning
			data.RearrangePuzzleMove = true
			rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
			pzlView := rearrangePzl.Get("rearrange_puzzle_view")
			move := world.TileSize * 24
			RearrangeSetNum(0)
			for i, ele := range pzlView.Elements {
				if ele.ElementType == ui.ContainerElement {
					switch ele.Key {
					case "puzzle_center":
						continue
					case "puzzle_float":
						floatX := world.TileSize * -16
						ele.ViewPort.PortPos.X = floatX
						ele.BorderVP.PortPos.X = floatX
						ele.Object.Pos.X = floatX
						CreatePuzzlePreview(ele, data.RearrangePuzzles[0])
					}
					interA := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.ViewPort.PortPos.X).
						SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
					interB := object.NewInterpolation(object.InterpolateCustom).
						SetValue(&ele.BorderVP.PortPos.X).
						SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X+move, constants.RearrangeMoveDur, ease.OutCubic)
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
							tmp := []int{data.RearrangePuzzles[data.RearrangePuzzleIndex]}
							tmp = append(tmp, append(data.RearrangePuzzles[:data.RearrangePuzzleIndex], data.RearrangePuzzles[data.RearrangePuzzleIndex+1:]...)...)
							data.RearrangePuzzles = tmp
							data.RearrangePuzzleIndex = 0
							ResetRearrangePuzzleView()
							data.RearrangePuzzleMove = false
						})
					}
					ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
				}
			}
		}
	}
}

func RearrangeBeginBump() {
	data.RearrangePuzzleMove = true
	rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
	pzlView := rearrangePzl.Get("rearrange_puzzle_view")
	move := 3.
	for i, ele := range pzlView.Elements {
		if ele.ElementType == ui.ContainerElement {
			interA := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.ViewPort.PortPos.X).
				SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.ViewPort.PortPos.X+move, ele.ViewPort.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interB := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.BorderVP.PortPos.X).
				SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.BorderVP.PortPos.X+move, ele.BorderVP.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interC := object.NewInterpolation(object.InterpolateX).
				SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.Object.Pos.X+move, ele.Object.Pos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			if i == len(pzlView.Elements)-1 {
				interC.SetOnComplete(func() {
					ResetRearrangePuzzleView()
					data.RearrangePuzzleMove = false
				})
			}
			ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
		}
	}
}

func RearrangeEndBump() {
	data.RearrangePuzzleMove = true
	rearrangePzl := ui.Dialogs[constants.DialogRearrangePuzzleSet]
	pzlView := rearrangePzl.Get("rearrange_puzzle_view")
	move := -3.
	for i, ele := range pzlView.Elements {
		if ele.ElementType == ui.ContainerElement {
			interA := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.ViewPort.PortPos.X).
				SetGween(ele.ViewPort.PortPos.X, ele.ViewPort.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.ViewPort.PortPos.X+move, ele.ViewPort.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interB := object.NewInterpolation(object.InterpolateCustom).
				SetValue(&ele.BorderVP.PortPos.X).
				SetGween(ele.BorderVP.PortPos.X, ele.BorderVP.PortPos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.BorderVP.PortPos.X+move, ele.BorderVP.PortPos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			interC := object.NewInterpolation(object.InterpolateX).
				SetGween(ele.Object.Pos.X, ele.Object.Pos.X+move, constants.RearrangeMoveDur*0.5, ease.Linear).
				AddGween(ele.Object.Pos.X+move, ele.Object.Pos.X, constants.RearrangeMoveDur*0.5, ease.OutCubic)
			if i == len(pzlView.Elements)-1 {
				interC.SetOnComplete(func() {
					ResetRearrangePuzzleView()
					data.RearrangePuzzleMove = false
				})
			}
			ele.Entity.AddComponent(myecs.Interpolation, []*object.Interpolation{interA, interB, interC})
		}
	}
}

func FillNext(ele *ui.Element) {
	floatX := world.TileSize * 16
	ele.ViewPort.PortPos.X = floatX
	ele.BorderVP.PortPos.X = floatX
	ele.Object.Pos.X = floatX
	if data.RearrangePuzzleIndex+2 >= len(data.RearrangePuzzles) {
		CreatePuzzlePreview(ele, -1)
	} else {
		CreatePuzzlePreview(ele, data.RearrangePuzzles[data.RearrangePuzzleIndex+2])
	}
}

func FillPrev(ele *ui.Element) {
	floatX := world.TileSize * -16
	ele.ViewPort.PortPos.X = floatX
	ele.BorderVP.PortPos.X = floatX
	ele.Object.Pos.X = floatX
	if data.RearrangePuzzleIndex-2 < 0 {
		CreatePuzzlePreview(ele, -1)
	} else {
		CreatePuzzlePreview(ele, data.RearrangePuzzles[data.RearrangePuzzleIndex-2])
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
			case data.BlockKeyBlue,
				data.BlockKeyGreen,
				data.BlockKeyPurple,
				data.BlockKeyBrown,
				data.BlockKeyYellow,
				data.BlockKeyOrange,
				data.BlockKeyCyan,
				data.BlockKeyGray:
				key = constants.PreviewKey
			case data.BlockPlayer1, data.BlockPlayer2,
				data.BlockPlayer3, data.BlockPlayer4:
				key = constants.PreviewPlayer
			case data.BlockDemon, data.BlockFly:
				key = constants.PreviewEnemy
			case data.BlockClosedBlue,
				data.BlockClosedGreen,
				data.BlockClosedPurple,
				data.BlockClosedBrown,
				data.BlockClosedYellow,
				data.BlockClosedOrange,
				data.BlockClosedCyan,
				data.BlockClosedGray,
				data.BlockLockBlue,
				data.BlockLockGreen,
				data.BlockLockPurple,
				data.BlockLockBrown,
				data.BlockLockYellow,
				data.BlockLockOrange,
				data.BlockLockCyan,
				data.BlockLockGray:
				key = constants.PreviewDoor
			case data.BlockGemBlue,
				data.BlockGemGreen,
				data.BlockGemPurple,
				data.BlockGemBrown,
				data.BlockGemYellow,
				data.BlockGemOrange,
				data.BlockGemCyan,
				data.BlockGemGray:
				key = constants.PreviewGem
			case data.BlockBar:
				key = constants.PreviewBar
			case data.BlockBomb, data.BlockBombLit:
				key = constants.PreviewBomb
			case data.BlockJetpack:
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

func RearrangeSetNameAndNum(index int) {
	pzlView := ui.Dialogs[constants.DialogRearrangePuzzleSet].Get("rearrange_puzzle_view")
	// name and number
	num := pzlView.Get("rearrange_puzzle_num")
	name := pzlView.Get("rearrange_puzzle_name")
	num.Text.SetText(fmt.Sprintf("%04d", index+1))
	theName := data.CurrPuzzleSet.Puzzles[data.RearrangePuzzles[index]].Metadata.Name
	width := name.Text.Text.BoundsOf(theName).W() * name.Text.RelativeSize
	name.Text.Obj.Pos.X = width * -0.5
	name.Text.SetText(theName)
}

func RearrangeSetNum(index int) {
	pzlView := ui.Dialogs[constants.DialogRearrangePuzzleSet].Get("rearrange_puzzle_view")
	num := pzlView.Get("rearrange_puzzle_num")
	num.Text.SetText(fmt.Sprintf("%04d", index+1))
}

func ResetRearrangePuzzleView() {
	pzlView := ui.Dialogs[constants.DialogRearrangePuzzleSet].Get("rearrange_puzzle_view")
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
	center.BorderVP.PortPos.X = 0
	center.Object.Pos.X = 0
	left.ViewPort.PortPos.X = data.RearrangeLeftX
	left.BorderVP.PortPos.X = data.RearrangeLeftX
	left.Object.Pos.X = data.RearrangeLeftX
	right.ViewPort.PortPos.X = data.RearrangeRightX
	right.BorderVP.PortPos.X = data.RearrangeRightX
	right.Object.Pos.X = data.RearrangeRightX
	float.ViewPort.PortPos.X = data.RearrangeFloatX
	float.BorderVP.PortPos.X = data.RearrangeFloatX
	float.Object.Pos.X = data.RearrangeFloatX
	// hide previews if at the end
	left.Object.Hidden = data.RearrangePuzzleIndex == 0
	right.Object.Hidden = data.RearrangePuzzleIndex == len(data.RearrangePuzzles)-1
}
