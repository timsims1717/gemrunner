package options

import (
	"gemrunner/pkg/viewport"
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	Updated         bool
	Focused         bool
	VSync           bool
	FullScreen      bool
	BilinearFilter  bool
	ResolutionIndex int
	Resolutions     []pixel.Vec
	CurrResolution  pixel.Vec

	fullscreen bool
	resIndex   int
)

func RegisterResolution(res pixel.Vec) {
	Resolutions = append(Resolutions, res)
}

func WindowUpdate(win *pixelgl.Window) {
	Updated = false
	Focused = win.Focused()
	if win.Focused() {
		win.SetVSync(VSync)
		win.SetSmooth(BilinearFilter)
		if FullScreen != fullscreen || ResolutionIndex != resIndex {
			// get window position (center)
			pos := win.GetPos()
			pos.X += win.Bounds().W() * 0.5
			pos.Y += win.Bounds().H() * 0.5

			// find current monitor
			var picked *pixelgl.Monitor
			if len(pixelgl.Monitors()) > 1 {
				for _, m := range pixelgl.Monitors() {
					x, y := m.Position()
					w, h := m.Size()
					if pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h {
						picked = m
						break
					}
				}
				if picked == nil {
					pos = win.GetPos()
					for _, m := range pixelgl.Monitors() {
						x, y := m.Position()
						w, h := m.Size()
						if pos.X >= x && pos.X <= x+w && pos.Y >= y && pos.Y <= y+h {
							picked = m
							break
						}
					}
				}
			}
			if picked == nil {
				picked = pixelgl.PrimaryMonitor()
			}
			if FullScreen {
				win.SetMonitor(picked)
				x, y := picked.Size()
				CurrResolution = pixel.V(x, y)
			} else {
				win.SetMonitor(nil)
				res := pixel.V(1600, 900)
				if len(Resolutions) > ResolutionIndex {
					res = Resolutions[ResolutionIndex]
				}
				CurrResolution = pixel.V(res.X, res.Y)
			}
			r := pixel.R(0, 0, CurrResolution.X, CurrResolution.Y)
			win.SetBounds(r)
			viewport.MainCamera.SetRect(r)
			fullscreen = FullScreen
			resIndex = ResolutionIndex
			Updated = true
		}
	}
}
