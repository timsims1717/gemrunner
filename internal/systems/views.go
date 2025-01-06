package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func UpdateViews() {
	wRatio := viewport.MainCamera.Rect.W() / (constants.PuzzleWidth * world.TileSize)
	hRatio := viewport.MainCamera.Rect.H() / (constants.PuzzleHeight * world.TileSize)
	maxRatio := wRatio
	if hRatio < wRatio {
		maxRatio = hRatio
	}
	maxRatio *= constants.ScreenRatioLimit

	constants.PickedRatio = 1.
	for constants.PickedRatio+1 < maxRatio {
		constants.PickedRatio += 1
	}

	if data.PuzzleView != nil {
		data.PuzzleView.PortPos = viewport.MainCamera.PostCamPos
		data.PuzzleViewNoShader.PortPos = viewport.MainCamera.PostCamPos
		data.BorderView.PortPos = viewport.MainCamera.PostCamPos
		data.PuzzleView.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
		data.PuzzleViewNoShader.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
		data.BorderView.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	if data.WorldView != nil {
		data.WorldView.PortPos = viewport.MainCamera.PostCamPos
		xWidth := constants.PuzzleWidth * world.TileSize * constants.PickedRatio
		yHeight := constants.PuzzleHeight * world.TileSize * constants.PickedRatio
		data.WorldView.SetRect(pixel.R(0, 0, xWidth, yHeight))
	}
	data.CursorObj.Sca = pixel.V(constants.PickedRatio, constants.PickedRatio)
	data.CursorObj.Offset = pixel.V(9, -9).Scaled(constants.PickedRatio)
	for _, dialog := range ui.Dialogs {
		UpdateDialogView(dialog)
	}
	if data.Editor != nil {
		if data.Editor.PosTop {
			panel := ui.Dialogs[constants.DialogEditorPanelTop]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (((panel.ViewPort.Canvas.Bounds().W() + data.Editor.BlockSelect.Canvas.Bounds().W()) * 0.5) - world.TileSize) * data.Editor.BlockSelect.PortSize.X
			data.Editor.BlockSelect.PortPos.Y -= ((data.Editor.BlockSelect.Canvas.Bounds().H() + world.TileSize) * 0.5) * data.Editor.BlockSelect.PortSize.Y
		} else {
			panel := ui.Dialogs[constants.DialogEditorPanelLeft]
			data.Editor.BlockSelect.PortPos = panel.ViewPort.PortPos
			data.Editor.BlockSelect.PortPos.X += (data.Editor.BlockSelect.Canvas.Bounds().W()*0.5 + world.HalfSize) * data.Editor.BlockSelect.PortSize.X
			//data.Editor.BlockSelect.PortPos.Y -= (panel.ViewPort.Canvas.Bounds().H()*0.5 - constants.BlockSelectHeight*world.TileSize*0.5 - world.HalfSize) * data.Editor.BlockSelect.PortSize.Y
		}
	}
	//viewport.MainCamera.CamPos.Y -= world.TileSize * 2
}

func UpdateDialogView(dialog *ui.Dialog) {
	posRatX := viewport.MainCamera.Rect.W() / constants.WinWidth
	posRatY := viewport.MainCamera.Rect.H() / constants.WinHeight
	nPos := pixel.V(dialog.Pos.X*posRatX, dialog.Pos.Y*posRatY)
	if !dialog.NoBorder {
		dialog.BorderVP.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
		dialog.BorderVP.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	dialog.ViewPort.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
	dialog.ViewPort.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
}
