package state

import "github.com/gopxl/pixel/pixelgl"

type loadingScreen struct {
	Init   func()
	Update func(*pixelgl.Window)
	Draw   func(*pixelgl.Window)
}

var LoadingScreen *loadingScreen
