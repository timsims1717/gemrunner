package systems

import (
	"gemrunner/internal/constants"
	"gemrunner/internal/data"
	"gemrunner/internal/ui"
	"gemrunner/pkg/options"
	"gemrunner/pkg/viewport"
	"gemrunner/pkg/world"
	"github.com/gopxl/pixel"
)

func UpdateViews() {
	data.ScreenView.SetRect(pixel.R(options.CurrResolution.X*-0.5, options.CurrResolution.Y*-0.5, options.CurrResolution.X*0.5, options.CurrResolution.Y*0.5))
	data.ScreenView.PortPos = viewport.MainCamera.CamPos
	data.ScreenView.CamPos = pixel.V(options.CurrResolution.X*0.5, options.CurrResolution.Y*0.5)

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

	if data.CurrentPlayArea != nil {
		UpdatePlayAreaView(data.CurrentPlayArea)
	}
	for _, fp := range data.OtherPlayAreas {
		UpdatePlayAreaView(fp)
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
	SetMainBorder(constants.PuzzleWidth, constants.PuzzleHeight)
	//SetMainBorder(int(width), int(height))
	//viewport.MainCamera.CamPos.Y -= world.TileSize * 2
}

func UpdatePlayAreaView(fp *data.PlayArea) {
	width := float64(constants.PuzzleWidth)
	height := float64(constants.PuzzleHeight)
	if fp.Level != nil {
		width = float64(fp.Level.Metadata.Width)
		height = float64(fp.Level.Metadata.Height)
	} else if fp.Puzzle != nil {
		width = float64(fp.Puzzle.Metadata.Width)
		height = float64(fp.Puzzle.Metadata.Height)
	}

	if fp.PuzzleView != nil {
		fp.PuzzleView.CamPos = pixel.V(world.TileSize*0.5*width, world.TileSize*0.5*height)
		fp.PuzzleView.SetRect(pixel.R(0, 0, world.TileSize*width, world.TileSize*height))
		//data.PuzzleView.PortPos = viewport.MainCamera.PostCamPos
		//data.PuzzleViewNoShader.PortPos = viewport.MainCamera.PostCamPos
		fp.PuzzleView.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	if fp.BorderView != nil {
		fp.BorderView.SetRect(pixel.R(0, 0, world.TileSize*(width+1), world.TileSize*(height+1)))
		fp.BorderView.CamPos = pixel.V(world.TileSize*0.5*width, world.TileSize*0.5*height)
		fp.BorderView.PortPos = viewport.MainCamera.PostCamPos
		fp.BorderView.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	if fp.PuzzleViewNoShader != nil {
		fp.PuzzleViewNoShader.CamPos = pixel.V(world.TileSize*0.5*width, world.TileSize*0.5*height)
		fp.PuzzleViewNoShader.SetRect(pixel.R(0, 0, world.TileSize*width, world.TileSize*height))
		fp.PuzzleViewNoShader.PortPos = viewport.MainCamera.PostCamPos
		fp.PuzzleViewNoShader.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	if fp.WorldView != nil {
		fp.WorldView.PortPos = viewport.MainCamera.PostCamPos
		xWidth := width * world.TileSize * constants.PickedRatio
		yHeight := height * world.TileSize * constants.PickedRatio
		fp.WorldView.SetRect(pixel.R(0, 0, xWidth, yHeight))
	}
	//viewport.MainCamera.CamPos.Y -= world.TileSize * 2
}

func UpdateDialogView(dialog *ui.Dialog) {
	//posRatX := data.ScreenView.Rect.W() / options.CurrResolution.X
	//posRatY := data.ScreenView.Rect.H() / options.CurrResolution.Y
	//nPos := pixel.V(dialog.Pos.X*posRatX, dialog.Pos.Y*posRatY)
	nPos := dialog.Pos.Scaled(constants.PickedRatio)
	if !dialog.NoBorder {
		dialog.BorderVP.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
		dialog.BorderVP.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
	}
	dialog.ViewPort.PortPos = viewport.MainCamera.PostCamPos.Add(nPos)
	dialog.ViewPort.PortSize = pixel.V(constants.PickedRatio, constants.PickedRatio)
}
