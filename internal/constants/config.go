package constants

import (
	"gemrunner/internal/data/config"
)

var DefaultConfiguration = config.Configuration{
	Gameplay: config.GameplayConf{
		FrameRate:    30,
		ShowTimer:    false,
		ScreenShake:  true,
		AlwaysRecord: false,
	},
	Graphics: config.GraphicsConf{
		Scanlines:      true,
		ShaderDetail:   2,
		BilinearFilter: true,
		VSync:          true,
		Fullscreen:     true,
		Resolution:     0,
		SetColorMode:   false,
	},
	Audio: config.AudioConf{
		MusicVolume:  100,
		MusicOn:      true,
		SfxVolume:    80,
		SfxOn:        true,
		MasterVolume: 100,
		MasterOn:     true,
		MuteUnfocus:  true,
	},
}
