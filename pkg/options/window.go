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
	Resolutions     = []pixel.Vec{
		pixel.V(1600, 900),
	}

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
		if FullScreen != fullscreen {
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
				viewport.MainCamera.SetRect(pixel.R(0, 0, x, y))
			} else {
				win.SetMonitor(nil)
				viewport.MainCamera.SetRect(pixel.R(0, 0, Resolutions[ResolutionIndex].X, Resolutions[ResolutionIndex].Y))
			}
			fullscreen = FullScreen
			Updated = true
		}
	}
}
