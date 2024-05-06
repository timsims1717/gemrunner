package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func UpdateViews() {
	pickedRatio := 1.
	if data.PuzzleView != nil {
		data.PuzzleView.PortPos = viewport.MainCamera.PostCamPos
		data.PuzzleViewNoShader.PortPos = viewport.MainCamera.PostCamPos
		data.BorderView.PortPos = viewport.MainCamera.PostCamPos
		wRatio := viewport.MainCamera.Rect.W() / data.PuzzleView.Rect.W()
		hRatio := viewport.MainCamera.Rect.H() / data.PuzzleView.Rect.H()
		pickedRatio = wRatio
		if hRatio < wRatio {
			pickedRatio = hRatio
		}
		pickedRatio *= 0.8
		data.PuzzleView.PortSize = pixel.V(pickedRatio, pickedRatio)
		data.PuzzleViewNoShader.PortSize = pixel.V(pickedRatio, pickedRatio)
		data.BorderView.PortSize = pixel.V(pickedRatio, pickedRatio)
	}
	data.CursorObj.Sca = pixel.V(pickedRatio, pickedRatio)
	data.CursorObj.Offset = pixel.V(9, -9).Scaled(pickedRatio)
	for _, dialog := range data.Dialogs {
		posRatX := viewport.MainCamera.Rect.W() / constants.WinWidth
		posRatY := viewport.MainCamera.Rect.H() / constants.WinHeight
		nPos := pixel.V(dialog.Pos.X*posRatX, dialog.Pos.Y*posRatY)
		if !dialog.NoBorder {
			dialog.BorderVP.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
			dialog.BorderVP.PortSize = pixel.V(pickedRatio, pickedRatio)
		}
		dialog.ViewPort.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
		dialog.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
	}
	if data.Editor != nil {
		if data.Editor.PosTop {
			panel := data.Dialogs[constants.DialogEditorPanelTop]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (((panel.ViewPort.Canvas.Bounds().W() + data.Editor.BlockSelect.Canvas.Bounds().W()) * 0.5) - world.TileSize) * data.Editor.BlockSelect.PortSize.X
			data.Editor.BlockSelect.PortPos.Y -= ((data.Editor.BlockSelect.Canvas.Bounds().H() + world.TileSize) * 0.5) * data.Editor.BlockSelect.PortSize.Y
		} else {
			panel := data.Dialogs[constants.DialogEditorPanelLeft]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (data.Editor.BlockSelect.Canvas.Bounds().W()*0.5 + world.HalfSize) * data.Editor.BlockSelect.PortSize.X
			//data.Editor.BlockSelect.PortPos.Y -= (panel.ViewPort.Canvas.Bounds().H()*0.5 - constants.BlockSelectHeight*world.TileSize*0.5 - world.HalfSize) * data.Editor.BlockSelect.PortSize.Y
		}
	}
}
