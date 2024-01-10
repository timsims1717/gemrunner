package reanimator

import (
	"gemrunner/pkg/img"
	"gemrunner/pkg/util"
	"github.com/gopxl/pixel"
	"image/color"
	"time"
)

var (
	Timer       time.Time
	FRate       int
	inter       float64
	FrameSwitch bool
)

type Tree struct {
	Root    *Switch
	spr     *pixel.Sprite
	anim    *Anim
	animKey string
	frame   int
	update  bool
	Done    bool
	Default string
}

func SetFrameRate(fRate int) {
	FRate = fRate
	inter = 1. / float64(fRate)
}

func Reset() {
	Timer = time.Now()
}

func Update() {
	FrameSwitch = time.Since(Timer).Seconds() > inter
	if FrameSwitch {
		Reset()
	}
}

func NewSimple(anim *Anim) *Tree {
	t := &Tree{
		Root: NewSwitch().
			AddAnimation(anim).
			SetChooseFn(func() int {
				return 0
			}),
	}
	t.Update()
	return t
}

func New(root *Switch, def string) *Tree {
	t := &Tree{
		Root:    root,
		update:  true,
		Default: def,
	}
	t.Update()
	return t
}

func (t *Tree) ForceUpdate() {
	t.update = true
}

func (t *Tree) Update() {
	if !t.Done {
		if FrameSwitch || t.update {
			t.update = false
			t.anim = t.Root.choose()
			if t.anim == nil {
				t.spr = nil
				t.animKey = ""
				t.frame = 0
			} else {
				pKey := t.animKey
				pFrame := t.frame
				var trigger int
				if t.anim.Key != t.animKey {
					t.anim.Step = 0
					trigger = 0
				} else {
					t.anim.Step++
					trigger = t.anim.Step
					if t.anim.Step%len(t.anim.S) == 0 {
						switch t.anim.Finish {
						case Loop:
							t.anim.Step = 0
							trigger = 0
						case Hold:
							t.anim.Step = len(t.anim.S) - 1
						case Tran:
							t.anim.Step = len(t.anim.S) - 1
							t.update = true
						case Done:
							t.anim.Step = len(t.anim.S) - 1
							t.Done = true
						}
					}
				}
				if t.anim.Triggers != nil {
					if fn, ok := t.anim.Triggers[trigger]; ok {
						fn(t.anim, pKey, pFrame)
					}
				}
				t.spr = t.anim.S[t.anim.Step]
				if t.update {
					t.animKey = t.Default
					t.frame = t.anim.Step
				} else {
					t.animKey = t.anim.Key
					t.frame = t.anim.Step
				}
			}
		}
	}
}

type Result struct {
	Spr   *pixel.Sprite
	Off   pixel.Vec
	Col   pixel.RGBA
	Batch string
}

func (t *Tree) CurrentSprite() *Result {
	if t.spr == nil {
		return nil
	}
	return &Result{
		Spr:   t.spr,
		Off:   t.anim.Offset,
		Col:   t.anim.Color,
		Batch: t.anim.Batch,
	}
}

func (t *Tree) Draw(target pixel.Target, mat pixel.Matrix) {
	if t.spr != nil {
		t.spr.Draw(target, mat)
	}
}

func (t *Tree) DrawColorMask(target pixel.Target, mat pixel.Matrix, col color.RGBA) {
	if t.spr != nil {
		t.spr.DrawColorMask(target, mat, col)
	}
}

type switchEl struct {
	Switch *Switch
	Anim   *Anim
}

type Switch struct {
	Elements []*switchEl
	Choose   func() int
}

func NewSwitch() *Switch {
	return &Switch{}
}

func (s *Switch) AddNull() *Switch {
	s.Elements = append(s.Elements, &switchEl{})
	return s
}

func (s *Switch) AddAnimation(anim *Anim) *Switch {
	s.Elements = append(s.Elements, &switchEl{
		Anim: anim,
	})
	return s
}

func (s *Switch) AddSubSwitch(ss *Switch) *Switch {
	s.Elements = append(s.Elements, &switchEl{
		Switch: ss,
	})
	return s
}

func (s *Switch) SetChooseFn(fn func() int) *Switch {
	s.Choose = fn
	return s
}

func (s *Switch) choose() *Anim {
	el := s.Elements[s.Choose()]
	if el.Switch != nil {
		return el.Switch.choose()
	} else if el.Anim != nil {
		return el.Anim
	} else {
		return nil
	}
}

type Anim struct {
	Key      string
	S        []*pixel.Sprite
	Step     int
	Finish   Finish
	Triggers map[int]func(*Anim, string, int)

	Offset pixel.Vec
	Color  pixel.RGBA
	Batch  string
}

type Finish int

const (
	Loop = iota
	Hold
	Tran
	Done
)

func (a *Anim) WithColor(col pixel.RGBA) *Anim {
	a.Color = col
	return a
}

func (a *Anim) WithBatch(batch string) *Anim {
	a.Batch = batch
	return a
}

func (a *Anim) WithOffset(offset pixel.Vec) *Anim {
	a.Offset = offset
	return a
}

func NewAnimFromSprite(key string, spr *pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      []*pixel.Sprite{spr},
		Step:   0,
		Finish: f,
		Color:  util.White,
	}
}

func NewAnimFromSprites(key string, spr []*pixel.Sprite, f Finish) *Anim {
	return &Anim{
		Key:    key,
		S:      spr,
		Step:   0,
		Finish: f,
		Color:  util.White,
	}
}

func NewAnimFromSheet(key string, spriteSheet *img.SpriteSheet, rs []int, f Finish) *Anim {
	var spr []*pixel.Sprite
	if len(rs) > 0 {
		for _, r := range rs {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, spriteSheet.Sprites[r]))
		}
	} else {
		for _, s := range spriteSheet.Sprites {
			spr = append(spr, pixel.NewSprite(spriteSheet.Img, s))
		}
	}
	return &Anim{
		Key:    key,
		S:      spr,
		Step:   0,
		Finish: f,
		Color:  util.White,
	}
}

func (anim *Anim) SetTriggerC(i int, fn func(*Anim, string, int)) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[i] = fn
	return anim
}

func (anim *Anim) SetTrigger(i int, fn func()) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	anim.Triggers[i] = func(*Anim, string, int) {
		fn()
	}
	return anim
}

func (anim *Anim) SetTriggerAll(fn func()) *Anim {
	if anim.Triggers == nil {
		anim.Triggers = map[int]func(*Anim, string, int){}
	}
	for i := range anim.S {
		anim.SetTrigger(i, fn)
	}
	anim.SetTrigger(len(anim.S), fn)
	return anim
}

func (anim *Anim) Copy() *Anim {
	return &Anim{
		Key:      anim.Key,
		S:        anim.S,
		Step:     anim.Step,
		Finish:   anim.Finish,
		Triggers: anim.Triggers,
		Color:    anim.Color,
		Batch:    anim.Batch,
		Offset:   anim.Offset,
	}
}
