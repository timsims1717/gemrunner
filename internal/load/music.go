package load

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/sfx"
)

func Music() {
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/beach.wav", constants.TrackBeach)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/dark.wav", constants.TrackDark)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/desert.wav", constants.TrackDesert)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/fungus.wav", constants.TrackFungus)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/ice.wav", constants.TrackIce)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/jungle.wav", constants.TrackJungle)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/lava.wav", constants.TrackLava)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/mech.wav", constants.TrackMech)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/reef.wav", constants.TrackReef)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/urban.wav", constants.TrackUrban)

	sfx.MusicPlayer.NewStream("game", sfx.Repeat, 0., 2.)
}
