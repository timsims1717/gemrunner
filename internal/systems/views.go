package systems

import (
	"gemrunner/internal/data"
	"gemrunner/pkg/viewport"
	"github.com/faiface/pixel"
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
	if data.EditorPanel != nil {
		data.EditorPanel.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
		data.EditorPanel.ViewPort.PortPos = pixel.V((viewport.MainCamera.Rect.W()-(data.BorderView.Rect.W()-data.EditorPanel.ViewPort.Rect.W())*pickedRatio)*0.5, (viewport.MainCamera.Rect.H()+(data.BorderView.Rect.H()-data.EditorPanel.ViewPort.Rect.H())*pickedRatio)*0.5)
	}
	for _, dialog := range data.Dialogs {
		dialog.BorderVP.PortPos = viewport.MainCamera.PostCamPos
		dialog.ViewPort.PortPos = viewport.MainCamera.PostCamPos
		dialog.BorderVP.PortSize = pixel.V(pickedRatio, pickedRatio)
		dialog.ViewPort.PortSize = pixel.V(pickedRatio, pickedRatio)
	}
}
