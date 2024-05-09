package load

import (
	"gemrunner/internal/constants"
	"gemrunner/pkg/sfx"
)

func Music() {
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/beach.ogg", constants.TrackBeach)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/dark.ogg", constants.TrackDark)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/desert.ogg", constants.TrackDesert)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/fungus.ogg", constants.TrackFungus)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/gilded.ogg", constants.TrackGilded)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/ice.ogg", constants.TrackIce)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/jungle.ogg", constants.TrackJungle)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/lava.ogg", constants.TrackLava)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/mech.ogg", constants.TrackMech)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/reef.ogg", constants.TrackReef)
	sfx.MusicPlayer.RegisterMusicTrack("assets/music/urban.ogg", constants.TrackUrban)

	sfx.MusicPlayer.NewStream("game", sfx.Repeat, 0., 2.)
}
