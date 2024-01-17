package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func UpdateViews() {
	var pickedRatio float64
	if data.PuzzleView != nil {
		data.PuzzleView.PortPos = viewport.MainCamera.PostCamPos
		data.BorderView.PortPos = viewport.MainCamera.PostCamPos
		wRatio := viewport.MainCamera.Rect.W() / data.PuzzleView.Rect.W()
		hRatio := viewport.MainCamera.Rect.H() / data.PuzzleView.Rect.H()
		pickedRatio = wRatio
		if hRatio < wRatio {
			pickedRatio = hRatio
		}
		pickedRatio *= 0.8
		data.PuzzleView.PortSize = pixel.V(pickedRatio, pickedRatio)
		data.BorderView.PortSize = pixel.V(pickedRatio, pickedRatio)
	}
	for _, dialog := range data.Dialogs {
		if !dialog.NoBorder {
			dialog.BorderVP.PortPos = viewport.MainCamera.PostCamPos.Add(dialog.Pos)
			dialog.BorderVP.PortSize = pixel.V(pickedRatio, pickedRatio)
		}
		dialog.ViewPort.PortPos = viewport.MainCamera.PostCamPos.Add(dialog.Pos)
		dialog.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
	}
	if data.Editor != nil {
		if data.Editor.PosTop {
			panel := data.Dialogs["editor_panel_top"]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (((panel.ViewPort.Canvas.Bounds().W() + data.Editor.BlockSelect.Canvas.Bounds().W()) * 0.5) - world.TileSize) * data.Editor.BlockSelect.PortSize.X
			data.Editor.BlockSelect.PortPos.Y -= ((data.Editor.BlockSelect.Canvas.Bounds().H() + world.TileSize) * 0.5) * data.Editor.BlockSelect.PortSize.Y
		} else {
			panel := data.Dialogs["editor_panel_left"]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (data.Editor.BlockSelect.Canvas.Bounds().W()*0.5 + world.TileSize*4) * data.Editor.BlockSelect.PortSize.X * 0.5
			data.Editor.BlockSelect.PortPos.Y += (panel.ViewPort.Canvas.Bounds().H()*0.5 + world.TileSize) * data.Editor.BlockSelect.PortSize.Y * 0.5
		}
	}
}
