package sfx

import (
	"fmt"
	gween "gemrunner/pkg/gween64"
	"gemrunner/pkg/gween64/ease"
	"gemrunner/pkg/timing"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"math"
)

type Mode int

const (
	Sequential = iota // play through tracks in order
	Random            // when current track ends, play a random other track (if possible)
	Repeat            // play the specified track again
	Single            // play only one track, then stop
)

// MusicStream holds the functions to play individual and sequences of
// audio tracks.
// To start playing music, create the MusicStream using the MusicPlayer,
// then call one of the following:
//
//	PlayTrack
//	SingleTrack
//	RepeatTrack
//
// To play multiple tracks, set the track list with SetTracks, then call
// Play.
type MusicStream struct {
	key    string
	tracks []string
	curr   string
	cId    int
	next   string

	mode     Mode
	playNext bool
	fade     float64
	vol      float64

	paused  bool
	stopped bool

	stream beep.StreamSeekCloser
	ctrl   *beep.Ctrl
	volume *effects.Volume
	interV *gween.Tween
	format beep.Format
}

// PlayTrack plays the requested track, but does not change the
// mode or change the tracks of the stream.
func (s *MusicStream) PlayTrack(track string) {
	if MusicPlayer.HasTrack(track) {
		s.pause(true)
		if track == "" {
			s.stopped = true
			return
		}
		s.next = track
		s.playNext = true
		s.stopped = false
	} else {
		fmt.Printf("MUSIC WARNING: track %s not registered", track)
	}
}

// SingleTrack plays the requested track and changes the stream's mode
// to Single.
func (s *MusicStream) SingleTrack(track string) {
	if MusicPlayer.HasTrack(track) {
		s.pause(true)
		s.next = track
		s.tracks = []string{}
		s.playNext = true
		s.mode = Single
		s.stopped = false
	} else {
		fmt.Printf("MUSIC WARNING: track %s not registered", track)
	}
}

// RepeatTrack plays the requested track and changes the stream's mode
// to Repeat.
func (s *MusicStream) RepeatTrack(track string) {
	if MusicPlayer.HasTrack(track) {
		s.pause(true)
		s.next = track
		s.tracks = []string{track}
		s.playNext = true
		s.mode = Repeat
		s.stopped = false
	} else {
		fmt.Printf("MUSIC WARNING: track %s not registered", track)
	}
}

// SetTracks sets the stream's track list. If the stream's mode
// is Single or Repeat, it is set to Sequential.
func (s *MusicStream) SetTracks(keys []string) {
	s.tracks = keys
	s.cId = 0
	s.next = ""
	if s.mode == Repeat || s.mode == Single {
		s.mode = Sequential
	}
}

// Play starts music if it is stopped, or unpauses it if
// it is paused.
func (s *MusicStream) Play() {
	if s.stopped {
		s.stopped = false
		s.playNext = true
	} else if s.paused {
		s.pause(false)
	}
}

// Resume unpauses the stream if it is paused.
func (s *MusicStream) Resume() {
	if !s.stopped {
		if s.paused {
			s.pause(false)
		}
	}
}

// Pause pauses the stream.
func (s *MusicStream) Pause() {
	if !s.stopped {
		s.pause(true)
	}
}

// NextTrack forces the stream to start playing the next track.
func (s *MusicStream) NextTrack() {
	if !s.stopped {
		if s.paused {
			s.pause(false)
		}
		s.playNext = true
	}
}

//func (s *MusicStream) chooseTrack(keys []string) {
//	if !s.stopped {
//		for _, k := range keys {
//			if k == s.curr {
//				return
//			}
//		}
//	}
//	s.stopped = false
//	s.next = keys[random.Intn(len(keys))]
//	s.playNext = true
//}

//func (s *MusicStream) setTrack(key string) {
//	s.next = key
//	s.playNext = s.next != s.curr
//	s.stopped = false
//}

func (s *MusicStream) pause(pause bool) {
	if pause && s.fade > 0. && s.volume != nil {
		s.interV = gween.New(s.volume.Volume, -8., s.fade, ease.Linear)
	} else if !pause && s.fade > 0. && s.volume != nil {
		s.interV = gween.New(s.volume.Volume, getMusicVolume()+s.vol, s.fade, ease.Linear)
	} else {
		s.interV = nil
	}
	s.paused = pause
}

func (s *MusicStream) Stop() {
	s.interV = nil
	s.paused = true
	s.stopped = true
	s.next = ""
	s.playNext = false
}

func (s *MusicStream) SetVolume(vol float64) {
	s.vol = vol
	if s.interV != nil && s.fade > 0. && s.volume != nil {
		s.interV = gween.New(s.volume.Volume, getMusicVolume()+s.vol, s.fade, ease.Linear)
	}
}

func (s *MusicStream) SetFade(fade float64) {
	s.fade = fade
}

func (s *MusicStream) update() {
	if s.playNext && len(s.tracks) > 0 {
		if !MusicPlayer.loading &&
			(s.ctrl == nil ||
				s.volume == nil ||
				s.ctrl.Paused ||
				s.volume.Silent ||
				s.mode == Repeat) {
			if s.next == "" {
				if len(s.tracks) > 1 {
					switch s.mode {
					case Random:
						t := -1
						for s.next == "" || s.next == s.curr {
							t = random.Intn(len(s.tracks))
							s.next = s.tracks[t]
						}
						s.cId = t
					case Sequential:
						s.cId++
						s.cId %= len(s.tracks)
						s.next = s.tracks[s.cId]
					}
				} else {
					s.next = s.tracks[0]
				}
			}
			go MusicPlayer.loadTrack(s)
		}
		if !s.paused {
			s.pause(true)
		}
	}
	if s.volume != nil {
		speaker.Lock()
		if !s.stopped && s.interV != nil {
			v, fin := s.interV.Update(timing.DT)
			if fin {
				s.volume.Silent = s.paused || getMusicMuted()
				s.volume.Volume = getMusicVolume() + s.vol
				s.ctrl.Paused = s.paused
				s.interV = nil
			} else {
				s.volume.Volume = math.Min(v, getMusicVolume()+s.vol)
				s.volume.Silent = getMusicMuted()
				s.ctrl.Paused = false
			}
		} else {
			if s.stopped || s.paused {
				s.volume.Volume = -8.
			} else {
				s.volume.Volume = getMusicVolume() + s.vol
			}
			s.volume.Silent = s.stopped || s.paused || getMusicMuted()
			s.ctrl.Paused = s.paused
		}
		speaker.Unlock()
	}
}
