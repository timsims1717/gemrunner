package sfx

import (
	"fmt"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/pkg/errors"
)

var MusicPlayer *musicPlayer

type musicPlayer struct {
	tracks  map[string]string
	streams map[string]*MusicStream
	loading bool
}

func init() {
	MusicPlayer = &musicPlayer{
		tracks:  make(map[string]string),
		streams: make(map[string]*MusicStream),
	}
}

func (p *musicPlayer) Update() {
	for _, s := range p.streams {
		s.update()
	}
}

func (p *musicPlayer) RegisterMusicTrack(path, key string) {
	p.tracks[key] = path
}

func (p *musicPlayer) NewStream(key string, mode Mode, vol, fade float64) {
	p.streams[key] = &MusicStream{
		key:  key,
		mode: mode,
		fade: fade,
		vol:  vol,
	}
}

func (p *musicPlayer) GetStream(key string) *MusicStream {
	s, ok := p.streams[key]
	if !ok {
		panic(fmt.Sprintf("fatal music player error: no stream '%s'", key))
	}
	return s
}

func (p *musicPlayer) HasStream(key string) bool {
	_, ok := p.streams[key]
	return ok
}

func (p *musicPlayer) loadTrack(set *MusicStream) {
	p.loading = true
	if err := p.loadTrackInner(set); err != nil {
		fmt.Printf("music player error: %s\n", err)
	} else {
		set.playNext = false
	}
	p.loading = false
}

func (p *musicPlayer) loadTrackInner(set *MusicStream) error {
	errMsg := fmt.Sprintf("load track %s", set.next)
	if path, ok := p.tracks[set.next]; ok {
		streamer, format, err := loadSoundFile(path)
		if err != nil {
			return errors.Wrap(err, errMsg)
		}
		speaker.Lock()
		if set.stream != nil {
			err = set.stream.Close()
			if err != nil {
				fmt.Println(errors.Wrap(err, errMsg))
			}
		}
		if set.ctrl != nil {
			set.ctrl.Paused = true
		}
		if set.volume != nil {
			set.volume.Silent = true
		}
		set.stream = streamer
		set.ctrl = &beep.Ctrl{
			Streamer: set.stream,
			Paused:   false,
		}
		set.volume = &effects.Volume{
			Streamer: set.ctrl,
			Base:     2,
			Volume:   getMusicVolume(),
			Silent:   false,
		}
		set.paused = false
		set.interV = nil
		fmt.Printf("playing track %s\n", set.next)
		set.curr = set.next
		if set.mode != Repeat {
			set.next = ""
		}
		speaker.Unlock()
		speaker.Play(beep.Seq(
			beep.Resample(4, format.SampleRate, sampleRate, set.volume),
			beep.Callback(func() {
				set.NextTrack()
			}),
		))
		return nil
	} else {
		set.next = ""
		return errors.Wrap(fmt.Errorf("key %s is not a registered track", set.next), errMsg)
	}
}

func (p *musicPlayer) stopAllMusic() {
	speaker.Clear()
	for _, s := range p.streams {
		if s.stream != nil {
			s.stream.Close()
		}
		s.ctrl = nil
		s.volume = nil
		s.interV = nil
		s.paused = true
	}
}

func (p *musicPlayer) PauseAllMusic() {
	for _, s := range p.streams {
		s.pause(true)
	}
}

func (p *musicPlayer) StopAllMusic() {
	for _, s := range p.streams {
		s.Stop()
	}
}

func (p *musicPlayer) HasTrack(key string) bool {
	_, ok := p.tracks[key]
	return ok
}
