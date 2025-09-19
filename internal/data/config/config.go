package config

type Configuration struct {
	Gameplay GameplayConf `toml:"gameplay"`
	Graphics GraphicsConf `toml:"graphics"`
	Audio    AudioConf    `toml:"audio"`
}

type GameplayConf struct {
	FrameRate    int  `toml:"framerate"`
	ShowTimer    bool `toml:"showTimer"`
	ScreenShake  bool `toml:"screenshake"`
	AlwaysRecord bool `toml:"alwaysRecord"`
}

type GraphicsConf struct {
	Scanlines      bool `toml:"scanlines"`
	ShaderDetail   int  `toml:"shaderDetail"`
	BilinearFilter bool `toml:"bilinearFilter"`
	VSync          bool `toml:"vsync"`
	Fullscreen     bool `toml:"fullscreen"`
	Resolution     int  `toml:"resolution"`
	SetColorMode   bool `toml:"setColorMode"`
}

type AudioConf struct {
	MusicVolume  int  `toml:"musicVol"`
	MusicOn      bool `toml:"musicOn"`
	SfxVolume    int  `toml:"sfxVol"`
	SfxOn        bool `toml:"sfxOn"`
	MasterVolume int  `toml:"masterVol"`
	MasterOn     bool `toml:"masterOn"`
	MuteUnfocus  bool `toml:"muteUnfocus"`
}

func (c Configuration) Copy() Configuration {
	return Configuration{
		Gameplay: c.Gameplay.Copy(),
		Graphics: c.Graphics.Copy(),
		Audio:    c.Audio.Copy(),
	}
}

func (gc GameplayConf) Copy() GameplayConf {
	return GameplayConf{
		FrameRate:    gc.FrameRate,
		ShowTimer:    gc.ShowTimer,
		ScreenShake:  gc.ScreenShake,
		AlwaysRecord: gc.AlwaysRecord,
	}
}

func (gc GraphicsConf) Copy() GraphicsConf {
	return GraphicsConf{
		Scanlines:      gc.Scanlines,
		ShaderDetail:   gc.ShaderDetail,
		BilinearFilter: gc.BilinearFilter,
		VSync:          gc.VSync,
		Fullscreen:     gc.Fullscreen,
		Resolution:     gc.Resolution,
		SetColorMode:   gc.SetColorMode,
	}
}

func (ac AudioConf) Copy() AudioConf {
	return AudioConf{
		MusicVolume:  ac.MusicVolume,
		MusicOn:      ac.MusicOn,
		SfxVolume:    ac.SfxVolume,
		SfxOn:        ac.SfxOn,
		MasterVolume: ac.MasterVolume,
		MasterOn:     ac.MasterOn,
		MuteUnfocus:  ac.MuteUnfocus,
	}
}
